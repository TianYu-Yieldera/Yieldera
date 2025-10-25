/**
 * Phase 5 Deployment Script
 * Deploys plugin system infrastructure and example plugins
 */

const { ethers } = require("hardhat");

async function main() {
  console.log("=".repeat(60));
  console.log("Phase 5: Plugin System Deployment");
  console.log("=".repeat(60));

  const [deployer] = await ethers.getSigners();
  console.log("\nDeploying with account:", deployer.address);

  const balance = await ethers.provider.getBalance(deployer.address);
  console.log("Account balance:", ethers.formatEther(balance), "ETH");

  // ============ Deploy Plugin Infrastructure ============
  console.log("\n🔌 Deploying Plugin Infrastructure...");

  // 1. PluginRegistry
  const PluginRegistry = await ethers.getContractFactory("PluginRegistry");
  const pluginRegistry = await PluginRegistry.deploy();
  await pluginRegistry.waitForDeployment();
  console.log("✓ PluginRegistry:", await pluginRegistry.getAddress());

  // 2. PluginPermissionManager
  const PluginPermissionManager = await ethers.getContractFactory("PluginPermissionManager");
  const pluginPermissionManager = await PluginPermissionManager.deploy();
  await pluginPermissionManager.waitForDeployment();
  console.log("✓ PluginPermissionManager:", await pluginPermissionManager.getAddress());

  // ============ Deploy Example Plugins ============
  console.log("\n📦 Deploying Example Plugins...");

  // 1. RewardMultiplierPlugin
  const RewardMultiplierPlugin = await ethers.getContractFactory("RewardMultiplierPlugin");
  const rewardMultiplierPlugin = await RewardMultiplierPlugin.deploy(deployer.address);
  await rewardMultiplierPlugin.waitForDeployment();
  console.log("✓ RewardMultiplierPlugin:", await rewardMultiplierPlugin.getAddress());

  // ============ Register Plugins ============
  console.log("\n📝 Registering Plugins...");

  const rewardMultiplierPluginAddr = await rewardMultiplierPlugin.getAddress();

  // Initialize plugin
  await rewardMultiplierPlugin.initialize("0x");
  console.log("✓ RewardMultiplierPlugin initialized");

  // Register plugin
  const pluginId = await rewardMultiplierPlugin.getPluginId();
  await pluginRegistry.registerPlugin(rewardMultiplierPluginAddr);
  console.log("✓ RewardMultiplierPlugin registered with ID:", pluginId);

  // ============ Configure Permissions ============
  console.log("\n🔐 Configuring Permissions...");

  // Grant READ_VAULT permission to RewardMultiplierPlugin
  const READ_PERMISSION = 1; // Permission.READ
  const VAULT_MODULE_ADDRESS = ethers.ZeroAddress; // Global permission

  await pluginPermissionManager.grantPermission(
    pluginId,
    READ_PERMISSION,
    VAULT_MODULE_ADDRESS,
    0 // Permanent permission
  );
  console.log("✓ Granted READ permission to RewardMultiplierPlugin");

  // ============ Verify Registration ============
  console.log("\n✅ Verifying Registration...");

  const isRegistered = await pluginRegistry.isPluginRegistered(pluginId);
  console.log("Plugin registered:", isRegistered);

  const isActive = await pluginRegistry.isPluginActive(pluginId);
  console.log("Plugin active:", isActive);

  const registration = await pluginRegistry.getPluginRegistration(pluginId);
  console.log("\nPlugin Registration:");
  console.log("  Name:", await rewardMultiplierPlugin.getPluginName());
  console.log("  Version:", await rewardMultiplierPlugin.getPluginVersion());
  console.log("  Type:", registration.pluginType.toString());
  console.log("  Author:", registration.author);

  // ============ Test Plugin Functionality ============
  console.log("\n🧪 Testing Plugin Functionality...");

  try {
    // Test: Record activity for a test user
    const testUser = deployer.address;
    const activityData = ethers.AbiCoder.defaultAbiCoder().encode(
      ["address", "uint256"],
      [testUser, 150] // 150 activity points
    );

    const executeData = ethers.AbiCoder.defaultAbiCoder().encode(
      ["string", "bytes"],
      ["recordActivity", activityData]
    );

    const result = await rewardMultiplierPlugin.execute(executeData);
    const [activityScore, multiplier] = ethers.AbiCoder.defaultAbiCoder().decode(
      ["uint256", "uint256"],
      result
    );

    console.log("✓ Activity recorded:");
    console.log("  Activity Score:", activityScore.toString());
    console.log("  Multiplier:", (Number(multiplier) / 100).toFixed(2) + "x");

    // Test: Calculate reward
    const rewardData = ethers.AbiCoder.defaultAbiCoder().encode(
      ["address", "uint256"],
      [testUser, ethers.parseEther("100")] // 100 tokens base reward
    );

    const rewardExecuteData = ethers.AbiCoder.defaultAbiCoder().encode(
      ["string", "bytes"],
      ["calculateReward", rewardData]
    );

    const rewardResult = await rewardMultiplierPlugin.execute(rewardExecuteData);
    const [finalReward] = ethers.AbiCoder.defaultAbiCoder().decode(
      ["uint256"],
      rewardResult
    );

    console.log("✓ Reward calculated:");
    console.log("  Base Reward: 100 tokens");
    console.log("  Final Reward:", ethers.formatEther(finalReward), "tokens");

    // Health check
    const [healthy, message] = await rewardMultiplierPlugin.healthCheck();
    console.log("\n🏥 Plugin Health Check:");
    console.log("  Status:", healthy ? "✓ HEALTHY" : "✗ UNHEALTHY");
    console.log("  Message:", message);

  } catch (error) {
    console.log("⚠️  Plugin test warning:", error.message);
  }

  // ============ Summary ============
  console.log("\n" + "=".repeat(60));
  console.log("Deployment Summary");
  console.log("=".repeat(60));

  const summary = {
    "Plugin Infrastructure": {
      "PluginRegistry": await pluginRegistry.getAddress(),
      "PluginPermissionManager": await pluginPermissionManager.getAddress()
    },
    "Example Plugins": {
      "RewardMultiplierPlugin": {
        "Address": rewardMultiplierPluginAddr,
        "Plugin ID": pluginId,
        "Registered": isRegistered,
        "Active": isActive
      }
    },
    "Statistics": {
      "Total Plugins": await pluginRegistry.getTotalPluginCount(),
      "Active Plugins": await pluginRegistry.getActivePluginCount(),
      "Verified Plugins": await pluginRegistry.getVerifiedPluginCount()
    }
  };

  console.log(JSON.stringify(summary, null, 2));

  console.log("\n" + "=".repeat(60));
  console.log("✅ Phase 5 Deployment Complete!");
  console.log("=".repeat(60));

  console.log("\n📚 Next Steps:");
  console.log("1. Verify plugin on Etherscan (if on mainnet/testnet)");
  console.log("2. Test plugin with actual vault module");
  console.log("3. Create additional plugins as needed");
  console.log("4. Set up plugin monitoring and analytics");

  return summary;
}

// Execute deployment
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ Deployment failed:");
    console.error(error);
    process.exit(1);
  });
