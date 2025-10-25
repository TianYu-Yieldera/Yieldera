import { expect } from "chai";
import { ethers, upgrades } from "hardhat";
import { loadFixture, time } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Phase 3 - Upgradeable Modules", function () {
  // Fixture for deploying upgradeable system
  async function deployUpgradeableSystemFixture() {
    const [owner, user1, user2, upgrader] = await ethers.getSigners();

    // Deploy tokens
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    const loyaltyToken = await MockERC20.deploy("LP", "LP", ethers.parseEther("1000000"));
    const lusdToken = await MockERC20.deploy("LUSD", "LUSD", ethers.parseEther("1000000"));

    // Deploy legacy contracts
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const collateralVault = await CollateralVault.deploy(await loyaltyToken.getAddress());

    // Deploy ProxyAdmin
    const ProxyAdmin = await ethers.getContractFactory("ProxyAdmin");
    const proxyAdmin = await ProxyAdmin.deploy();

    // Deploy ModuleRegistry
    const ModuleRegistry = await ethers.getContractFactory("ModuleRegistry");
    const moduleRegistry = await ModuleRegistry.deploy();

    // Deploy VaultModuleV2 implementation
    const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
    const vaultImplementationV2 = await VaultModuleV2.deploy();

    // Encode init data
    const initData = VaultModuleV2.interface.encodeFunctionData("initialize", [
      await collateralVault.getAddress(),
      await loyaltyToken.getAddress(),
      await lusdToken.getAddress()
    ]);

    // Deploy proxy
    const ModuleProxy = await ethers.getContractFactory("ModuleProxy");
    const vaultProxy = await ModuleProxy.deploy(
      await vaultImplementationV2.getAddress(),
      initData
    );

    // Attach V2 interface to proxy
    const vaultModule = VaultModuleV2.attach(await vaultProxy.getAddress());

    // Register proxy
    await proxyAdmin.registerProxy(await vaultProxy.getAddress(), "VaultModule");

    // Register with module registry
    await moduleRegistry.registerModule(await vaultProxy.getAddress());
    await moduleRegistry.enableModule(await vaultModule.MODULE_ID());

    // Setup tokens for users
    await loyaltyToken.transfer(user1.address, ethers.parseEther("10000"));
    await loyaltyToken.transfer(user2.address, ethers.parseEther("10000"));
    await loyaltyToken.connect(user1).approve(await collateralVault.getAddress(), ethers.MaxUint256);
    await loyaltyToken.connect(user2).approve(await collateralVault.getAddress(), ethers.MaxUint256);

    return {
      vaultModule,
      vaultProxy,
      vaultImplementationV2,
      proxyAdmin,
      moduleRegistry,
      collateralVault,
      loyaltyToken,
      lusdToken,
      owner,
      user1,
      user2,
      upgrader
    };
  }

  describe("Proxy Deployment", function () {
    it("Should deploy proxy with correct implementation", async function () {
      const { vaultModule, vaultImplementationV2 } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const implementation = await vaultModule.getImplementation();
      expect(implementation).to.equal(await vaultImplementationV2.getAddress());
    });

    it("Should have correct version", async function () {
      const { vaultModule } = await loadFixture(deployUpgradeableSystemFixture);

      const version = await vaultModule.getImplementationVersion();
      expect(version).to.equal("2.0.0");
    });

    it("Should be registered in ModuleRegistry", async function () {
      const { vaultModule, moduleRegistry } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const moduleId = await vaultModule.MODULE_ID();
      const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);

      expect(moduleInfo.isRegistered).to.be.true;
      expect(moduleInfo.isEnabled).to.be.true;
    });

    it("Should be registered in ProxyAdmin", async function () {
      const { vaultProxy, proxyAdmin } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const proxyInfo = await proxyAdmin.getProxyInfo(await vaultProxy.getAddress());
      expect(proxyInfo.isRegistered).to.be.true;
      expect(proxyInfo.moduleName).to.equal("VaultModule");
    });
  });

  describe("Diamond Storage", function () {
    it("Should store data in Diamond Storage", async function () {
      const { vaultModule, user1 } = await loadFixture(deployUpgradeableSystemFixture);

      // Deposit collateral
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));

      // Check balance
      const balance = await vaultModule.getCollateralBalance(user1.address);
      expect(balance).to.equal(ethers.parseEther("1000"));
    });

    it("Should maintain data across function calls", async function () {
      const { vaultModule, user1 } = await loadFixture(deployUpgradeableSystemFixture);

      // Multiple operations
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("500"));

      const balance = await vaultModule.getCollateralBalance(user1.address);
      expect(balance).to.equal(ethers.parseEther("1500"));
    });

    it("Should track configuration in Diamond Storage", async function () {
      const { vaultModule, owner } = await loadFixture(deployUpgradeableSystemFixture);

      // Update config
      await vaultModule.connect(owner).setMinCollateralRatio(160);

      const config = await vaultModule.getVaultConfig();
      expect(config.minCollateralRatio).to.equal(160);
    });
  });

  describe("Module Functionality", function () {
    it("Should support all vault operations", async function () {
      const { vaultModule, collateralVault, user1 } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      // Deposit
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("2000"));

      // Increase debt
      await collateralVault.connect(user1).increaseDebt(user1.address, ethers.parseEther("1000"));

      // Check position
      const position = await vaultModule.getPosition(user1.address);
      expect(position.collateralAmount).to.equal(ethers.parseEther("2000"));
      expect(position.debtAmount).to.equal(ethers.parseEther("1000"));
      expect(position.isActive).to.be.true;

      // Check ratio
      const ratio = await vaultModule.getCollateralRatio(user1.address);
      expect(ratio).to.equal(200);
    });

    it("Should perform health checks", async function () {
      const { vaultModule } = await loadFixture(deployUpgradeableSystemFixture);

      const [healthy, message] = await vaultModule.healthCheck();
      expect(healthy).to.be.true;
      expect(message).to.include("healthy");
    });

    it("Should track active positions", async function () {
      const { vaultModule, user1, user2 } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));
      await vaultModule.connect(user2).depositCollateral(ethers.parseEther("2000"));

      const count = await vaultModule.getActivePositionCount();
      expect(count).to.equal(2);
    });
  });

  describe("Upgrade Authorization", function () {
    it("Should allow owner to authorize upgraders", async function () {
      const { vaultModule, owner, upgrader } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      await vaultModule.connect(owner).authorizeUpgrader(upgrader.address);

      expect(await vaultModule.canUpgrade(upgrader.address)).to.be.true;
    });

    it("Should allow owner to revoke authorization", async function () {
      const { vaultModule, owner, upgrader } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      await vaultModule.connect(owner).authorizeUpgrader(upgrader.address);
      await vaultModule.connect(owner).revokeUpgradeAuthorization(upgrader.address);

      expect(await vaultModule.canUpgrade(upgrader.address)).to.be.false;
    });

    it("Should validate upgrade implementations", async function () {
      const { vaultModule } = await loadFixture(deployUpgradeableSystemFixture);

      // Valid contract
      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      const [valid, reason] = await vaultModule.validateUpgrade(await newImpl.getAddress());
      expect(valid).to.be.true;

      // Invalid - zero address
      const [invalid, reason2] = await vaultModule.validateUpgrade(ethers.ZeroAddress);
      expect(invalid).to.be.false;
      expect(reason2).to.equal("Zero address");
    });
  });

  describe("Upgrade Process", function () {
    it("Should track upgrade history", async function () {
      const { vaultModule } = await loadFixture(deployUpgradeableSystemFixture);

      const [implementations, timestamps, versions] = await vaultModule.getUpgradeHistory();

      expect(implementations.length).to.equal(1);
      expect(versions[0]).to.equal("2.0.0");
    });

    it("Should allow pausing upgrades", async function () {
      const { vaultModule, owner } = await loadFixture(deployUpgradeableSystemFixture);

      await vaultModule.connect(owner).pauseUpgrades();

      expect(await vaultModule.upgradesPaused()).to.be.true;
    });

    it("Should prevent upgrades when paused", async function () {
      const { vaultModule, owner } = await loadFixture(deployUpgradeableSystemFixture);

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await vaultModule.connect(owner).pauseUpgrades();

      await expect(
        vaultModule.connect(owner).upgradeTo(await newImpl.getAddress())
      ).to.be.revertedWith("Upgrades paused");
    });
  });

  describe("ProxyAdmin Management", function () {
    it("Should schedule upgrade with timelock", async function () {
      const { vaultProxy, proxyAdmin, owner } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await proxyAdmin.connect(owner).scheduleUpgrade(
        await vaultProxy.getAddress(),
        await newImpl.getAddress()
      );

      const scheduledUpgrade = await proxyAdmin.scheduledUpgrades(await vaultProxy.getAddress());
      expect(scheduledUpgrade.newImplementation).to.equal(await newImpl.getAddress());
    });

    it("Should enforce timelock", async function () {
      const { vaultProxy, proxyAdmin, owner } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await proxyAdmin.connect(owner).scheduleUpgrade(
        await vaultProxy.getAddress(),
        await newImpl.getAddress()
      );

      // Try to execute immediately (should fail)
      await expect(
        proxyAdmin.connect(owner).executeUpgrade(await vaultProxy.getAddress())
      ).to.be.revertedWith("Timelock not expired");
    });

    it("Should execute upgrade after timelock", async function () {
      const { vaultProxy, proxyAdmin, owner } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await proxyAdmin.connect(owner).scheduleUpgrade(
        await vaultProxy.getAddress(),
        await newImpl.getAddress()
      );

      // Fast forward past timelock (2 days)
      await time.increase(2 * 24 * 60 * 60 + 1);

      // Execute upgrade
      await proxyAdmin.connect(owner).executeUpgrade(await vaultProxy.getAddress());

      const proxyInfo = await proxyAdmin.getProxyInfo(await vaultProxy.getAddress());
      expect(proxyInfo.currentImplementation).to.equal(await newImpl.getAddress());
    });

    it("Should allow cancelling scheduled upgrade", async function () {
      const { vaultProxy, proxyAdmin, owner } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await proxyAdmin.connect(owner).scheduleUpgrade(
        await vaultProxy.getAddress(),
        await newImpl.getAddress()
      );

      await proxyAdmin.connect(owner).cancelUpgrade(await vaultProxy.getAddress());

      const scheduledUpgrade = await proxyAdmin.scheduledUpgrades(await vaultProxy.getAddress());
      expect(scheduledUpgrade.cancelled).to.be.true;
    });

    it("Should support emergency upgrades", async function () {
      const { vaultProxy, proxyAdmin, owner } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      // Emergency upgrade bypasses timelock
      await proxyAdmin.connect(owner).emergencyUpgrade(
        await vaultProxy.getAddress(),
        await newImpl.getAddress()
      );

      const proxyInfo = await proxyAdmin.getProxyInfo(await vaultProxy.getAddress());
      expect(proxyInfo.currentImplementation).to.equal(await newImpl.getAddress());
    });

    it("Should pause all upgrades in emergency", async function () {
      const { proxyAdmin, owner } = await loadFixture(deployUpgradeableSystemFixture);

      await proxyAdmin.connect(owner).pause();

      expect(await proxyAdmin.paused()).to.be.true;
    });
  });

  describe("Storage Persistence Across Upgrades", function () {
    it("Should maintain state after upgrade", async function () {
      const { vaultProxy, vaultModule, proxyAdmin, user1, owner } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      // Deposit before upgrade
      await vaultModule.connect(user1).depositCollateral(ethers.parseEther("1000"));

      const balanceBefore = await vaultModule.getCollateralBalance(user1.address);

      // Deploy new implementation
      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      // Emergency upgrade (to skip timelock for testing)
      await proxyAdmin.connect(owner).emergencyUpgrade(
        await vaultProxy.getAddress(),
        await newImpl.getAddress()
      );

      // Check balance after upgrade
      const balanceAfter = await vaultModule.getCollateralBalance(user1.address);
      expect(balanceAfter).to.equal(balanceBefore);
    });
  });

  describe("Access Control", function () {
    it("Should only allow owner to upgrade", async function () {
      const { vaultModule, user1 } = await loadFixture(deployUpgradeableSystemFixture);

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await expect(
        vaultModule.connect(user1).upgradeTo(await newImpl.getAddress())
      ).to.be.reverted;
    });

    it("Should only allow owner to manage ProxyAdmin", async function () {
      const { vaultProxy, proxyAdmin, user1 } = await loadFixture(
        deployUpgradeableSystemFixture
      );

      const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");
      const newImpl = await VaultModuleV2.deploy();

      await expect(
        proxyAdmin.connect(user1).scheduleUpgrade(
          await vaultProxy.getAddress(),
          await newImpl.getAddress()
        )
      ).to.be.revertedWith("Not authorized");
    });
  });
});
