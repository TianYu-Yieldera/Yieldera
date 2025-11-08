/**
 * 监控系统测试脚本
 * 用于验证所有监听器的功能
 */

import dotenv from 'dotenv';
import { ethers } from 'ethers';
import { MarketplaceListener } from './listeners/treasury/MarketplaceListener';
import { AssetFactoryListener } from './listeners/treasury/AssetFactoryListener';

dotenv.config();

interface TestResults {
  passed: number;
  failed: number;
  errors: string[];
}

class MonitoringSystemTest {
  private results: TestResults = {
    passed: 0,
    failed: 0,
    errors: [],
  };

  async runAllTests() {
    console.log('🧪 Starting Monitoring System Tests\n');
    console.log('=' .repeat(60));

    try {
      // Test 1: Environment Configuration
      await this.testEnvironmentConfig();

      // Test 2: WebSocket Connection
      await this.testWebSocketConnection();

      // Test 3: Treasury Listeners
      await this.testTreasuryListeners();

      // Test 4: Historical Events
      await this.testHistoricalEvents();

      // Test 5: Listener Stats
      await this.testListenerStats();

      // Print Summary
      this.printSummary();

    } catch (error: any) {
      console.error('\n❌ Test suite failed:', error.message);
      this.results.errors.push(`Test suite error: ${error.message}`);
    }
  }

  /**
   * Test 1: 验证环境配置
   */
  private async testEnvironmentConfig() {
    console.log('\n📋 Test 1: Environment Configuration');
    console.log('-'.repeat(60));

    const requiredVars = [
      'ARBITRUM_SEPOLIA_WS',
      'ARBITRUM_SEPOLIA_RPC',
      'TREASURY_MARKETPLACE_ADDRESS',
      'TREASURY_ASSET_FACTORY_ADDRESS',
      'USDC_ADDRESS',
    ];

    let allPresent = true;
    for (const varName of requiredVars) {
      const value = process.env[varName];
      if (!value) {
        console.log(`  ❌ Missing: ${varName}`);
        this.results.errors.push(`Missing environment variable: ${varName}`);
        allPresent = false;
      } else {
        console.log(`  ✅ ${varName}: ${value.substring(0, 50)}${value.length > 50 ? '...' : ''}`);
      }
    }

    if (allPresent) {
      console.log('\n  ✅ All required environment variables are set');
      this.results.passed++;
    } else {
      console.log('\n  ❌ Some environment variables are missing');
      this.results.failed++;
    }
  }

  /**
   * Test 2: 测试WebSocket连接
   */
  private async testWebSocketConnection() {
    console.log('\n🔌 Test 2: WebSocket Connection');
    console.log('-'.repeat(60));

    const wsUrl = process.env.ARBITRUM_SEPOLIA_WS!;

    try {
      const provider = new ethers.WebSocketProvider(wsUrl);

      // Test connection with timeout
      const timeoutPromise = new Promise((_, reject) =>
        setTimeout(() => reject(new Error('Connection timeout after 10s')), 10000)
      );

      const blockPromise = provider.getBlockNumber();

      const blockNumber = await Promise.race([blockPromise, timeoutPromise]) as number;

      console.log(`  ✅ Connected to Arbitrum Sepolia`);
      console.log(`  📦 Current block: ${blockNumber}`);

      await provider.destroy();
      this.results.passed++;

    } catch (error: any) {
      console.log(`  ❌ WebSocket connection failed: ${error.message}`);
      console.log(`  ℹ️  This may be due to:`);
      console.log(`     - Invalid or demo API key`);
      console.log(`     - Network issues`);
      console.log(`     - Rate limiting`);
      this.results.errors.push(`WebSocket connection: ${error.message}`);
      this.results.failed++;
    }
  }

  /**
   * Test 3: 测试Treasury监听器初始化
   */
  private async testTreasuryListeners() {
    console.log('\n🏛️ Test 3: Treasury Listeners Initialization');
    console.log('-'.repeat(60));

    const wsUrl = process.env.ARBITRUM_SEPOLIA_WS!;
    const marketplaceAddress = process.env.TREASURY_MARKETPLACE_ADDRESS!;
    const assetFactoryAddress = process.env.TREASURY_ASSET_FACTORY_ADDRESS!;

    let marketplaceListener: MarketplaceListener | null = null;
    let assetFactoryListener: AssetFactoryListener | null = null;

    try {
      // Test Marketplace Listener
      console.log('\n  Testing MarketplaceListener...');
      marketplaceListener = new MarketplaceListener(wsUrl, marketplaceAddress);

      const startPromise = marketplaceListener.start();
      const timeoutPromise = new Promise((_, reject) =>
        setTimeout(() => reject(new Error('Listener start timeout')), 15000)
      );

      await Promise.race([startPromise, timeoutPromise]);

      console.log(`    ✅ MarketplaceListener initialized`);
      console.log(`    📍 Contract: ${marketplaceAddress}`);

      const stats = marketplaceListener.getStats();
      console.log(`    📊 Initial stats:`, stats);

      await marketplaceListener.stop();

      // Wait a bit before next listener to avoid rate limit
      await new Promise(resolve => setTimeout(resolve, 2000));

      this.results.passed++;

    } catch (error: any) {
      console.log(`    ❌ MarketplaceListener failed: ${error.message}`);
      this.results.errors.push(`MarketplaceListener: ${error.message}`);
      this.results.failed++;
    }

    try {
      // Test AssetFactory Listener
      console.log('\n  Testing AssetFactoryListener...');
      assetFactoryListener = new AssetFactoryListener(wsUrl, assetFactoryAddress);

      const startPromise = assetFactoryListener.start();
      const timeoutPromise = new Promise((_, reject) =>
        setTimeout(() => reject(new Error('Listener start timeout')), 15000)
      );

      await Promise.race([startPromise, timeoutPromise]);

      console.log(`    ✅ AssetFactoryListener initialized`);
      console.log(`    📍 Contract: ${assetFactoryAddress}`);

      const stats = assetFactoryListener.getStats();
      console.log(`    📊 Initial stats:`, stats);

      await assetFactoryListener.stop();
      this.results.passed++;

    } catch (error: any) {
      console.log(`    ❌ AssetFactoryListener failed: ${error.message}`);

      // Check if it's a rate limit issue
      if (error.message.includes('provider destroyed') || error.message.includes('400')) {
        console.log(`    ℹ️  This may be due to using demo API key - please use your own API key`);
        this.results.errors.push(`AssetFactoryListener: Rate limited (demo API key)`);
      } else {
        this.results.errors.push(`AssetFactoryListener: ${error.message}`);
      }
      this.results.failed++;
    }
  }

  /**
   * Test 4: 测试历史事件查询
   */
  private async testHistoricalEvents() {
    console.log('\n📜 Test 4: Historical Events Query');
    console.log('-'.repeat(60));

    try {
      const httpUrl = process.env.ARBITRUM_SEPOLIA_RPC!;
      const marketplaceAddress = process.env.TREASURY_MARKETPLACE_ADDRESS!;

      const provider = new ethers.JsonRpcProvider(httpUrl);
      const currentBlock = await provider.getBlockNumber();

      console.log(`  Current block: ${currentBlock}`);

      // Check contract deployment
      const code = await provider.getCode(marketplaceAddress);
      if (code === '0x') {
        console.log(`  ⚠️  No contract found at Marketplace address`);
        console.log(`     This is expected if contracts haven't been deployed yet`);
      } else {
        console.log(`  ✅ Contract exists at Marketplace address`);
      }

      // Query recent events
      const abi = [
        'event OrderCreated(uint256 indexed orderId, address indexed seller, uint256 indexed assetId, uint256 amount, uint256 pricePerToken, uint8 orderType)',
        'event OrderFilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amount, uint256 totalPrice)',
      ];

      const marketplace = new ethers.Contract(marketplaceAddress, abi, provider);
      const fromBlock = Math.max(0, currentBlock - 10000);

      console.log(`  Searching events from block ${fromBlock} to ${currentBlock}...`);

      const [orderCreatedEvents, orderFilledEvents] = await Promise.all([
        marketplace.queryFilter(marketplace.filters.OrderCreated(), fromBlock, currentBlock),
        marketplace.queryFilter(marketplace.filters.OrderFilled(), fromBlock, currentBlock),
      ]);

      console.log(`  📊 Found ${orderCreatedEvents.length} OrderCreated events`);
      console.log(`  📊 Found ${orderFilledEvents.length} OrderFilled events`);

      if (orderCreatedEvents.length > 0) {
        console.log(`\n  Recent OrderCreated events:`);
        orderCreatedEvents.slice(-3).forEach((event: any) => {
          console.log(`    Block ${event.blockNumber}: Order #${event.args.orderId}`);
        });
      }

      this.results.passed++;

    } catch (error: any) {
      console.log(`  ❌ Historical events query failed: ${error.message}`);
      this.results.errors.push(`Historical events: ${error.message}`);
      this.results.failed++;
    }
  }

  /**
   * Test 5: 测试监听器统计功能
   */
  private async testListenerStats() {
    console.log('\n📊 Test 5: Listener Statistics');
    console.log('-'.repeat(60));

    try {
      const wsUrl = process.env.ARBITRUM_SEPOLIA_WS!;
      const marketplaceAddress = process.env.TREASURY_MARKETPLACE_ADDRESS!;

      const listener = new MarketplaceListener(wsUrl, marketplaceAddress);

      // Test stats before start
      const statsBeforeStart = listener.getStats();
      console.log(`  Stats before start:`, statsBeforeStart);

      // Verify initial state
      if (statsBeforeStart.totalOrders === 0 &&
          statsBeforeStart.totalTrades === 0 &&
          statsBeforeStart.totalVolume === '$0.0') {
        console.log(`  ✅ Initial stats are correctly zeroed`);
      } else {
        console.log(`  ⚠️  Unexpected initial stats`);
      }

      // Test resetStats
      listener.resetStats();
      const statsAfterReset = listener.getStats();
      console.log(`  Stats after reset:`, statsAfterReset);

      if (statsAfterReset.totalOrders === 0) {
        console.log(`  ✅ resetStats() works correctly`);
        this.results.passed++;
      } else {
        console.log(`  ❌ resetStats() did not reset stats`);
        this.results.failed++;
      }

    } catch (error: any) {
      console.log(`  ❌ Listener stats test failed: ${error.message}`);
      this.results.errors.push(`Listener stats: ${error.message}`);
      this.results.failed++;
    }
  }

  /**
   * 打印测试摘要
   */
  private printSummary() {
    console.log('\n' + '='.repeat(60));
    console.log('📊 Test Summary');
    console.log('='.repeat(60));

    const total = this.results.passed + this.results.failed;
    const passRate = total > 0 ? ((this.results.passed / total) * 100).toFixed(1) : '0.0';

    console.log(`\nTotal Tests: ${total}`);
    console.log(`✅ Passed: ${this.results.passed}`);
    console.log(`❌ Failed: ${this.results.failed}`);
    console.log(`📈 Pass Rate: ${passRate}%`);

    if (this.results.errors.length > 0) {
      console.log('\n❌ Errors encountered:');
      this.results.errors.forEach((error, index) => {
        console.log(`  ${index + 1}. ${error}`);
      });
    }

    console.log('\n' + '='.repeat(60));

    if (this.results.failed === 0) {
      console.log('🎉 All tests passed! Monitoring system is ready.\n');
    } else if (this.results.passed > 0) {
      console.log('⚠️  Some tests failed. Review errors above.\n');
    } else {
      console.log('❌ All tests failed. Check configuration.\n');
    }

    // Next steps
    console.log('📝 Next Steps:');
    if (this.results.errors.some(e => e.includes('WebSocket'))) {
      console.log('  1. Obtain API key from Alchemy/Infura/Quicknode');
      console.log('  2. Update ARBITRUM_SEPOLIA_WS in .env file');
      console.log('     Example: wss://arb-sepolia.g.alchemy.com/v2/YOUR_API_KEY');
    }
    if (this.results.errors.some(e => e.includes('environment variable'))) {
      console.log('  3. Ensure all required environment variables are set in .env');
    }
    if (this.results.passed > 0 && this.results.failed === 0) {
      console.log('  ✅ Run the main monitoring system: npm run start');
      console.log('  ✅ Deploy DeFi adapters and update their addresses in .env');
    }
    console.log('');
  }
}

// Run tests
const tester = new MonitoringSystemTest();
tester.runAllTests().catch(console.error);
