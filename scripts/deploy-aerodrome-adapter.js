const hre = require("hardhat");

async function main() {
  console.log("\n🚀 Deploying Aerodrome Adapter to Base Sepolia\n");
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

  if (chainId !== "84532") {
    console.log("\n❌ Wrong network! This script is for Base Sepolia (84532)");
    console.log("Current network:", chainId);
    process.exit(1);
  }

  // Contract addresses on Base (Mainnet - need to find Sepolia addresses)
  // Note: Aerodrome might not be on Base Sepolia testnet
  // For production deployment on Base Mainnet, use these:
  const AERODROME_ROUTER_MAINNET = "0xcF77a3Ba9A5CA399B7c97c74d54e5b1Beb874E43";
  const AERO_TOKEN_MAINNET = "0x940181a94A35A4569E4529A3CDfB74e38FD98631";

  // For Sepolia, we'll need to deploy mocks or use testnet addresses
  console.log("\n⚠️  Note: Aerodrome is primarily on Base Mainnet");
  console.log("For testnet, we need mock contracts or testnet deployment addresses");
  console.log("Proceeding with placeholder addresses for compilation testing...\n");

  // Placeholder for testing (replace with actual testnet addresses when available)
  const ROUTER_ADDRESS = process.env.AERODROME_ROUTER_SEPOLIA || AERODROME_ROUTER_MAINNET;
  const AERO_TOKEN = process.env.AERO_TOKEN_SEPOLIA || AERO_TOKEN_MAINNET;

  // Use existing L2StateAggregator if deployed, or deploy a mock
  const STATE_AGGREGATOR = process.env.L2_STATE_AGGREGATOR || "0x0000000000000000000000000000000000000001";

  console.log("Configuration:");
  console.log("Router:", ROUTER_ADDRESS);
  console.log("AERO Token:", AERO_TOKEN);
  console.log("State Aggregator:", STATE_AGGREGATOR);

  console.log("\n📝 Deploying AerodromeAdapter...");

  const AerodromeAdapter = await hre.ethers.getContractFactory("AerodromeAdapter");

  // Estimate gas
  const deployTx = await AerodromeAdapter.getDeployTransaction(
    ROUTER_ADDRESS,
    AERO_TOKEN,
    STATE_AGGREGATOR
  );
  const gasEstimate = await deployer.provider.estimateGas(deployTx);
  console.log("Estimated gas:", gasEstimate.toString());

  const adapter = await AerodromeAdapter.deploy(
    ROUTER_ADDRESS,
    AERO_TOKEN,
    STATE_AGGREGATOR,
    { gasLimit: gasEstimate * 120n / 100n }
  );

  await adapter.waitForDeployment();
  const adapterAddress = await adapter.getAddress();

  console.log("\n" + "=".repeat(70));
  console.log("✅ Aerodrome Adapter Deployed Successfully!");
  console.log("=".repeat(70));
  console.log("\nDeployed Contract:");
  console.log("  AerodromeAdapter:", adapterAddress);

  console.log("\nConfiguration:");
  console.log("  Router:", ROUTER_ADDRESS);
  console.log("  AERO Token:", AERO_TOKEN);
  console.log("  State Aggregator:", STATE_AGGREGATOR);
  console.log("  Max Slippage:", "1% (100 bps)");

  console.log("\n🔗 View on BaseScan:");
  console.log(`  https://sepolia.basescan.org/address/${adapterAddress}`);

  console.log("\n💡 Key Features:");
  console.log("  ✅ Add/remove liquidity to Aerodrome pools");
  console.log("  ✅ Swap tokens through optimized routes");
  console.log("  ✅ Stake LP tokens in gauges");
  console.log("  ✅ Earn AERO rewards");
  console.log("  ✅ Support for stable and volatile pools");

  console.log("\n📊 Next Steps:");
  console.log("  1. Update .env with deployed address:");
  console.log(`     AERODROME_ADAPTER=${adapterAddress}`);
  console.log("  2. Verify contract on BaseScan (if on mainnet)");
  console.log("  3. Configure state aggregator integration");
  console.log("  4. Test liquidity operations");
  console.log("  5. Set up gauge staking for AERO rewards");

  console.log("\n⚠️  Important Notes:");
  console.log("  - Aerodrome primarily operates on Base Mainnet");
  console.log("  - Testnet functionality requires testnet deployment or mocks");
  console.log("  - Always test with small amounts first");
  console.log("  - Monitor slippage on volatile pairs");

  // Save deployment info
  const deployment = {
    network: "baseSepolia",
    chainId: chainId,
    contracts: {
      AerodromeAdapter: adapterAddress,
      Router: ROUTER_ADDRESS,
      AeroToken: AERO_TOKEN,
      StateAggregator: STATE_AGGREGATOR
    },
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    gasUsed: gasEstimate.toString()
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
