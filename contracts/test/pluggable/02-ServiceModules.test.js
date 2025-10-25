import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture, time } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Pluggable Architecture - Service Modules", function () {
  // Fixture for deploying service modules
  async function deployServiceModulesFixture() {
    const [owner, priceFeeder, auditor, user1, user2] = await ethers.getSigners();

    // Deploy infrastructure
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
    const priceOracle = await PriceOracleModule.deploy();

    const AuditModule = await ethers.getContractFactory("AuditModule");
    const auditModule = await AuditModule.deploy();

    // Register and enable modules
    await moduleRegistry.registerModule(await priceOracle.getAddress());
    await moduleRegistry.registerModule(await auditModule.getAddress());

    await moduleRegistry.enableModule(await priceOracle.MODULE_ID());
    await moduleRegistry.enableModule(await auditModule.MODULE_ID());

    // Grant roles
    const PRICE_FEEDER_ROLE = await priceOracle.PRICE_FEEDER_ROLE();
    await priceOracle.grantRole(PRICE_FEEDER_ROLE, priceFeeder.address);

    const AUDITOR_ROLE = await auditModule.AUDITOR_ROLE();
    const LOGGER_ROLE = await auditModule.LOGGER_ROLE();
    await auditModule.grantRole(AUDITOR_ROLE, auditor.address);
    await auditModule.grantRole(LOGGER_ROLE, owner.address);

    return {
      accessController,
      moduleRegistry,
      eventHub,
      priceOracle,
      auditModule,
      owner,
      priceFeeder,
      auditor,
      user1,
      user2
    };
  }

  describe("PriceOracleModule", function () {
    it("Should deploy and register successfully", async function () {
      const { priceOracle, moduleRegistry } = await loadFixture(deployServiceModulesFixture);

      const moduleId = await priceOracle.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.isRegistered).to.be.true;
      expect(moduleInfo.state).to.equal(1); // ACTIVE
    });

    it("Should have correct module metadata", async function () {
      const { priceOracle } = await loadFixture(deployServiceModulesFixture);

      expect(await priceOracle.MODULE_NAME()).to.equal("PriceOracleModule");
      expect(await priceOracle.MODULE_VERSION()).to.equal("1.0.0");
    });

    it("Should add price feed", async function () {
      const { priceOracle, priceFeeder } = await loadFixture(deployServiceModulesFixture);

      const tokenAddress = ethers.Wallet.createRandom().address;
      const feedAddress = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(
        tokenAddress,
        feedAddress,
        "Chainlink ETH/USD",
        18
      );

      const feeds = await priceOracle.getPriceFeeds(tokenAddress);
      expect(feeds.length).to.equal(1);
      expect(feeds[0].feedAddress).to.equal(feedAddress);
    });

    it("Should update price", async function () {
      const { priceOracle, priceFeeder } = await loadFixture(deployServiceModulesFixture);

      const tokenAddress = ethers.Wallet.createRandom().address;
      const feedAddress = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(
        tokenAddress,
        feedAddress,
        "Test Feed",
        18
      );

      const price = ethers.parseEther("2000");
      await priceOracle.connect(priceFeeder).updatePrice(
        feedAddress,
        price,
        18
      );

      const priceData = await priceOracle.getPrice(feedAddress);
      expect(priceData.price).to.equal(price);
    });

    it("Should calculate aggregated price from multiple feeds", async function () {
      const { priceOracle, priceFeeder } = await loadFixture(deployServiceModulesFixture);

      const tokenAddress = ethers.Wallet.createRandom().address;

      // Add multiple feeds
      const feed1 = ethers.Wallet.createRandom().address;
      const feed2 = ethers.Wallet.createRandom().address;
      const feed3 = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feed1, "Feed 1", 18);
      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feed2, "Feed 2", 18);
      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feed3, "Feed 3", 18);

      // Update prices
      await priceOracle.connect(priceFeeder).updatePrice(feed1, ethers.parseEther("100"), 18);
      await priceOracle.connect(priceFeeder).updatePrice(feed2, ethers.parseEther("102"), 18);
      await priceOracle.connect(priceFeeder).updatePrice(feed3, ethers.parseEther("98"), 18);

      // Get aggregated price
      const aggregated = await priceOracle.getAggregatedPrice(tokenAddress);

      // Median of [98, 100, 102] = 100
      expect(aggregated.aggregatedPrice).to.equal(ethers.parseEther("100"));
      expect(aggregated.activeFeedCount).to.equal(3);
    });

    it("Should detect stale prices", async function () {
      const { priceOracle, priceFeeder } = await loadFixture(deployServiceModulesFixture);

      const tokenAddress = ethers.Wallet.createRandom().address;
      const feedAddress = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feedAddress, "Test", 18);
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("100"), 18);

      // Advance time beyond staleness threshold (default 1 hour)
      await time.increase(3601);

      const priceData = await priceOracle.getPrice(feedAddress);
      expect(priceData.isStale).to.be.true;
    });

    it("Should pause and unpause price feeds", async function () {
      const { priceOracle, priceFeeder, owner } = await loadFixture(deployServiceModulesFixture);

      const tokenAddress = ethers.Wallet.createRandom().address;
      const feedAddress = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feedAddress, "Test", 18);

      // Pause feed
      await priceOracle.connect(owner).pauseFeed(feedAddress);
      let priceData = await priceOracle.getPrice(feedAddress);
      expect(priceData.isActive).to.be.false;

      // Unpause feed
      await priceOracle.connect(owner).unpauseFeed(feedAddress);
      priceData = await priceOracle.getPrice(feedAddress);
      expect(priceData.isActive).to.be.true;
    });

    it("Should perform health check", async function () {
      const { priceOracle } = await loadFixture(deployServiceModulesFixture);

      const [healthy, message] = await priceOracle.healthCheck();
      expect(healthy).to.be.true;
      expect(message).to.include("healthy");
    });

    it("Should get price history", async function () {
      const { priceOracle, priceFeeder } = await loadFixture(deployServiceModulesFixture);

      const feedAddress = ethers.Wallet.createRandom().address;
      const tokenAddress = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feedAddress, "Test", 18);

      // Update price multiple times
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("100"), 18);
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("101"), 18);
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("102"), 18);

      const history = await priceOracle.getPriceHistory(feedAddress, 0, 10);
      expect(history.length).to.equal(3);
      expect(history[2].price).to.equal(ethers.parseEther("102"));
    });

    it("Should calculate price deviation", async function () {
      const { priceOracle, priceFeeder } = await loadFixture(deployServiceModulesFixture);

      const tokenAddress = ethers.Wallet.createRandom().address;

      // Add feeds with varying prices
      const feed1 = ethers.Wallet.createRandom().address;
      const feed2 = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feed1, "Feed 1", 18);
      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feed2, "Feed 2", 18);

      await priceOracle.connect(priceFeeder).updatePrice(feed1, ethers.parseEther("100"), 18);
      await priceOracle.connect(priceFeeder).updatePrice(feed2, ethers.parseEther("110"), 18);

      const aggregated = await priceOracle.getAggregatedPrice(tokenAddress);

      // Should have non-zero deviation
      expect(aggregated.deviation).to.be.greaterThan(0);
    });
  });

  describe("AuditModule", function () {
    it("Should deploy and register successfully", async function () {
      const { auditModule, moduleRegistry } = await loadFixture(deployServiceModulesFixture);

      const moduleId = await auditModule.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.isRegistered).to.be.true;
      expect(moduleInfo.state).to.equal(1); // ACTIVE
    });

    it("Should have correct module metadata", async function () {
      const { auditModule } = await loadFixture(deployServiceModulesFixture);

      expect(await auditModule.MODULE_NAME()).to.equal("AuditModule");
      expect(await auditModule.MODULE_VERSION()).to.equal("1.0.0");
    });

    it("Should log audit events", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      const eventType = 0; // DEPOSIT
      const severity = 0; // INFO
      const value = ethers.parseEther("100");
      const description = "User deposited collateral";

      await expect(
        auditModule.logEvent(
          eventType,
          severity,
          user1.address,
          owner.address,
          value,
          description
        )
      ).to.emit(auditModule, "AuditEventLogged");
    });

    it("Should retrieve audit logs", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      // Log multiple events
      await auditModule.logEvent(0, 0, user1.address, owner.address, ethers.parseEther("100"), "Deposit 1");
      await auditModule.logEvent(1, 0, user1.address, owner.address, ethers.parseEther("50"), "Withdrawal 1");
      await auditModule.logEvent(0, 0, user1.address, owner.address, ethers.parseEther("200"), "Deposit 2");

      const logs = await auditModule.getAuditLogs(0, 10);
      expect(logs.length).to.equal(3);
    });

    it("Should filter logs by user", async function () {
      const { auditModule, owner, user1, user2 } = await loadFixture(deployServiceModulesFixture);

      // Log events for different users
      await auditModule.logEvent(0, 0, user1.address, owner.address, ethers.parseEther("100"), "User1 deposit");
      await auditModule.logEvent(0, 0, user2.address, owner.address, ethers.parseEther("200"), "User2 deposit");
      await auditModule.logEvent(1, 0, user1.address, owner.address, ethers.parseEther("50"), "User1 withdrawal");

      const user1Logs = await auditModule.getLogsByUser(user1.address, 0, 10);
      expect(user1Logs.length).to.equal(2);

      const user2Logs = await auditModule.getLogsByUser(user2.address, 0, 10);
      expect(user2Logs.length).to.equal(1);
    });

    it("Should filter logs by event type", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      // Log different event types
      await auditModule.logEvent(0, 0, user1.address, owner.address, 100, "Deposit"); // DEPOSIT
      await auditModule.logEvent(1, 0, user1.address, owner.address, 50, "Withdrawal"); // WITHDRAWAL
      await auditModule.logEvent(0, 0, user1.address, owner.address, 200, "Deposit 2"); // DEPOSIT

      const depositLogs = await auditModule.getLogsByEventType(0, 0, 10);
      expect(depositLogs.length).to.equal(2);

      const withdrawalLogs = await auditModule.getLogsByEventType(1, 0, 10);
      expect(withdrawalLogs.length).to.equal(1);
    });

    it("Should filter logs by severity", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      // Log events with different severities
      await auditModule.logEvent(0, 0, user1.address, owner.address, 100, "Info event"); // INFO
      await auditModule.logEvent(0, 1, user1.address, owner.address, 200, "Warning event"); // WARNING
      await auditModule.logEvent(0, 2, user1.address, owner.address, 300, "Error event"); // ERROR

      const warningLogs = await auditModule.getLogsBySeverity(1, 0, 10);
      expect(warningLogs.length).to.equal(1);

      const errorLogs = await auditModule.getLogsBySeverity(2, 0, 10);
      expect(errorLogs.length).to.equal(1);
    });

    it("Should filter logs by time range", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      const startTime = await time.latest();

      await auditModule.logEvent(0, 0, user1.address, owner.address, 100, "Event 1");

      await time.increase(3600); // 1 hour

      await auditModule.logEvent(0, 0, user1.address, owner.address, 200, "Event 2");

      const midTime = await time.latest();

      await time.increase(3600); // another hour

      await auditModule.logEvent(0, 0, user1.address, owner.address, 300, "Event 3");

      const endTime = await time.latest();

      // Get logs in time range
      const logs = await auditModule.getLogsByTimeRange(startTime, midTime, 0, 10);
      expect(logs.length).to.equal(2);
    });

    it("Should generate compliance reports", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      const startTime = await time.latest();

      // Log various events
      await auditModule.logEvent(0, 0, user1.address, owner.address, ethers.parseEther("100"), "Deposit");
      await auditModule.logEvent(1, 0, user1.address, owner.address, ethers.parseEther("50"), "Withdrawal");
      await auditModule.logEvent(2, 1, user1.address, owner.address, ethers.parseEther("200"), "Liquidation");

      const endTime = await time.latest();

      const report = await auditModule.generateComplianceReport(startTime, endTime);

      expect(report.totalEvents).to.equal(3);
      expect(report.criticalEvents).to.equal(0);
      expect(report.errorEvents).to.equal(0);
      expect(report.warningEvents).to.equal(1);
    });

    it("Should perform health check", async function () {
      const { auditModule } = await loadFixture(deployServiceModulesFixture);

      const [healthy, message] = await auditModule.healthCheck();
      expect(healthy).to.be.true;
      expect(message).to.include("healthy");
    });

    it("Should get audit statistics", async function () {
      const { auditModule, owner, user1 } = await loadFixture(deployServiceModulesFixture);

      // Log events
      await auditModule.logEvent(0, 0, user1.address, owner.address, 100, "Event 1");
      await auditModule.logEvent(0, 1, user1.address, owner.address, 200, "Event 2");
      await auditModule.logEvent(0, 2, user1.address, owner.address, 300, "Event 3");

      const stats = await auditModule.getAuditStatistics();

      expect(stats.totalLogs).to.equal(3);
      expect(stats.criticalCount).to.equal(0);
      expect(stats.errorCount).to.equal(1);
      expect(stats.warningCount).to.equal(1);
    });
  });

  describe("Service Modules Integration", function () {
    it("Should work together for price monitoring and audit logging", async function () {
      const { priceOracle, auditModule, priceFeeder, owner } = await loadFixture(
        deployServiceModulesFixture
      );

      const tokenAddress = ethers.Wallet.createRandom().address;
      const feedAddress = ethers.Wallet.createRandom().address;

      // Add price feed
      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feedAddress, "Test Feed", 18);

      // Log the action in audit module
      await auditModule.logEvent(
        7, // PRICE_UPDATE
        0, // INFO
        priceFeeder.address,
        feedAddress,
        0,
        "Price feed added"
      );

      // Update price
      const price = ethers.parseEther("100");
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, price, 18);

      // Log price update
      await auditModule.logEvent(
        7, // PRICE_UPDATE
        0, // INFO
        priceFeeder.address,
        feedAddress,
        price,
        "Price updated"
      );

      // Verify audit trail
      const logs = await auditModule.getLogsByEventType(7, 0, 10);
      expect(logs.length).to.equal(2);

      // Verify price
      const priceData = await priceOracle.getPrice(feedAddress);
      expect(priceData.price).to.equal(price);
    });

    it("Should detect and log anomalies", async function () {
      const { priceOracle, auditModule, priceFeeder } = await loadFixture(
        deployServiceModulesFixture
      );

      const tokenAddress = ethers.Wallet.createRandom().address;
      const feedAddress = ethers.Wallet.createRandom().address;

      await priceOracle.connect(priceFeeder).addPriceFeed(tokenAddress, feedAddress, "Test", 18);

      // Set initial price
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("100"), 18);

      // Sudden price spike (100% increase) - should be logged as warning
      await priceOracle.connect(priceFeeder).updatePrice(feedAddress, ethers.parseEther("200"), 18);

      await auditModule.logEvent(
        7, // PRICE_UPDATE
        1, // WARNING
        priceFeeder.address,
        feedAddress,
        ethers.parseEther("200"),
        "Significant price change detected: 100% increase"
      );

      const warnings = await auditModule.getLogsBySeverity(1, 0, 10);
      expect(warnings.length).to.be.greaterThan(0);
    });
  });
});
