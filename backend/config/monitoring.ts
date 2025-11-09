/**
 * 监控系统配置
 */

export const MONITORING_CONFIG = {
  // 区块链连接
  blockchain: {
    arbitrumSepoliaWs: process.env.ARBITRUM_SEPOLIA_WS || 'wss://sepolia-rollup.arbitrum.io/rpc',
    arbitrumSepoliaHttp: process.env.ARBITRUM_SEPOLIA_RPC || 'https://sepolia-rollup.arbitrum.io/rpc',
  },

  // 合约地址
  contracts: {
    // DeFi Adapters (从部署记录获取)
    uniswapAdapter: process.env.UNISWAP_ADAPTER_ADDRESS || '',
    aaveAdapter: process.env.AAVE_ADAPTER_ADDRESS || '',
    compoundAdapter: process.env.COMPOUND_ADAPTER_ADDRESS || '',

    // GMX V2 Adapter (新增)
    gmxv2Adapter: process.env.GMXV2_ADAPTER_ADDRESS || '',

    // Treasury Contracts
    treasuryMarketplace: process.env.TREASURY_MARKETPLACE_ADDRESS || '0x90708d3663C3BE0DF3002dC293Bb06c45b67a334',
    treasuryAssetFactory: process.env.TREASURY_ASSET_FACTORY_ADDRESS || '0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B',
    treasuryYieldDistributor: process.env.TREASURY_YIELD_DISTRIBUTOR_ADDRESS || '0x0BE14D40188FCB5924c36af46630faBD76698A80',

    // RWA Contracts
    rwaYieldDistributor: process.env.RWA_YIELD_DISTRIBUTOR_ADDRESS || '',
  },

  // 风险阈值
  thresholds: {
    // Uniswap
    highSlippage: 0.02, // 2%
    largeSwapAmount: 100, // ETH

    // Aave
    criticalUtilization: 0.9, // 90%
    warningUtilization: 0.8, // 80%
    largeWithdrawal: 50, // ETH
    largeFlashLoan: 100, // ETH

    // GMX V2 (新增)
    leverageWarning: 30, // 杠杆警告线
    leverageCritical: 40, // 杠杆危险线
    largePositionUsd: 50000, // 大额仓位阈值 ($50k)

    // Treasury
    priceDeviationPercent: 0.05, // 5%
    volumeChangePercent: 0.3, // 30%
    lowLiquidity: 10000, // USD

    // YieldDistributor
    largeYieldDeposit: 100000, // $100,000
    shortClaimPeriodDays: 7, // 7 days
    lowClaimRate: 0.5, // 50%
  },

  // 告警配置
  alerts: {
    slack: {
      enabled: process.env.SLACK_ENABLED === 'true',
      webhookUrl: process.env.SLACK_WEBHOOK_URL || '',
      channelName: process.env.SLACK_CHANNEL || '#monitoring',
      botName: process.env.SLACK_BOT_NAME || 'Loyalty Points Monitor',
      minLevel: (process.env.SLACK_MIN_LEVEL || 'WARNING') as 'CRITICAL' | 'WARNING' | 'INFO',
    },
    email: {
      enabled: false,
      // 后续配置
    },
    // 告警冷却时间 (防止重复告警)
    cooldownMs: parseInt(process.env.ALERT_COOLDOWN_MS || '300000'), // 5分钟
  },

  // 数据库
  database: {
    host: process.env.DB_HOST || 'localhost',
    port: parseInt(process.env.DB_PORT || '5432'),
    database: process.env.DB_NAME || 'loyalty_monitoring',
    user: process.env.DB_USER || 'postgres',
    password: process.env.DB_PASSWORD || '',
  },

  // 监听器配置
  listeners: {
    reconnectAttempts: 10,
    reconnectDelayMs: 1000,
    maxReconnectDelayMs: 30000,
  },

  // 性能配置
  performance: {
    eventProcessingDelayMs: 5000,
    riskCalculationIntervalMs: 60000, // 1分钟
    statsReportIntervalMs: 300000, // 5分钟
  },

  // GMX 专用配置
  gmx: {
    monitoring: {
      leverageWarning: 30,
      leverageCritical: 40,
      largePositionUsd: '50000',
    },
    advisorMode: {
      enabled: true,
      autoHedge: false,
      requireUserConfirmation: true,
      trackRecommendations: true,
    },
  },
};

// 验证配置
export function validateConfig(): void {
  const required = [
    'ARBITRUM_SEPOLIA_WS',
    'DB_HOST',
    'DB_NAME',
    'DB_USER',
    'DB_PASSWORD',
  ];

  const missing = required.filter(key => !process.env[key]);

  if (missing.length > 0) {
    console.warn(`Warning: Missing environment variables: ${missing.join(', ')}`);
    console.warn('Some features may not work correctly.');
  }

  // 验证合约地址
  if (!MONITORING_CONFIG.contracts.uniswapAdapter) {
    console.warn('Warning: UNISWAP_ADAPTER_ADDRESS not set');
  }
  if (!MONITORING_CONFIG.contracts.aaveAdapter) {
    console.warn('Warning: AAVE_ADAPTER_ADDRESS not set');
  }
  if (!MONITORING_CONFIG.contracts.treasuryMarketplace) {
    console.warn('Warning: TREASURY_MARKETPLACE_ADDRESS not set');
  }
}
