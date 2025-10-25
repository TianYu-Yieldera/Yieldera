import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Pluggable Architecture - Infrastructure", function () {
  // Fixture for deploying core infrastructure
  async function deployInfrastructureFixture() {
    const [owner, admin, user1, user2] = await ethers.getSigners();

    // Deploy AccessController
    const AccessController = await ethers.getContractFactory("AccessController");
    const accessController = await AccessController.deploy();

    // Deploy ModuleRegistry
    const ModuleRegistry = await ethers.getContractFactory("ModuleRegistry");
    const moduleRegistry = await ModuleRegistry.deploy();

    // Deploy EventHub
    const EventHub = await ethers.getContractFactory("EventHub");
    const eventHub = await EventHub.deploy();

    // Configure connections
    await moduleRegistry.setAccessController(await accessController.getAddress());
    await eventHub.setModuleRegistry(await moduleRegistry.getAddress());

    return {
      accessController,
      moduleRegistry,
      eventHub,
      owner,
      admin,
      user1,
      user2
    };
  }

  describe("AccessController", function () {
    it("Should deploy successfully", async function () {
      const { accessController } = await loadFixture(deployInfrastructureFixture);
      expect(await accessController.getAddress()).to.not.equal(ethers.ZeroAddress);
    });

    it("Should grant and check roles", async function () {
      const { accessController, owner, admin } = await loadFixture(deployInfrastructureFixture);

      const ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ADMIN_ROLE"));

      await accessController.createRole(ADMIN_ROLE, "Admin Role", "Administrative access");
      await accessController.grantRole(ADMIN_ROLE, admin.address);

      expect(await accessController.hasRole(ADMIN_ROLE, admin.address)).to.be.true;
      expect(await accessController.hasRole(ADMIN_ROLE, owner.address)).to.be.false;
    });

    it("Should create and manage permissions", async function () {
      const { accessController, owner } = await loadFixture(deployInfrastructureFixture);

      const MINT_PERMISSION = ethers.keccak256(ethers.toUtf8Bytes("MINT_PERMISSION"));
      const ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ADMIN_ROLE"));

      await accessController.createRole(ADMIN_ROLE, "Admin Role", "Administrative access");
      await accessController.createPermission(
        MINT_PERMISSION,
        "Mint Permission",
        "Permission to mint tokens",
        0 // No timelock
      );

      await accessController.grantPermissionToRole(MINT_PERMISSION, ADMIN_ROLE);

      expect(await accessController.roleHasPermission(ADMIN_ROLE, MINT_PERMISSION)).to.be.true;
    });

    it("Should enforce timelock on operations", async function () {
      const { accessController, admin } = await loadFixture(deployInfrastructureFixture);

      const CRITICAL_PERMISSION = ethers.keccak256(ethers.toUtf8Bytes("CRITICAL_PERMISSION"));
      const ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ADMIN_ROLE"));

      // Create permission with 1 hour timelock
      await accessController.createRole(ADMIN_ROLE, "Admin Role", "Admin");
      await accessController.grantRole(ADMIN_ROLE, admin.address);
      await accessController.createPermission(
        CRITICAL_PERMISSION,
        "Critical Permission",
        "Critical operation",
        3600 // 1 hour timelock
      );
      await accessController.grantPermissionToRole(CRITICAL_PERMISSION, ADMIN_ROLE);

      // Schedule operation
      const targetAddress = ethers.ZeroAddress;
      const callData = "0x";

      await accessController.connect(admin).scheduleOperation(
        CRITICAL_PERMISSION,
        targetAddress,
        callData
      );

      // Verify operation is scheduled
      const operationId = await accessController.computeOperationId(
        CRITICAL_PERMISSION,
        targetAddress,
        callData,
        await ethers.provider.getBlockNumber()
      );

      const operation = await accessController.getOperation(operationId);
      expect(operation.isScheduled).to.be.true;
    });

    it("Should support emergency pause", async function () {
      const { accessController, owner } = await loadFixture(deployInfrastructureFixture);

      await accessController.emergencyPause();
      expect(await accessController.paused()).to.be.true;

      await accessController.unpause();
      expect(await accessController.paused()).to.be.false;
    });
  });

  describe("ModuleRegistry", function () {
    it("Should deploy successfully", async function () {
      const { moduleRegistry } = await loadFixture(deployInfrastructureFixture);
      expect(await moduleRegistry.getAddress()).to.not.equal(ethers.ZeroAddress);
    });

    it("Should connect to AccessController", async function () {
      const { moduleRegistry, accessController } = await loadFixture(deployInfrastructureFixture);

      expect(await moduleRegistry.accessController()).to.equal(
        await accessController.getAddress()
      );
    });

    it("Should register a module", async function () {
      const { moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy a test module
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();

      await moduleRegistry.registerModule(await priceOracle.getAddress());

      const moduleId = await priceOracle.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.isRegistered).to.be.true;
      expect(moduleInfo.moduleAddress).to.equal(await priceOracle.getAddress());
    });

    it("Should validate module dependencies", async function () {
      const { moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy service modules (no dependencies)
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();

      const AuditModule = await ethers.getContractFactory("AuditModule");
      const auditModule = await AuditModule.deploy();

      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.registerModule(await auditModule.getAddress());

      const priceOracleId = await priceOracle.MODULE_ID();
      const auditModuleId = await auditModule.MODULE_ID();

      // Enable service modules
      await moduleRegistry.enableModule(priceOracleId);
      await moduleRegistry.enableModule(auditModuleId);

      // Now deploy a business module with dependencies
      const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
      const lusd = await LoyaltyUSD.deploy();
      const LPToken = await ethers.getContractFactory("LoyaltyUSD");
      const lpToken = await LPToken.deploy();
      const CollateralVault = await ethers.getContractFactory("CollateralVault");
      const vault = await CollateralVault.deploy(await lpToken.getAddress());

      const VaultModule = await ethers.getContractFactory("VaultModule");
      const vaultModule = await VaultModule.deploy(
        await vault.getAddress(),
        await lpToken.getAddress(),
        await lusd.getAddress()
      );

      await moduleRegistry.registerModule(await vaultModule.getAddress());

      const vaultModuleId = await vaultModule.MODULE_ID();

      // Dependencies should be satisfied
      expect(await moduleRegistry.validateDependencies(vaultModuleId)).to.be.true;
    });

    it("Should prevent enabling module with unsatisfied dependencies", async function () {
      const { moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy business module without deploying its dependencies
      const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
      const lusd = await LoyaltyUSD.deploy();
      const LPToken = await ethers.getContractFactory("LoyaltyUSD");
      const lpToken = await LPToken.deploy();
      const CollateralVault = await ethers.getContractFactory("CollateralVault");
      const vault = await CollateralVault.deploy(await lpToken.getAddress());

      const VaultModule = await ethers.getContractFactory("VaultModule");
      const vaultModule = await VaultModule.deploy(
        await vault.getAddress(),
        await lpToken.getAddress(),
        await lusd.getAddress()
      );

      await moduleRegistry.registerModule(await vaultModule.getAddress());

      const vaultModuleId = await vaultModule.MODULE_ID();

      // Should fail because dependencies are not registered/enabled
      await expect(
        moduleRegistry.enableModule(vaultModuleId)
      ).to.be.revertedWith("Dependencies not satisfied");
    });

    it("Should perform system health check", async function () {
      const { moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy and enable modules
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();

      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.enableModule(await priceOracle.MODULE_ID());

      const [healthyCount, totalCount, unhealthyIds] = await moduleRegistry.systemHealthCheck();

      expect(healthyCount).to.equal(1);
      expect(totalCount).to.equal(1);
      expect(unhealthyIds.length).to.equal(0);
    });

    it("Should track module states correctly", async function () {
      const { moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();
      const moduleId = await priceOracle.MODULE_ID();

      // Register module
      await moduleRegistry.registerModule(await priceOracle.getAddress());
      let moduleInfo = await moduleRegistry.getModuleInfo(moduleId);
      expect(moduleInfo.state).to.equal(0); // UNINITIALIZED

      // Enable module
      await moduleRegistry.enableModule(moduleId);
      moduleInfo = await moduleRegistry.getModuleInfo(moduleId);
      expect(moduleInfo.state).to.equal(1); // ACTIVE

      // Disable module
      await moduleRegistry.disableModule(moduleId);
      moduleInfo = await moduleRegistry.getModuleInfo(moduleId);
      expect(moduleInfo.state).to.equal(2); // PAUSED
    });
  });

  describe("EventHub", function () {
    it("Should deploy successfully", async function () {
      const { eventHub } = await loadFixture(deployInfrastructureFixture);
      expect(await eventHub.getAddress()).to.not.equal(ethers.ZeroAddress);
    });

    it("Should connect to ModuleRegistry", async function () {
      const { eventHub, moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      expect(await eventHub.moduleRegistry()).to.equal(
        await moduleRegistry.getAddress()
      );
    });

    it("Should publish events", async function () {
      const { eventHub, moduleRegistry, owner } = await loadFixture(deployInfrastructureFixture);

      // Deploy and register a module
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();
      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.enableModule(await priceOracle.MODULE_ID());

      // Publish event from module
      const eventType = "PRICE_UPDATE";
      const eventData = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address", "uint256"],
        [ethers.ZeroAddress, ethers.parseEther("100")]
      );

      await expect(
        eventHub.publishEvent(
          1, // MODULE category
          0, // INFO severity
          eventType,
          eventData
        )
      ).to.emit(eventHub, "EventPublished");
    });

    it("Should support event subscriptions", async function () {
      const { eventHub, moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy modules
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();
      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.enableModule(await priceOracle.MODULE_ID());

      const AuditModule = await ethers.getContractFactory("AuditModule");
      const auditModule = await AuditModule.deploy();
      await moduleRegistry.registerModule(await auditModule.getAddress());
      await moduleRegistry.enableModule(await auditModule.MODULE_ID());

      // Subscribe to events
      const priceOracleId = await priceOracle.MODULE_ID();
      const eventType = "PRICE_UPDATE";
      const callbackSelector = "0x12345678"; // Mock callback

      await eventHub.subscribe(
        priceOracleId,
        1, // MODULE category
        eventType,
        await auditModule.getAddress(),
        callbackSelector
      );

      // Verify subscription
      const subscriptionId = ethers.keccak256(
        ethers.AbiCoder.defaultAbiCoder().encode(
          ["bytes32", "uint8", "string", "address"],
          [priceOracleId, 1, eventType, await auditModule.getAddress()]
        )
      );

      const subscription = await eventHub.getSubscription(subscriptionId);
      expect(subscription.isActive).to.be.true;
    });

    it("Should filter events by category", async function () {
      const { eventHub, moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy and register module
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();
      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.enableModule(await priceOracle.MODULE_ID());

      // Publish events of different categories
      await eventHub.publishEvent(0, 0, "SYSTEM_START", "0x");
      await eventHub.publishEvent(1, 0, "MODULE_ENABLED", "0x");
      await eventHub.publishEvent(2, 0, "TRANSACTION_COMPLETE", "0x");

      // Get events by category
      const moduleEvents = await eventHub.getEventsByCategory(1, 0, 10);
      expect(moduleEvents.length).to.be.greaterThan(0);
    });

    it("Should filter events by severity", async function () {
      const { eventHub, moduleRegistry } = await loadFixture(deployInfrastructureFixture);

      // Deploy and register module
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();
      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.enableModule(await priceOracle.MODULE_ID());

      // Publish events of different severities
      await eventHub.publishEvent(0, 0, "INFO_EVENT", "0x"); // INFO
      await eventHub.publishEvent(0, 1, "WARNING_EVENT", "0x"); // WARNING
      await eventHub.publishEvent(0, 2, "ERROR_EVENT", "0x"); // ERROR

      // Get critical events
      const criticalEvents = await eventHub.getEventsBySeverity(3, 0, 10);
      expect(criticalEvents.length).to.equal(0); // No critical events published

      const warningEvents = await eventHub.getEventsBySeverity(1, 0, 10);
      expect(warningEvents.length).to.be.greaterThan(0);
    });
  });

  describe("Infrastructure Integration", function () {
    it("Should work together as a complete system", async function () {
      const {
        accessController,
        moduleRegistry,
        eventHub,
        admin,
        owner
      } = await loadFixture(deployInfrastructureFixture);

      // Create admin role
      const ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ADMIN_ROLE"));
      await accessController.createRole(ADMIN_ROLE, "Admin Role", "Admin");
      await accessController.grantRole(ADMIN_ROLE, admin.address);

      // Deploy modules
      const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
      const priceOracle = await PriceOracleModule.deploy();

      const AuditModule = await ethers.getContractFactory("AuditModule");
      const auditModule = await AuditModule.deploy();

      // Register modules
      await moduleRegistry.registerModule(await priceOracle.getAddress());
      await moduleRegistry.registerModule(await auditModule.getAddress());

      // Enable modules
      await moduleRegistry.enableModule(await priceOracle.MODULE_ID());
      await moduleRegistry.enableModule(await auditModule.MODULE_ID());

      // Publish event
      await eventHub.publishEvent(1, 0, "MODULE_SYSTEM_READY", "0x");

      // Verify system health
      const [healthyCount, totalCount] = await moduleRegistry.systemHealthCheck();
      expect(healthyCount).to.equal(2);
      expect(totalCount).to.equal(2);
    });
  });
});
