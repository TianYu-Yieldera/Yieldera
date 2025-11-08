/**
 * 监控系统启动入口
 */

import dotenv from 'dotenv';
import { UniswapListener } from './listeners/adapters/UniswapListener';
import { AaveListener } from './listeners/adapters/AaveListener';
import { MONITORING_CONFIG, validateConfig } from './config/monitoring';

// 加载环境变量
dotenv.config();

// 验证配置
validateConfig();

class MonitoringSystem {
  private listeners: Map<string, any> = new Map();

  async start() {
    console.log('🚀 Starting Loyalty Points Monitoring System...\n');

    try {
      // 启动DeFi适配器监听
      await this.startAdapterListeners();

      // 设置告警处理
      this.setupAlertHandlers();

      console.log('\n✅ Monitoring System Started Successfully!');
      console.log('📊 Monitoring the following contracts:');
      this.listeners.forEach((listener, name) => {
        console.log(`  - ${name}: ${listener.contractAddress}`);
      });

    } catch (error) {
      console.error('❌ Failed to start monitoring system:', error);
      process.exit(1);
    }
  }

  /**
   * 启动DeFi适配器监听器
   */
  private async startAdapterListeners() {
    const { blockchain, contracts } = MONITORING_CONFIG;

    // Uniswap监听
    if (contracts.uniswapAdapter) {
      const uniswapListener = new UniswapListener(
        blockchain.arbitrumSepoliaWs,
        contracts.uniswapAdapter
      );

      // 监听事件
      uniswapListener.on('swap', (data) => this.handleSwapEvent(data));
      uniswapListener.on('alert', (alert) => this.handleAlert(alert));
      uniswapListener.on('error', (error) => this.handleError('Uniswap', error));

      await uniswapListener.start();
      this.listeners.set('Uniswap', uniswapListener);
    }

    // Aave监听
    if (contracts.aaveAdapter) {
      const aaveListener = new AaveListener(
        blockchain.arbitrumSepoliaWs,
        contracts.aaveAdapter
      );

      // 监听事件
      aaveListener.on('supply', (data) => this.handleAaveEvent('supply', data));
      aaveListener.on('borrow', (data) => this.handleAaveEvent('borrow', data));
      aaveListener.on('flashLoan', (data) => this.handleFlashLoan(data));
      aaveListener.on('alert', (alert) => this.handleAlert(alert));
      aaveListener.on('error', (error) => this.handleError('Aave', error));

      await aaveListener.start();
      this.listeners.set('Aave', aaveListener);
    }

    console.log(`Started ${this.listeners.size} adapter listeners`);
  }

  /**
   * 设置告警处理
   */
  private setupAlertHandlers() {
    // 定时输出统计数据
    setInterval(() => {
      this.printStats();
    }, MONITORING_CONFIG.performance.statsReportIntervalMs);
  }

  /**
   * 处理Swap事件
   */
  private handleSwapEvent(data: any) {
    console.log('💱 Swap Event:', {
      user: data.user.substring(0, 10) + '...',
      amountIn: data.amountIn,
      slippage: data.slippage,
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理Aave事件
   */
  private handleAaveEvent(type: string, data: any) {
    console.log(`🏦 Aave ${type}:`, {
      user: data.user.substring(0, 10) + '...',
      amount: data.amount,
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理闪电贷
   */
  private handleFlashLoan(data: any) {
    console.log('⚡ Flash Loan Detected:', {
      initiator: data.initiator.substring(0, 10) + '...',
      amount: data.amount,
      premium: data.premium,
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理告警
   */
  private handleAlert(alert: any) {
    const emoji = alert.severity === 'critical' ? '🚨' :
                  alert.severity === 'warning' ? '⚠️' : 'ℹ️';

    console.log(`${emoji} ALERT [${alert.severity.toUpperCase()}]: ${alert.type}`);
    console.log(`   ${alert.message}`);

    // TODO: 发送到Slack
    if (MONITORING_CONFIG.alerts.slack.enabled) {
      this.sendSlackAlert(alert);
    }
  }

  /**
   * 处理错误
   */
  private handleError(source: string, error: any) {
    console.error(`❌ Error from ${source}:`, error.message);
  }

  /**
   * 发送Slack告警
   */
  private async sendSlackAlert(alert: any) {
    // TODO: 实现Slack webhook
    console.log('[Slack] Would send:', alert.message);
  }

  /**
   * 打印统计数据
   */
  private printStats() {
    console.log('\n📊 === Monitoring Stats ===');

    this.listeners.forEach((listener, name) => {
      if (typeof listener.getStats === 'function') {
        const stats = listener.getStats();
        console.log(`${name}:`, stats);
      }
    });

    console.log('========================\n');
  }

  /**
   * 优雅关闭
   */
  async shutdown() {
    console.log('\n🛑 Shutting down monitoring system...');

    for (const [name, listener] of this.listeners) {
      console.log(`Stopping ${name}...`);
      await listener.stop();
    }

    console.log('✅ Shutdown complete');
    process.exit(0);
  }
}

// 启动系统
const monitoringSystem = new MonitoringSystem();
monitoringSystem.start();

// 处理退出信号
process.on('SIGINT', () => monitoringSystem.shutdown());
process.on('SIGTERM', () => monitoringSystem.shutdown());
