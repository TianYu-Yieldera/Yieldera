const hre = require("hardhat");

async function main() {
  console.log("\n📊 Checking Base Sepolia Balance\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Deployer address:", deployer.address);

  const balance = await deployer.provider.getBalance(deployer.address);
  console.log("Balance:", hre.ethers.formatEther(balance), "ETH");

  const network = await deployer.provider.getNetwork();
  console.log("Network:", network.name);
  console.log("Chain ID:", network.chainId.toString());

  const requiredBalance = "0.01";
  if (parseFloat(hre.ethers.formatEther(balance)) < parseFloat(requiredBalance)) {
    console.log("\n⚠️  Insufficient balance for deployment!");
    console.log(`Required: ${requiredBalance} ETH`);
    console.log("\n💡 Get Base Sepolia testnet ETH from:");
    console.log("   https://www.coinbase.com/faucets/base-ethereum-sepolia-faucet");
    console.log("   (Requires Coinbase account)");
  } else {
    console.log("\n✅ Sufficient balance for deployment");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
