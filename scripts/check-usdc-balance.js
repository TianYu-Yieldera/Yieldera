/**
 * Check USDC Balance Script
 * 检查钱包中的 USDC 余额
 */

const { ethers } = require('ethers');
require('dotenv').config();

// USDC 合约地址 (不同网络不同)
const USDC_ADDRESSES = {
  // Ethereum Mainnet
  mainnet: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',

  // Sepolia Testnet
  sepolia: '0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238',

  // Arbitrum One
  arbitrum: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831',

  // Arbitrum Sepolia
  arbitrumSepolia: '0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d',

  // Base
  base: '0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913',

  // Base Sepolia
  baseSepolia: '0x036CbD53842c5426634e7929541eC2318f3dCF7e'
};

// ERC20 ABI (只需要 balanceOf 和 decimals)
const ERC20_ABI = [
  'function balanceOf(address owner) view returns (uint256)',
  'function decimals() view returns (uint8)',
  'function symbol() view returns (string)',
  'function name() view returns (string)'
];

async function checkBalance(walletAddress, networkName, rpcUrl) {
  console.log(`\n🔍 检查 ${networkName}...`);
  console.log(`   RPC: ${rpcUrl}`);

  try {
    const provider = new ethers.JsonRpcProvider(rpcUrl);

    // 检查网络连接
    const network = await provider.getNetwork();
    console.log(`   ✅ 网络连接成功: Chain ID ${network.chainId}`);

    // 检查 ETH 余额
    const ethBalance = await provider.getBalance(walletAddress);
    console.log(`   💰 ETH 余额: ${ethers.formatEther(ethBalance)} ETH`);

    // 检查 USDC 余额
    const usdcAddress = USDC_ADDRESSES[networkName];
    if (!usdcAddress) {
      console.log(`   ⚠️  未配置 USDC 地址`);
      return;
    }

    console.log(`   📄 USDC 合约: ${usdcAddress}`);

    const usdcContract = new ethers.Contract(usdcAddress, ERC20_ABI, provider);

    try {
      const [balance, decimals, symbol, name] = await Promise.all([
        usdcContract.balanceOf(walletAddress),
        usdcContract.decimals(),
        usdcContract.symbol(),
        usdcContract.name()
      ]);

      const formattedBalance = ethers.formatUnits(balance, decimals);

      console.log(`   📌 代币信息: ${name} (${symbol})`);
      console.log(`   💵 USDC 余额: ${formattedBalance} ${symbol}`);
      console.log(`   🔢 原始余额: ${balance.toString()}`);

      if (balance > 0n) {
        console.log(`   ✅ 找到 USDC!`);
      } else {
        console.log(`   ❌ USDC 余额为 0`);
      }

    } catch (contractError) {
      console.log(`   ❌ 无法读取 USDC 合约: ${contractError.message}`);
      console.log(`   💡 可能原因: USDC 合约地址错误，或此网络无 USDC`);
    }

  } catch (error) {
    console.log(`   ❌ 错误: ${error.message}`);
  }
}

async function main() {
  console.log('=' .repeat(70));
  console.log('💰 USDC 余额检查工具');
  console.log('=' .repeat(70));

  // 获取钱包地址
  let walletAddress = process.argv[2];

  if (!walletAddress) {
    // 尝试从私钥生成地址
    const privateKey = process.env.PRIVATE_KEY;
    if (privateKey) {
      const wallet = new ethers.Wallet(privateKey);
      walletAddress = wallet.address;
      console.log('\n📍 使用环境变量中的钱包地址');
    } else {
      console.error('\n❌ 错误: 请提供钱包地址作为参数');
      console.log('用法: node check-usdc-balance.js <钱包地址>');
      console.log('或者: 在 .env 文件中设置 PRIVATE_KEY');
      process.exit(1);
    }
  }

  console.log(`\n👛 钱包地址: ${walletAddress}`);
  console.log('=' .repeat(70));

  // 配置要检查的网络
  const networks = [
    {
      name: 'sepolia',
      rpc: process.env.SEPOLIA_RPC_URL || 'https://rpc.sepolia.org'
    },
    {
      name: 'arbitrumSepolia',
      rpc: process.env.ARBITRUM_SEPOLIA_RPC_URL || 'https://sepolia-rollup.arbitrum.io/rpc'
    },
    {
      name: 'baseSepolia',
      rpc: process.env.BASE_SEPOLIA_RPC_URL || 'https://sepolia.base.org'
    },
    {
      name: 'arbitrum',
      rpc: process.env.ARBITRUM_RPC_URL || 'https://arb1.arbitrum.io/rpc'
    },
    {
      name: 'base',
      rpc: process.env.BASE_RPC_URL || 'https://mainnet.base.org'
    }
  ];

  // 检查所有网络
  for (const network of networks) {
    await checkBalance(walletAddress, network.name, network.rpc);
  }

  console.log('\n' + '=' .repeat(70));
  console.log('✅ 检查完成!');
  console.log('=' .repeat(70));

  console.log('\n💡 常见问题排查:');
  console.log('   1. 检查你发送 USDC 的交易哈希在区块链浏览器上确认');
  console.log('   2. 确认你发送到了正确的网络 (Sepolia/Arbitrum/Base?)');
  console.log('   3. 确认你使用的是测试网 USDC 还是主网 USDC');
  console.log('   4. 交易可能还在 pending,等待几分钟后重试');
  console.log('\n📱 区块链浏览器:');
  console.log(`   Sepolia: https://sepolia.etherscan.io/address/${walletAddress}`);
  console.log(`   Arbitrum Sepolia: https://sepolia.arbiscan.io/address/${walletAddress}`);
  console.log(`   Base Sepolia: https://sepolia.basescan.org/address/${walletAddress}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error('\n❌ 发生错误:', error);
    process.exit(1);
  });
