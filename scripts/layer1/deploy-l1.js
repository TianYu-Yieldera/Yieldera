/**
 * L1 Deployment Script
 * Deploys core L1 contracts to Ethereum Mainnet or Testnet
 *
 * Deploys:
 * 1. LoyaltyUSDL1 (stablecoin)
 * 2. CollateralVaultL1 (asset custody)
 * 3. L1StateRegistry (state verification)
 * 4. L1Gateway (bridge)
 */

import hre from "hardhat";
import { writeFileSync } from "fs";

async function main() {
  console.log("====================================");
  console.log("🚀 L1 Deployment Script");
  console.log("====================================\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);
  console.log("Account balance:", hre.ethers.formatEther(await hre.ethers.provider.getBalance(deployer.address)), "ETH\n");

  const network = hre.network.name;
  console.log("Network:", network);
  console.log("Chain ID:", (await hre.ethers.provider.getNetwork()).chainId);
  console.log("====================================\n");

  // Deployment results
  const deployedContracts = {
    network,
    chainId: Number((await hre.ethers.provider.getNetwork()).chainId),
    deployedAt: new Date().toISOString(),
    deployer: deployer.address,
    contracts: {}
  };

  // ============ Step 1: Deploy LoyaltyUSDL1 ============
  console.log("📄 Step 1: Deploying LoyaltyUSDL1...");
  const LoyaltyUSDL1 = await hre.ethers.getContractFactory("LoyaltyUSDL1");
  const loyaltyUSD = await LoyaltyUSDL1.deploy();
  await loyaltyUSD.waitForDeployment();
  const loyaltyUSDAddress = await loyaltyUSD.getAddress();
  console.log("✅ LoyaltyUSDL1 deployed to:", loyaltyUSDAddress);
  deployedContracts.contracts.loyaltyUSD = loyaltyUSDAddress;
  console.log("");

  // ============ Step 2: Deploy Mock Collateral Token (for testing) ============
  console.log("📄 Step 2: Deploying Mock Collateral Token...");
  const MockERC20 = await hre.ethers.getContractFactory("MockERC20");
  const collateralToken = await MockERC20.deploy(
    "Loyalty Points",
    "LP",
    hre.ethers.parseEther("1000000") // 1M tokens initial supply
  );
  await collateralToken.waitForDeployment();
  const collateralTokenAddress = await collateralToken.getAddress();
  console.log("✅ Collateral Token deployed to:", collateralTokenAddress);
  deployedContracts.contracts.collateralToken = collateralTokenAddress;
  console.log("");

  // ============ Step 3: Deploy CollateralVaultL1 ============
  console.log("📄 Step 3: Deploying CollateralVaultL1...");
  const CollateralVaultL1 = await hre.ethers.getContractFactory("CollateralVaultL1");
  const collateralVault = await CollateralVaultL1.deploy(collateralTokenAddress);
  await collateralVault.waitForDeployment();
  const collateralVaultAddress = await collateralVault.getAddress();
  console.log("✅ CollateralVaultL1 deployed to:", collateralVaultAddress);
  deployedContracts.contracts.collateralVault = collateralVaultAddress;
  console.log("");

  // ============ Step 4: Deploy L1StateRegistry ============
  console.log("📄 Step 4: Deploying L1StateRegistry...");
  // Note: L2 aggregator address will be set later after L2 deployment
  const tempL2Aggregator = deployer.address; // Temporary, will update later
  const L1StateRegistry = await hre.ethers.getContractFactory("L1StateRegistry");
  const stateRegistry = await L1StateRegistry.deploy(tempL2Aggregator);
  await stateRegistry.waitForDeployment();
  const stateRegistryAddress = await stateRegistry.getAddress();
  console.log("✅ L1StateRegistry deployed to:", stateRegistryAddress);
  console.log("⚠️  Note: L2 aggregator set to deployer address temporarily");
  deployedContracts.contracts.stateRegistry = stateRegistryAddress;
  console.log("");

  // ============ Step 5: Get Arbitrum Bridge Addresses ============
  console.log("📄 Step 5: Getting Arbitrum Bridge Addresses...");
  let inboxAddress, outboxAddress;

  if (network === "mainnet") {
    inboxAddress = process.env.ARBITRUM_INBOX || "0x4Dbd4fc535Ac27206064B68FfCf827b0A60BAB3f";
    outboxAddress = process.env.ARBITRUM_OUTBOX || "0x0B9857ae2D4A3DBe74ffE1d7DF045bb7F96E4840";
  } else if (network === "sepolia") {
    inboxAddress = process.env.ARBITRUM_SEPOLIA_INBOX || "0xaAe29B0366299461418F5324a79Afc425BE5ae21";
    outboxAddress = process.env.ARBITRUM_SEPOLIA_OUTBOX || "0x65f07C7D521164a4d5DaC6eB8Fac8DA067A3B78F";
  } else {
    // For local testing, use mock addresses
    console.log("⚠️  Using mock bridge addresses for local network");
    inboxAddress = "0x0000000000000000000000000000000000000001";
    outboxAddress = "0x0000000000000000000000000000000000000002";
  }

  console.log("Inbox address:", inboxAddress);
  console.log("Outbox address:", outboxAddress);
  deployedContracts.arbitrumBridge = { inbox: inboxAddress, outbox: outboxAddress };
  console.log("");

  // ============ Step 6: Deploy L1Gateway ============
  console.log("📄 Step 6: Deploying L1Gateway...");
  const L1Gateway = await hre.ethers.getContractFactory("L1Gateway");
  const l1Gateway = await L1Gateway.deploy(
    collateralVaultAddress,
    loyaltyUSDAddress,
    collateralTokenAddress,
    inboxAddress,
    outboxAddress
  );
  await l1Gateway.waitForDeployment();
  const l1GatewayAddress = await l1Gateway.getAddress();
  console.log("✅ L1Gateway deployed to:", l1GatewayAddress);
  deployedContracts.contracts.l1Gateway = l1GatewayAddress;
  console.log("");

  // ============ Step 7: Setup Permissions ============
  console.log("📄 Step 7: Setting up permissions...");

  // Grant BRIDGE_ROLE to L1Gateway in LoyaltyUSD
  console.log("Granting BRIDGE_ROLE to L1Gateway in LoyaltyUSD...");
  const BRIDGE_ROLE = await loyaltyUSD.BRIDGE_ROLE();
  await loyaltyUSD.grantRole(BRIDGE_ROLE, l1GatewayAddress);
  console.log("✅ BRIDGE_ROLE granted");

  // Set L2Bridge in CollateralVault
  console.log("Setting L2Bridge in CollateralVault...");
  await collateralVault.setL2Bridge(l1GatewayAddress);
  console.log("✅ L2Bridge set");

  // Set StateRegistry in CollateralVault
  console.log("Setting StateRegistry in CollateralVault...");
  await collateralVault.setStateRegistry(stateRegistryAddress);
  console.log("✅ StateRegistry set");

  console.log("");

  // ============ Step 8: Verify Configuration ============
  console.log("📄 Step 8: Verifying configuration...");
  const canMint = await loyaltyUSD.canMint(l1GatewayAddress);
  const vaultBridge = await collateralVault.l2Bridge();
  console.log("L1Gateway can mint LUSD:", canMint);
  console.log("CollateralVault L2 bridge:", vaultBridge);
  console.log("");

  // ============ Save Deployment Info ============
  const timestamp = new Date().toISOString().replace(/:/g, "-").split(".")[0];
  const filename = `deployments/l1-${network}-${timestamp}.json`;

  try {
    writeFileSync(filename, JSON.stringify(deployedContracts, null, 2));
    console.log("✅ Deployment info saved to:", filename);
  } catch (error) {
    console.log("⚠️  Could not save to file, creating deployments directory...");
    const fs = await import("fs");
    fs.mkdirSync("deployments", { recursive: true });
    writeFileSync(filename, JSON.stringify(deployedContracts, null, 2));
    console.log("✅ Deployment info saved to:", filename);
  }

  // ============ Summary ============
  console.log("\n====================================");
  console.log("🎉 L1 Deployment Complete!");
  console.log("====================================");
  console.log("\n📋 Deployed Contracts:");
  console.log("├─ LoyaltyUSDL1:", loyaltyUSDAddress);
  console.log("├─ Collateral Token:", collateralTokenAddress);
  console.log("├─ CollateralVaultL1:", collateralVaultAddress);
  console.log("├─ L1StateRegistry:", stateRegistryAddress);
  console.log("└─ L1Gateway:", l1GatewayAddress);

  console.log("\n⚠️  Next Steps:");
  console.log("1. Deploy L2 contracts using: npx hardhat run scripts/layer2/deploy-l2.js --network arbitrumSepolia");
  console.log("2. Update L1StateRegistry with L2 aggregator address");
  console.log("3. Update L1Gateway with L2 gateway address");
  console.log("4. Update .env file with deployed addresses");
  console.log("\n====================================\n");

  // Return addresses for programmatic use
  return deployedContracts;
}

// Execute deployment
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
