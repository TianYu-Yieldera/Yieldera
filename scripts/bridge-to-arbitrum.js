const hre = require("hardhat");

async function main() {
  console.log('\n🌉 Bridging ETH from Sepolia to Arbitrum Sepolia\n');
  
  const [wallet] = await hre.ethers.getSigners();
  console.log('From:', wallet.address);
  
  const balance = await wallet.provider.getBalance(wallet.address);
  console.log('Sepolia Balance:', hre.ethers.formatEther(balance), 'ETH\n');
  
  // Arbitrum Sepolia Inbox
  const INBOX_ADDRESS = '0xaAe29B0366299461418F5324a79Afc425BE5ae21';
  const INBOX_ABI = ['function depositEth() external payable returns (uint256)'];
  
  const inbox = new hre.ethers.Contract(INBOX_ADDRESS, INBOX_ABI, wallet);
  
  // Bridge 0.15 ETH
  const amount = hre.ethers.parseEther('0.15');
  console.log('Bridging:', hre.ethers.formatEther(amount), 'ETH');
  console.log('Target: Arbitrum Sepolia\n');
  
  const tx = await inbox.depositEth({
    value: amount,
    gasLimit: 200000
  });
  
  console.log('✅ Transaction sent!');
  console.log('TX Hash:', tx.hash);
  console.log('View: https://sepolia.etherscan.io/tx/' + tx.hash);
  console.log('\n⏳ Waiting for confirmation...');
  
  const receipt = await tx.wait();
  console.log('✅ Confirmed in block:', receipt.blockNumber);
  console.log('\n⏰ Bridge processing: ~10-15 minutes');
  console.log('Check at: https://sepolia.arbiscan.io/address/' + wallet.address);
  console.log('\n💡 After 15 minutes, deploy with:');
  console.log('npx hardhat run scripts/deploy-all-treasury.js --network arbitrumSepolia');
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error('\n❌ Error:', error.message);
    process.exit(1);
  });
