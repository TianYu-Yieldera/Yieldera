const hre = require("hardhat");

async function main() {
  const provider = new hre.ethers.JsonRpcProvider('https://sepolia-rollup.arbitrum.io/rpc');
  const address = '0x3C07226A3f1488320426eB5FE9976f72E5712346';
  
  const balance = await provider.getBalance(address);
  const balanceEth = hre.ethers.formatEther(balance);
  
  console.log('\n📊 Arbitrum Sepolia Balance Check\n');
  console.log('Address:', address);
  console.log('Balance:', balanceEth, 'ETH');
  
  if (parseFloat(balanceEth) >= 0.05) {
    console.log('\n✅ Sufficient balance for deployment!');
    console.log('Ready to deploy with:');
    console.log('npx hardhat run scripts/deploy-all-treasury.js --network arbitrumSepolia');
  } else if (parseFloat(balanceEth) > 0) {
    console.log('\n⏳ Bridge in progress...');
    console.log('Please wait a few more minutes.');
  } else {
    console.log('\n⏳ Waiting for bridge...');
    console.log('Expected arrival: ~10-15 minutes from bridge transaction');
  }
  
  console.log('\nView on Arbiscan: https://sepolia.arbiscan.io/address/' + address + '\n');
}

main().catch(console.error);
