# GMX V2 监控系统集成计划

## 概述

GMX V2 是去中心化永续合约和现货交易平台，需要实时监控交易活动、仓位风险和流动性状态。

## GMX V2 架构理解

### 核心合约组件
1. **OrderBook** - 订单管理
2. **PositionManager** - 仓位管理
3. **Reader** - 数据查询
4. **Vault** - 资金池管理
5. **RewardRouter** - 奖励分配

### 关键事件
```solidity
// 订单事件
event CreateIncreaseOrder(address indexed account, uint256 orderIndex, ...)
event CreateDecreaseOrder(address indexed account, uint256 orderIndex, ...)
event ExecuteIncreaseOrder(address indexed account, uint256 orderIndex, ...)
event ExecuteDecreaseOrder(address indexed account, uint256 orderIndex, ...)
event CancelOrder(address indexed account, uint256 orderIndex, ...)

// 仓位事件
event IncreasePosition(bytes32 key, address account, ...)
event DecreasePosition(bytes32 key, address account, ...)
event LiquidatePosition(bytes32 key, address account, ...)
event UpdatePosition(bytes32 key, uint256 size, ...)

// 流动性事件
event BuyUSDG(address account, address token, uint256 amount, ...)
event SellUSDG(address account, address token, uint256 amount, ...)
```

## 监控指标设计

### 1. 订单监控
**指标：**
- 订单创建频率
- 订单执行率
- 订单取消率
- 大额订单追踪（>$10k）
- 订单类型分布（做多/做空）

**告警条件：**
- 订单执行失败率 > 20%
- 大额订单（>$100k）
- 订单堆积（待执行 > 100）

### 2. 仓位风险监控
**指标：**
- 总持仓量（Long/Short）
- 平均杠杆倍数
- 接近清算的仓位数量
- 清算事件频率
- 最大单一仓位规模

**告警条件：**
- 杠杆倍数 > 50x
- 保证金率 < 1.5%（接近清算）
- 批量清算（1分钟内 > 5个）
- 单一仓位 > 总TVL的10%

### 3. 流动性监控
**指标：**
- GLP 池总价值（TVL）
- 各资产利用率
- 池子深度变化
- 大额存取款
- APY 变化

**告警条件：**
- 单一资产利用率 > 90%
- TVL 单日变化 > 30%
- 大额提款（> $1M）
- 流动性枯竭风险

### 4. 价格和滑点监控
**指标：**
- 执行价格 vs 预期价格
- 价格影响（滑点）
- 资金费率
- 标记价格 vs 指数价格偏差

**告警条件：**
- 滑点 > 5%
- 价格偏差 > 3%
- 资金费率异常（> 0.1% per hour）

## 实现计划

### Phase 1: 核心监听器开发（2-3天）

#### 1.1 GMXOrderListener
```typescript
// listeners/gmx/GMXOrderListener.ts
export class GMXOrderListener extends BaseListener {
  private orderStats = {
    totalOrders: 0,
    executedOrders: 0,
    cancelledOrders: 0,
    largeOrders: 0,
    longOrders: 0,
    shortOrders: 0,
  };

  // 监听事件
  - CreateIncreaseOrder
  - CreateDecreaseOrder
  - ExecuteIncreaseOrder
  - ExecuteDecreaseOrder
  - CancelOrder
}
```

#### 1.2 GMXPositionListener
```typescript
// listeners/gmx/GMXPositionListener.ts
export class GMXPositionListener extends BaseListener {
  private positionStats = {
    totalLongSize: BigInt(0),
    totalShortSize: BigInt(0),
    totalPositions: 0,
    liquidationCount: 0,
    highLeveragePositions: 0,
  };

  // 监听事件
  - IncreasePosition
  - DecreasePosition
  - LiquidatePosition
  - UpdatePosition
}
```

#### 1.3 GMXVaultListener
```typescript
// listeners/gmx/GMXVaultListener.ts
export class GMXVaultListener extends BaseListener {
  private vaultStats = {
    totalTVL: BigInt(0),
    buyVolume: BigInt(0),
    sellVolume: BigInt(0),
    utilizationRate: 0,
  };

  // 监听事件
  - BuyUSDG
  - SellUSDG
  - Swap
}
```

### Phase 2: 风险计算引擎（1-2天）

```typescript
// services/gmx/GMXRiskCalculator.ts
export class GMXRiskCalculator {
  /**
   * 计算清算风险
   */
  calculateLiquidationRisk(position: Position): number {
    // (抵押品价值 - 仓位损失) / 仓位规模
    const marginRatio = ...;
    return marginRatio < 0.015 ? 'CRITICAL' :
           marginRatio < 0.03 ? 'WARNING' : 'SAFE';
  }

  /**
   * 计算流动性风险
   */
  calculateLiquidityRisk(asset: string): number {
    // 已用流动性 / 总流动性
    const utilization = ...;
    return utilization;
  }

  /**
   * 计算价格影响
   */
  calculatePriceImpact(size: bigint, liquidity: bigint): number {
    // size / liquidity
    return Number(size) / Number(liquidity);
  }
}
```

### Phase 3: 告警系统（1天）

```typescript
// services/alerts/GMXAlertService.ts
export class GMXAlertService {
  /**
   * 检查并发送告警
   */
  checkAndAlert(type: AlertType, data: any) {
    const alerts = [
      this.checkLiquidationRisk(data),
      this.checkLeverageRisk(data),
      this.checkLiquidityRisk(data),
      this.checkPriceImpact(data),
    ];

    alerts.filter(Boolean).forEach(alert => {
      this.sendAlert(alert);
    });
  }

  private sendAlert(alert: Alert) {
    // Slack notification
    // Email notification
    // Database log
  }
}
```

### Phase 4: 集成到主系统（0.5天）

```typescript
// index.ts
import { GMXOrderListener } from './listeners/gmx/GMXOrderListener';
import { GMXPositionListener } from './listeners/gmx/GMXPositionListener';
import { GMXVaultListener } from './listeners/gmx/GMXVaultListener';

// 启动 GMX 监听器
private async startGMXListeners() {
  const { blockchain, contracts } = MONITORING_CONFIG;

  // Order Book
  const orderListener = new GMXOrderListener(
    blockchain.arbitrumSepoliaWs,
    contracts.gmxOrderBook
  );
  await orderListener.start();

  // Position Manager
  const positionListener = new GMXPositionListener(
    blockchain.arbitrumSepoliaWs,
    contracts.gmxPositionManager
  );
  await positionListener.start();

  // Vault
  const vaultListener = new GMXVaultListener(
    blockchain.arbitrumSepoliaWs,
    contracts.gmxVault
  );
  await vaultListener.start();
}
```

## 环境配置

### 合约地址（Arbitrum One）

```env
# GMX V2 Contract Addresses (Mainnet)
GMX_ORDER_BOOK=0x09f77E8A13De2a0E6d26f17Ab5eF9e60dE0Fa4E3
GMX_POSITION_MANAGER=0x75E42e6f5b8FA8AefC7c1Ff8C9B1e9A4BDAe2b88
GMX_VAULT=0x489ee077994B6658eAfA855C308275EAd8097C4A
GMX_READER=0x22199a49A999c351eF7927602CFB187ec3cae489
GMX_REWARD_ROUTER=0xA906F338CB21815cBc4Bc87ace9e68c87eF8d8F1

# 测试网（如果有）
GMX_TESTNET_ORDER_BOOK=
GMX_TESTNET_POSITION_MANAGER=
GMX_TESTNET_VAULT=
```

## 测试策略

### 1. 单元测试
```typescript
// tests/listeners/GMXOrderListener.test.ts
describe('GMXOrderListener', () => {
  it('should track order creation', async () => {
    const listener = new GMXOrderListener(wsUrl, contractAddress);
    // Mock event
    // Verify stats update
  });

  it('should alert on large orders', async () => {
    // Test alert threshold
  });
});
```

### 2. 集成测试
```typescript
// tests/integration/gmx-monitoring.test.ts
describe('GMX Monitoring Integration', () => {
  it('should monitor complete trade lifecycle', async () => {
    // Create order -> Execute -> Update position
    // Verify all events captured
  });
});
```

### 3. 负载测试
- 模拟高频交易场景
- 测试批量清算处理
- 验证告警不重复发送

## 监控面板设计

### Dashboard Metrics

```
┌─────────────────────────────────────────────────────────┐
│ GMX V2 实时监控面板                                       │
├─────────────────────────────────────────────────────────┤
│ 订单统计                                                 │
│   总订单: 1,234  |  执行: 1,100  |  取消: 134           │
│   执行率: 89.1%  |  大额订单: 23                         │
├─────────────────────────────────────────────────────────┤
│ 仓位概览                                                 │
│   总持仓: $12.5M  |  Long: $7.2M  |  Short: $5.3M      │
│   平均杠杆: 15.2x |  高风险仓位: 5                       │
├─────────────────────────────────────────────────────────┤
│ 流动性状态                                               │
│   TVL: $45.2M    |  利用率: 67.3%                       │
│   24h 流入: $2.1M |  24h 流出: $1.8M                    │
├─────────────────────────────────────────────────────────┤
│ 风险告警                                                 │
│   🚨 高杠杆仓位: 3个 (>50x)                              │
│   ⚠️  接近清算: 2个 (保证金<2%)                          │
│   ℹ️  大额订单: $150k Long BTC                          │
└─────────────────────────────────────────────────────────┘
```

## 时间估算

| 阶段 | 任务 | 预计时间 |
|------|------|----------|
| Phase 1 | GMX 监听器开发 | 2-3天 |
| Phase 2 | 风险计算引擎 | 1-2天 |
| Phase 3 | 告警系统 | 1天 |
| Phase 4 | 系统集成 | 0.5天 |
| 测试 | 单元测试 + 集成测试 | 1天 |
| **总计** | | **5.5-7.5天** |

## 技术挑战

1. **高频事件处理** - GMX 交易频繁，需要优化事件处理性能
2. **复杂风险计算** - 需要准确计算清算价格、保证金率等
3. **多合约协调** - GMX V2 由多个合约组成，需要同步监听
4. **历史数据回溯** - 需要查询链上历史数据计算初始状态

## 优化建议

1. **事件批处理** - 批量处理事件减少计算开销
2. **缓存机制** - 缓存仓位数据减少链上查询
3. **异步告警** - 告警发送异步化避免阻塞
4. **数据库索引** - 优化历史数据查询性能

## 参考资源

- GMX V2 官方文档: https://docs.gmx.io/
- GMX V2 合约代码: https://github.com/gmx-io/gmx-contracts
- GMX V2 Subgraph: https://thegraph.com/hosted-service/subgraph/gmx-io/gmx-stats
- Arbitrum RPC: https://arbitrum.io/

## 下一步行动

1. ✅ 研究 GMX V2 合约结构和事件
2. ⬜ 创建 GMXOrderListener 监听器
3. ⬜ 创建 GMXPositionListener 监听器
4. ⬜ 创建 GMXVaultListener 监听器
5. ⬜ 实现风险计算引擎
6. ⬜ 集成到主监控系统
7. ⬜ 编写测试用例
8. ⬜ 部署到生产环境
