const hre = require("hardhat");
const fs = require('fs');
const path = require('path');

async function main() {
  console.log("\n🚀 ============================================");
  console.log("   Deploying Treasury Module to Base Sepolia");
  console.log("   Chain ID: 84532");
  console.log("   Purpose: Conservative RWA & US Treasury Bonds");
  console.log("============================================\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);

  const balance = await deployer.provider.getBalance(deployer.address);
  console.log("Account balance:", hre.ethers.formatEther(balance), "ETH");

  const network = await deployer.provider.getNetwork();
  const chainId = network.chainId.toString();

  console.log("Network:", hre.network.name);
  console.log("Chain ID:", chainId);

  if (chainId !== "84532") {
    console.log("\n❌ Wrong network! This script is for Base Sepolia (84532)");
    console.log("Current network:", chainId);
    console.log("\nPlease run: npx hardhat run scripts/deploy-treasury-base.js --network baseSepolia");
    process.exit(1);
  }

  if (parseFloat(hre.ethers.formatEther(balance)) < 0.01) {
    console.log("\n❌ Insufficient balance!");
    console.log("Minimum required: 0.01 ETH");
    console.log("Please get test ETH from: https://www.coinbase.com/faucets/base-ethereum-sepolia-faucet");
    process.exit(1);
  }

  console.log("\n📝 Configuration:");
  const admin = deployer.address;
  const feeCollector = deployer.address;
  console.log("Admin:", admin);
  console.log("Fee Collector:", feeCollector);

  // Base Sepolia USDC address
  // If not available, we'll deploy a mock
  let usdcAddress;
  const USDC_BASE_SEPOLIA = "0x036CbD53842c5426634e7929541eC2318f3dCF7e"; // Official Base Sepolia USDC

  try {
    const code = await deployer.provider.getCode(USDC_BASE_SEPOLIA);
    if (code !== '0x') {
      console.log("\n✅ Using existing USDC at:", USDC_BASE_SEPOLIA);
      usdcAddress = USDC_BASE_SEPOLIA;
    } else {
      throw new Error("USDC not found");
    }
  } catch (error) {
    console.log("\n📝 Deploying Mock USDC...");
    const MockERC20 = await hre.ethers.getContractFactory("MockERC20");
    const mockUsdc = await MockERC20.deploy("USD Coin", "USDC", 6);
    await mockUsdc.waitForDeployment();
    usdcAddress = await mockUsdc.getAddress();
    console.log("✅ Mock USDC deployed at:", usdcAddress);
  }

  const deployments = {
    network: hre.network.name,
    chainId: chainId,
    timestamp: new Date().toISOString(),
    deployer: deployer.address,
    purpose: "Conservative RWA & US Treasury on Base",
    contracts: {}
  };

  // Deploy TreasuryAssetFactory
  console.log("\n1️⃣ Deploying TreasuryAssetFactory...");
  const TreasuryAssetFactory = await hre.ethers.getContractFactory("TreasuryAssetFactory");
  const factory = await TreasuryAssetFactory.deploy(admin);
  await factory.waitForDeployment();
  const factoryAddress = await factory.getAddress();
  deployments.contracts.TreasuryAssetFactory = factoryAddress;
  console.log("✅ Deployed at:", factoryAddress);

  // Deploy TreasuryPriceOracle
  console.log("\n2️⃣ Deploying TreasuryPriceOracle...");
  const TreasuryPriceOracle = await hre.ethers.getContractFactory("TreasuryPriceOracle");
  const oracle = await TreasuryPriceOracle.deploy(admin, factoryAddress);
  await oracle.waitForDeployment();
  const oracleAddress = await oracle.getAddress();
  deployments.contracts.TreasuryPriceOracle = oracleAddress;
  console.log("✅ Deployed at:", oracleAddress);

  // Deploy TreasuryYieldDistributor
  console.log("\n3️⃣ Deploying TreasuryYieldDistributor...");
  const TreasuryYieldDistributor = await hre.ethers.getContractFactory("TreasuryYieldDistributor");
  const yieldDistributor = await TreasuryYieldDistributor.deploy(admin, factoryAddress, usdcAddress);
  await yieldDistributor.waitForDeployment();
  const yieldAddress = await yieldDistributor.getAddress();
  deployments.contracts.TreasuryYieldDistributor = yieldAddress;
  console.log("✅ Deployed at:", yieldAddress);

  // Deploy TreasuryMarketplace
  console.log("\n4️⃣ Deploying TreasuryMarketplace...");
  const TreasuryMarketplace = await hre.ethers.getContractFactory("TreasuryMarketplace");
  const marketplace = await TreasuryMarketplace.deploy(admin, factoryAddress, usdcAddress, feeCollector);
  await marketplace.waitForDeployment();
  const marketplaceAddress = await marketplace.getAddress();
  deployments.contracts.TreasuryMarketplace = marketplaceAddress;
  deployments.contracts.USDC = usdcAddress;
  console.log("✅ Deployed at:", marketplaceAddress);

  // Create sample Treasury assets (Conservative strategy for Base)
  console.log("\n5️⃣ Creating sample Treasury assets (Conservative)...");
  const currentTime = Math.floor(Date.now() / 1000);
  const oneDay = 24 * 60 * 60;

  // Conservative portfolio: More short-term and mid-term treasuries
  const assets = [
    { type: 0, term: "4W", cusip: "912796TW4", maturity: 28, coupon: 530, name: "T-Bill 4W (Ultra Safe)" },
    { type: 0, term: "13W", cusip: "912796TB1", maturity: 91, coupon: 525, name: "T-Bill 13W" },
    { type: 0, term: "26W", cusip: "912796TC9", maturity: 182, coupon: 520, name: "T-Bill 26W" },
    { type: 1, term: "2Y", cusip: "91282CHX6", maturity: 730, coupon: 450, name: "T-Note 2Y" },
    { type: 1, term: "5Y", cusip: "912828YL2", maturity: 1825, coupon: 435, name: "T-Note 5Y" },
    { type: 1, term: "10Y", cusip: "912828YK4", maturity: 3650, coupon: 425, name: "T-Note 10Y" }
  ];

  deployments.sampleAssets = [];

  for (const asset of assets) {
    console.log(`  Creating ${asset.name}...`);
    const tx = await factory.createTreasuryAsset(
      asset.type,
      asset.term,
      asset.cusip,
      currentTime,
      currentTime + (asset.maturity * oneDay),
      hre.ethers.parseEther("10000"), // Higher liquidity for conservative assets
      asset.coupon
    );
    await tx.wait();
    deployments.sampleAssets.push({
      name: asset.name,
      cusip: asset.cusip,
      term: asset.term,
      type: asset.type === 0 ? "Treasury Bill" : "Treasury Note"
    });
    console.log(`  ✅ ${asset.name} created`);
  }

  // Set initial prices (Conservative pricing)
  console.log("\n6️⃣ Setting initial prices...");
  const prices = [
    { assetId: 1, price: "995", yield: 535 },  // 4W T-Bill (near par)
    { assetId: 2, price: "980", yield: 540 },  // 13W T-Bill
    { assetId: 3, price: "975", yield: 535 },  // 26W T-Bill
    { assetId: 4, price: "985", yield: 465 },  // 2Y T-Note
    { assetId: 5, price: "970", yield: 455 },  // 5Y T-Note
    { assetId: 6, price: "950", yield: 475 }   // 10Y T-Note
  ];

  for (const p of prices) {
    await oracle.updatePrice(
      p.assetId,
      hre.ethers.parseEther(p.price),
      p.yield,
      "BASE_INITIAL_SETUP"
    );
  }
  console.log("✅ All prices set");

  // Save deployment info
  const deploymentsDir = path.join(__dirname, '../deployments');
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  const timestamp = Date.now();
  const filename = `treasury-base-sepolia-${timestamp}.json`;
  const filepath = path.join(deploymentsDir, filename);
  fs.writeFileSync(filepath, JSON.stringify(deployments, null, 2));

  // Also create a "latest" symlink for easy access
  const latestPath = path.join(deploymentsDir, 'treasury-base-sepolia-latest.json');
  fs.writeFileSync(latestPath, JSON.stringify(deployments, null, 2));

  console.log("\n" + "=".repeat(70));
  console.log("🎉 Base Treasury Deployment Complete!");
  console.log("=".repeat(70));
  console.log("\n📋 Contract Addresses (Base Sepolia - Chain ID 84532):\n");
  console.log("USDC (Payment Token):     ", usdcAddress);
  console.log("TreasuryAssetFactory:     ", factoryAddress);
  console.log("TreasuryPriceOracle:      ", oracleAddress);
  console.log("TreasuryYieldDistributor: ", yieldAddress);
  console.log("TreasuryMarketplace:      ", marketplaceAddress);

  console.log("\n📊 Sample Assets (Conservative Portfolio):\n");
  deployments.sampleAssets.forEach((a, i) => {
    console.log(`Asset ${i+1}: ${a.name} (CUSIP: ${a.cusip}, Type: ${a.type})`);
  });

  console.log("\n📁 Deployment saved to:", filename);

  console.log("\n🔧 Environment Variables to Add:\n");
  console.log("# Base Treasury Contracts (Chain ID: 84532)");
  console.log(`BASE_TREASURY_FACTORY=${factoryAddress}`);
  console.log(`BASE_TREASURY_MARKETPLACE=${marketplaceAddress}`);
  console.log(`BASE_TREASURY_YIELD_DISTRIBUTOR=${yieldAddress}`);
  console.log(`BASE_TREASURY_ORACLE=${oracleAddress}`);
  console.log(`BASE_USDC=${usdcAddress}`);

  console.log("\n💡 Next Steps:");
  console.log("1. Add above environment variables to .env");
  console.log("2. Update database treasury_holdings with chain_id=84532");
  console.log("3. Create Base event listeners (backend/listeners/treasury/BaseTreasuryListener.ts)");
  console.log("4. Update frontend to show Base treasury assets");
  console.log("5. Test AI adapter integration with MultiChainManager\n");

  console.log("🌐 Base Sepolia Block Explorer:");
  console.log(`   https://sepolia.basescan.org/address/${factoryAddress}\n`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ Deployment failed:");
    console.error(error);
    process.exit(1);
  });
