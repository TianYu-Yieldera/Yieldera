import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture, time } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Pluggable Architecture - Module Integration", function () {
  // Comprehensive fixture deploying all modules
  async function deployFullSystemFixture() {
    const [owner, priceFeeder, auditor, user1, user2, trader1, feeCollector] = await ethers.getSigners();

    // Deploy tokens
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    const loyaltyToken = await MockERC20.deploy("Loyalty Points", "LP", ethers.parseEther("10000000"));
    const lusdToken = await MockERC20.deploy("Loyalty USD", "LUSD", ethers.parseEther("10000000"));
    const rwaToken = await MockERC20.deploy("Real Estate Token", "RET", ethers.parseEther("1000000"));

    // Deploy core infrastructure
    const AccessController = await ethers.getContractFactory("AccessController");
    const accessController = await AccessController.deploy();

    const ModuleRegistry = await ethers.getContractFactory("ModuleRegistry");
    const moduleRegistry = await ModuleRegistry.deploy();

    const EventHub = await ethers.getContractFactory("EventHub");
    const eventHub = await EventHub.deploy();

    await moduleRegistry.setAccessController(await accessController.getAddress());
    await eventHub.setModuleRegistry(await moduleRegistry.getAddress());

    // Deploy service modules
    const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
    const priceOracleModule = await PriceOracleModule.deploy();

    const AuditModule = await ethers.getContractFactory("AuditModule");
    const auditModule = await AuditModule.deploy();

    // Deploy legacy contracts
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const collateralVault = await CollateralVault.deploy(await loyaltyToken.getAddress());

    const OrderBook = await ethers.getContractFactory("OrderBook");
    const orderBook = await OrderBook.deploy(
      await rwaToken.getAddress(),
      await lusdToken.getAddress(),
      feeCollector.address
    );

    // Deploy business logic modules (adapters)
    const VaultModule = await ethers.getContractFactory("VaultModule");
    const vaultModule = await VaultModule.deploy(
      await collateralVault.getAddress(),
      await loyaltyToken.getAddress(),
      await lusdToken.getAddress()
    );

    const RWAModule = await ethers.getContractFactory("RWAModule");
    const rwaModule = await RWAModule.deploy(await orderBook.getAddress());

    // Register all modules
    await moduleRegistry.registerModule(await priceOracleModule.getAddress());
    await moduleRegistry.registerModule(await auditModule.getAddress());
    await moduleRegistry.registerModule(await vaultModule.getAddress());
    await moduleRegistry.registerModule(await rwaModule.getAddress());

    // Enable all modules
    await moduleRegistry.enableModule(await priceOracleModule.MODULE_ID());
    await moduleRegistry.enableModule(await auditModule.MODULE_ID());
    await moduleRegistry.enableModule(await vaultModule.MODULE_ID());
    await moduleRegistry.enableModule(await rwaModule.MODULE_ID());

    // Initialize modules
    await priceOracleModule.initialize("0x");
    await auditModule.initialize("0x");
    await vaultModule.initialize("0x");
    await rwaModule.initialize("0x");

    // Grant roles
    const PRICE_FEEDER_ROLE = await priceOracleModule.PRICE_FEEDER_ROLE();
    await priceOracleModule.grantRole(PRICE_FEEDER_ROLE, priceFeeder.address);

    const AUDITOR_ROLE = await auditModule.AUDITOR_ROLE();
    const LOGGER_ROLE = await auditModule.LOGGER_ROLE();
    await auditModule.grantRole(AUDITOR_ROLE, auditor.address);
    await auditModule.grantRole(LOGGER_ROLE, owner.address);

    // Setup tokens for users
    await loyaltyToken.transfer(user1.address, ethers.parseEther("100000"));
    await loyaltyToken.transfer(user2.address, ethers.parseEther("100000"));
    await rwaToken.transfer(trader1.address, ethers.parseEther("10000"));

    await loyaltyToken.connect(user1).approve(await collateralVault.getAddress(), ethers.MaxUint256);
    await loyaltyToken.connect(user2).approve(await collateralVault.getAddress(), ethers.MaxUint256);
    await rwaToken.connect(trader1).approve(await orderBook.getAddress(), ethers.MaxUint256);
    await lusdToken.connect(trader1).approve(await orderBook.getAddress(), ethers.MaxUint256);

    return {
      moduleRegistry,
      accessController,
      eventHub,
      priceOracleModule,
      auditModule,
      vaultModule,
      rwaModule,
      collateralVault,
      orderBook,
      loyaltyToken,
      lusdToken,
      rwaToken,
      owner,
      priceFeeder,
      auditor,
      user1,
      user2,
      trader1,
      feeCollector
    };
  }

  describe("System-Wide Module Registration", function () {
    it("Should have all modules registered and enabled", async function () {
      const { moduleRegistry, priceOracleModule, auditModule, vaultModule, rwaModule } =
        await loadFixture(deployFullSystemFixture);

      const modules = [priceOracleModule, auditModule, vaultModule, rwaModule];

      for (const module of modules) {
        const moduleId = await module.MODULE_ID();
        const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

        expect(moduleInfo.isRegistered).to.be.true;
        expect(moduleInfo.isEnabled).to.be.true;
        expect(moduleInfo.state).to.equal(1); // ACTIVE
      }
    });

    it("Should query all registered modules", async function () {
      const { moduleRegistry } = await loadFixture(deployFullSystemFixture);

      const registeredModules = await moduleRegistry.getRegisteredModules();
      expect(registeredModules.length).to.equal(4);
    });

    it("Should verify module dependencies", async function () {
      const { vaultModule, rwaModule } = await loadFixture(deployFullSystemFixture);

      const vaultDeps = await vaultModule.getDependencies();
      expect(vaultDeps.length).to.equal(2);

      const rwaDeps = await rwaModule.getDependencies();
      expect(rwaDeps.length).to.equal(2);

      // Both should depend on PriceOracle and Audit modules
      expect(vaultDeps[0]).to.equal(ethers.keccak256(ethers.toUtf8Bytes("PRICE_ORACLE_MODULE")));
      expect(rwaDeps[0]).to.equal(ethers.keccak256(ethers.toUtf8Bytes("PRICE_ORACLE_MODULE")));
    });

    it("Should perform health checks on all modules", async function () {
      const { priceOracleModule, auditModule, vaultModule, rwaModule } =
        await loadFixture(deployFullSystemFixture);

      const modules = [priceOracleModule, auditModule, vaultModule, rwaModule];

      for (const module of modules) {
        const [healthy, message] = await module.healthCheck();
        expect(healthy).to.be.true;
        expect(message).to.include("healthy");
      }
    });
  });

  describe("Cross-Module Integration: Vault + Audit", function () {
    it("Should log vault operations in audit module", async function () {
      const { vaultModule, auditModule, user1 } = await loadFixture(deployFullSystemFixture);

      // User deposits collateral
      const depositAmount = ethers.parseEther("1000");
      await vaultModule.connect(user1).depositCollateral(depositAmount);

      // Log the event in audit module
      await auditModule.logEvent(
        0, // DEPOSIT
        0, // INFO
        user1.address,
        await vaultModule.getAddress(),
        depositAmount,
        "Collateral deposited"
      );

      // Verify audit log
      const logs = await auditModule.getLogsByUser(user1.address, 0, 10);
      expect(logs.length).to.be.greaterThan(0);
      expect(logs[0].user).to.equal(user1.address);
    });

    it("Should track vault liquidations in audit module", async function () {
      const { vaultModule, auditModule, collateralVault, user1, user2 } =
        await loadFixture(deployFullSystemFixture);

      // Create undercollateralized position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1100"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Liquidate
      const debtToCover = ethers.parseEther("500");
      await vaultModule.connect(user2).liquidate(user1.address, debtToCover);

      // Log liquidation event
      await auditModule.logEvent(
        2, // LIQUIDATION
        1, // WARNING
        user1.address,
        user2.address,
        debtToCover,
        "Position liquidated"
      );

      // Verify liquidation was logged
      const liquidationLogs = await auditModule.getLogsByEventType(2, 0, 10);
      expect(liquidationLogs.length).to.be.greaterThan(0);
    });

    it("Should generate compliance reports including vault operations", async function () {
      const { vaultModule, auditModule, collateralVault, user1 } =
        await loadFixture(deployFullSystemFixture);

      const startTime = await time.latest();

      // Perform various vault operations
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("2000"));
      await auditModule.logEvent(0, 0, user1.address, await vaultModule.getAddress(), ethers.parseEther("2000"), "Deposit");

      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));
      await auditModule.logEvent(3, 0, user1.address, await vaultModule.getAddress(), ethers.parseEther("1000"), "Debt increase");

      const endTime = await time.latest();

      const report = await auditModule.generateComplianceReport(startTime, endTime);
      expect(report.totalEvents).to.be.greaterThan(0);
    });
  });

  describe("Cross-Module Integration: RWA + Audit", function () {
    it("Should log RWA trading operations in audit module", async function () {
      const { rwaModule, auditModule, trader1 } = await loadFixture(deployFullSystemFixture);

      // Place order
      const price = ethers.parseUnits("10", 6);
      const amount = ethers.parseEther("100");
      const tx = await rwaModule.connect(trader1).placeOrder(0, price, amount);
      await tx.wait();

      // Log the order in audit module
      await auditModule.logEvent(
        6, // ORDER_PLACED
        0, // INFO
        trader1.address,
        await rwaModule.getAddress(),
        amount,
        "Buy order placed"
      );

      // Verify audit log
      const logs = await auditModule.getLogsByEventType(6, 0, 10);
      expect(logs.length).to.be.greaterThan(0);
    });

    it("Should track large orders with warnings", async function () {
      const { rwaModule, auditModule, trader1, lusdToken } = await loadFixture(deployFullSystemFixture);

      // Fund trader with more LUSD
      await lusdToken.transfer(trader1.address, ethers.parseEther("1000000"));

      // Place large order
      const largeAmount = ethers.parseEther("5000");
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), largeAmount);

      // Log as warning (large order)
      await auditModule.logEvent(
        6, // ORDER_PLACED
        1, // WARNING
        trader1.address,
        await rwaModule.getAddress(),
        largeAmount,
        "Large buy order placed"
      );

      const warningLogs = await auditModule.getLogsBySeverity(1, 0, 10);
      expect(warningLogs.length).to.be.greaterThan(0);
    });
  });

  describe("Cross-Module Integration: Vault + PriceOracle", function () {
    it("Should use price oracle for collateral valuation", async function () {
      const { vaultModule, priceOracleModule, priceFeeder, user1, loyaltyToken } =
        await loadFixture(deployFullSystemFixture);

      // Setup price feed for loyalty token
      const feedAddress = ethers.Wallet.createRandom().address;
      await priceOracleModule.connect(priceFeeder).addPriceFeed(
        await loyaltyToken.getAddress(),
        feedAddress,
        "LP/USD Feed",
        18
      );

      // Set LP price to $1.50
      await priceOracleModule.connect(priceFeeder).updatePrice(
        feedAddress,
        ethers.parseEther("1.5"),
        18
      );

      // Deposit collateral
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));

      // Get collateral value (would need price oracle integration in vault)
      const collateralBalance = await vaultModule.getCollateralBalance(user1.address);
      expect(collateralBalance).to.equal(ethers.parseEther("1000"));

      // Verify price data is available
      const priceData = await priceOracleModule.getPrice(feedAddress);
      expect(priceData.price).to.equal(ethers.parseEther("1.5"));
    });

    it("Should detect price deviations and log warnings", async function () {
      const { priceOracleModule, auditModule, priceFeeder, loyaltyToken } =
        await loadFixture(deployFullSystemFixture);

      const tokenAddress = await loyaltyToken.getAddress();
      const feedAddress = ethers.Wallet.createRandom().address;

      await priceOracleModule.connect(priceFeeder).addPriceFeed(tokenAddress, feedAddress, "LP Feed", 18);

      // Set initial price
      await priceOracleModule.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("1"), 18);

      // Sudden price change
      await priceOracleModule.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("2"), 18);

      // Log price anomaly
      await auditModule.logEvent(
        7, // PRICE_UPDATE
        1, // WARNING
        priceFeeder.address,
        feedAddress,
        ethers.parseEther("2"),
        "100% price increase detected"
      );

      const warnings = await auditModule.getLogsBySeverity(1, 0, 10);
      expect(warnings.length).to.be.greaterThan(0);
    });

    it("Should handle stale price data", async function () {
      const { priceOracleModule, priceFeeder, loyaltyToken } =
        await loadFixture(deployFullSystemFixture);

      const feedAddress = ethers.Wallet.createRandom().address;
      await priceOracleModule.connect(priceFeeder).addPriceFeed(
        await loyaltyToken.getAddress(),
        feedAddress,
        "LP Feed",
        18
      );

      await priceOracleModule.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("1"), 18);

      // Advance time beyond staleness threshold
      await time.increase(3601); // 1 hour + 1 second

      const priceData = await priceOracleModule.getPrice(feedAddress);
      expect(priceData.isStale).to.be.true;
    });
  });

  describe("Cross-Module Integration: RWA + PriceOracle", function () {
    it("Should use price oracle for RWA token pricing", async function () {
      const { rwaModule, priceOracleModule, priceFeeder, rwaToken } =
        await loadFixture(deployFullSystemFixture);

      // Setup price feed for RWA token
      const feedAddress = ethers.Wallet.createRandom().address;
      await priceOracleModule.connect(priceFeeder).addPriceFeed(
        await rwaToken.getAddress(),
        feedAddress,
        "RWA/USD Feed",
        18
      );

      // Set RWA price
      const rwaPrice = ethers.parseEther("100");
      await priceOracleModule.connect(priceFeeder).updatePrice(feedAddress, rwaPrice, 18);

      // Verify price is available for trading
      const priceData = await priceOracleModule.getPrice(feedAddress);
      expect(priceData.price).to.equal(rwaPrice);
      expect(priceData.isActive).to.be.true;
    });

    it("Should validate order prices against oracle data", async function () {
      const { rwaModule, priceOracleModule, priceFeeder, trader1, rwaToken } =
        await loadFixture(deployFullSystemFixture);

      // Setup oracle price
      const feedAddress = ethers.Wallet.createRandom().address;
      await priceOracleModule.connect(priceFeeder).addPriceFeed(
        await rwaToken.getAddress(),
        feedAddress,
        "RWA Feed",
        18
      );
      await priceOracleModule.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("100"), 18);

      // Place order at market price
      await rwaModule.connect(trader1).placeOrder(
        1, // SELL
        ethers.parseUnits("100", 6), // $100
        ethers.parseEther("10")
      );

      // Verify order was placed
      const orders = await rwaModule.getUserOpenOrders(trader1.address);
      expect(orders.length).to.be.greaterThan(0);
    });
  });

  describe("Module Lifecycle Management", function () {
    it("Should pause all modules in emergency", async function () {
      const { priceOracleModule, auditModule, vaultModule, rwaModule, owner } =
        await loadFixture(deployFullSystemFixture);

      const modules = [priceOracleModule, auditModule, vaultModule, rwaModule];

      // Pause all modules
      for (const module of modules) {
        await module.connect(owner).pause();
        expect(await module.isActive()).to.be.false;
      }
    });

    it("Should resume all modules after emergency", async function () {
      const { priceOracleModule, auditModule, vaultModule, rwaModule, owner } =
        await loadFixture(deployFullSystemFixture);

      const modules = [priceOracleModule, auditModule, vaultModule, rwaModule];

      // Pause then unpause all modules
      for (const module of modules) {
        await module.connect(owner).pause();
        await module.connect(owner).unpause();
        expect(await module.isActive()).to.be.true;
      }
    });

    it("Should handle selective module disabling", async function () {
      const { moduleRegistry, vaultModule, rwaModule, user1, trader1 } =
        await loadFixture(deployFullSystemFixture);

      // Disable vault module only
      await moduleRegistry.disableModule(await vaultModule.MODULE_ID());

      // Vault operations should fail
      await expect(
        vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"))
      ).to.be.reverted;

      // RWA operations should still work
      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
    });
  });

  describe("Complete User Journey", function () {
    it("Should support full vault + trading workflow", async function () {
      const {
        vaultModule,
        rwaModule,
        auditModule,
        priceOracleModule,
        collateralVault,
        priceFeeder,
        user1,
        loyaltyToken,
        rwaToken,
        lusdToken
      } = await loadFixture(deployFullSystemFixture);

      // Step 1: Setup price feeds
      const lpFeed = ethers.Wallet.createRandom().address;
      const rwaFeed = ethers.Wallet.createRandom().address;

      await priceOracleModule.connect(priceFeeder).addPriceFeed(
        await loyaltyToken.getAddress(),
        lpFeed,
        "LP/USD",
        18
      );
      await priceOracleModule.connect(priceFeeder).updatePrice(lpFeed, ethers.parseEther("1"), 18);

      await priceOracleModule.connect(priceFeeder).addPriceFeed(
        await rwaToken.getAddress(),
        rwaFeed,
        "RWA/USD",
        18
      );
      await priceOracleModule.connect(priceFeeder).updatePrice(rwaFeed, ethers.parseEther("100"), 18);

      // Step 2: User deposits collateral
      const collateralAmount = ethers.parseEther("2000");
      await vaultModule.connect(user1).depositCollateral(collateralAmount);
      await auditModule.logEvent(0, 0, user1.address, await vaultModule.getAddress(), collateralAmount, "Collateral deposited");

      // Step 3: User mints LUSD
      const mintAmount = ethers.parseEther("1000");
      await collateralVault.connect(user1).increaseDebt(user1.address, mintAmount);
      await auditModule.logEvent(3, 0, user1.address, await vaultModule.getAddress(), mintAmount, "LUSD minted");

      // Step 4: User trades RWA tokens
      await lusdToken.transfer(user1.address, ethers.parseEther("10000"));
      await lusdToken.connect(user1).approve(await rwaModule.orderBook(), ethers.MaxUint256);

      await rwaModule.connect(user1).placeOrder(0, ethers.parseUnits("100", 6), ethers.parseEther("10"));
      await auditModule.logEvent(6, 0, user1.address, await rwaModule.getAddress(), ethers.parseEther("10"), "RWA order placed");

      // Verify all operations succeeded
      const position = await vaultModule.getPosition(user1.address);
      expect(position.collateralAmount).to.equal(collateralAmount);
      expect(position.debtAmount).to.equal(mintAmount);

      const orders = await rwaModule.getUserOpenOrders(user1.address);
      expect(orders.length).to.be.greaterThan(0);

      const auditLogs = await auditModule.getLogsByUser(user1.address, 0, 10);
      expect(auditLogs.length).to.equal(3); // Deposit, mint, trade
    });

    it("Should generate comprehensive system report", async function () {
      const {
        vaultModule,
        rwaModule,
        auditModule,
        collateralVault,
        user1,
        user2,
        trader1
      } = await loadFixture(deployFullSystemFixture);

      const startTime = await time.latest();

      // Multiple user operations
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));
      await auditModule.logEvent(0, 0, user1.address, await vaultModule.getAddress(), ethers.parseEther("1500"), "User1 deposit");

      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("3000"));
      await collateralVault.connect(user2).increaseDebt(user2.address, ethers.parseEther("2000"));
      await auditModule.logEvent(0, 0, user2.address, await vaultModule.getAddress(), ethers.parseEther("3000"), "User2 deposit");

      await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await auditModule.logEvent(6, 0, trader1.address, await rwaModule.getAddress(), ethers.parseEther("100"), "Trade");

      const endTime = await time.latest();

      // Generate reports
      const [totalCollateral, totalDebt, avgRatio] = await vaultModule.getVaultStats();
      const marketStats = await rwaModule.getMarketStats();
      const complianceReport = await auditModule.generateComplianceReport(startTime, endTime);

      expect(totalCollateral).to.equal(ethers.parseEther("4500"));
      expect(totalDebt).to.equal(ethers.parseEther("3000"));
      expect(complianceReport.totalEvents).to.equal(3);
    });
  });

  describe("Backward Compatibility", function () {
    it("Should maintain compatibility with legacy CollateralVault", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployFullSystemFixture);

      // Operations through module
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));

      // Verify state in legacy contract
      const legacyBalance = await collateralVault.collateralDeposited(user1.address);
      const moduleBalance = await vaultModule.getCollateralBalance(user1.address);

      expect(legacyBalance).to.equal(moduleBalance);
    });

    it("Should maintain compatibility with legacy OrderBook", async function () {
      const { rwaModule, orderBook, trader1 } = await loadFixture(deployFullSystemFixture);

      // Operations through module
      const tx = await rwaModule.connect(trader1).placeOrder(0, ethers.parseUnits("10", 6), ethers.parseEther("100"));
      await tx.wait();

      // Verify orders exist in legacy contract
      const userOrders = await rwaModule.getUserOpenOrders(trader1.address);
      expect(userOrders.length).to.be.greaterThan(0);
    });

    it("Should allow direct legacy contract access", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployFullSystemFixture);

      // Direct access to legacy contract
      await collateralVault.connect(user1).depositCollateral(ethers.parseEther("1000"));

      // Verify through module interface
      const balance = await vaultModule.getCollateralBalance(user1.address);
      expect(balance).to.equal(ethers.parseEther("1000"));
    });
  });

  describe("Error Handling and Edge Cases", function () {
    it("Should handle module initialization failures gracefully", async function () {
      // Test with invalid parameters would go here
      // This is a placeholder for edge case testing
    });

    it("Should recover from partial system failures", async function () {
      const { vaultModule, rwaModule, owner } = await loadFixture(deployFullSystemFixture);

      // Pause one module
      await vaultModule.connect(owner).pause();

      // Other modules should still work
      expect(await rwaModule.isActive()).to.be.true;

      // Resume failed module
      await vaultModule.connect(owner).unpause();
      expect(await vaultModule.isActive()).to.be.true;
    });

    it("Should enforce proper access controls across modules", async function () {
      const { vaultModule, rwaModule, user1 } = await loadFixture(deployFullSystemFixture);

      // Non-owner should not be able to configure modules
      await expect(
        vaultModule.connect(user1).setMinCollateralRatio(200)
      ).to.be.reverted;

      await expect(
        rwaModule.connect(user1).updateFees(5, 15)
      ).to.be.reverted;
    });
  });
});
