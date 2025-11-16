const hre = require("hardhat");

async function main() {
  console.log("\n📊 Deploying GMX Performance Tracker\n");
  console.log("=".repeat(70));

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deployer:", deployer.address);

  const balance = await deployer.provider.getBalance(deployer.address);
  console.log("Balance:", hre.ethers.formatEther(balance), "ETH");

  const network = await deployer.provider.getNetwork();
  console.log("Network:", network.name);
  console.log("Chain ID:", network.chainId.toString());

  console.log("\n📝 Deploying GMXPerformanceTracker...");

  const GMXPerformanceTracker = await hre.ethers.getContractFactory("GMXPerformanceTracker");

  const deployTx = await GMXPerformanceTracker.getDeployTransaction();
  const gasEstimate = await deployer.provider.estimateGas(deployTx);
  console.log("Estimated gas:", gasEstimate.toString());

  const tracker = await GMXPerformanceTracker.deploy({
    gasLimit: gasEstimate * 120n / 100n
  });

  await tracker.waitForDeployment();
  const trackerAddress = await tracker.getAddress();

  console.log("\n" + "=".repeat(70));
  console.log("✅ GMX Performance Tracker Deployed!");
  console.log("=".repeat(70));

  console.log("\nContract Address:", trackerAddress);

  console.log("\n📊 Tracked Metrics:");
  console.log("  • Execution prices and slippage");
  console.log("  • Funding rates");
  console.log("  • Liquidity depth and utilization");
  console.log("  • Gas costs");
  console.log("  • Trade execution time");
  console.log("  • PnL performance");

  console.log("\n🔄 Comparison Capabilities:");
  console.log("  Arbitrum GMX vs Base GMX:");
  console.log("    - Average slippage");
  console.log("    - Average gas costs");
  console.log("    - Execution speed");
  console.log("    - Liquidity availability");
  console.log("    - Overall profitability");

  console.log("\n🔗 View on Explorer:");
  if (network.chainId.toString() === "84532") {
    console.log(`  https://sepolia.basescan.org/address/${trackerAddress}`);
  } else if (network.chainId.toString() === "8453") {
    console.log(`  https://basescan.org/address/${trackerAddress}`);
  } else if (network.chainId.toString() === "421614") {
    console.log(`  https://sepolia.arbiscan.io/address/${trackerAddress}`);
  } else {
    console.log(`  Chain ID ${network.chainId}: Check appropriate block explorer`);
  }

  console.log("\n⚙️  Configuration Needed:");
  console.log("  1. Authorize GMX adapters as reporters:");
  console.log("     tracker.setReporterAuthorization(arbitrumAdapter, true)");
  console.log("     tracker.setReporterAuthorization(baseAdapter, true)");
  console.log("  2. Integrate with backend monitoring service");
  console.log("  3. Set up automated data collection");

  console.log("\n📈 Integration with AI System:");
  console.log("  Phase 4: AI will use this data to:");
  console.log("    • Recommend optimal chain for new positions");
  console.log("    • Detect arbitrage opportunities");
  console.log("    • Optimize execution strategies");
  console.log("    • Predict market conditions");

  console.log("\n🎯 Use Cases:");
  console.log("  1. Real-time chain selection for trades");
  console.log("  2. Historical performance analysis");
  console.log("  3. Cost optimization (gas + slippage)");
  console.log("  4. Liquidity monitoring and alerts");
  console.log("  5. Performance dashboards (Phase 7)");

  // Save deployment info
  const deployment = {
    contract: "GMXPerformanceTracker",
    address: trackerAddress,
    network: network.name,
    chainId: network.chainId.toString(),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    gasUsed: gasEstimate.toString()
  };

  console.log("\n📄 Deployment Info:");
  console.log(JSON.stringify(deployment, null, 2));

  console.log("\n💡 Next Steps:");
  console.log("  1. Update .env:");
  console.log(`     GMX_PERFORMANCE_TRACKER=${trackerAddress}`);
  console.log("  2. Deploy/configure GMX adapters on both chains");
  console.log("  3. Start collecting performance data");
  console.log("  4. Build comparison dashboard");
  console.log();
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ Deployment failed:");
    console.error(error);
    process.exit(1);
  });
