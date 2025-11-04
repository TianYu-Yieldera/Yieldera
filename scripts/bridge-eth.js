const hre = require("hardhat");

async function main() {
  console.log('\n🌉 Bridging ETH from Sepolia to Arbitrum Sepolia...\n');
  
  const [wallet] = await hre.ethers.getSigners();
  
  // Arbitrum Sepolia Inbox address
  const INBOX_ADDRESS = '0xaAe29B0366299461418F5324a79Afc425BE5ae21';
  
  // Amount to bridge (0.2 ETH)
  const amount = hre.ethers.parseEther('0.2');
  
  console.log('From:', wallet.address);
  console.log('To (Inbox):', INBOX_ADDRESS);
  console.log('Amount:', hre.ethers.formatEther(amount), 'ETH\n');
  
  // Simple ETH transfer to Inbox (will be automatically bridged)
  const tx = await wallet.sendTransaction({
    to: INBOX_ADDRESS,
    value: amount,
    gasLimit: 100000
  });
  
  console.log('📤 Transaction sent:', tx.hash);
  console.log('🔗 View on Etherscan:', `https://sepolia.etherscan.io/tx/${tx.hash}\n`);
  console.log('⏳ Waiting for confirmation...');
  
  const receipt = await tx.wait();
  console.log('✅ Transaction confirmed in block:', receipt.blockNumber);
  console.log('\n⏳ Bridge processing...');
  console.log('   This takes 10-15 minutes. ETH will appear on Arbitrum Sepolia.');
  console.log('   Check status: https://sepolia.arbiscan.io/address/' + wallet.address);
  console.log('\n💡 We can start deploying to Sepolia (L1) now while waiting!\n');
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
