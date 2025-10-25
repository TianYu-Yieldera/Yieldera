/**
 * Deployment Script 2: Deploy Service Modules
 *
 * Deploys the independent service modules:
 * 1. PriceOracleModule - Price oracle service
 * 2. AuditModule - Audit logging service
 *
 * Prerequisites: Must run 01-deploy-infrastructure.js first
 */

async function main() {
  const [deployer] = await ethers.getSigners();
  const fs = require('fs');
  const path = require('path');

  console.log("=".repeat(60));
  console.log("Deploying Service Modules");
  console.log("=".repeat(60));
  console.log("Deployer address:", deployer.address);
  console.log("");

  // ============ Load Infrastructure Addresses ============
  const infraPath = path.join(__dirname, '../../deployments/infrastructure.json');
  if (!fs.existsSync(infraPath)) {
    throw new Error("Infrastructure not deployed! Run 01-deploy-infrastructure.js first");
  }

  const infrastructure = JSON.parse(fs.readFileSync(infraPath, 'utf8'));
  const moduleRegistryAddress = infrastructure.contracts.ModuleRegistry;

  console.log("Using ModuleRegistry at:", moduleRegistryAddress);
  console.log("");

  // ============ Deploy PriceOracleModule ============
  console.log("1️⃣  Deploying PriceOracleModule...");
  const PriceOracleModule = await ethers.getContractFactory("PriceOracleModule");
  const priceOracle = await PriceOracleModule.deploy();
  await priceOracle.waitForDeployment();
  const priceOracleAddress = await priceOracle.getAddress();
  console.log("✅ PriceOracleModule deployed to:", priceOracleAddress);
  console.log("");

  // ============ Deploy AuditModule ============
  console.log("2️⃣  Deploying AuditModule...");
  const AuditModule = await ethers.getContractFactory("AuditModule");
  const auditModule = await AuditModule.deploy();
  await auditModule.waitForDeployment();
  const auditModuleAddress = await auditModule.getAddress();
  console.log("✅ AuditModule deployed to:", auditModuleAddress);
  console.log("");

  // ============ Register Modules ============
  console.log("3️⃣  Registering modules in ModuleRegistry...");
  const moduleRegistry = await ethers.getContractAt("ModuleRegistry", moduleRegistryAddress);

  // Register PriceOracleModule
  console.log("   - Registering PriceOracleModule...");
  const tx1 = await moduleRegistry.registerModule(priceOracleAddress);
  await tx1.wait();
  const priceOracleModuleId = await priceOracle.MODULE_ID();
  console.log("   ✓ PriceOracleModule registered with ID:", priceOracleModuleId);

  // Register AuditModule
  console.log("   - Registering AuditModule...");
  const tx2 = await moduleRegistry.registerModule(auditModuleAddress);
  await tx2.wait();
  const auditModuleId = await auditModule.MODULE_ID();
  console.log("   ✓ AuditModule registered with ID:", auditModuleId);
  console.log("");

  // ============ Enable Modules ============
  console.log("4️⃣  Enabling modules...");

  console.log("   - Enabling PriceOracleModule...");
  const tx3 = await moduleRegistry.enableModule(priceOracleModuleId);
  await tx3.wait();
  console.log("   ✓ PriceOracleModule enabled");

  console.log("   - Enabling AuditModule...");
  const tx4 = await moduleRegistry.enableModule(auditModuleId);
  await tx4.wait();
  console.log("   ✓ AuditModule enabled");
  console.log("");

  // ============ Save Deployment Info ============
  const deploymentInfo = {
    network: (await ethers.provider.getNetwork()).name,
    chainId: (await ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    modules: {
      PriceOracleModule: {
        address: priceOracleAddress,
        moduleId: priceOracleModuleId,
        status: "enabled"
      },
      AuditModule: {
        address: auditModuleAddress,
        moduleId: auditModuleId,
        status: "enabled"
      }
    }
  };

  const deployDir = path.join(__dirname, '../../deployments');
  fs.writeFileSync(
    path.join(deployDir, 'service-modules.json'),
    JSON.stringify(deploymentInfo, null, 2)
  );

  console.log("=".repeat(60));
  console.log("✅ Service Modules Deployment Complete!");
  console.log("=".repeat(60));
  console.log("Deployment info saved to: deployments/service-modules.json");
  console.log("");
  console.log("Module Addresses:");
  console.log("  PriceOracleModule:", priceOracleAddress);
  console.log("  AuditModule:      ", auditModuleAddress);
  console.log("");
  console.log("Next step: Run 03-deploy-business-modules.js");
  console.log("=".repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
