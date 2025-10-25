/**
 * Deployment Script 3: Deploy Business Modules
 *
 * Deploys the business logic modules (adapters):
 * 1. VaultModule - Wraps CollateralVault
 * 2. RWAModule - Wraps OrderBook
 *
 * Prerequisites:
 * - Must run 01-deploy-infrastructure.js first
 * - Must run 02-deploy-service-modules.js first
 * - Existing CollateralVault and OrderBook must be deployed
 */

async function main() {
  const [deployer] = await ethers.getSigners();
  const fs = require('fs');
  const path = require('path');

  console.log("=".repeat(60));
  console.log("Deploying Business Modules");
  console.log("=".repeat(60));
  console.log("Deployer address:", deployer.address);
  console.log("");

  // ============ Load Previous Deployments ============
  const infraPath = path.join(__dirname, '../../deployments/infrastructure.json');
  const servicesPath = path.join(__dirname, '../../deployments/service-modules.json');

  if (!fs.existsSync(infraPath) || !fs.existsSync(servicesPath)) {
    throw new Error("Previous deployments not found! Run scripts 01 and 02 first");
  }

  const infrastructure = JSON.parse(fs.readFileSync(infraPath, 'utf8'));
  const services = JSON.parse(fs.readFileSync(servicesPath, 'utf8'));

  const moduleRegistryAddress = infrastructure.contracts.ModuleRegistry;

  console.log("Using ModuleRegistry at:", moduleRegistryAddress);
  console.log("");

  // ============ Get Existing Contract Addresses ============
  // NOTE: Update these addresses with your actual deployed contracts
  console.log("⚠️  Configuration Required:");
  console.log("   Please update the following addresses in this script:");
  console.log("");

  const COLLATERAL_VAULT_ADDRESS = process.env.COLLATERAL_VAULT_ADDRESS || "0x0000000000000000000000000000000000000000";
  const COLLATERAL_TOKEN_ADDRESS = process.env.COLLATERAL_TOKEN_ADDRESS || "0x0000000000000000000000000000000000000000";
  const DEBT_TOKEN_ADDRESS = process.env.DEBT_TOKEN_ADDRESS || "0x0000000000000000000000000000000000000000";
  const ORDER_BOOK_ADDRESS = process.env.ORDER_BOOK_ADDRESS || "0x0000000000000000000000000000000000000000";

  if (COLLATERAL_VAULT_ADDRESS === "0x0000000000000000000000000000000000000000") {
    console.log("❌ Error: Please set environment variables:");
    console.log("   - COLLATERAL_VAULT_ADDRESS");
    console.log("   - COLLATERAL_TOKEN_ADDRESS");
    console.log("   - DEBT_TOKEN_ADDRESS");
    console.log("   - ORDER_BOOK_ADDRESS");
    console.log("");
    console.log("Example:");
    console.log("   export COLLATERAL_VAULT_ADDRESS=0x...");
    console.log("   npx hardhat run scripts/deploy/03-deploy-business-modules.js");
    process.exit(1);
  }

  console.log("Using existing contracts:");
  console.log("  CollateralVault:", COLLATERAL_VAULT_ADDRESS);
  console.log("  CollateralToken:", COLLATERAL_TOKEN_ADDRESS);
  console.log("  DebtToken:      ", DEBT_TOKEN_ADDRESS);
  console.log("  OrderBook:      ", ORDER_BOOK_ADDRESS);
  console.log("");

  // ============ Deploy VaultModule ============
  console.log("1️⃣  Deploying VaultModule...");
  const VaultModule = await ethers.getContractFactory("VaultModule");
  const vaultModule = await VaultModule.deploy(
    COLLATERAL_VAULT_ADDRESS,
    COLLATERAL_TOKEN_ADDRESS,
    DEBT_TOKEN_ADDRESS
  );
  await vaultModule.waitForDeployment();
  const vaultModuleAddress = await vaultModule.getAddress();
  console.log("✅ VaultModule deployed to:", vaultModuleAddress);
  console.log("");

  // ============ Deploy RWAModule ============
  console.log("2️⃣  Deploying RWAModule...");
  const RWAModule = await ethers.getContractFactory("RWAModule");
  const rwaModule = await RWAModule.deploy(ORDER_BOOK_ADDRESS);
  await rwaModule.waitForDeployment();
  const rwaModuleAddress = await rwaModule.getAddress();
  console.log("✅ RWAModule deployed to:", rwaModuleAddress);
  console.log("");

  // ============ Register Modules ============
  console.log("3️⃣  Registering modules in ModuleRegistry...");
  const moduleRegistry = await ethers.getContractAt("ModuleRegistry", moduleRegistryAddress);

  // Register VaultModule
  console.log("   - Registering VaultModule...");
  const tx1 = await moduleRegistry.registerModule(vaultModuleAddress);
  await tx1.wait();
  const vaultModuleId = await vaultModule.MODULE_ID();
  console.log("   ✓ VaultModule registered with ID:", vaultModuleId);

  // Register RWAModule
  console.log("   - Registering RWAModule...");
  const tx2 = await moduleRegistry.registerModule(rwaModuleAddress);
  await tx2.wait();
  const rwaModuleId = await rwaModule.MODULE_ID();
  console.log("   ✓ RWAModule registered with ID:", rwaModuleId);
  console.log("");

  // ============ Validate Dependencies ============
  console.log("4️⃣  Validating module dependencies...");

  const vaultDepsValid = await moduleRegistry.validateDependencies(vaultModuleId);
  console.log("   - VaultModule dependencies:", vaultDepsValid ? "✓ Satisfied" : "✗ Not satisfied");

  const rwaDepsValid = await moduleRegistry.validateDependencies(rwaModuleId);
  console.log("   - RWAModule dependencies:", rwaDepsValid ? "✓ Satisfied" : "✗ Not satisfied");
  console.log("");

  // ============ Enable Modules ============
  if (vaultDepsValid && rwaDepsValid) {
    console.log("5️⃣  Enabling modules...");

    console.log("   - Enabling VaultModule...");
    const tx3 = await moduleRegistry.enableModule(vaultModuleId);
    await tx3.wait();
    console.log("   ✓ VaultModule enabled");

    console.log("   - Enabling RWAModule...");
    const tx4 = await moduleRegistry.enableModule(rwaModuleId);
    await tx4.wait();
    console.log("   ✓ RWAModule enabled");
    console.log("");
  } else {
    console.log("⚠️  Warning: Cannot enable modules - dependencies not satisfied");
    console.log("   Make sure PriceOracleModule and AuditModule are enabled");
    console.log("");
  }

  // ============ System Health Check ============
  console.log("6️⃣  Performing system health check...");
  const [healthyCount, totalCount, unhealthyIds] = await moduleRegistry.systemHealthCheck();
  console.log(`   - Healthy modules: ${healthyCount}/${totalCount}`);
  if (unhealthyIds.length > 0) {
    console.log("   - Unhealthy modules:", unhealthyIds);
  }
  console.log("");

  // ============ Save Deployment Info ============
  const deploymentInfo = {
    network: (await ethers.provider.getNetwork()).name,
    chainId: (await ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    existingContracts: {
      CollateralVault: COLLATERAL_VAULT_ADDRESS,
      CollateralToken: COLLATERAL_TOKEN_ADDRESS,
      DebtToken: DEBT_TOKEN_ADDRESS,
      OrderBook: ORDER_BOOK_ADDRESS
    },
    modules: {
      VaultModule: {
        address: vaultModuleAddress,
        moduleId: vaultModuleId,
        status: vaultDepsValid ? "enabled" : "registered"
      },
      RWAModule: {
        address: rwaModuleAddress,
        moduleId: rwaModuleId,
        status: rwaDepsValid ? "enabled" : "registered"
      }
    },
    systemHealth: {
      healthyModules: healthyCount.toString(),
      totalModules: totalCount.toString()
    }
  };

  const deployDir = path.join(__dirname, '../../deployments');
  fs.writeFileSync(
    path.join(deployDir, 'business-modules.json'),
    JSON.stringify(deploymentInfo, null, 2)
  );

  console.log("=".repeat(60));
  console.log("✅ Business Modules Deployment Complete!");
  console.log("=".repeat(60));
  console.log("Deployment info saved to: deployments/business-modules.json");
  console.log("");
  console.log("Module Addresses:");
  console.log("  VaultModule:", vaultModuleAddress);
  console.log("  RWAModule:  ", rwaModuleAddress);
  console.log("");
  console.log("System Status:");
  console.log(`  Active Modules: ${healthyCount}/${totalCount}`);
  console.log("");
  console.log("Deployment Complete! 🎉");
  console.log("=".repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
