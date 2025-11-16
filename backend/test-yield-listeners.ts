/**
 * YieldDistributor 监听器测试套件
 */

import dotenv from 'dotenv';
import { TreasuryYieldDistributorListener } from './listeners/treasury/TreasuryYieldDistributorListener';
import { RWAYieldDistributorListener } from './listeners/rwa/RWAYieldDistributorListener';

dotenv.config();

const WS_URL = process.env.ARBITRUM_SEPOLIA_WS || 'wss://arb-sepolia.g.alchemy.com/v2/YOUR_API_KEY';
const TREASURY_YIELD_ADDRESS = process.env.TREASURY_YIELD_DISTRIBUTOR_ADDRESS || '0x0BE14D40188FCB5924c36af46630faBD76698A80';
const RWA_YIELD_ADDRESS = process.env.RWA_YIELD_DISTRIBUTOR_ADDRESS || '';

console.log('🧪 Starting YieldDistributor Listener Tests\n');
console.log('============================================================\n');

// Test 1: 环境配置检查
async function testEnvironmentConfig() {
  console.log('📋 Test 1: Environment Configuration');
  console.log('------------------------------------------------------------');

  if (!WS_URL || WS_URL.includes('YOUR_API_KEY')) {
    console.log('  ❌ ARBITRUM_SEPOLIA_WS not configured');
    console.log('  ℹ️  Please set ARBITRUM_SEPOLIA_WS in .env file');
    return false;
  }

  console.log(`  ✅ ARBITRUM_SEPOLIA_WS: ${WS_URL.substring(0, 50)}...`);
  console.log(`  ✅ TREASURY_YIELD_DISTRIBUTOR_ADDRESS: ${TREASURY_YIELD_ADDRESS}`);

  if (RWA_YIELD_ADDRESS) {
    console.log(`  ✅ RWA_YIELD_DISTRIBUTOR_ADDRESS: ${RWA_YIELD_ADDRESS}`);
  } else {
    console.log('  ⚠️  RWA_YIELD_DISTRIBUTOR_ADDRESS not set (optional)');
  }

  console.log('');
  return true;
}

// Test 2: TreasuryYieldDistributor监听器连接测试
async function testTreasuryYieldListener() {
  console.log('🔌 Test 2: TreasuryYieldDistributor Listener Connection');
  console.log('------------------------------------------------------------');

  try {
    const listener = new TreasuryYieldDistributorListener(WS_URL, TREASURY_YIELD_ADDRESS);

    // 设置事件监听
    listener.on('yieldDeposited', (data) => {
      console.log('  📥 YieldDeposited event received:', {
        distributionId: data.distributionId,
        assetId: data.assetId,
        totalYield: data.totalYield,
        distributionType: data.distributionType,
      });
    });

    listener.on('yieldClaimed', (data) => {
      console.log('  💰 YieldClaimed event received:', {
        user: data.user.substring(0, 10) + '...',
        amount: data.amount,
        distributionId: data.distributionId,
      });
    });

    listener.on('batchDistributed', (data) => {
      console.log('  📦 BatchDistributed event received:', {
        distributionId: data.distributionId,
        recipientsCount: data.recipientsCount,
        totalAmount: data.totalAmount,
      });
    });

    listener.on('alert', (alert) => {
      console.log(`  ⚠️  Alert [${alert.level}]: ${alert.message}`);
    });

    listener.on('error', (error) => {
      console.error('  ❌ Listener error:', error.message);
    });

    await listener.start();
    console.log('  ✅ TreasuryYieldDistributor listener connected');
    console.log(`  📍 Monitoring contract: ${TREASURY_YIELD_ADDRESS}`);
    console.log('');

    return listener;
  } catch (error: any) {
    console.error('  ❌ Failed to connect:', error.message);
    console.log('');
    return null;
  }
}

// Test 3: RWAYieldDistributor监听器连接测试
async function testRWAYieldListener() {
  console.log('🔌 Test 3: RWAYieldDistributor Listener Connection');
  console.log('------------------------------------------------------------');

  if (!RWA_YIELD_ADDRESS) {
    console.log('  ⚠️  RWA_YIELD_DISTRIBUTOR_ADDRESS not configured, skipping');
    console.log('');
    return null;
  }

  try {
    const listener = new RWAYieldDistributorListener(WS_URL, RWA_YIELD_ADDRESS);

    // 设置事件监听
    listener.on('yieldDeposited', (data) => {
      console.log('  📥 YieldDeposited event received:', {
        distributionId: data.distributionId,
        assetId: data.assetId,
        amount: data.amount,
        paymentToken: data.paymentToken,
      });
    });

    listener.on('yieldClaimed', (data) => {
      console.log('  💰 YieldClaimed event received:', {
        distributionId: data.distributionId,
        user: data.user.substring(0, 10) + '...',
        amount: data.amount,
      });
    });

    listener.on('distributionFinalized', (data) => {
      console.log('  ✅ DistributionFinalized event received:', {
        distributionId: data.distributionId,
        totalClaimed: data.totalClaimed,
        unclaimed: data.unclaimed,
      });
    });

    listener.on('unclaimedYieldReclaimed', (data) => {
      console.log('  🔙 UnclaimedYieldReclaimed event received:', {
        distributionId: data.distributionId,
        amount: data.amount,
        recipient: data.recipient.substring(0, 10) + '...',
      });
    });

    listener.on('alert', (alert) => {
      console.log(`  ⚠️  Alert [${alert.level}]: ${alert.message}`);
    });

    listener.on('error', (error) => {
      console.error('  ❌ Listener error:', error.message);
    });

    await listener.start();
    console.log('  ✅ RWAYieldDistributor listener connected');
    console.log(`  📍 Monitoring contract: ${RWA_YIELD_ADDRESS}`);
    console.log('');

    return listener;
  } catch (error: any) {
    console.error('  ❌ Failed to connect:', error.message);
    console.log('');
    return null;
  }
}

// Test 4: 统计数据测试
async function testStatistics(treasuryListener: any, rwaListener: any) {
  console.log('📊 Test 4: Statistics Functionality');
  console.log('------------------------------------------------------------');

  if (treasuryListener) {
    console.log('  TreasuryYieldDistributor Stats:');
    const stats = treasuryListener.getStats();
    console.log('    - Total Distributions:', stats.totalDistributions);
    console.log('    - Total Yield Distributed:', stats.totalYieldDistributed);
    console.log('    - Total Claims:', stats.totalClaims);
    console.log('    - Coupon Payments:', stats.couponPayments);
    console.log('    - Maturity Payments:', stats.maturityPayments);
    console.log('    - Claim Rate:', stats.claimRate);
    console.log('');
  }

  if (rwaListener) {
    console.log('  RWAYieldDistributor Stats:');
    const stats = rwaListener.getStats();
    console.log('    - Total Distributions:', stats.totalDistributions);
    console.log('    - Total Yield Deposited:', stats.totalYieldDeposited);
    console.log('    - Total Claimed:', stats.totalClaimed);
    console.log('    - Total Finalized:', stats.totalFinalized);
    console.log('    - Overall Claim Rate:', stats.overallClaimRate);
    console.log('    - Active Distributions:', stats.activeDistributions);
    console.log('');
  }
}

// Test 5: 实时监控测试
async function testLiveMonitoring(treasuryListener: any, rwaListener: any) {
  console.log('🎧 Test 5: Live Event Monitoring');
  console.log('------------------------------------------------------------');
  console.log('  Listening for events... (Press Ctrl+C to stop)');
  console.log('  ⏰ Will report stats every 30 seconds');
  console.log('');

  // 每30秒输出一次统计
  const statsInterval = setInterval(() => {
    console.log('\n📊 === Current Stats ===');
    testStatistics(treasuryListener, rwaListener);
    console.log('========================\n');
  }, 30000);

  // 处理退出
  process.on('SIGINT', async () => {
    console.log('\n\n🛑 Stopping listeners...');
    clearInterval(statsInterval);

    if (treasuryListener) {
      await treasuryListener.stop();
      console.log('  ✅ TreasuryYieldDistributor listener stopped');
    }

    if (rwaListener) {
      await rwaListener.stop();
      console.log('  ✅ RWAYieldDistributor listener stopped');
    }

    console.log('\n✅ All tests completed\n');
    process.exit(0);
  });
}

// 运行所有测试
async function runAllTests() {
  try {
    // Test 1: 环境配置
    const configOk = await testEnvironmentConfig();
    if (!configOk) {
      console.log('❌ Configuration failed, exiting...\n');
      process.exit(1);
    }

    // Test 2 & 3: 监听器连接
    const treasuryListener = await testTreasuryYieldListener();
    const rwaListener = await testRWAYieldListener();

    if (!treasuryListener && !rwaListener) {
      console.log('❌ No listeners connected, exiting...\n');
      process.exit(1);
    }

    // Test 4: 统计数据
    await testStatistics(treasuryListener, rwaListener);

    // Test 5: 实时监控
    await testLiveMonitoring(treasuryListener, rwaListener);

  } catch (error: any) {
    console.error('❌ Test suite failed:', error.message);
    process.exit(1);
  }
}

// 启动测试
runAllTests();
