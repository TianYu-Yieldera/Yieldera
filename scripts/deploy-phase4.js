/**
 * Phase 4 Deployment Script
 * Deploys modular vault sub-modules and VaultModuleV3 coordinator
 */

const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("=".repeat(60));
  console.log("Phase 4: Modular Architecture Deployment");
  console.log("=".repeat(60));

  const [deployer] = await ethers.getSigners();
  console.log("\nDeploying with account:", deployer.address);

  const balance = await ethers.provider.getBalance(deployer.address);
  console.log("Account balance:", ethers.formatEther(balance), "ETH");

  // ============ Deploy Mock Tokens (for testing) ============
  console.log("\n📦 Deploying mock tokens...");

  const MockERC20 = await ethers.getContractFactory("MockERC20");

  const collateralToken = await MockERC20.deploy("Mock Collateral", "mCOL");
  await collateralToken.waitForDeployment();
  console.log("✓ Collateral Token:", await collateralToken.getAddress());

  const debtToken = await MockERC20.deploy("Mock Debt", "mDEBT");
  await debtToken.waitForDeployment();
  console.log("✓ Debt Token:", await debtToken.getAddress());

  // ============ Deploy Vault Sub-Modules ============
  console.log("\n🔧 Deploying Vault Sub-Modules...");

  // 1. CollateralManager
  const CollateralManager = await ethers.getContractFactory("CollateralManager");
  const collateralManager = await CollateralManager.deploy(await collateralToken.getAddress());
  await collateralManager.waitForDeployment();
  console.log("✓ CollateralManager:", await collateralManager.getAddress());

  // 2. PositionManager
  const PositionManager = await ethers.getContractFactory("PositionManager");
  const positionManager = await PositionManager.deploy();
  await positionManager.waitForDeployment();
  console.log("✓ PositionManager:", await positionManager.getAddress());

  // 3. DebtManager
  const DebtManager = await ethers.getContractFactory("DebtManager");
  const debtManager = await DebtManager.deploy();
  await debtManager.waitForDeployment();
  console.log("✓ DebtManager:", await debtManager.getAddress());

  // 4. InterestCalculator (with 2% annual fee = 200 basis points)
  const InterestCalculator = await ethers.getContractFactory("InterestCalculator");
  const interestCalculator = await InterestCalculator.deploy(200);
  await interestCalculator.waitForDeployment();
  console.log("✓ InterestCalculator:", await interestCalculator.getAddress());

  // 5. LiquidationEngine (threshold: 120%, penalty: 10%)
  const LiquidationEngine = await ethers.getContractFactory("LiquidationEngine");
  const liquidationEngine = await LiquidationEngine.deploy(120, 10);
  await liquidationEngine.waitForDeployment();
  console.log("✓ LiquidationEngine:", await liquidationEngine.getAddress());

  // ============ Deploy Legacy Vault (mock) ============
  console.log("\n📦 Deploying legacy vault (mock)...");

  // For this demo, we'll use CollateralManager address as legacy vault
  const legacyVault = await collateralManager.getAddress();
  console.log("✓ Legacy Vault (mock):", legacyVault);

  // ============ Deploy VaultModuleV3 with Proxy ============
  console.log("\n🚀 Deploying VaultModuleV3 Coordinator...");

  const VaultModuleV3 = await ethers.getContractFactory("VaultModuleV3");

  const vaultModuleV3 = await upgrades.deployProxy(
    VaultModuleV3,
    [
      await collateralManager.getAddress(),
      await positionManager.getAddress(),
      await debtManager.getAddress(),
      await interestCalculator.getAddress(),
      await liquidationEngine.getAddress(),
      legacyVault,
      await debtToken.getAddress()
    ],
    {
      kind: "uups",
      initializer: "initialize"
    }
  );
  await vaultModuleV3.waitForDeployment();

  const vaultProxyAddress = await vaultModuleV3.getAddress();
  const vaultImplAddress = await upgrades.erc1967.getImplementationAddress(vaultProxyAddress);

  console.log("✓ VaultModuleV3 Proxy:", vaultProxyAddress);
  console.log("✓ VaultModuleV3 Implementation:", vaultImplAddress);

  // ============ Configure Sub-Modules ============
  console.log("\n⚙️  Configuring sub-modules...");

  // Set vault module address in each sub-module
  await collateralManager.setVaultModule(vaultProxyAddress);
  console.log("✓ CollateralManager configured");

  await positionManager.setVaultModule(vaultProxyAddress);
  console.log("✓ PositionManager configured");

  await debtManager.setVaultModule(vaultProxyAddress);
  console.log("✓ DebtManager configured");

  await interestCalculator.setVaultModule(vaultProxyAddress);
  console.log("✓ InterestCalculator configured");

  await liquidationEngine.setVaultModule(vaultProxyAddress);
  console.log("✓ LiquidationEngine configured");

  // ============ Verify Health ============
  console.log("\n🏥 Verifying system health...");

  const [healthy, message] = await vaultModuleV3.healthCheck();
  console.log("Health Status:", healthy ? "✓ HEALTHY" : "✗ UNHEALTHY");
  console.log("Health Message:", message);

  const moduleInfo = await vaultModuleV3.getModuleInfo();
  console.log("\nModule Info:");
  console.log("  Name:", moduleInfo.name);
  console.log("  Version:", moduleInfo.version);
  console.log("  State:", moduleInfo.state === 0 ? "INACTIVE" : moduleInfo.state === 1 ? "ACTIVE" : "PAUSED");

  // ============ Deploy RWA Sub-Modules ============
  console.log("\n🔧 Deploying RWA Sub-Modules...");

  // 1. OrderManager
  const OrderManager = await ethers.getContractFactory("OrderManager");
  const orderManager = await OrderManager.deploy();
  await orderManager.waitForDeployment();
  console.log("✓ OrderManager:", await orderManager.getAddress());

  // 2. MatchingEngine
  const MatchingEngine = await ethers.getContractFactory("MatchingEngine");
  const matchingEngine = await MatchingEngine.deploy(await orderManager.getAddress());
  await matchingEngine.waitForDeployment();
  console.log("✓ MatchingEngine:", await matchingEngine.getAddress());

  // 3. MarketDataProvider
  const MarketDataProvider = await ethers.getContractFactory("MarketDataProvider");
  const marketDataProvider = await MarketDataProvider.deploy();
  await marketDataProvider.waitForDeployment();
  console.log("✓ MarketDataProvider:", await marketDataProvider.getAddress());

  // 4. FeeCalculator (makerFee: 25 bp = 0.25%, takerFee: 50 bp = 0.50%)
  const FeeCalculator = await ethers.getContractFactory("FeeCalculator");
  const feeCalculator = await FeeCalculator.deploy(
    await orderManager.getAddress(),
    await collateralToken.getAddress(),
    25,
    50
  );
  await feeCalculator.waitForDeployment();
  console.log("✓ FeeCalculator:", await feeCalculator.getAddress());

  // 5. LiquidityAnalyzer
  const LiquidityAnalyzer = await ethers.getContractFactory("LiquidityAnalyzer");
  const liquidityAnalyzer = await LiquidityAnalyzer.deploy(await orderManager.getAddress());
  await liquidityAnalyzer.waitForDeployment();
  console.log("✓ LiquidityAnalyzer:", await liquidityAnalyzer.getAddress());

  // ============ Deploy RWAModuleV3 with Proxy ============
  console.log("\n🚀 Deploying RWAModuleV3 Coordinator...");

  // Mock tokens for RWA trading
  const rwaToken = await MockERC20.deploy("Mock RWA Token", "mRWA");
  await rwaToken.waitForDeployment();
  console.log("✓ RWA Token:", await rwaToken.getAddress());

  const legacyOrderBook = await orderManager.getAddress(); // Use OrderManager as legacy

  const RWAModuleV3 = await ethers.getContractFactory("RWAModuleV3");

  const rwaModuleV3 = await upgrades.deployProxy(
    RWAModuleV3,
    [
      await orderManager.getAddress(),
      await matchingEngine.getAddress(),
      await marketDataProvider.getAddress(),
      await feeCalculator.getAddress(),
      await liquidityAnalyzer.getAddress(),
      await rwaToken.getAddress(),
      await collateralToken.getAddress(),
      legacyOrderBook
    ],
    {
      kind: "uups",
      initializer: "initialize"
    }
  );
  await rwaModuleV3.waitForDeployment();

  const rwaProxyAddress = await rwaModuleV3.getAddress();
  const rwaImplAddress = await upgrades.erc1967.getImplementationAddress(rwaProxyAddress);

  console.log("✓ RWAModuleV3 Proxy:", rwaProxyAddress);
  console.log("✓ RWAModuleV3 Implementation:", rwaImplAddress);

  // ============ Configure RWA Sub-Modules ============
  console.log("\n⚙️  Configuring RWA sub-modules...");

  await orderManager.setRWAModule(rwaProxyAddress);
  console.log("✓ OrderManager configured");

  await matchingEngine.setRWAModule(rwaProxyAddress);
  console.log("✓ MatchingEngine configured");

  await marketDataProvider.setRWAModule(rwaProxyAddress);
  console.log("✓ MarketDataProvider configured");

  await feeCalculator.setRWAModule(rwaProxyAddress);
  console.log("✓ FeeCalculator configured");

  await liquidityAnalyzer.setRWAModule(rwaProxyAddress);
  console.log("✓ LiquidityAnalyzer configured");

  // ============ Verify RWA Module Health ============
  console.log("\n🏥 Verifying RWA module health...");

  const [rwaHealthy, rwaMessage] = await rwaModuleV3.healthCheck();
  console.log("Health Status:", rwaHealthy ? "✓ HEALTHY" : "✗ UNHEALTHY");
  console.log("Health Message:", rwaMessage);

  const rwaModuleInfo = await rwaModuleV3.getModuleInfo();
  console.log("\nRWA Module Info:");
  console.log("  Name:", rwaModuleInfo.name);
  console.log("  Version:", rwaModuleInfo.version);
  console.log("  State:", rwaModuleInfo.state === 0 ? "INACTIVE" : rwaModuleInfo.state === 1 ? "ACTIVE" : "PAUSED");

  // ============ Summary ============
  console.log("\n" + "=".repeat(60));
  console.log("Deployment Summary");
  console.log("=".repeat(60));

  const summary = {
    "Tokens": {
      "Collateral Token": await collateralToken.getAddress(),
      "Debt Token": await debtToken.getAddress(),
      "RWA Token": await rwaToken.getAddress()
    },
    "Vault Sub-Modules": {
      "CollateralManager": await collateralManager.getAddress(),
      "PositionManager": await positionManager.getAddress(),
      "DebtManager": await debtManager.getAddress(),
      "InterestCalculator": await interestCalculator.getAddress(),
      "LiquidationEngine": await liquidationEngine.getAddress()
    },
    "VaultModuleV3": {
      "Proxy": vaultProxyAddress,
      "Implementation": vaultImplAddress
    },
    "RWA Sub-Modules": {
      "OrderManager": await orderManager.getAddress(),
      "MatchingEngine": await matchingEngine.getAddress(),
      "MarketDataProvider": await marketDataProvider.getAddress(),
      "FeeCalculator": await feeCalculator.getAddress(),
      "LiquidityAnalyzer": await liquidityAnalyzer.getAddress()
    },
    "RWAModuleV3": {
      "Proxy": rwaProxyAddress,
      "Implementation": rwaImplAddress
    }
  };

  console.log(JSON.stringify(summary, null, 2));

  console.log("\n" + "=".repeat(60));
  console.log("✅ Phase 4 Deployment Complete!");
  console.log("=".repeat(60));

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
