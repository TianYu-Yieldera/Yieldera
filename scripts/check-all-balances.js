const hre = require("hardhat");

async function checkNetwork(networkName, rpcUrl, chainId) {
  try {
    const provider = new hre.ethers.JsonRpcProvider(rpcUrl);
    const wallet = new hre.ethers.Wallet(process.env.PRIVATE_KEY, provider);
    const balance = await provider.getBalance(wallet.address);

    return {
      network: networkName,
      chainId: chainId,
      balance: hre.ethers.formatEther(balance),
      hasBalance: parseFloat(hre.ethers.formatEther(balance)) > 0
    };
  } catch (error) {
    return {
      network: networkName,
      chainId: chainId,
      balance: "Error",
      hasBalance: false,
      error: error.message
    };
  }
}

async function main() {
  console.log("\n💰 Checking All Network Balances\n");
  console.log("Wallet:", "0x3C07226A3f1488320426eB5FE9976f72E5712346");
  console.log("=".repeat(70));

  const networks = [
    { name: "Ethereum Sepolia", rpc: "https://eth-sepolia.g.alchemy.com/v2/demo", chainId: 11155111 },
    { name: "Arbitrum Sepolia", rpc: "https://sepolia-rollup.arbitrum.io/rpc", chainId: 421614 },
    { name: "Base Sepolia", rpc: "https://sepolia.base.org", chainId: 84532 },
  ];

  const results = [];

  for (const net of networks) {
    const result = await checkNetwork(net.name, net.rpc, net.chainId);
    results.push(result);
  }

  console.log("\n📊 Balance Summary:\n");
  results.forEach(r => {
    const icon = r.hasBalance ? "✅" : "❌";
    const balanceStr = r.error ? `Error: ${r.error.substring(0, 50)}` : `${r.balance} ETH`;
    console.log(`${icon} ${r.network.padEnd(25)} (${r.chainId}) - ${balanceStr}`);
  });

  const totalBalance = results
    .filter(r => !r.error)
    .reduce((sum, r) => sum + parseFloat(r.balance), 0);

  console.log("\n" + "=".repeat(70));
  console.log(`Total across all networks: ${totalBalance.toFixed(6)} ETH`);
  console.log("=".repeat(70));

  const hasEnough = results.filter(r => r.hasBalance && parseFloat(r.balance) > 0.01);
  if (hasEnough.length > 0) {
    console.log("\n💡 Networks with bridgeable balance:");
    hasEnough.forEach(r => {
      console.log(`   - ${r.network}: ${r.balance} ETH`);
    });
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
