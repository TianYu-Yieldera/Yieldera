const hre = require("hardhat");

async function main() {
  console.log("\n🚀 Deploying GMX V2 Adapter to Base Sepolia\n");
  console.log("=".repeat(70));

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deployer address:", deployer.address);

  const balance = await deployer.provider.getBalance(deployer.address);
  console.log("Balance:", hre.ethers.formatEther(balance), "ETH");

  // Verify we're on Base Sepolia
  const network = await deployer.provider.getNetwork();
  const chainId = network.chainId.toString();
  console.log("Network:", network.name);
  console.log("Chain ID:", chainId);

  if (chainId !== "84532" && chainId !== "8453") {
    console.log("\n❌ Wrong network! This script is for Base networks");
    console.log("Current network:", chainId);
    process.exit(1);
  }

  console.log("\n⚠️  GMX V2 Deployment Status on Base:");
  console.log("GMX V2 is expanding to Base, but official contract addresses may vary.");
  console.log("This adapter is designed to be compatible once GMX deploys to Base.");
  console.log("\n");

  // GMX V2 Contract Addresses
  // NOTE: These are placeholder addresses for Arbitrum
  // Replace with actual Base addresses when GMX deploys to Base
  const GMX_EXCHANGE_ROUTER = process.env.GMX_EXCHANGE_ROUTER_BASE ||
    "0x7C68C7866A64FA2160F78EEaE12217FFbf871fa8"; // Arbitrum address as reference

  const GMX_READER = process.env.GMX_READER_BASE ||
    "0x60a0fF4cDaF0f6D496d35   e1b6e2e8b6a82bb"); // Arbitrum address as reference

  const GMX_DATASTORE = process.env.GMX_DATASTORE_BASE ||
    "0xFD70de6b91282D8017aA4E741e9Ae325CAb992d8"; // Arbitrum address as reference

  console.log("Configuration (using Arbitrum references for Base deployment):");
  console.log("ExchangeRouter:", GMX_EXCHANGE_ROUTER);
  console.log("Reader:", GMX_READER);
  console.log("DataStore:", GMX_DATASTORE);

  console.log("\n📝 Deploying GMXV2Adapter...");

  const GMXV2Adapter = await hre.ethers.getContractFactory("GMXV2Adapter");

  // Estimate gas
  const deployTx = await GMXV2Adapter.getDeployTransaction(
    GMX_EXCHANGE_ROUTER,
    GMX_READER,
    GMX_DATASTORE
  );
  const gasEstimate = await deployer.provider.estimateGas(deployTx);
  console.log("Estimated gas:", gasEstimate.toString());

  const adapter = await GMXV2Adapter.deploy(
    GMX_EXCHANGE_ROUTER,
    GMX_READER,
    GMX_DATASTORE,
    { gasLimit: gasEstimate * 120n / 100n }
  );

  await adapter.waitForDeployment();
  const adapterAddress = await adapter.getAddress();

  console.log("\n" + "=".repeat(70));
  console.log("✅ GMX V2 Adapter Deployed Successfully!");
  console.log("=".repeat(70));

  console.log("\nDeployed Contract:");
  console.log("  GMXV2Adapter:", adapterAddress);

  console.log("\nConfiguration:");
  console.log("  ExchangeRouter:", GMX_EXCHANGE_ROUTER);
  console.log("  Reader:", GMX_READER);
  console.log("  DataStore:", GMX_DATASTORE);
  console.log("  Max Leverage:", "50x");
  console.log("  Max Slippage:", "2% (200 bps)");

  console.log("\n🔗 View on BaseScan:");
  if (chainId === "84532") {
    console.log(`  https://sepolia.basescan.org/address/${adapterAddress}`);
  } else {
    console.log(`  https://basescan.org/address/${adapterAddress}`);
  }

  console.log("\n💡 Key Features:");
  console.log("  ✅ Open/close perpetual positions (long/short)");
  console.log("  ✅ Emergency hedge execution for risk management");
  console.log("  ✅ Position tracking and PnL monitoring");
  console.log("  ✅ Support for multiple collateral types");
  console.log("  ✅ Leverage up to 50x");
  console.log("  ✅ Integration with AI risk engine");

  console.log("\n📊 Use Cases:");
  console.log("  • Hedge DeFi portfolio risks");
  console.log("  • Arbitrage between Arbitrum and Base GMX markets");
  console.log("  • Performance comparison: GMX on Base vs Arbitrum");
  console.log("  • Automated risk management strategies");

  console.log("\n⚠️  Important Notes:");
  console.log("  - GMX V2 on Base is in early deployment phase");
  console.log("  - Update contract addresses when official Base deployment is live");
  console.log("  - Test with small positions first");
  console.log("  - Monitor liquidation risks with high leverage");
  console.log("  - Always use stop-loss orders");

  console.log("\n📈 Performance Comparison Setup:");
  console.log("  Arbitrum GMX (existing): High liquidity, established market");
  console.log("  Base GMX (new): Lower fees, faster finality");
  console.log("  → AI engine will track and compare:");
  console.log("     • Execution prices and slippage");
  console.log("     • Funding rates");
  console.log("     • Liquidity depth");
  console.log("     • Gas costs");

  console.log("\n📄 Next Steps:");
  console.log("  1. Update .env with deployed address:");
  console.log(`     GMX_ADAPTER_BASE=${adapterAddress}`);
  console.log("  2. Configure supported markets:");
  console.log("     • ETH-USD");
  console.log("     • BTC-USD");
  console.log("     • Add more as GMX expands on Base");
  console.log("  3. Set up supported collateral tokens:");
  console.log("     • USDC");
  console.log("     • WETH");
  console.log("  4. Integrate with Phase 4 risk engine");
  console.log("  5. Build performance comparison dashboard (Phase 7)");

  // Save deployment info
  const deployment = {
    network: chainId === "84532" ? "baseSepolia" : "base",
    chainId: chainId,
    contracts: {
      GMXV2Adapter: adapterAddress,
      ExchangeRouter: GMX_EXCHANGE_ROUTER,
      Reader: GMX_READER,
      DataStore: GMX_DATASTORE
    },
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    gasUsed: gasEstimate.toString(),
    note: "Using Arbitrum reference addresses - update when Base deployment is official"
  };

  console.log("\n📄 Deployment Summary:");
  console.log(JSON.stringify(deployment, null, 2));
  console.log();
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ Deployment Failed:");
    console.error(error);
    process.exit(1);
  });
