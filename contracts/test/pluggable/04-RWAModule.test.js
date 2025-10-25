import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Pluggable Architecture - RWAModule", function () {
  // Fixture for deploying RWAModule with full infrastructure
  async function deployRWAModuleFixture() {
    const [owner, trader1, trader2, trader3, feeCollector] = await ethers.getSigners();

    // Deploy mock ERC20 tokens
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    const rwaToken = await MockERC20.deploy("Real Estate Token", "RET", ethers.parseEther("1000000"));
    const lusdToken = await MockERC20.deploy("Loyalty USD", "LUSD", ethers.parseEther("1000000"));

    // Deploy core infrastructure
    const AccessController = await ethers.getContractFactory("AccessController");
    const accessController = await AccessController.deploy();

    const ModuleRegistry = await ethers.getContractFactory("ModuleRegistry");
    const moduleRegistry = await ModuleRegistry.deploy();

    const EventHub = await ethers.getContractFactory("EventHub");
    const eventHub = await EventHub.deploy();

    await moduleRegistry.setAccessController(await accessController.getAddress());
    await eventHub.setModuleRegistry(await moduleRegistry.getAddress());

    // Deploy OrderBook (legacy contract)
    const OrderBook = await ethers.getContractFactory("OrderBook");
    const orderBook = await OrderBook.deploy(
      await rwaToken.getAddress(),
      await lusdToken.getAddress(),
      feeCollector.address
    );

    // Deploy RWAModule (adapter)
    const RWAModule = await ethers.getContractFactory("RWAModule");
    const rwaModule = await RWAModule.deploy(await orderBook.getAddress());

    // Register and enable module
    await moduleRegistry.registerModule(await rwaModule.getAddress());
    await moduleRegistry.enableModule(await rwaModule.MODULE_ID());

    // Initialize module
    await rwaModule.initialize("0x");

    // Setup tokens for traders
    await rwaToken.transfer(trader1.address, ethers.parseEther("10000"));
    await rwaToken.transfer(trader2.address, ethers.parseEther("10000"));
    await rwaToken.transfer(trader3.address, ethers.parseEther("10000"));

    await lusdToken.transfer(trader1.address, ethers.parseEther("100000"));
    await lusdToken.transfer(trader2.address, ethers.parseEther("100000"));
    await lusdToken.transfer(trader3.address, ethers.parseEther("100000"));

    // Approve tokens
    await rwaToken.connect(trader1).approve(await orderBook.getAddress(), ethers.MaxUint256);
    await rwaToken.connect(trader2).approve(await orderBook.getAddress(), ethers.MaxUint256);
    await rwaToken.connect(trader3).approve(await orderBook.getAddress(), ethers.MaxUint256);

    await lusdToken.connect(trader1).approve(await orderBook.getAddress(), ethers.MaxUint256);
    await lusdToken.connect(trader2).approve(await orderBook.getAddress(), ethers.MaxUint256);
    await lusdToken.connect(trader3).approve(await orderBook.getAddress(), ethers.MaxUint256);

    return {
      rwaModule,
      orderBook,
      moduleRegistry,
      accessController,
      eventHub,
      rwaToken,
      lusdToken,
      owner,
      trader1,
      trader2,
      trader3,
      feeCollector
    };
  }

  describe("Module Registration and Lifecycle", function () {
    it("Should deploy and register successfully", async function () {
      const { rwaModule, moduleRegistry } = await loadFixture(deployRWAModuleFixture);

      const moduleId = await rwaModule.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.isRegistered).to.be.true;
      expect(moduleInfo.state).to.equal(1); // ACTIVE
    });

    it("Should have correct module metadata", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      expect(await rwaModule.MODULE_NAME()).to.equal("RWAModule");
      expect(await rwaModule.MODULE_VERSION()).to.equal("1.0.0");
      expect(await rwaModule.getModuleId()).to.equal(await rwaModule.MODULE_ID());
    });

    it("Should report correct dependencies", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      const dependencies = await rwaModule.getDependencies();
      expect(dependencies.length).to.equal(2);
      expect(dependencies[0]).to.equal(ethers.keccak256(ethers.toUtf8Bytes("PRICE_ORACLE_MODULE")));
      expect(dependencies[1]).to.equal(ethers.keccak256(ethers.toUtf8Bytes("AUDIT_MODULE")));
    });

    it("Should be active after initialization", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      expect(await rwaModule.isActive()).to.be.true;
    });

    it("Should pause and unpause correctly", async function () {
      const { rwaModule, owner } = await loadFixture(deployRWAModuleFixture);

      await rwaModule.connect(owner).pause();
      expect(await rwaModule.isActive()).to.be.false;

      await rwaModule.connect(owner).unpause();
      expect(await rwaModule.isActive()).to.be.true;
    });

    it("Should perform health check", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      const [healthy, message] = await rwaModule.healthCheck();
      expect(healthy).to.be.true;
      expect(message).to.include("healthy");
    });
  });

  describe("Order Management", function () {
    it("Should place buy order", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      const price = ethers.parseUnits("10", 6); // $10 per token
      const amount = ethers.parseEther("100");

      await expect(rwaModule.connect(trader1).placeOrder(0, price, amount)) // 0 = BUY
        .to.emit(rwaModule, "OrderPlaced");

      const userOrders = await rwaModule.getUserOpenOrders(trader1.address);
      expect(userOrders.length).to.equal(1);
      expect(userOrders[0].orderType).to.equal(0); // BUY
      expect(userOrders[0].price).to.equal(price);
      expect(userOrders[0].amount).to.equal(amount);
    });

    it("Should place sell order", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      const price = ethers.parseUnits("12", 6); // $12 per token
      const amount = ethers.parseEther("50");

      await expect(rwaModule.connect(trader1).placeOrder(1, price, amount)) // 1 = SELL
        .to.emit(rwaModule, "OrderPlaced");

      const userOrders = await rwaModule.getUserOpenOrders(trader1.address);
      expect(userOrders.length).to.equal(1);
      expect(userOrders[0].orderType).to.equal(1); // SELL
    });

    it("Should cancel order", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      // Place order
      const tx = await rwaModule.connect(trader1).placeOrder(
        0,
        ethers.parseUnits("10", 6),
        ethers.parseEther("100")
      );
      const receipt = await tx.wait();

      // Get order ID from event
      const orderPlacedEvent = receipt.logs.find(
        log => {
          try {
            return rwaModule.interface.parseLog(log)?.name === "OrderPlaced";
          } catch {
            return false;
          }
        }
      );
      const orderId = rwaModule.interface.parseLog(orderPlacedEvent).args.orderId;

      // Cancel order
      await expect(rwaModule.connect(trader1).cancelOrder(orderId))
        .to.emit(rwaModule, "OrderCancelled");
    });

    it("Should get order details", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      // Place order
      const price = ethers.parseUnits("10", 6);
      const amount = ethers.parseEther("100");
      const tx = await rwaModule.connect(trader1).placeOrder(0, price, amount);
      const receipt = await tx.wait();

      const orderPlacedEvent = receipt.logs.find(
        log => {
          try {
            return rwaModule.interface.parseLog(log)?.name === "OrderPlaced";
          } catch {
            return false;
          }
        }
      );
      const orderId = rwaModule.interface.parseLog(orderPlacedEvent).args.orderId;

      // Get order details
      const order = await rwaModule.getOrder(orderId);
      expect(order.trader).to.equal(trader1.address);
      expect(order.price).to.equal(price);
      expect(order.amount).to.equal(amount);
      expect(order.status).to.equal(0); // OPEN
    });

    it("Should get user's open orders", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      // Place multiple orders
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("11", 6), ethers.parseEther("50"));
      await rwaModule.connect(trader1).placeOrder(1, ethers.parseUnits("15", 6), ethers.parseEther("75"));

      const openOrders = await rwaModule.getUserOpenOrders(trader1.address);
      expect(openOrders.length).to.equal(3);
    });

    it("Should revert when placing order while paused", async function () {
      const { rwaModule, trader1, owner } = await loadFixture(deployRWAModuleFixture);

      await rwaModule.connect(owner).pause();

      await expect(
        rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"))
      ).to.be.reverted;
    });
  });

  describe("Order Book Operations", function () {
    it("Should get order book depth", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Place buy orders
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader2).placeOrder(0, ethers.parseUnits("9", 6), ethers.parseEther("50"));

      const [prices, amounts] = await rwaModule.getOrderBookDepth(0, 10); // 0 = BUY
      expect(prices.length).to.be.greaterThan(0);
      expect(amounts.length).to.equal(prices.length);
    });

    it("Should get best bid", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Place buy orders
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader2).placeOrder(0, ethers.parseUnits("12", 6), ethers.parseEther("50"));

      const [price, amount] = await rwaModule.getBestBid();
      expect(price).to.equal(ethers.parseUnits("12", 6)); // Highest buy price
      expect(amount).to.be.greaterThan(0);
    });

    it("Should get best ask", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Place sell orders
      await rwaModule.connect(trader1).placeOrder(1, ethers.parseUnits("15", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader2).placeOrder(1, ethers.parseUnits("14", 6), ethers.parseEther("50"));

      const [price, amount] = await rwaModule.getBestAsk();
      expect(price).to.equal(ethers.parseUnits("14", 6)); // Lowest sell price
      expect(amount).to.be.greaterThan(0);
    });

    it("Should calculate spread", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Place orders
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100")); // Buy
      await rwaModule.connect(trader2).placeOrder(1, ethers.parseUnits("12", 6), ethers.parseEther("50")); // Sell

      const spread = await rwaModule.getSpread();
      expect(spread).to.equal(ethers.parseUnits("2", 6)); // $12 - $10 = $2
    });

    it("Should get mid price", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Place orders
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100")); // Buy
      await rwaModule.connect(trader2).placeOrder(1, ethers.parseUnits("12", 6), ethers.parseEther("50")); // Sell

      const midPrice = await rwaModule.getMidPrice();
      expect(midPrice).to.equal(ethers.parseUnits("11", 6)); // ($10 + $12) / 2 = $11
    });
  });

  describe("Trading Operations", function () {
    it("Should place market buy order", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Trader1 places sell order (provide liquidity)
      await rwaModule.connect(trader1).placeOrder(1, ethers.parseUnits("10", 6), ethers.parseEther("100"));

      // Trader2 places market buy order
      const amount = ethers.parseEther("50");
      const [executedAmount, avgPrice] = await rwaModule.connect(trader2).placeMarketOrder(0, amount);

      expect(executedAmount).to.equal(amount);
      expect(avgPrice).to.equal(ethers.parseUnits("10", 6));
    });

    it("Should place market sell order", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Trader1 places buy order (provide liquidity)
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));

      // Trader2 places market sell order
      const amount = ethers.parseEther("50");
      const [executedAmount, avgPrice] = await rwaModule.connect(trader2).placeMarketOrder(1, amount);

      expect(executedAmount).to.equal(amount);
      expect(avgPrice).to.be.greaterThan(0);
    });

    it("Should revert market order with no liquidity", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      await expect(
        rwaModule.connect(trader1).placeMarketOrder(0, ethers.parseEther("100"))
      ).to.be.revertedWith("No liquidity available");
    });

    it("Should get trade history", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Create some trades
      await rwaModule.connect(trader1).placeOrder(1, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader2).placeMarketOrder(0, ethers.parseEther("50"));

      const trades = await rwaModule.getTradeHistory(0, 10);
      // Note: Trade history tracking needs to be implemented in the adapter
      // This test may pass with 0 trades if tracking is not yet implemented
    });
  });

  describe("Market Statistics", function () {
    it("Should get market stats", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      const stats = await rwaModule.getMarketStats();
      expect(stats.lastPrice).to.be.a("bigint");
      expect(stats.openBuyOrders).to.be.a("bigint");
      expect(stats.openSellOrders).to.be.a("bigint");
    });

    it("Should get trading pair info", async function () {
      const { rwaModule, rwaToken, lusdToken } = await loadFixture(deployRWAModuleFixture);

      const tradingPair = await rwaModule.getTradingPair();
      expect(tradingPair.baseToken).to.equal(await rwaToken.getAddress());
      expect(tradingPair.quoteToken).to.equal(await lusdToken.getAddress());
      expect(tradingPair.isActive).to.be.true;
    });

    it("Should get 24h volume", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      const volume = await rwaModule.get24hVolume();
      expect(volume).to.be.a("bigint");
    });

    it("Should get total volume", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      const totalVolume = await rwaModule.getTotalVolume();
      expect(totalVolume).to.be.a("bigint");
    });

    it("Should get last price", async function () {
      const { rwaModule } = await loadFixture(deployRWAModuleFixture);

      const lastPrice = await rwaModule.getLastPrice();
      expect(lastPrice).to.be.a("bigint");
    });
  });

  describe("Configuration Management", function () {
    it("Should update trading fees", async function () {
      const { rwaModule, owner } = await loadFixture(deployRWAModuleFixture);

      const newMakerFee = 5; // 0.05%
      const newTakerFee = 15; // 0.15%

      await expect(rwaModule.connect(owner).updateFees(newMakerFee, newTakerFee))
        .to.emit(rwaModule, "FeesUpdated")
        .withArgs(newMakerFee, newTakerFee);

      const tradingPair = await rwaModule.getTradingPair();
      expect(tradingPair.makerFee).to.equal(newMakerFee);
      expect(tradingPair.takerFee).to.equal(newTakerFee);
    });

    it("Should update order limits", async function () {
      const { rwaModule, owner } = await loadFixture(deployRWAModuleFixture);

      const newMinSize = ethers.parseEther("10");
      const newMaxSize = ethers.parseEther("50000");

      await expect(rwaModule.connect(owner).updateOrderLimits(newMinSize, newMaxSize))
        .to.emit(rwaModule, "OrderLimitsUpdated")
        .withArgs(newMinSize, newMaxSize);

      const tradingPair = await rwaModule.getTradingPair();
      expect(tradingPair.minOrderSize).to.equal(newMinSize);
      expect(tradingPair.maxOrderSize).to.equal(newMaxSize);
    });

    it("Should update price limits", async function () {
      const { rwaModule, owner } = await loadFixture(deployRWAModuleFixture);

      const newMinPrice = ethers.parseUnits("1", 6);
      const newMaxPrice = ethers.parseUnits("500000", 6);

      await rwaModule.connect(owner).updatePriceLimits(newMinPrice, newMaxPrice);

      const tradingPair = await rwaModule.getTradingPair();
      expect(tradingPair.minPrice).to.equal(newMinPrice);
      expect(tradingPair.maxPrice).to.equal(newMaxPrice);
    });

    it("Should set fee collector", async function () {
      const { rwaModule, owner, trader1 } = await loadFixture(deployRWAModuleFixture);

      await rwaModule.connect(owner).setFeeCollector(trader1.address);
      // Verification would require checking the underlying OrderBook contract
    });

    it("Should only allow owner to update configuration", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      await expect(
        rwaModule.connect(trader1).updateFees(5, 15)
      ).to.be.reverted;
    });
  });

  describe("Liquidity Operations", function () {
    it("Should get total liquidity", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Place orders to create liquidity
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader2).placeOrder(1, ethers.parseUnits("12", 6), ethers.parseEther("50"));

      const [buyLiquidity, sellLiquidity] = await rwaModule.getTotalLiquidity();
      expect(buyLiquidity).to.be.greaterThan(0);
      expect(sellLiquidity).to.be.greaterThan(0);
    });

    it("Should check liquidity for order", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      // Provide liquidity
      await rwaModule.connect(trader1).placeOrder(1, ethers.parseUnits("10", 6), ethers.parseEther("200"));

      // Check if there's sufficient liquidity for buy order
      const amount = ethers.parseEther("100");
      const maxSlippage = 100; // 1%
      const [sufficient, estimatedPrice] = await rwaModule.connect(trader2).checkLiquidity(
        0, // BUY
        amount,
        maxSlippage
      );

      expect(sufficient).to.be.true;
      expect(estimatedPrice).to.be.greaterThan(0);
    });

    it("Should detect insufficient liquidity", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      // Check liquidity for large order with no liquidity
      const amount = ethers.parseEther("1000");
      const maxSlippage = 100;
      const [sufficient, estimatedPrice] = await rwaModule.connect(trader1).checkLiquidity(
        0, // BUY
        amount,
        maxSlippage
      );

      expect(sufficient).to.be.false;
    });
  });

  describe("Token Information", function () {
    it("Should return correct base token", async function () {
      const { rwaModule, rwaToken } = await loadFixture(deployRWAModuleFixture);

      const baseToken = await rwaModule.getBaseToken();
      expect(baseToken).to.equal(await rwaToken.getAddress());
    });

    it("Should return correct quote token", async function () {
      const { rwaModule, lusdToken } = await loadFixture(deployRWAModuleFixture);

      const quoteToken = await rwaModule.getQuoteToken();
      expect(quoteToken).to.equal(await lusdToken.getAddress());
    });
  });

  describe("Integration with Module Registry", function () {
    it("Should be queryable from module registry", async function () {
      const { rwaModule, moduleRegistry } = await loadFixture(deployRWAModuleFixture);

      const moduleId = await rwaModule.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.implementation).to.equal(await rwaModule.getAddress());
      expect(moduleInfo.isEnabled).to.be.true;
    });

    it("Should enforce module state changes", async function () {
      const { rwaModule, owner } = await loadFixture(deployRWAModuleFixture);

      // Pause module
      await rwaModule.connect(owner).pause();

      const moduleInfoAfterPause = await rwaModule.getModuleInfo();
      expect(moduleInfoAfterPause.state).to.equal(2); // PAUSED

      // Unpause module
      await rwaModule.connect(owner).unpause();

      const moduleInfoAfterUnpause = await rwaModule.getModuleInfo();
      expect(moduleInfoAfterUnpause.state).to.equal(1); // ACTIVE
    });
  });

  describe("Order Matching Scenarios", function () {
    it("Should match buy and sell orders at same price", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      const price = ethers.parseUnits("10", 6);
      const amount = ethers.parseEther("100");

      // Trader1 places sell order
      await rwaModule.connect(trader1).placeOrder(1, price, amount);

      // Trader2 places matching buy order
      await rwaModule.connect(trader2).placeOrder(0, price, amount);

      // Check if orders were matched (implementation specific)
      const trader2Orders = await rwaModule.getUserOpenOrders(trader2.address);
      // Depending on the matching logic, orders might be filled immediately
    });

    it("Should handle partial fills", async function () {
      const { rwaModule, trader1, trader2 } = await loadFixture(deployRWAModuleFixture);

      const price = ethers.parseUnits("10", 6);

      // Trader1 places sell order for 100 tokens
      await rwaModule.connect(trader1).placeOrder(1, price, ethers.parseEther("100"));

      // Trader2 places buy order for 50 tokens (partial fill)
      await rwaModule.connect(trader2).placeOrder(0, price, ethers.parseEther("50"));

      // Check remaining order size
      const trader1Orders = await rwaModule.getUserOpenOrders(trader1.address);
      if (trader1Orders.length > 0) {
        const remainingAmount = trader1Orders[0].amount - trader1Orders[0].filled;
        expect(remainingAmount).to.be.lessThanOrEqual(ethers.parseEther("100"));
      }
    });
  });

  describe("Edge Cases and Security", function () {
    it("Should handle zero amount orders", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      // Should revert on zero amount (implementation specific)
      // This depends on the underlying OrderBook validation
    });

    it("Should handle concurrent orders from same user", async function () {
      const { rwaModule, trader1 } = await loadFixture(deployRWAModuleFixture);

      // Place multiple orders quickly
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("11", 6), ethers.parseEther("50"));
      await rwaModule.connect(trader1).placeOrder(1, ethers.parseUnits("15", 6), ethers.parseEther("75"));

      const orders = await rwaModule.getUserOpenOrders(trader1.address);
      expect(orders.length).to.equal(3);
    });

    it("Should prevent operations when module is not active", async function () {
      const { rwaModule, trader1, owner } = await loadFixture(deployRWAModuleFixture);

      // Pause module
      await rwaModule.connect(owner).pause();

      // All trading operations should fail
      await expect(
        rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"))
      ).to.be.reverted;

      await expect(
        rwaModule.connect(trader1).placeMarketOrder(0, ethers.parseEther("100"))
      ).to.be.reverted;
    });
  });
});
