const hre = require("hardhat");

/**
 * Deploy Treasury Module Contracts
 *
 * Deployment Order:
 * 1. TreasuryAssetFactory
 * 2. TreasuryPriceOracle
 * 3. TreasuryYieldDistributor
 * 4. TreasuryMarketplace
 */

async function main() {
  console.log("🚀 Deploying Treasury Module to Arbitrum Sepolia...\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString(), "\n");

  // Configuration
  const admin = deployer.address;
  const feeCollector = deployer.address; // Can be changed later

  // Mock USDC for testnet (or use existing USDC if available)
  let usdcAddress;

  // Check if USDC is already deployed
  const USDC_ADDRESS_ARBITRUM_SEPOLIA = "0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d"; // Arbitrum Sepolia USDC

  try {
    const usdc = await hre.ethers.getContractAt("IERC20", USDC_ADDRESS_ARBITRUM_SEPOLIA);
    const symbol = await usdc.symbol();
    console.log(`✅ Using existing USDC at ${USDC_ADDRESS_ARBITRUM_SEPOLIA}`);
    usdcAddress = USDC_ADDRESS_ARBITRUM_SEPOLIA;
  } catch (error) {
    console.log("📝 Deploying Mock USDC for testing...");
    const MockERC20 = await hre.ethers.getContractFactory("MockERC20");
    const mockUsdc = await MockERC20.deploy("USD Coin", "USDC", 6); // 6 decimals like real USDC
    await mockUsdc.waitForDeployment();
    usdcAddress = await mockUsdc.getAddress();
    console.log(`✅ Mock USDC deployed at: ${usdcAddress}\n`);
  }

  // ============================================
  // 1. Deploy TreasuryAssetFactory
  // ============================================
  console.log("1️⃣ Deploying TreasuryAssetFactory...");
  const TreasuryAssetFactory = await hre.ethers.getContractFactory("TreasuryAssetFactory");
  const factory = await TreasuryAssetFactory.deploy(admin);
  await factory.waitForDeployment();
  const factoryAddress = await factory.getAddress();
  console.log(`✅ TreasuryAssetFactory deployed at: ${factoryAddress}\n`);

  // ============================================
  // 2. Deploy TreasuryPriceOracle
  // ============================================
  console.log("2️⃣ Deploying TreasuryPriceOracle...");
  const TreasuryPriceOracle = await hre.ethers.getContractFactory("TreasuryPriceOracle");
  const oracle = await TreasuryPriceOracle.deploy(admin, factoryAddress);
  await oracle.waitForDeployment();
  const oracleAddress = await oracle.getAddress();
  console.log(`✅ TreasuryPriceOracle deployed at: ${oracleAddress}\n`);

  // ============================================
  // 3. Deploy TreasuryYieldDistributor
  // ============================================
  console.log("3️⃣ Deploying TreasuryYieldDistributor...");
  const TreasuryYieldDistributor = await hre.ethers.getContractFactory("TreasuryYieldDistributor");
  const yieldDistributor = await TreasuryYieldDistributor.deploy(
    admin,
    factoryAddress,
    usdcAddress
  );
  await yieldDistributor.waitForDeployment();
  const yieldDistributorAddress = await yieldDistributor.getAddress();
  console.log(`✅ TreasuryYieldDistributor deployed at: ${yieldDistributorAddress}\n`);

  // ============================================
  // 4. Deploy TreasuryMarketplace
  // ============================================
  console.log("4️⃣ Deploying TreasuryMarketplace...");
  const TreasuryMarketplace = await hre.ethers.getContractFactory("TreasuryMarketplace");
  const marketplace = await TreasuryMarketplace.deploy(
    admin,
    factoryAddress,
    usdcAddress,
    feeCollector
  );
  await marketplace.waitForDeployment();
  const marketplaceAddress = await marketplace.getAddress();
  console.log(`✅ TreasuryMarketplace deployed at: ${marketplaceAddress}\n`);

  // ============================================
  // 5. Create Sample Treasury Assets
  // ============================================
  console.log("5️⃣ Creating sample treasury assets...\n");

  const currentTime = Math.floor(Date.now() / 1000);
  const oneDay = 24 * 60 * 60;

  // T-Bill 13 Week
  console.log("Creating T-Bill 13W...");
  const tBillTx = await factory.createTreasuryAsset(
    0, // TreasuryType.T_BILL
    "13W",
    "912796TB1",
    currentTime,
    currentTime + (91 * oneDay), // 13 weeks
    hre.ethers.parseEther("1000"), // $1000 face value
    525 // 5.25% coupon
  );
  await tBillTx.wait();
  console.log("✅ T-Bill 13W created\n");

  // T-Note 2 Year
  console.log("Creating T-Note 2Y...");
  const tNote2YTx = await factory.createTreasuryAsset(
    1, // TreasuryType.T_NOTE
    "2Y",
    "91282CHX6",
    currentTime,
    currentTime + (730 * oneDay), // 2 years
    hre.ethers.parseEther("1000"),
    450 // 4.50%
  );
  await tNote2YTx.wait();
  console.log("✅ T-Note 2Y created\n");

  // T-Note 10 Year
  console.log("Creating T-Note 10Y...");
  const tNote10YTx = await factory.createTreasuryAsset(
    1, // TreasuryType.T_NOTE
    "10Y",
    "912828YK4",
    currentTime,
    currentTime + (3650 * oneDay), // 10 years
    hre.ethers.parseEther("1000"),
    425 // 4.25%
  );
  await tNote10YTx.wait();
  console.log("✅ T-Note 10Y created\n");

  // T-Bond 30 Year
  console.log("Creating T-Bond 30Y...");
  const tBond30YTx = await factory.createTreasuryAsset(
    2, // TreasuryType.T_BOND
    "30Y",
    "912810TT4",
    currentTime,
    currentTime + (10950 * oneDay), // 30 years
    hre.ethers.parseEther("1000"),
    400 // 4.00%
  );
  await tBond30YTx.wait();
  console.log("✅ T-Bond 30Y created\n");

  // ============================================
  // 6. Initialize Sample Prices
  // ============================================
  console.log("6️⃣ Setting initial prices...\n");

  // Asset 1: T-Bill 13W
  await oracle.updatePrice(
    1,
    hre.ethers.parseEther("980"), // $980 price
    540, // 5.40% yield
    "INITIAL_SETUP"
  );
  console.log("✅ Price set for T-Bill 13W\n");

  // Asset 2: T-Note 2Y
  await oracle.updatePrice(
    2,
    hre.ethers.parseEther("985"),
    465, // 4.65%
    "INITIAL_SETUP"
  );
  console.log("✅ Price set for T-Note 2Y\n");

  // Asset 3: T-Note 10Y
  await oracle.updatePrice(
    3,
    hre.ethers.parseEther("950"),
    475, // 4.75%
    "INITIAL_SETUP"
  );
  console.log("✅ Price set for T-Note 10Y\n");

  // Asset 4: T-Bond 30Y
  await oracle.updatePrice(
    4,
    hre.ethers.parseEther("920"),
    450, // 4.50%
    "INITIAL_SETUP"
  );
  console.log("✅ Price set for T-Bond 30Y\n");

  // ============================================
  // 7. Summary
  // ============================================
  console.log("=" .repeat(60));
  console.log("🎉 Treasury Module Deployment Complete!");
  console.log("=" .repeat(60));
  console.log("\n📋 Contract Addresses:\n");
  console.log(`USDC (Payment Token):        ${usdcAddress}`);
  console.log(`TreasuryAssetFactory:        ${factoryAddress}`);
  console.log(`TreasuryPriceOracle:         ${oracleAddress}`);
  console.log(`TreasuryYieldDistributor:    ${yieldDistributorAddress}`);
  console.log(`TreasuryMarketplace:         ${marketplaceAddress}`);

  console.log("\n📊 Sample Assets Created:\n");
  console.log("Asset ID 1: T-Bill 13W  (CUSIP: 912796TB1)");
  console.log("Asset ID 2: T-Note 2Y   (CUSIP: 91282CHX6)");
  console.log("Asset ID 3: T-Note 10Y  (CUSIP: 912828YK4)");
  console.log("Asset ID 4: T-Bond 30Y  (CUSIP: 912810TT4)");

  console.log("\n💡 Next Steps:\n");
  console.log("1. Update .env with contract addresses");
  console.log("2. Run database migration: psql -f db/migrations/003_treasury_module.sql");
  console.log("3. Update backend services to listen to Treasury events");
  console.log("4. Configure price oracle data source");
  console.log("5. Deploy frontend with Treasury module\n");

  // Save deployment info to file
  const fs = require('fs');
  const deploymentInfo = {
    network: hre.network.name,
    timestamp: new Date().toISOString(),
    deployer: deployer.address,
    contracts: {
      USDC: usdcAddress,
      TreasuryAssetFactory: factoryAddress,
      TreasuryPriceOracle: oracleAddress,
      TreasuryYieldDistributor: yieldDistributorAddress,
      TreasuryMarketplace: marketplaceAddress
    },
    sampleAssets: [
      { id: 1, type: "T-BILL", term: "13W", cusip: "912796TB1" },
      { id: 2, type: "T-NOTE", term: "2Y", cusip: "91282CHX6" },
      { id: 3, type: "T-NOTE", term: "10Y", cusip: "912828YK4" },
      { id: 4, type: "T-BOND", term: "30Y", cusip: "912810TT4" }
    ]
  };

  fs.writeFileSync(
    'deployments/treasury-deployment.json',
    JSON.stringify(deploymentInfo, null, 2)
  );
  console.log("✅ Deployment info saved to deployments/treasury-deployment.json\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
