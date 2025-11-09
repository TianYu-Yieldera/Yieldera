/**
 * Treasury监控系统测试
 * 只监控Treasury合约（已部署）
 */

import dotenv from 'dotenv';
import { MarketplaceListener } from './listeners/treasury/MarketplaceListener';
import { MONITORING_CONFIG, validateConfig } from './config/monitoring';

// 加载环境变量
dotenv.config();

class TreasuryMonitorTest {
  private marketplaceListener?: MarketplaceListener;

  async start() {
    console.log('🎯 Testing Treasury Monitoring System...\n');

    try {
      const { blockchain, contracts } = MONITORING_CONFIG;

      if (!contracts.treasuryMarketplace) {
        console.error('❌ TREASURY_MARKETPLACE_ADDRESS not set in .env');
        console.log('\nPlease add to backend/.env:');
        console.log('TREASURY_MARKETPLACE_ADDRESS=0x90708d3663C3BE0DF3002dC293Bb06c45b67a334');
        process.exit(1);
      }

      // 启动Marketplace监听
      this.marketplaceListener = new MarketplaceListener(
        blockchain.arbitrumSepoliaWs,
        contracts.treasuryMarketplace
      );

      // 监听事件
      this.marketplaceListener.on('orderCreated', (data) => {
        console.log('📝 New Order:', data);
      });

      this.marketplaceListener.on('orderFilled', (data) => {
        console.log('✅ Order Filled:', data);
      });

      this.marketplaceListener.on('orderCancelled', (data) => {
        console.log('❌ Order Cancelled:', data);
      });

      this.marketplaceListener.on('alert', (alert) => {
        const emoji = alert.severity === 'critical' ? '🚨' :
                      alert.severity === 'warning' ? '⚠️' : 'ℹ️';
        console.log(`${emoji} ALERT: ${alert.message}`);
      });

      this.marketplaceListener.on('error', (error) => {
        console.error('❌ Error:', error.message);
      });

      await this.marketplaceListener.start();

      console.log('\n✅ Treasury Monitoring Started!');
      console.log('📊 Watching TreasuryMarketplace at:', contracts.treasuryMarketplace);
      console.log('\nWaiting for events...\n');

      // 定时输出统计
      setInterval(() => {
        this.printStats();
      }, 60000); // 1分钟

      // 获取历史事件
      await this.getHistoricalEvents();

    } catch (error) {
      console.error('❌ Failed to start:', error);
      process.exit(1);
    }
  }

  /**
   * 获取历史事件
   */
  private async getHistoricalEvents() {
    if (!this.marketplaceListener) return;

    try {
      console.log('🔍 Fetching historical events...\n');

      const currentBlock = await this.marketplaceListener['provider'].getBlockNumber();
      const fromBlock = Math.max(0, currentBlock - 10000); // 最近10000个区块

      const orderCreatedEvents = await this.marketplaceListener.getHistoricalEvents(
        'OrderCreated',
        fromBlock,
        'latest'
      );

      const orderFilledEvents = await this.marketplaceListener.getHistoricalEvents(
        'OrderFilled',
        fromBlock,
        'latest'
      );

      console.log(`📜 Historical Events (last 10000 blocks):`);
      console.log(`  - Orders Created: ${orderCreatedEvents.length}`);
      console.log(`  - Orders Filled: ${orderFilledEvents.length}`);
      console.log('');

    } catch (error) {
      console.error('Error fetching historical events:', error);
    }
  }

  /**
   * 打印统计数据
   */
  private printStats() {
    console.log('\n📊 === Treasury Marketplace Stats ===');

    if (this.marketplaceListener) {
      const stats = this.marketplaceListener.getStats();
      console.log('Marketplace:', stats);
    }

    console.log('====================================\n');
  }

  /**
   * 优雅关闭
   */
  async shutdown() {
    console.log('\n🛑 Shutting down...');

    if (this.marketplaceListener) {
      await this.marketplaceListener.stop();
    }

    console.log('✅ Shutdown complete');
    process.exit(0);
  }
}

// 启动测试
const monitor = new TreasuryMonitorTest();
monitor.start();

// 处理退出信号
process.on('SIGINT', () => monitor.shutdown());
process.on('SIGTERM', () => monitor.shutdown());
