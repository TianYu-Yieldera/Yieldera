import { ethers, upgrades } from "hardhat";

/**
 * Deploy Phase 3 - UUPS Proxy System
 *
 * This script deploys the upgradeable proxy system for all modules:
 * 1. Deploy ProxyAdmin
 * 2. Deploy VaultModuleV2 behind proxy
 * 3. Deploy RWAModuleV2 behind proxy (when ready)
 * 4. Register proxies with ModuleRegistry
 */
async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying Phase 3 with account:", deployer.address);

  // ============ Step 1: Deploy ProxyAdmin ============
  console.log("\n📋 Deploying ProxyAdmin...");
  const ProxyAdmin = await ethers.getContractFactory("ProxyAdmin");
  const proxyAdmin = await ProxyAdmin.deploy();
  await proxyAdmin.waitForDeployment();
  console.log("✅ ProxyAdmin deployed to:", await proxyAdmin.getAddress());

  // ============ Step 2: Deploy Legacy Contracts (if not already deployed) ============
  console.log("\n📋 Deploying legacy contracts...");

  // Deploy mock tokens for testing
  const MockERC20 = await ethers.getContractFactory("MockERC20");
  const loyaltyToken = await MockERC20.deploy(
    "Loyalty Points",
    "LP",
    ethers.parseEther("10000000")
  );
  await loyaltyToken.waitForDeployment();
  console.log("✅ LoyaltyToken deployed to:", await loyaltyToken.getAddress());

  const lusdToken = await MockERC20.deploy(
    "Loyalty USD",
    "LUSD",
    ethers.parseEther("10000000")
  );
  await lusdToken.waitForDeployment();
  console.log("✅ LUSD deployed to:", await lusdToken.getAddress());

  // Deploy CollateralVault
  const CollateralVault = await ethers.getContractFactory("CollateralVault");
  const collateralVault = await CollateralVault.deploy(await loyaltyToken.getAddress());
  await collateralVault.waitForDeployment();
  console.log("✅ CollateralVault deployed to:", await collateralVault.getAddress());

  // ============ Step 3: Deploy Module Registry ============
  console.log("\n📋 Deploying ModuleRegistry...");
  const ModuleRegistry = await ethers.getContractFactory("ModuleRegistry");
  const moduleRegistry = await ModuleRegistry.deploy();
  await moduleRegistry.waitForDeployment();
  console.log("✅ ModuleRegistry deployed to:", await moduleRegistry.getAddress());

  // ============ Step 4: Deploy VaultModuleV2 with UUPS Proxy ============
  console.log("\n📋 Deploying VaultModuleV2 with proxy...");

  const VaultModuleV2 = await ethers.getContractFactory("VaultModuleV2");

  // Deploy implementation
  const vaultImplementation = await VaultModuleV2.deploy();
  await vaultImplementation.waitForDeployment();
  console.log("✅ VaultModuleV2 implementation:", await vaultImplementation.getAddress());

  // Encode initialization data
  const initData = VaultModuleV2.interface.encodeFunctionData("initialize", [
    await collateralVault.getAddress(),
    await loyaltyToken.getAddress(),
    await lusdToken.getAddress()
  ]);

  // Deploy proxy
  const ModuleProxy = await ethers.getContractFactory("ModuleProxy");
  const vaultProxy = await ModuleProxy.deploy(
    await vaultImplementation.getAddress(),
    initData
  );
  await vaultProxy.waitForDeployment();
  console.log("✅ VaultModule Proxy deployed to:", await vaultProxy.getAddress());

  // Get proxy contract with V2 interface
  const vaultModule = VaultModuleV2.attach(await vaultProxy.getAddress());

  // Verify initialization
  const version = await vaultModule.getImplementationVersion();
  console.log("   Version:", version);

  // ============ Step 5: Register proxy with ProxyAdmin ============
  console.log("\n📋 Registering proxy with ProxyAdmin...");
  await proxyAdmin.registerProxy(await vaultProxy.getAddress(), "VaultModule");
  console.log("✅ VaultModule proxy registered");

  // ============ Step 6: Register module with ModuleRegistry ============
  console.log("\n📋 Registering module with ModuleRegistry...");
  await moduleRegistry.registerModule(await vaultProxy.getAddress());
  const moduleId = await vaultModule.MODULE_ID();
  await moduleRegistry.enableModule(moduleId);
  console.log("✅ VaultModule registered and enabled");

  // ============ Step 7: Verify Deployment ============
  console.log("\n📋 Verifying deployment...");

  const moduleInfo = await moduleRegistry.getModuleInfo(moduleId);
  console.log("Module Info:");
  console.log("  - Name:", moduleInfo.name);
  console.log("  - Version:", moduleInfo.version);
  console.log("  - State:", moduleInfo.state);
  console.log("  - Is Registered:", moduleInfo.isRegistered);
  console.log("  - Is Enabled:", moduleInfo.isEnabled);

  const [healthy, message] = await vaultModule.healthCheck();
  console.log("Health Check:");
  console.log("  - Healthy:", healthy);
  console.log("  - Message:", message);

  // ============ Summary ============
  console.log("\n" + "=".repeat(60));
  console.log("📊 Phase 3 Deployment Summary");
  console.log("=".repeat(60));
  console.log("ProxyAdmin:           ", await proxyAdmin.getAddress());
  console.log("ModuleRegistry:       ", await moduleRegistry.getAddress());
  console.log("LoyaltyToken:         ", await loyaltyToken.getAddress());
  console.log("LUSD:                 ", await lusdToken.getAddress());
  console.log("CollateralVault:      ", await collateralVault.getAddress());
  console.log("-".repeat(60));
  console.log("VaultModuleV2 Impl:   ", await vaultImplementation.getAddress());
  console.log("VaultModule Proxy:    ", await vaultProxy.getAddress());
  console.log("=".repeat(60));

  // Save deployment addresses
  const deployment = {
    network: (await ethers.provider.getNetwork()).name,
    timestamp: new Date().toISOString(),
    deployer: deployer.address,
    contracts: {
      proxyAdmin: await proxyAdmin.getAddress(),
      moduleRegistry: await moduleRegistry.getAddress(),
      loyaltyToken: await loyaltyToken.getAddress(),
      lusd: await lusdToken.getAddress(),
      collateralVault: await collateralVault.getAddress(),
      vaultModule: {
        implementation: await vaultImplementation.getAddress(),
        proxy: await vaultProxy.getAddress()
      }
    }
  };

  console.log("\n💾 Deployment data:");
  console.log(JSON.stringify(deployment, null, 2));

  return deployment;
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
