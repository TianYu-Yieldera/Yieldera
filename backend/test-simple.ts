/**
 * 简单测试 - 使用HTTP provider查询历史事件
 */

import dotenv from 'dotenv';
import { ethers } from 'ethers';

dotenv.config();

async function test() {
  console.log('🔍 Testing Treasury Marketplace contract...\n');

  const provider = new ethers.JsonRpcProvider('https://sepolia-rollup.arbitrum.io/rpc');
  const marketplaceAddress = '0x90708d3663C3BE0DF3002dC293Bb06c45b67a334';

  // ABI
  const abi = [
    'event OrderCreated(uint256 indexed orderId, address indexed seller, uint256 indexed assetId, uint256 amount, uint256 pricePerToken, uint8 orderType)',
    'event OrderFilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amount, uint256 totalPrice)',
  ];

  const marketplace = new ethers.Contract(marketplaceAddress, abi, provider);

  // 检查当前区块
  const currentBlock = await provider.getBlockNumber();
  console.log(`✅ Connected to Arbitrum Sepolia`);
  console.log(`Current block: ${currentBlock}\n`);

  // 检查合约代码
  const code = await provider.getCode(marketplaceAddress);
  if (code === '0x') {
    console.log('❌ No contract found at this address');
    return;
  }
  console.log('✅ Contract exists\n');

  // 查询历史事件
  const fromBlock = Math.max(0, currentBlock - 10000);
  console.log(`🔍 Searching for events from block ${fromBlock} to ${currentBlock}...\n`);

  try {
    const orderCreatedFilter = marketplace.filters.OrderCreated();
    const orderFilledFilter = marketplace.filters.OrderFilled();

    const [createdEvents, filledEvents] = await Promise.all([
      marketplace.queryFilter(orderCreatedFilter, fromBlock, currentBlock),
      marketplace.queryFilter(orderFilledFilter, fromBlock, currentBlock),
    ]);

    console.log(`📊 Results:`);
    console.log(`  Orders Created: ${createdEvents.length}`);
    console.log(`  Orders Filled: ${filledEvents.length}\n`);

    if (createdEvents.length > 0) {
      console.log('📝 Recent OrderCreated events:');
      createdEvents.slice(-3).forEach((event: any) => {
        console.log(`  Block ${event.blockNumber}: Order #${event.args.orderId} by ${event.args.seller.substring(0, 10)}...`);
      });
      console.log('');
    }

    if (filledEvents.length > 0) {
      console.log('✅ Recent OrderFilled events:');
      filledEvents.slice(-3).forEach((event: any) => {
        console.log(`  Block ${event.blockNumber}: Order #${event.args.orderId} - ${ethers.formatUnits(event.args.totalPrice, 6)} USDC`);
      });
    }

    if (createdEvents.length === 0 && filledEvents.length === 0) {
      console.log('ℹ️  No events found in the last 10,000 blocks');
      console.log('   This is normal if the contract was recently deployed or has no activity yet');
    }

  } catch (error: any) {
    console.error('❌ Error querying events:', error.message);
  }
}

test().catch(console.error);
