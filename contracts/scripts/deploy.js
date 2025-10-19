const hre = require("hardhat");

async function main() {
  console.log("🚀 Starting deployment to", hre.network.name);
  console.log("=".repeat(50));

  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 Deploying contracts with account:", deployer.address);

  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 Account balance:", hre.ethers.formatEther(balance), "ETH");
  console.log("=".repeat(50));

  // Deploy LoyaltyUSD (LUSD)
  console.log("\n📄 Deploying LoyaltyUSD...");
  const LoyaltyUSD = await hre.ethers.getContractFactory("LoyaltyUSD");
  const lusd = await LoyaltyUSD.deploy();
  await lusd.waitForDeployment();
  const lusdAddress = await lusd.getAddress();
  console.log("✅ LoyaltyUSD deployed to:", lusdAddress);

  // Deploy CollateralVault
  console.log("\n📄 Deploying CollateralVault...");
  const CollateralVault = await hre.ethers.getContractFactory("CollateralVault");
  const vault = await CollateralVault.deploy(lusdAddress); // Using LUSD as collateral for testing
  await vault.waitForDeployment();
  const vaultAddress = await vault.getAddress();
  console.log("✅ CollateralVault deployed to:", vaultAddress);

  // Grant roles
  console.log("\n🔐 Setting up roles...");
  const MINTER_ROLE = await lusd.MINTER_ROLE();
  const BURNER_ROLE = await lusd.BURNER_ROLE();

  console.log("   Granting MINTER_ROLE to vault...");
  const mintTx = await lusd.grantRole(MINTER_ROLE, vaultAddress);
  await mintTx.wait();
  console.log("   ✅ MINTER_ROLE granted");

  console.log("   Granting BURNER_ROLE to vault...");
  const burnTx = await lusd.grantRole(BURNER_ROLE, vaultAddress);
  await burnTx.wait();
  console.log("   ✅ BURNER_ROLE granted");

  // Verify roles
  console.log("\n🔍 Verifying roles...");
  const hasMinter = await lusd.hasRole(MINTER_ROLE, vaultAddress);
  const hasBurner = await lusd.hasRole(BURNER_ROLE, vaultAddress);
  console.log("   Vault has MINTER_ROLE:", hasMinter);
  console.log("   Vault has BURNER_ROLE:", hasBurner);

  // Display contract info
  console.log("\n" + "=".repeat(50));
  console.log("📋 DEPLOYMENT SUMMARY");
  console.log("=".repeat(50));
  console.log("Network:", hre.network.name);
  console.log("Deployer:", deployer.address);
  console.log("\nContracts:");
  console.log("  LoyaltyUSD (LUSD):", lusdAddress);
  console.log("  CollateralVault:  ", vaultAddress);
  console.log("\nNext steps:");
  console.log("  1. Verify contracts on Etherscan");
  console.log("  2. Update .env files with addresses");
  console.log("  3. Test mint/redeem flow");
  console.log("=".repeat(50));

  // Save addresses to file
  const fs = require('fs');
  const addresses = {
    network: hre.network.name,
    deployer: deployer.address,
    lusd: lusdAddress,
    vault: vaultAddress,
    timestamp: new Date().toISOString()
  };

  const outputPath = `deployments-${hre.network.name}.json`;
  fs.writeFileSync(outputPath, JSON.stringify(addresses, null, 2));
  console.log("\n💾 Addresses saved to:", outputPath);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Deployment failed:");
    console.error(error);
    process.exit(1);
  });
