/**
 * L2 Deployment Script
 * Deploys all business logic modules to Arbitrum L2
 *
 * Deploys:
 * 1. L2StateAggregator (state aggregation)
 * 2. All VaultModuleV3 sub-modules
 * 3. All RWAModuleV3 sub-modules
 * 4. Plugin system contracts
 */

import hre from "hardhat";
import { writeFileSync, readFileSync } from "fs";

async function main() {
  console.log("====================================");
  console.log("🚀 L2 Deployment Script (Arbitrum)");
  console.log("====================================\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);
  console.log("Account balance:", hre.ethers.formatEther(await hre.ethers.provider.getBalance(deployer.address)), "ETH\n");

  const network = hre.network.name;
  console.log("Network:", network);
  console.log("Chain ID:", (await hre.ethers.provider.getNetwork()).chainId);
  console.log("====================================\n");

  // Load L1 deployment info
  let l1StateRegistry;
  try {
    // Try to read from .env
    l1StateRegistry = process.env.L1_STATE_REGISTRY || deployer.address;
    console.log("L1 State Registry address:", l1StateRegistry);
  } catch (error) {
    console.log("⚠️  Could not find L1 deployment info, using deployer address");
    l1StateRegistry = deployer.address;
  }
  console.log("");

  // Deployment results
  const deployedContracts = {
    network,
    chainId: Number((await hre.ethers.provider.getNetwork()).chainId),
    deployedAt: new Date().toISOString(),
    deployer: deployer.address,
    l1StateRegistry,
    contracts: {}
  };

  // ============ Step 1: Deploy L2StateAggregator ============
  console.log("📄 Step 1: Deploying L2StateAggregator...");
  const L2StateAggregator = await hre.ethers.getContractFactory("L2StateAggregator");
  const stateAggregator = await L2StateAggregator.deploy(l1StateRegistry);
  await stateAggregator.waitForDeployment();
  const stateAggregatorAddress = await stateAggregator.getAddress();
  console.log("✅ L2StateAggregator deployed to:", stateAggregatorAddress);
  deployedContracts.contracts.stateAggregator = stateAggregatorAddress;
  console.log("");

  // ============ Step 2: Deploy Mock Tokens (for testing) ============
  console.log("📄 Step 2: Deploying Mock Tokens...");
  const MockERC20 = await hre.ethers.getContractFactory("MockERC20");

  const collateralToken = await MockERC20.deploy(
    "Loyalty Points L2",
    "LP",
    hre.ethers.parseEther("1000000")
  );
  await collateralToken.waitForDeployment();
  const collateralTokenAddress = await collateralToken.getAddress();
  console.log("✅ Collateral Token deployed to:", collateralTokenAddress);

  const debtToken = await MockERC20.deploy(
    "Loyalty USD L2",
    "LUSD",
    hre.ethers.parseUnits("1000000", 6) // 6 decimals for LUSD
  );
  await debtToken.waitForDeployment();
  const debtTokenAddress = await debtToken.getAddress();
  console.log("✅ Debt Token deployed to:", debtTokenAddress);

  deployedContracts.contracts.collateralToken = collateralTokenAddress;
  deployedContracts.contracts.debtToken = debtTokenAddress;
  console.log("");

  // ============ Step 3: Deploy Vault Sub-Modules ============
  console.log("📄 Step 3: Deploying Vault Sub-Modules...");

  console.log("  3.1 Deploying CollateralManager...");
  const CollateralManager = await hre.ethers.getContractFactory("CollateralManager");
  const collateralManager = await CollateralManager.deploy(collateralTokenAddress);
  await collateralManager.waitForDeployment();
  const collateralManagerAddress = await collateralManager.getAddress();
  console.log("  ✅ CollateralManager:", collateralManagerAddress);

  console.log("  3.2 Deploying PositionManager...");
  const PositionManager = await hre.ethers.getContractFactory("PositionManager");
  const positionManager = await PositionManager.deploy();
  await positionManager.waitForDeployment();
  const positionManagerAddress = await positionManager.getAddress();
  console.log("  ✅ PositionManager:", positionManagerAddress);

  console.log("  3.3 Deploying DebtManager...");
  const DebtManager = await hre.ethers.getContractFactory("DebtManager");
  const debtManager = await DebtManager.deploy();
  await debtManager.waitForDeployment();
  const debtManagerAddress = await debtManager.getAddress();
  console.log("  ✅ DebtManager:", debtManagerAddress);

  console.log("  3.4 Deploying InterestCalculator...");
  const InterestCalculator = await hre.ethers.getContractFactory("InterestCalculator");
  const interestCalculator = await InterestCalculator.deploy();
  await interestCalculator.waitForDeployment();
  const interestCalculatorAddress = await interestCalculator.getAddress();
  console.log("  ✅ InterestCalculator:", interestCalculatorAddress);

  console.log("  3.5 Deploying LiquidationEngine...");
  const LiquidationEngine = await hre.ethers.getContractFactory("LiquidationEngine");
  const liquidationEngine = await LiquidationEngine.deploy();
  await liquidationEngine.waitForDeployment();
  const liquidationEngineAddress = await liquidationEngine.getAddress();
  console.log("  ✅ LiquidationEngine:", liquidationEngineAddress);

  deployedContracts.contracts.vaultSubModules = {
    collateralManager: collateralManagerAddress,
    positionManager: positionManagerAddress,
    debtManager: debtManagerAddress,
    interestCalculator: interestCalculatorAddress,
    liquidationEngine: liquidationEngineAddress
  };
  console.log("");

  // ============ Step 4: Deploy VaultModuleV3 ============
  console.log("📄 Step 4: Deploying VaultModuleV3...");

  // Deploy implementation
  const VaultModuleV3 = await hre.ethers.getContractFactory("VaultModuleV3");
  const vaultModuleImpl = await VaultModuleV3.deploy();
  await vaultModuleImpl.waitForDeployment();
  const vaultModuleImplAddress = await vaultModuleImpl.getAddress();
  console.log("  Implementation deployed to:", vaultModuleImplAddress);

  // Deploy proxy
  const ERC1967Proxy = await hre.ethers.getContractFactory("@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol:ERC1967Proxy");

  // Create legacy vault mock
  const legacyVault = deployer.address; // Temporary

  const initData = vaultModuleImpl.interface.encodeFunctionData("initialize", [
    collateralManagerAddress,
    positionManagerAddress,
    debtManagerAddress,
    interestCalculatorAddress,
    liquidationEngineAddress,
    legacyVault,
    debtTokenAddress
  ]);

  const vaultProxy = await ERC1967Proxy.deploy(vaultModuleImplAddress, initData);
  await vaultProxy.waitForDeployment();
  const vaultProxyAddress = await vaultProxy.getAddress();
  console.log("  ✅ VaultModuleV3 Proxy:", vaultProxyAddress);

  deployedContracts.contracts.vaultModule = {
    implementation: vaultModuleImplAddress,
    proxy: vaultProxyAddress
  };
  console.log("");

  // ============ Step 5: Deploy RWA Sub-Modules ============
  console.log("📄 Step 5: Deploying RWA Sub-Modules...");

  console.log("  5.1 Deploying OrderManager...");
  const OrderManager = await hre.ethers.getContractFactory("OrderManager");
  const orderManager = await OrderManager.deploy();
  await orderManager.waitForDeployment();
  const orderManagerAddress = await orderManager.getAddress();
  console.log("  ✅ OrderManager:", orderManagerAddress);

  console.log("  5.2 Deploying MatchingEngine...");
  const MatchingEngine = await hre.ethers.getContractFactory("MatchingEngine");
  const matchingEngine = await MatchingEngine.deploy();
  await matchingEngine.waitForDeployment();
  const matchingEngineAddress = await matchingEngine.getAddress();
  console.log("  ✅ MatchingEngine:", matchingEngineAddress);

  console.log("  5.3 Deploying MarketDataProvider...");
  const MarketDataProvider = await hre.ethers.getContractFactory("MarketDataProvider");
  const marketDataProvider = await MarketDataProvider.deploy();
  await marketDataProvider.waitForDeployment();
  const marketDataProviderAddress = await marketDataProvider.getAddress();
  console.log("  ✅ MarketDataProvider:", marketDataProviderAddress);

  console.log("  5.4 Deploying FeeCalculator...");
  const FeeCalculator = await hre.ethers.getContractFactory("FeeCalculator");
  const feeCalculator = await FeeCalculator.deploy();
  await feeCalculator.waitForDeployment();
  const feeCalculatorAddress = await feeCalculator.getAddress();
  console.log("  ✅ FeeCalculator:", feeCalculatorAddress);

  console.log("  5.5 Deploying LiquidityAnalyzer...");
  const LiquidityAnalyzer = await hre.ethers.getContractFactory("LiquidityAnalyzer");
  const liquidityAnalyzer = await LiquidityAnalyzer.deploy();
  await liquidityAnalyzer.waitForDeployment();
  const liquidityAnalyzerAddress = await liquidityAnalyzer.getAddress();
  console.log("  ✅ LiquidityAnalyzer:", liquidityAnalyzerAddress);

  deployedContracts.contracts.rwaSubModules = {
    orderManager: orderManagerAddress,
    matchingEngine: matchingEngineAddress,
    marketDataProvider: marketDataProviderAddress,
    feeCalculator: feeCalculatorAddress,
    liquidityAnalyzer: liquidityAnalyzerAddress
  };
  console.log("");

  // ============ Step 6: Deploy RWAModuleV3 ============
  console.log("📄 Step 6: Deploying RWAModuleV3...");

  const RWAModuleV3 = await hre.ethers.getContractFactory("RWAModuleV3");
  const rwaModuleImpl = await RWAModuleV3.deploy();
  await rwaModuleImpl.waitForDeployment();
  const rwaModuleImplAddress = await rwaModuleImpl.getAddress();
  console.log("  Implementation deployed to:", rwaModuleImplAddress);

  const legacyOrderBook = deployer.address; // Temporary

  const rwaInitData = rwaModuleImpl.interface.encodeFunctionData("initialize", [
    orderManagerAddress,
    matchingEngineAddress,
    marketDataProviderAddress,
    feeCalculatorAddress,
    liquidityAnalyzerAddress,
    collateralTokenAddress, // base token
    debtTokenAddress, // quote token
    legacyOrderBook
  ]);

  const rwaProxy = await ERC1967Proxy.deploy(rwaModuleImplAddress, rwaInitData);
  await rwaProxy.waitForDeployment();
  const rwaProxyAddress = await rwaProxy.getAddress();
  console.log("  ✅ RWAModuleV3 Proxy:", rwaProxyAddress);

  deployedContracts.contracts.rwaModule = {
    implementation: rwaModuleImplAddress,
    proxy: rwaProxyAddress
  };
  console.log("");

  // ============ Step 7: Register Modules with State Aggregator ============
  console.log("📄 Step 7: Registering modules with StateAggregator...");

  const VAULT_MODULE_ID = hre.ethers.id("VAULT_MODULE");
  const RWA_MODULE_ID = hre.ethers.id("RWA_MODULE");

  await stateAggregator.registerModule(VAULT_MODULE_ID, vaultProxyAddress);
  console.log("  ✅ VaultModule registered");

  await stateAggregator.registerModule(RWA_MODULE_ID, rwaProxyAddress);
  console.log("  ✅ RWAModule registered");

  console.log("");

  // ============ Save Deployment Info ============
  const timestamp = new Date().toISOString().replace(/:/g, "-").split(".")[0];
  const filename = `deployments/l2-${network}-${timestamp}.json`;

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
  console.log("🎉 L2 Deployment Complete!");
  console.log("====================================");
  console.log("\n📋 Deployed Contracts:");
  console.log("├─ L2StateAggregator:", stateAggregatorAddress);
  console.log("├─ VaultModuleV3:", vaultProxyAddress);
  console.log("│  ├─ CollateralManager:", collateralManagerAddress);
  console.log("│  ├─ PositionManager:", positionManagerAddress);
  console.log("│  ├─ DebtManager:", debtManagerAddress);
  console.log("│  ├─ InterestCalculator:", interestCalculatorAddress);
  console.log("│  └─ LiquidationEngine:", liquidationEngineAddress);
  console.log("└─ RWAModuleV3:", rwaProxyAddress);
  console.log("   ├─ OrderManager:", orderManagerAddress);
  console.log("   ├─ MatchingEngine:", matchingEngineAddress);
  console.log("   ├─ MarketDataProvider:", marketDataProviderAddress);
  console.log("   ├─ FeeCalculator:", feeCalculatorAddress);
  console.log("   └─ LiquidityAnalyzer:", liquidityAnalyzerAddress);

  console.log("\n⚠️  Next Steps:");
  console.log("1. Update L1StateRegistry with L2StateAggregator address:", stateAggregatorAddress);
  console.log("2. Test cross-chain communication");
  console.log("3. Update .env file with all deployed addresses");
  console.log("\n====================================\n");

  return deployedContracts;
}

// Execute deployment
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
