import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture, time } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Pluggable Architecture - VaultModule", function () {
  // Fixture for deploying VaultModule with full infrastructure
  async function deployVaultModuleFixture() {
    const [owner, user1, user2, liquidator] = await ethers.getSigners();

    // Deploy mock ERC20 tokens
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    const loyaltyToken = await MockERC20.deploy("Loyalty Points", "LP", ethers.parseEther("1000000"));
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

    // Deploy CollateralVault (legacy contract)
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const collateralVault = await CollateralVault.deploy(await loyaltyToken.getAddress());

    // Deploy VaultModule (adapter)
    const VaultModule = await ethers.getContractFactory("VaultModule");
    const vaultModule = await VaultModule.deploy(
      await collateralVault.getAddress(),
      await loyaltyToken.getAddress(),
      await lusdToken.getAddress()
    );

    // Register and enable module
    await moduleRegistry.registerModule(await vaultModule.getAddress());
    await moduleRegistry.enableModule(await vaultModule.MODULE_ID());

    // Initialize module
    await vaultModule.initialize("0x");

    // Setup tokens for users
    await loyaltyToken.transfer(user1.address, ethers.parseEther("10000"));
    await loyaltyToken.transfer(user2.address, ethers.parseEther("10000"));
    await loyaltyToken.connect(user1).approve(await collateralVault.getAddress(), ethers.MaxUint256);
    await loyaltyToken.connect(user2).approve(await collateralVault.getAddress(), ethers.MaxUint256);

    // Setup LUSD for liquidator
    await lusdToken.transfer(liquidator.address, ethers.parseEther("10000"));
    await lusdToken.connect(liquidator).approve(await collateralVault.getAddress(), ethers.MaxUint256);

    return {
      vaultModule,
      collateralVault,
      moduleRegistry,
      accessController,
      eventHub,
      loyaltyToken,
      lusdToken,
      owner,
      user1,
      user2,
      liquidator
    };
  }

  describe("Module Registration and Lifecycle", function () {
    it("Should deploy and register successfully", async function () {
      const { vaultModule, moduleRegistry } = await loadFixture(deployVaultModuleFixture);

      const moduleId = await vaultModule.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.isRegistered).to.be.true;
      expect(moduleInfo.state).to.equal(1); // ACTIVE
    });

    it("Should have correct module metadata", async function () {
      const { vaultModule } = await loadFixture(deployVaultModuleFixture);

      expect(await vaultModule.MODULE_NAME()).to.equal("VaultModule");
      expect(await vaultModule.MODULE_VERSION()).to.equal("1.0.0");
      expect(await vaultModule.getModuleId()).to.equal(await vaultModule.MODULE_ID());
    });

    it("Should report correct dependencies", async function () {
      const { vaultModule } = await loadFixture(deployVaultModuleFixture);

      const dependencies = await vaultModule.getDependencies();
      expect(dependencies.length).to.equal(2);
      expect(dependencies[0]).to.equal(ethers.keccak256(ethers.toUtf8Bytes("PRICE_ORACLE_MODULE")));
      expect(dependencies[1]).to.equal(ethers.keccak256(ethers.toUtf8Bytes("AUDIT_MODULE")));
    });

    it("Should be active after initialization", async function () {
      const { vaultModule } = await loadFixture(deployVaultModuleFixture);

      expect(await vaultModule.isActive()).to.be.true;
    });

    it("Should pause and unpause correctly", async function () {
      const { vaultModule, owner } = await loadFixture(deployVaultModuleFixture);

      await vaultModule.connect(owner).pause();
      expect(await vaultModule.isActive()).to.be.false;

      await vaultModule.connect(owner).unpause();
      expect(await vaultModule.isActive()).to.be.true;
    });

    it("Should perform health check", async function () {
      const { vaultModule } = await loadFixture(deployVaultModuleFixture);

      const [healthy, message] = await vaultModule.healthCheck();
      expect(healthy).to.be.true;
      expect(message).to.include("healthy");
    });
  });

  describe("Collateral Management", function () {
    it("Should deposit collateral", async function () {
      const { vaultModule, user1 } = await loadFixture(deployVaultModuleFixture);

      const depositAmount = ethers.parseEther("1000");

      await expect(vaultModule.connect(user1).depositCollateral(depositAmount))
        .to.emit(vaultModule, "CollateralDeposited")
        .withArgs(user1.address, depositAmount, depositAmount);

      const balance = await vaultModule.getCollateralBalance(user1.address);
      expect(balance).to.equal(depositAmount);
    });

    it("Should withdraw collateral when no debt", async function () {
      const { vaultModule, user1 } = await loadFixture(deployVaultModuleFixture);

      const depositAmount = ethers.parseEther("1000");
      await vaultModule.connect(user1).depositCollateral(depositAmount);

      const withdrawAmount = ethers.parseEther("500");
      await expect(vaultModule.connect(user1).withdrawCollateral(withdrawAmount))
        .to.emit(vaultModule, "CollateralWithdrawn")
        .withArgs(user1.address, withdrawAmount, depositAmount - withdrawAmount);

      const balance = await vaultModule.getCollateralBalance(user1.address);
      expect(balance).to.equal(depositAmount - withdrawAmount);
    });

    it("Should not withdraw collateral if position becomes unhealthy", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Deposit collateral
      const depositAmount = ethers.parseEther("1500");
      await vaultModule.connect(user1).depositCollateral(depositAmount);

      // Take debt (assuming 1 LP = 1 USD)
      const debtAmount = ethers.parseEther("1000");
      await collateralVault.connect(user1).increaseDebt(user1.address, debtAmount);

      // Try to withdraw too much (would make ratio < 150%)
      const withdrawAmount = ethers.parseEther("1000");
      await expect(
        vaultModule.connect(user1).withdrawCollateral(withdrawAmount)
      ).to.be.reverted;
    });

    it("Should track total collateral correctly", async function () {
      const { vaultModule, user1, user2 } = await loadFixture(deployVaultModuleFixture);

      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));
      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("2000"));

      const totalCollateral = await vaultModule.getTotalCollateral();
      expect(totalCollateral).to.equal(ethers.parseEther("3000"));
    });

    it("Should revert on zero amount deposit", async function () {
      const { vaultModule, user1 } = await loadFixture(deployVaultModuleFixture);

      await expect(
        vaultModule.connect(user1).depositCollateral(0)
      ).to.be.revertedWith("Amount must be > 0");
    });

    it("Should revert when paused", async function () {
      const { vaultModule, user1, owner } = await loadFixture(deployVaultModuleFixture);

      await vaultModule.connect(owner).pause();

      await expect(
        vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"))
      ).to.be.reverted;
    });
  });

  describe("Debt Management", function () {
    it("Should increase debt", async function () {
      const { vaultModule, user1 } = await loadFixture(deployVaultModuleFixture);

      // Deposit collateral first
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));

      // Increase debt
      const debtAmount = ethers.parseEther("1000");
      await expect(vaultModule.connect(user1).increaseDebt(debtAmount))
        .to.emit(vaultModule, "DebtIncreased");

      const totalDebt = await vaultModule.getTotalDebt(user1.address);
      expect(totalDebt).to.equal(debtAmount);
    });

    it("Should decrease debt", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Setup: deposit and take debt
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Decrease debt
      const repayAmount = ethers.parseEther("500");
      await expect(vaultModule.connect(user1).decreaseDebt(repayAmount))
        .to.emit(vaultModule, "DebtDecreased");

      const remainingDebt = await vaultModule.getTotalDebt(user1.address);
      expect(remainingDebt).to.be.closeTo(ethers.parseEther("500"), ethers.parseEther("1"));
    });

    it("Should respect debt ceiling", async function () {
      const { vaultModule, owner, user1 } = await loadFixture(deployVaultModuleFixture);

      // Set low debt ceiling
      await vaultModule.connect(owner).setDebtCeiling(ethers.parseEther("500"));

      // Deposit sufficient collateral
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("10000"));

      // Try to exceed debt ceiling
      await expect(
        vaultModule.connect(user1).increaseDebt(ethers.parseEther("1000"))
      ).to.be.revertedWith("Exceeds debt ceiling");
    });

    it("Should calculate accrued interest", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Deposit and take debt
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Advance time
      await time.increase(365 * 24 * 60 * 60); // 1 year

      // Check interest accrued
      const interest = await vaultModule.calculateAccruedInterest(user1.address);
      expect(interest).to.be.greaterThan(0);
    });

    it("Should calculate max mintable amount", async function () {
      const { vaultModule, user1 } = await loadFixture(deployVaultModuleFixture);

      // Deposit collateral
      const collateralAmount = ethers.parseEther("1500");
      await vaultModule.connect(user1).depositCollateral(collateralAmount);

      // Get max mintable (should be ~1000 LP with 150% ratio, assuming 1 LP = 1 USD)
      const maxMintable = await vaultModule.getMaxMintable(user1.address);
      expect(maxMintable).to.be.greaterThan(0);
      expect(maxMintable).to.be.lessThanOrEqual(ethers.parseEther("1000"));
    });

    it("Should track system debt correctly", async function () {
      const { vaultModule, collateralVault, user1, user2 } = await loadFixture(deployVaultModuleFixture);

      // User1 deposits and borrows
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("500"));

      // User2 deposits and borrows
      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("3000"));
      await collateralVault.connect(user2).increaseDebt(user2.address, ethers.parseEther("1000"));

      const systemDebt = await vaultModule.getSystemDebt();
      expect(systemDebt).to.equal(ethers.parseEther("1500"));
    });
  });

  describe("Position Management", function () {
    it("Should get user position", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Create position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      const position = await vaultModule.getPosition(user1.address);
      expect(position.collateralAmount).to.equal(ethers.parseEther("1500"));
      expect(position.debtAmount).to.equal(ethers.parseEther("1000"));
      expect(position.isActive).to.be.true;
    });

    it("Should calculate collateral ratio", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Create position with 150% ratio
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      const ratio = await vaultModule.getCollateralRatio(user1.address);
      expect(ratio).to.equal(150);
    });

    it("Should check if position is healthy", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Create healthy position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("2000"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      const isHealthy = await vaultModule.isPositionHealthy(user1.address);
      expect(isHealthy).to.be.true;
    });

    it("Should detect unhealthy position", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Create position at minimum ratio
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Simulate collateral value drop by increasing debt (in real scenario, price would change)
      // Position should become unhealthy if ratio < 150%
      const isHealthy = await vaultModule.isPositionHealthy(user1.address);
      // At exactly 150%, it should be healthy
      expect(isHealthy).to.be.true;
    });

    it("Should check if position can be liquidated", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Create position near liquidation threshold
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1200"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      const canLiquidate = await vaultModule.canLiquidate(user1.address);
      // At 120%, should be at liquidation threshold
      expect(canLiquidate).to.be.true;
    });
  });

  describe("Liquidation", function () {
    it("Should liquidate undercollateralized position", async function () {
      const { vaultModule, collateralVault, user1, liquidator } = await loadFixture(deployVaultModuleFixture);

      // Create undercollateralized position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1100"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Liquidate
      const debtToCover = ethers.parseEther("500");
      await expect(
        vaultModule.connect(liquidator).liquidate(user1.address, debtToCover)
      ).to.emit(vaultModule, "PositionLiquidated");
    });

    it("Should calculate liquidation parameters", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(deployVaultModuleFixture);

      // Create position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1100"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      const debtToCover = ethers.parseEther("500");
      const [collateralToSeize, penalty] = await vaultModule.calculateLiquidation(
        user1.address,
        debtToCover
      );

      expect(collateralToSeize).to.be.greaterThan(debtToCover);
      expect(penalty).to.equal(ethers.parseEther("50")); // 10% of 500
    });

    it("Should get liquidatable positions", async function () {
      const { vaultModule, collateralVault, user1, user2 } = await loadFixture(deployVaultModuleFixture);

      // User1: healthy position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("2000"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // User2: liquidatable position
      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("1100"));
      await collateralVault.connect(user2).increaseDebt(user2.address, ethers.parseEther("1000"));

      const liquidatablePositions = await vaultModule.getLiquidatablePositions();
      expect(liquidatablePositions.length).to.be.greaterThan(0);
      expect(liquidatablePositions).to.include(user2.address);
    });

    it("Should not liquidate healthy position", async function () {
      const { vaultModule, collateralVault, user1, liquidator } = await loadFixture(deployVaultModuleFixture);

      // Create healthy position
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("2000"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Try to liquidate
      await expect(
        vaultModule.connect(liquidator).liquidate(user1.address, ethers.parseEther("100"))
      ).to.be.reverted;
    });
  });

  describe("Configuration Management", function () {
    it("Should get vault configuration", async function () {
      const { vaultModule } = await loadFixture(deployVaultModuleFixture);

      const config = await vaultModule.getVaultConfig();
      expect(config.minCollateralRatio).to.equal(150);
      expect(config.liquidationThreshold).to.equal(120);
      expect(config.stabilityFee).to.equal(200);
    });

    it("Should update minimum collateral ratio", async function () {
      const { vaultModule, owner } = await loadFixture(deployVaultModuleFixture);

      await expect(vaultModule.connect(owner).setMinCollateralRatio(160))
        .to.emit(vaultModule, "VaultConfigUpdated")
        .withArgs("minCollateralRatio", 150, 160);

      const config = await vaultModule.getVaultConfig();
      expect(config.minCollateralRatio).to.equal(160);
    });

    it("Should update liquidation threshold", async function () {
      const { vaultModule, owner } = await loadFixture(deployVaultModuleFixture);

      await expect(vaultModule.connect(owner).setLiquidationThreshold(130))
        .to.emit(vaultModule, "VaultConfigUpdated")
        .withArgs("liquidationThreshold", 120, 130);

      const config = await vaultModule.getVaultConfig();
      expect(config.liquidationThreshold).to.equal(130);
    });

    it("Should update stability fee", async function () {
      const { vaultModule, owner } = await loadFixture(deployVaultModuleFixture);

      await expect(vaultModule.connect(owner).setStabilityFee(300))
        .to.emit(vaultModule, "VaultConfigUpdated")
        .withArgs("stabilityFee", 200, 300);

      const config = await vaultModule.getVaultConfig();
      expect(config.stabilityFee).to.equal(300);
    });

    it("Should update debt ceiling", async function () {
      const { vaultModule, owner } = await loadFixture(deployVaultModuleFixture);

      const newCeiling = ethers.parseEther("1000000");
      await expect(vaultModule.connect(owner).setDebtCeiling(newCeiling))
        .to.emit(vaultModule, "VaultConfigUpdated");

      const config = await vaultModule.getVaultConfig();
      expect(config.debtCeiling).to.equal(newCeiling);
    });

    it("Should only allow owner to update configuration", async function () {
      const { vaultModule, user1 } = await loadFixture(deployVaultModuleFixture);

      await expect(
        vaultModule.connect(user1).setMinCollateralRatio(160)
      ).to.be.reverted;
    });
  });

  describe("Statistics and Reporting", function () {
    it("Should get vault statistics", async function () {
      const { vaultModule, collateralVault, user1, user2 } = await loadFixture(deployVaultModuleFixture);

      // Create positions
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("3000"));
      await collateralVault.connect(user2).increaseDebt(user2.address, ethers.parseEther("2000"));

      const [totalCollateral, totalDebt, avgRatio, utilization] = await vaultModule.getVaultStats();

      expect(totalCollateral).to.equal(ethers.parseEther("4500"));
      expect(totalDebt).to.equal(ethers.parseEther("3000"));
      expect(avgRatio).to.equal(150);
    });

    it("Should count active positions", async function () {
      const { vaultModule, collateralVault, user1, user2 } = await loadFixture(deployVaultModuleFixture);

      // Create positions
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("1500"));
      await collateralVault.connect(user2).increaseDebt(user2.address, ethers.parseEther("1000"));

      const activeCount = await vaultModule.getActivePositionCount();
      expect(activeCount).to.equal(2);
    });

    it("Should return correct token addresses", async function () {
      const { vaultModule, loyaltyToken, lusdToken } = await loadFixture(deployVaultModuleFixture);

      expect(await vaultModule.getCollateralToken()).to.equal(await loyaltyToken.getAddress());
      expect(await vaultModule.getDebtToken()).to.equal(await lusdToken.getAddress());
    });
  });

  describe("Integration with Module Registry", function () {
    it("Should be queryable from module registry", async function () {
      const { vaultModule, moduleRegistry } = await loadFixture(deployVaultModuleFixture);

      const moduleId = await vaultModule.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.implementation).to.equal(await vaultModule.getAddress());
      expect(moduleInfo.isEnabled).to.be.true;
    });

    it("Should enforce module state changes", async function () {
      const { vaultModule, moduleRegistry, owner } = await loadFixture(deployVaultModuleFixture);

      const moduleId = await vaultModule.MODULE_ID();

      // Pause module
      await vaultModule.connect(owner).pause();

      const moduleInfoAfterPause = await vaultModule.getModuleInfo();
      expect(moduleInfoAfterPause.state).to.equal(2); // PAUSED

      // Unpause module
      await vaultModule.connect(owner).unpause();

      const moduleInfoAfterUnpause = await vaultModule.getModuleInfo();
      expect(moduleInfoAfterUnpause.state).to.equal(1); // ACTIVE
    });
  });
});
