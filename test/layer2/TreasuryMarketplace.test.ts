import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  TreasuryMarketplace,
  TreasuryAssetFactory,
  TreasuryToken,
} from "../../typechain-types/index.js";

describe("TreasuryMarketplace", function () {
  let marketplace: TreasuryMarketplace;
  let factory: TreasuryAssetFactory;
  let token: TreasuryToken;

  let admin: SignerWithAddress;
  let seller: SignerWithAddress;
  let buyer: SignerWithAddress;
  let feeCollector: SignerWithAddress;

  const OPERATOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("OPERATOR_ROLE"));
  const MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MINTER_ROLE"));

  let assetId: number;
  const FACE_VALUE = ethers.parseEther("1000");
  const INITIAL_SUPPLY = ethers.parseEther("1000"); // 1000 tokens
  const SELL_AMOUNT = ethers.parseEther("100");
  const PRICE_PER_TOKEN = ethers.parseEther("995"); // $995 per token (below face value)

  beforeEach(async function () {
    [admin, seller, buyer, feeCollector] = await ethers.getSigners();

    // Deploy TreasuryAssetFactory
    const FactoryContract = await ethers.getContractFactory("TreasuryAssetFactory");
    factory = await FactoryContract.deploy(admin.address);

    // Deploy TreasuryMarketplace
    const MarketplaceFactory = await ethers.getContractFactory("TreasuryMarketplace");
    marketplace = await MarketplaceFactory.deploy(
      admin.address,
      await factory.getAddress(),
      feeCollector.address
    );

    // Grant roles
    await factory.grantRole(OPERATOR_ROLE, admin.address);
    await marketplace.grantRole(OPERATOR_ROLE, admin.address);

    // Create a Treasury asset
    const issueDate = Math.floor(Date.now() / 1000);
    const maturityDate = issueDate + 90 * 24 * 60 * 60;

    await factory.createTreasuryAsset(
      0, // T-BILL
      "13W",
      "912796YZ1",
      issueDate,
      maturityDate,
      FACE_VALUE,
      525
    );

    assetId = 1;

    // Mint tokens to seller
    await factory.mintTokens(assetId, seller.address, INITIAL_SUPPLY);

    // Get token contract
    const asset = await factory.getTreasuryAsset(assetId);
    const TreasuryToken = await ethers.getContractFactory("TreasuryToken");
    token = TreasuryToken.attach(asset.tokenAddress) as any;
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      const DEFAULT_ADMIN_ROLE = await marketplace.DEFAULT_ADMIN_ROLE();
      expect(await marketplace.hasRole(DEFAULT_ADMIN_ROLE, admin.address)).to.be.true;
    });

    it("Should set correct factory and fee collector", async function () {
      expect(await marketplace.treasuryFactory()).to.equal(await factory.getAddress());
      expect(await marketplace.feeCollector()).to.equal(feeCollector.address);
    });
  });

  describe("Sell Orders", function () {
    it("Should create a sell order", async function () {
      // Approve marketplace
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);

      const tx = await marketplace.connect(seller).createSellOrder(
        assetId,
        SELL_AMOUNT,
        PRICE_PER_TOKEN
      );

      await expect(tx).to.emit(marketplace, "SellOrderCreated");

      const orderId = 1;
      const order = await marketplace.getSellOrder(orderId);
      expect(order.seller).to.equal(seller.address);
      expect(order.amount).to.equal(SELL_AMOUNT);
      expect(order.pricePerToken).to.equal(PRICE_PER_TOKEN);
      expect(order.status).to.equal(0); // Active
    });

    it("Should fail if insufficient balance", async function () {
      const tooMuch = ethers.parseEther("2000");

      await token.connect(seller).approve(await marketplace.getAddress(), tooMuch);

      await expect(
        marketplace.connect(seller).createSellOrder(assetId, tooMuch, PRICE_PER_TOKEN)
      ).to.be.revertedWith("Insufficient balance");
    });

    it("Should fail if not approved", async function () {
      await expect(
        marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN)
      ).to.be.revertedWith("Insufficient allowance");
    });

    it("Should cancel a sell order", async function () {
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);

      const orderId = 1;
      await marketplace.connect(seller).cancelSellOrder(orderId);

      const order = await marketplace.getSellOrder(orderId);
      expect(order.status).to.equal(2); // Cancelled
    });

    it("Should prevent non-owner from cancelling", async function () {
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);

      const orderId = 1;
      await expect(
        marketplace.connect(buyer).cancelSellOrder(orderId)
      ).to.be.revertedWith("Not order owner");
    });
  });

  describe("Buy Orders", function () {
    it("Should create a buy order with payment", async function () {
      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");

      const tx = await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );

      await expect(tx).to.emit(marketplace, "BuyOrderCreated");

      const orderId = 1;
      const order = await marketplace.getBuyOrder(orderId);
      expect(order.buyer).to.equal(buyer.address);
      expect(order.amount).to.equal(buyAmount);
      expect(order.status).to.equal(0); // Active
    });

    it("Should fail with insufficient payment", async function () {
      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");
      const insufficientPayment = totalCost / 2n;

      await expect(
        marketplace.connect(buyer).createBuyOrder(
          assetId,
          buyAmount,
          PRICE_PER_TOKEN,
          { value: insufficientPayment }
        )
      ).to.be.revertedWith("Insufficient payment");
    });

    it("Should cancel buy order and refund", async function () {
      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");

      await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );

      const orderId = 1;
      const balanceBefore = await ethers.provider.getBalance(buyer.address);

      const tx = await marketplace.connect(buyer).cancelBuyOrder(orderId);
      const receipt = await tx.wait();
      const gasUsed = receipt!.gasUsed * receipt!.gasPrice;

      const balanceAfter = await ethers.provider.getBalance(buyer.address);
      expect(balanceAfter).to.be.closeTo(balanceBefore + totalCost - gasUsed, ethers.parseEther("0.001"));
    });
  });

  describe("Order Matching", function () {
    it("Should match buy and sell orders", async function () {
      // Create sell order
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);
      const sellOrderId = 1;

      // Create buy order
      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");
      await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );
      const buyOrderId = 1;

      // Match orders
      const tx = await marketplace.connect(admin).matchOrders(buyOrderId, sellOrderId, buyAmount);

      await expect(tx).to.emit(marketplace, "OrderMatched");

      // Verify token transfer
      expect(await token.balanceOf(buyer.address)).to.equal(buyAmount);
      expect(await token.balanceOf(seller.address)).to.equal(INITIAL_SUPPLY - buyAmount);
    });

    it("Should handle partial fills", async function () {
      // Create sell order for 100 tokens
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);
      const sellOrderId = 1;

      // Create buy order for 30 tokens
      const buyAmount = ethers.parseEther("30");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");
      await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );
      const buyOrderId = 1;

      // Match orders
      await marketplace.connect(admin).matchOrders(buyOrderId, sellOrderId, buyAmount);

      // Check sell order is partially filled
      const sellOrder = await marketplace.getSellOrder(sellOrderId);
      expect(sellOrder.filledAmount).to.equal(buyAmount);
      expect(sellOrder.status).to.equal(0); // Still active

      // Check buy order is fully filled
      const buyOrder = await marketplace.getBuyOrder(buyOrderId);
      expect(buyOrder.status).to.equal(1); // Filled
    });

    it("Should charge trading fees", async function () {
      // Set fee to 1% (100 basis points)
      await marketplace.connect(admin).setTradingFee(100);

      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);
      const sellOrderId = 1;

      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");
      await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );
      const buyOrderId = 1;

      const feeCollectorBalanceBefore = await ethers.provider.getBalance(feeCollector.address);

      await marketplace.connect(admin).matchOrders(buyOrderId, sellOrderId, buyAmount);

      const feeCollectorBalanceAfter = await ethers.provider.getBalance(feeCollector.address);
      const expectedFee = totalCost / 100n; // 1%

      expect(feeCollectorBalanceAfter - feeCollectorBalanceBefore).to.equal(expectedFee);
    });
  });

  describe("Order Book", function () {
    it("Should return active sell orders", async function () {
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);

      const orders = await marketplace.getActiveSellOrders(assetId);
      expect(orders.length).to.equal(1);
      expect(orders[0].seller).to.equal(seller.address);
    });

    it("Should return active buy orders", async function () {
      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");

      await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );

      const orders = await marketplace.getActiveBuyOrders(assetId);
      expect(orders.length).to.equal(1);
      expect(orders[0].buyer).to.equal(buyer.address);
    });
  });

  describe("Access Control", function () {
    it("Should allow only operator to match orders", async function () {
      await token.connect(seller).approve(await marketplace.getAddress(), SELL_AMOUNT);
      await marketplace.connect(seller).createSellOrder(assetId, SELL_AMOUNT, PRICE_PER_TOKEN);

      const buyAmount = ethers.parseEther("50");
      const totalCost = (buyAmount * PRICE_PER_TOKEN) / ethers.parseEther("1");
      await marketplace.connect(buyer).createBuyOrder(
        assetId,
        buyAmount,
        PRICE_PER_TOKEN,
        { value: totalCost }
      );

      await expect(
        marketplace.connect(buyer).matchOrders(1, 1, buyAmount)
      ).to.be.reverted;
    });
  });
});
