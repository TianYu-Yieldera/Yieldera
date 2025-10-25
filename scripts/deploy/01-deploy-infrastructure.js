/**
 * Deployment Script 1: Deploy Core Infrastructure
 *
 * Deploys the foundational contracts for the pluggable architecture:
 * 1. AccessController - Unified access control
 * 2. ModuleRegistry - Module management
 * 3. EventHub - Event bus for inter-module communication
 */

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("=".repeat(60));
  console.log("Deploying Core Infrastructure");
  console.log("=".repeat(60));
  console.log("Deployer address:", deployer.address);
  console.log("Account balance:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)));
  console.log("");

  // ============ Deploy AccessController ============
  console.log("1️⃣  Deploying AccessController...");
  const AccessController = await ethers.getContractFactory("AccessController");
  const accessController = await AccessController.deploy();
  await accessController.waitForDeployment();
  const accessControllerAddress = await accessController.getAddress();
  console.log("✅ AccessController deployed to:", accessControllerAddress);
  console.log("");

  // ============ Deploy ModuleRegistry ============
  console.log("2️⃣  Deploying ModuleRegistry...");
  const ModuleRegistry = await ethers.getContractFactory("ModuleRegistry");
  const moduleRegistry = await ModuleRegistry.deploy();
  await moduleRegistry.waitForDeployment();
  const moduleRegistryAddress = await moduleRegistry.getAddress();
  console.log("✅ ModuleRegistry deployed to:", moduleRegistryAddress);
  console.log("");

  // ============ Deploy EventHub ============
  console.log("3️⃣  Deploying EventHub...");
  const EventHub = await ethers.getContractFactory("EventHub");
  const eventHub = await EventHub.deploy();
  await eventHub.waitForDeployment();
  const eventHubAddress = await eventHub.getAddress();
  console.log("✅ EventHub deployed to:", eventHubAddress);
  console.log("");

  // ============ Configure Connections ============
  console.log("4️⃣  Configuring infrastructure connections...");

  // Set AccessController in ModuleRegistry
  console.log("   - Setting AccessController in ModuleRegistry...");
  const tx1 = await moduleRegistry.setAccessController(accessControllerAddress);
  await tx1.wait();
  console.log("   ✓ AccessController set");

  // Set ModuleRegistry in EventHub
  console.log("   - Setting ModuleRegistry in EventHub...");
  const tx2 = await eventHub.setModuleRegistry(moduleRegistryAddress);
  await tx2.wait();
  console.log("   ✓ ModuleRegistry set");
  console.log("");

  // ============ Save Deployment Info ============
  const deploymentInfo = {
    network: (await ethers.provider.getNetwork()).name,
    chainId: (await ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    contracts: {
      AccessController: accessControllerAddress,
      ModuleRegistry: moduleRegistryAddress,
      EventHub: eventHubAddress
    }
  };

  const fs = require('fs');
  const path = require('path');
  const deployDir = path.join(__dirname, '../../deployments');

  if (!fs.existsSync(deployDir)) {
    fs.mkdirSync(deployDir, { recursive: true });
  }

  fs.writeFileSync(
    path.join(deployDir, 'infrastructure.json'),
    JSON.stringify(deploymentInfo, null, 2)
  );

  console.log("=".repeat(60));
  console.log("✅ Core Infrastructure Deployment Complete!");
  console.log("=".repeat(60));
  console.log("Deployment info saved to: deployments/infrastructure.json");
  console.log("");
  console.log("Contract Addresses:");
  console.log("  AccessController:", accessControllerAddress);
  console.log("  ModuleRegistry:  ", moduleRegistryAddress);
  console.log("  EventHub:        ", eventHubAddress);
  console.log("");
  console.log("Next step: Run 02-deploy-service-modules.js");
  console.log("=".repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
