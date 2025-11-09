const hre = require("hardhat");
const fs = require('fs');
const path = require('path');

async function main() {
  console.log("\n🚀 ============================================");
  console.log("   Deploying Treasury Module to Arbitrum Sepolia");
  console.log("============================================\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);

  const balance = await deployer.provider.getBalance(deployer.address);
  console.log("Account balance:", hre.ethers.formatEther(balance), "ETH");

  if (parseFloat(hre.ethers.formatEther(balance)) < 0.05) {
    console.log("\n❌ Insufficient balance!");
    console.log("Minimum required: 0.05 ETH");
    console.log("Please get test ETH from: https://faucet.quicknode.com/arbitrum/sepolia");
    process.exit(1);
  }

  console.log("\n📝 Configuration:");
  const admin = deployer.address;
  const feeCollector = deployer.address;
  console.log("Admin:", admin);
  console.log("Fee Collector:", feeCollector);

  // Use Arbitrum Sepolia USDC or deploy mock
  let usdcAddress;
  const USDC_ARB_SEPOLIA = "0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d";

  try {
    const code = await deployer.provider.getCode(USDC_ARB_SEPOLIA);
    if (code !== '0x') {
      console.log("\n✅ Using existing USDC at:", USDC_ARB_SEPOLIA);
      usdcAddress = USDC_ARB_SEPOLIA;
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
    chainId: (await deployer.provider.getNetwork()).chainId.toString(),
    timestamp: new Date().toISOString(),
    deployer: deployer.address,
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

  // Create sample assets
  console.log("\n5️⃣ Creating sample Treasury assets...");
  const currentTime = Math.floor(Date.now() / 1000);
  const oneDay = 24 * 60 * 60;

  const assets = [
    { type: 0, term: "13W", cusip: "912796TB1", maturity: 91, coupon: 525, name: "T-Bill 13W" },
    { type: 1, term: "2Y", cusip: "91282CHX6", maturity: 730, coupon: 450, name: "T-Note 2Y" },
    { type: 1, term: "10Y", cusip: "912828YK4", maturity: 3650, coupon: 425, name: "T-Note 10Y" },
    { type: 2, term: "30Y", cusip: "912810TT4", maturity: 10950, coupon: 400, name: "T-Bond 30Y" }
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
      hre.ethers.parseEther("1000"),
      asset.coupon
    );
    await tx.wait();
    deployments.sampleAssets.push({
      name: asset.name,
      cusip: asset.cusip,
      term: asset.term
    });
    console.log(`  ✅ ${asset.name} created`);
  }

  // Set initial prices
  console.log("\n6️⃣ Setting initial prices...");
  const prices = [
    { assetId: 1, price: "980", yield: 540 },
    { assetId: 2, price: "985", yield: 465 },
    { assetId: 3, price: "950", yield: 475 },
    { assetId: 4, price: "920", yield: 450 }
  ];

  for (const p of prices) {
    await oracle.updatePrice(
      p.assetId,
      hre.ethers.parseEther(p.price),
      p.yield,
      "INITIAL_SETUP"
    );
  }
  console.log("✅ All prices set");

  // Save deployment info
  const deploymentsDir = path.join(__dirname, '../deployments');
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  const timestamp = Date.now();
  const filename = `treasury-${hre.network.name}-${timestamp}.json`;
  const filepath = path.join(deploymentsDir, filename);
  fs.writeFileSync(filepath, JSON.stringify(deployments, null, 2));

  console.log("\n" + "=".repeat(60));
  console.log("🎉 Deployment Complete!");
  console.log("=".repeat(60));
  console.log("\n📋 Contract Addresses:\n");
  console.log("USDC (Payment Token):     ", usdcAddress);
  console.log("TreasuryAssetFactory:     ", factoryAddress);
  console.log("TreasuryPriceOracle:      ", oracleAddress);
  console.log("TreasuryYieldDistributor: ", yieldAddress);
  console.log("TreasuryMarketplace:      ", marketplaceAddress);

  console.log("\n📊 Sample Assets:\n");
  deployments.sampleAssets.forEach((a, i) => {
    console.log(`Asset ${i+1}: ${a.name} (CUSIP: ${a.cusip})`);
  });

  console.log("\n📁 Deployment saved to:", filename);
  console.log("\n💡 Next Steps:");
  console.log("1. Update .env with contract addresses");
  console.log("2. Update frontend configuration");
  console.log("3. Start backend services");
  console.log("4. Test the complete flow\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ Deployment failed:");
    console.error(error);
    process.exit(1);
  });
