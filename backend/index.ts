/**
 * 监控系统启动入口
 */

import dotenv from 'dotenv';
import { UniswapListener } from './listeners/adapters/UniswapListener';
import { AaveListener } from './listeners/adapters/AaveListener';
import { CompoundListener } from './listeners/adapters/CompoundListener';
import { GMXPositionListener } from './listeners/adapters/GMXPositionListener';
import { MarketplaceListener } from './listeners/treasury/MarketplaceListener';
import { AssetFactoryListener } from './listeners/treasury/AssetFactoryListener';
import { TreasuryYieldDistributorListener } from './listeners/treasury/TreasuryYieldDistributorListener';
import { RWAYieldDistributorListener } from './listeners/rwa/RWAYieldDistributorListener';
import { SlackAlertService } from './services/alerts/SlackAlertService';
import { MONITORING_CONFIG, validateConfig } from './config/monitoring';

// 加载环境变量
dotenv.config();

// 验证配置
validateConfig();

class MonitoringSystem {
  private listeners: Map<string, any> = new Map();
  private slackAlertService: SlackAlertService;

  constructor() {
    // 初始化 Slack 告警服务
    this.slackAlertService = new SlackAlertService({
      webhookUrl: MONITORING_CONFIG.alerts.slack.webhookUrl,
      enabled: MONITORING_CONFIG.alerts.slack.enabled,
      channelName: MONITORING_CONFIG.alerts.slack.channelName,
      botName: MONITORING_CONFIG.alerts.slack.botName,
      minLevel: MONITORING_CONFIG.alerts.slack.minLevel,
    });
  }

  async start() {
    console.log('🚀 Starting Loyalty Points Monitoring System...\n');

    try {
      // 启动DeFi适配器监听
      await this.startAdapterListeners();

      // 启动Treasury监听
      await this.startTreasuryListeners();

      // 启动GMX监听 (新增)
      await this.startGMXListeners();

      // 设置告警处理
      this.setupAlertHandlers();

      console.log('\n✅ Monitoring System Started Successfully!');
      console.log('📊 Monitoring the following contracts:');
      this.listeners.forEach((listener, name) => {
        console.log(`  - ${name}: ${listener.contractAddress}`);
      });

      // 发送启动通知到 Slack
      await this.slackAlertService.sendStartupNotification();

      console.log('\n📱 Slack Alerts:', this.slackAlertService.getStats());

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

    // Compound监听
    if (contracts.compoundAdapter) {
      const compoundListener = new CompoundListener(
        blockchain.arbitrumSepoliaWs,
        contracts.compoundAdapter
      );

      // 监听事件
      compoundListener.on('supply', (data) => this.handleCompoundEvent('supply', data));
      compoundListener.on('withdraw', (data) => this.handleCompoundEvent('withdraw', data));
      compoundListener.on('supplyRateUpdated', (data) => this.handleRateUpdate('supply', data));
      compoundListener.on('borrowRateUpdated', (data) => this.handleRateUpdate('borrow', data));
      compoundListener.on('alert', (alert) => this.handleAlert(alert));
      compoundListener.on('error', (error) => this.handleError('Compound', error));

      await compoundListener.start();
      this.listeners.set('Compound', compoundListener);
    }

    console.log(`Started ${this.listeners.size} adapter listeners`);
  }

  /**
   * 启动Treasury监听器
   */
  private async startTreasuryListeners() {
    const { blockchain, contracts } = MONITORING_CONFIG;

    // Marketplace监听
    if (contracts.treasuryMarketplace) {
      const marketplaceListener = new MarketplaceListener(
        blockchain.arbitrumSepoliaWs,
        contracts.treasuryMarketplace
      );

      // 监听事件
      marketplaceListener.on('orderCreated', (data) => this.handleMarketplaceEvent('orderCreated', data));
      marketplaceListener.on('orderFilled', (data) => this.handleMarketplaceEvent('orderFilled', data));
      marketplaceListener.on('orderCancelled', (data) => this.handleMarketplaceEvent('orderCancelled', data));
      marketplaceListener.on('alert', (alert) => this.handleAlert(alert));
      marketplaceListener.on('error', (error) => this.handleError('Marketplace', error));

      await marketplaceListener.start();
      this.listeners.set('TreasuryMarketplace', marketplaceListener);
    }

    // AssetFactory监听
    if (contracts.treasuryAssetFactory) {
      const assetFactoryListener = new AssetFactoryListener(
        blockchain.arbitrumSepoliaWs,
        contracts.treasuryAssetFactory
      );

      // 监听事件
      assetFactoryListener.on('assetCreated', (data) => this.handleAssetEvent('created', data));
      assetFactoryListener.on('assetVerified', (data) => this.handleAssetEvent('verified', data));
      assetFactoryListener.on('assetStatusUpdated', (data) => this.handleAssetEvent('statusUpdated', data));
      assetFactoryListener.on('assetMatured', (data) => this.handleAssetEvent('matured', data));
      assetFactoryListener.on('alert', (alert) => this.handleAlert(alert));
      assetFactoryListener.on('error', (error) => this.handleError('AssetFactory', error));

      await assetFactoryListener.start();
      this.listeners.set('TreasuryAssetFactory', assetFactoryListener);
    }

    // TreasuryYieldDistributor监听
    if (contracts.treasuryYieldDistributor) {
      const treasuryYieldListener = new TreasuryYieldDistributorListener(
        blockchain.arbitrumSepoliaWs,
        contracts.treasuryYieldDistributor
      );

      // 监听事件
      treasuryYieldListener.on('yieldDeposited', (data) => this.handleYieldEvent('deposited', data));
      treasuryYieldListener.on('yieldClaimed', (data) => this.handleYieldEvent('claimed', data));
      treasuryYieldListener.on('batchDistributed', (data) => this.handleYieldEvent('batchDistributed', data));
      treasuryYieldListener.on('alert', (alert) => this.handleAlert(alert));
      treasuryYieldListener.on('error', (error) => this.handleError('TreasuryYieldDistributor', error));

      await treasuryYieldListener.start();
      this.listeners.set('TreasuryYieldDistributor', treasuryYieldListener);
    }

    // RWAYieldDistributor监听
    if (contracts.rwaYieldDistributor) {
      const rwaYieldListener = new RWAYieldDistributorListener(
        blockchain.arbitrumSepoliaWs,
        contracts.rwaYieldDistributor
      );

      // 监听事件
      rwaYieldListener.on('yieldDeposited', (data) => this.handleRWAYieldEvent('deposited', data));
      rwaYieldListener.on('yieldClaimed', (data) => this.handleRWAYieldEvent('claimed', data));
      rwaYieldListener.on('distributionFinalized', (data) => this.handleRWAYieldEvent('finalized', data));
      rwaYieldListener.on('unclaimedYieldReclaimed', (data) => this.handleRWAYieldEvent('reclaimed', data));
      rwaYieldListener.on('alert', (alert) => this.handleAlert(alert));
      rwaYieldListener.on('error', (error) => this.handleError('RWAYieldDistributor', error));

      await rwaYieldListener.start();
      this.listeners.set('RWAYieldDistributor', rwaYieldListener);
    }

    const adapterCount = [contracts.uniswapAdapter, contracts.aaveAdapter, contracts.compoundAdapter].filter(Boolean).length;
    console.log(`Started ${this.listeners.size - adapterCount} Treasury listeners`);
  }

  /**
   * 启动GMX监听器
   */
  private async startGMXListeners() {
    const { blockchain, contracts } = MONITORING_CONFIG;

    if (!contracts.gmxv2Adapter) {
      console.log('⚠️  GMX V2 Adapter not configured, skipping...\n');
      return;
    }

    console.log('🎯 Starting GMX V2 Listeners...\n');

    // GMX Position Listener
    const gmxPositionListener = new GMXPositionListener(
      blockchain.arbitrumSepoliaWs,
      contracts.gmxv2Adapter
    );

    // 监听事件
    gmxPositionListener.on('positionOpened', (data) => {
      console.log('📈 GMX Position Opened:', {
        user: data.user.substring(0, 10) + '...',
        market: data.market.substring(0, 10) + '...',
        direction: data.isLong ? 'LONG' : 'SHORT',
        size: data.sizeUsd,
        leverage: data.leverage,
        isHedge: data.isHedge,
      });
    });

    gmxPositionListener.on('positionClosed', (data) => {
      console.log('📉 GMX Position Closed:', {
        user: data.user.substring(0, 10) + '...',
        market: data.market.substring(0, 10) + '...',
        pnl: data.pnl,
        profitable: data.profitable,
      });
    });

    gmxPositionListener.on('emergencyHedge', (data) => {
      console.log('🚨 GMX Emergency Hedge Executed:', {
        user: data.user.substring(0, 10) + '...',
        market: data.market.substring(0, 10) + '...',
        hedgeSize: data.hedgeSize,
        reason: data.reason,
      });
    });

    gmxPositionListener.on('alert', async (alert) => {
      await this.handleGMXAlert(alert);
    });

    gmxPositionListener.on('error', (error) => {
      this.handleError('GMXPosition', error);
    });

    await gmxPositionListener.start();
    this.listeners.set('GMX-Position', gmxPositionListener);

    console.log('✅ GMX V2 Listeners started\n');
  }

  /**
   * 处理GMX告警 (建议式)
   */
  private async handleGMXAlert(alert: any) {
    const { level, type, message, user, recommendation } = alert;

    console.log(`\n💡 GMX Risk Advisory - ${type}`);
    console.log(`  Level: ${level}`);
    console.log(`  User: ${user}`);
    console.log(`  Message: ${message}`);

    if (recommendation) {
      console.log(`  Action: ${recommendation.action}`);
      console.log(`  Priority: ${recommendation.priority}`);
      console.log(`  Reason: ${recommendation.reason}`);
      console.log(`  User Decision Required: ${recommendation.userDecision ? 'YES' : 'NO'}`);
    }

    // 发送到 Slack (建议式告警)
    if (recommendation && recommendation.userDecision) {
      await this.slackAlertService.send({
        title: `💡 GMX 风险建议 - ${type}`,
        level: level,
        message: message,
        fields: [
          { label: '用户', value: user, short: true },
          { label: '优先级', value: recommendation.priority, short: true },
          { label: '建议行动', value: recommendation.action },
          { label: '原因', value: recommendation.reason },
          { label: '预期效果', value: recommendation.expectedOutcome },
        ],
        footer: '⚠️ 需要用户自主决策 - Advisory Mode',
      });
    } else {
      // 普通告警 (如紧急对冲已执行)
      await this.slackAlertService.send({
        title: `🚨 GMX ${type}`,
        level: level,
        message: message,
        fields: [
          { label: '用户', value: user, short: true },
        ],
      });
    }
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
   * 处理Compound事件
   */
  private handleCompoundEvent(type: string, data: any) {
    console.log(`🏛️ Compound ${type}:`, {
      user: data.user.substring(0, 10) + '...',
      amount: data.amount,
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理利率更新
   */
  private handleRateUpdate(rateType: string, data: any) {
    console.log(`📈 ${rateType} Rate Updated:`, {
      newRate: data.newRate,
      timestamp: data.timestamp,
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理Marketplace事件
   */
  private handleMarketplaceEvent(type: string, data: any) {
    const eventEmojis: Record<string, string> = {
      orderCreated: '📝',
      orderFilled: '✅',
      orderCancelled: '❌',
    };

    console.log(`${eventEmojis[type] || '📊'} Marketplace ${type}:`, {
      orderId: data.orderId,
      seller: data.seller?.substring(0, 10) + '...' || 'N/A',
      buyer: data.buyer?.substring(0, 10) + '...' || 'N/A',
      amount: data.amount || 'N/A',
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理Asset事件
   */
  private handleAssetEvent(type: string, data: any) {
    const eventEmojis: Record<string, string> = {
      created: '🆕',
      verified: '✔️',
      statusUpdated: '🔄',
      matured: '💰',
    };

    console.log(`${eventEmojis[type] || '📄'} Asset ${type}:`, {
      assetId: data.assetId,
      symbol: data.symbol || 'N/A',
      value: data.totalValue || data.finalValue || 'N/A',
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理Treasury Yield事件
   */
  private handleYieldEvent(type: string, data: any) {
    const eventEmojis: Record<string, string> = {
      deposited: '💰',
      claimed: '🎁',
      batchDistributed: '📦',
    };

    console.log(`${eventEmojis[type] || '💵'} Treasury Yield ${type}:`, {
      distributionId: data.distributionId,
      assetId: data.assetId || 'N/A',
      amount: data.totalYield || data.amount || data.totalAmount || 'N/A',
      type: data.distributionType || 'N/A',
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理RWA Yield事件
   */
  private handleRWAYieldEvent(type: string, data: any) {
    const eventEmojis: Record<string, string> = {
      deposited: '💰',
      claimed: '🎁',
      finalized: '✅',
      reclaimed: '🔙',
    };

    console.log(`${eventEmojis[type] || '💵'} RWA Yield ${type}:`, {
      distributionId: data.distributionId,
      assetId: data.assetId || 'N/A',
      amount: data.amount || data.totalClaimed || 'N/A',
      unclaimed: data.unclaimed || 'N/A',
      txHash: data.transactionHash,
    });
  }

  /**
   * 处理告警
   */
  private async handleAlert(alert: any) {
    const emoji = alert.level === 'CRITICAL' ? '🚨' :
                  alert.level === 'WARNING' ? '⚠️' : 'ℹ️';

    console.log(`${emoji} ALERT [${alert.level}]: ${alert.type}`);
    console.log(`   ${alert.message}`);

    // 发送到Slack
    await this.slackAlertService.sendAlert(alert);
  }

  /**
   * 处理错误
   */
  private handleError(source: string, error: any) {
    console.error(`❌ Error from ${source}:`, error.message);

    // 发送错误告警到 Slack
    this.slackAlertService.sendAlert({
      level: 'WARNING',
      type: 'LISTENER_ERROR',
      message: `Error in ${source}: ${error.message}`,
      data: { source, error: error.message },
    });
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

    // 发送关闭通知
    await this.slackAlertService.sendShutdownNotification();

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
