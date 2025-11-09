# GMX V2 Adapter - 使用文档

## 📖 概述

GMXV2Adapter 是用于风险对冲和衍生品交易的智能合约适配器，提供：
- ✅ 开仓/平仓永续合约
- ✅ 紧急对冲功能
- ✅ 多市场支持（ETH/BTC/等）
- ✅ 杠杆限制和滑点保护
- ✅ 完整的权限管理

## 🏗️ 架构设计

```
User/RiskSystem
       ↓
   GMXV2Adapter (本合约)
       ↓
   GMX V2 ExchangeRouter
       ↓
   GMX V2 Markets
```

## 🚀 快速开始

### 1. 部署合约

```bash
# 部署到 Arbitrum One (主网)
npx hardhat run scripts/deploy-gmx-adapter.js --network arbitrumOne

# 部署到 Arbitrum Sepolia (测试网)
npx hardhat run scripts/deploy-gmx-adapter.js --network arbitrumSepolia

# 本地测试 (使用 Mock 合约)
npx hardhat run scripts/deploy-gmx-adapter.js --network localhost
```

### 2. 配置环境变量

```bash
# .env
GMXV2_ADAPTER_ADDRESS=0x...  # 部署后的合约地址
```

### 3. 授予角色

```javascript
const adapter = await ethers.getContractAt("GMXV2Adapter", ADAPTER_ADDRESS);

// 授予风控系统 RISK_MANAGER_ROLE
const RISK_MANAGER_ROLE = await adapter.RISK_MANAGER_ROLE();
await adapter.grantRole(RISK_MANAGER_ROLE, RISK_SYSTEM_ADDRESS);

// 授予操作员 OPERATOR_ROLE
const OPERATOR_ROLE = await adapter.OPERATOR_ROLE();
await adapter.grantRole(OPERATOR_ROLE, OPERATOR_ADDRESS);
```

## 📚 核心功能

### 1. 开仓 (Open Position)

**做多 (Long)**:
```javascript
const tx = await adapter.openPosition(
  marketAddress,        // 市场地址 (如 ETH/USD 市场)
  collateralToken,      // 抵押品代币 (USDC/USDT/WETH)
  collateralAmount,     // 抵押品数量
  sizeInUsd,            // 仓位大小 (USD, 18 decimals)
  true,                 // isLong = true (做多)
  acceptablePrice,      // 可接受价格
  executionFee,         // 执行费用 (ETH)
  { value: executionFee }
);
```

**做空 (Short)**:
```javascript
const tx = await adapter.openPosition(
  marketAddress,
  collateralToken,
  collateralAmount,
  sizeInUsd,
  false,                // isLong = false (做空)
  acceptablePrice,
  executionFee,
  { value: executionFee }
);
```

**参数说明**:
- `marketAddress`: GMX 市场地址
  - ETH/USD: `0x70d95587d40A2caf56bd97485aB3Eec10Bee6336`
  - BTC/USD: `0x47c031236e19d024b42f8AE6780E44A573170703`
- `collateralToken`: 抵押品代币地址
  - USDC: `0xaf88d065e77c8cC2239327C5EDb3A432268e5831`
  - USDT: `0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9`
- `collateralAmount`: 抵押品数量（注意精度）
  - USDC: 6 decimals → 1000 USDC = `1000_000000`
  - WETH: 18 decimals → 1 ETH = `1000000000000000000`
- `sizeInUsd`: 仓位大小 (18 decimals)
  - 10,000 USD = `ethers.parseEther("10000")`
- `acceptablePrice`: 可接受的执行价格 (18 decimals)
  - 防止滑点过大
  - 2000 USD = `ethers.parseEther("2000")`
- `executionFee`: 执行费用 (ETH, 最低 0.0001)
  - 建议: `ethers.parseEther("0.001")` (0.001 ETH)

**示例 - 开 10x 杠杆多单**:
```javascript
const USDC = "0xaf88d065e77c8cC2239327C5EDb3A432268e5831";
const ETH_USD_MARKET = "0x70d95587d40A2caf56bd97485aB3Eec10Bee6336";

// 用 1000 USDC 开 10,000 USD 的多单 (10x 杠杆)
const collateralAmount = ethers.parseUnits("1000", 6); // 1000 USDC
const sizeInUsd = ethers.parseEther("10000");          // 10k USD
const executionFee = ethers.parseEther("0.001");

// 批准 USDC
const usdcContract = await ethers.getContractAt("IERC20", USDC);
await usdcContract.approve(adapterAddress, collateralAmount);

// 开仓
const tx = await adapter.openPosition(
  ETH_USD_MARKET,
  USDC,
  collateralAmount,
  sizeInUsd,
  true,  // 做多
  ethers.parseEther("2000"), // 接受最高 2000 USD/ETH
  executionFee,
  { value: executionFee }
);

await tx.wait();
console.log("Position opened!");
```

### 2. 平仓 (Close Position)

```javascript
const tx = await adapter.closePosition(
  marketAddress,
  collateralToken,
  sizeInUsd,            // 平仓大小 (可以部分平仓)
  isLong,               // 与开仓时一致
  acceptablePrice,
  executionFee,
  { value: executionFee }
);
```

**示例 - 平掉全部多单**:
```javascript
const tx = await adapter.closePosition(
  ETH_USD_MARKET,
  USDC,
  ethers.parseEther("10000"), // 平掉全部 10k USD
  true,  // 多单
  ethers.parseEther("1800"), // 接受最低 1800 USD/ETH
  executionFee,
  { value: executionFee }
);
```

### 3. 紧急对冲 (Emergency Hedge)

**只有 RISK_MANAGER 可以调用**

```javascript
// 由风控系统调用
const tx = await adapter.emergencyHedge(
  userAddress,          // 用户地址
  marketAddress,        // 市场
  collateralToken,      // 抵押品
  hedgeSize,            // 对冲规模 (USD)
  "Liquidation protection", // 对冲原因
  { value: executionFee }
);
```

**示例 - 对冲用户仓位**:
```javascript
// 用户有 20k USD 的 ETH 多单，价格下跌，Health Factor 降低
// 风控系统自动对冲 10k USD (开空单)

const hedgeSize = ethers.parseEther("10000");
const executionFee = ethers.parseEther("0.001");

const tx = await adapter.connect(riskManager).emergencyHedge(
  userAddress,
  ETH_USD_MARKET,
  USDC,
  hedgeSize,
  "Health Factor < 1.3, auto-hedge triggered",
  { value: executionFee }
);

// 结果: 用户现在有 20k 多单 + 10k 空单 = 净敞口 10k
```

### 4. 查询仓位

```javascript
// 查询用户在 GMX 的实时仓位
const position = await adapter.getPosition(
  userAddress,
  marketAddress,
  collateralToken,
  isLong
);

console.log("Size (USD):", ethers.formatEther(position.sizeInUsd));
console.log("Collateral:", ethers.formatUnits(position.collateralAmount, 6));
console.log("Average Price:", ethers.formatEther(position.averagePrice));

// 查询用户所有仓位记录
const positions = await adapter.getUserPositions(userAddress);
for (const pos of positions) {
  console.log("Market:", pos.market);
  console.log("Is Long:", pos.isLong);
  console.log("Leverage:", pos.leverage);
  console.log("Is Hedge:", pos.isHedge);
}
```

### 5. 统计数据

```javascript
const stats = await adapter.getStatistics();

console.log("Total Orders:", stats.totalOrders);
console.log("Total Hedges:", stats.totalHedges);
console.log("Total Volume:", ethers.formatEther(stats.totalVolume), "USD");
console.log("Successful Orders:", stats.successfulOrders);
```

## 🔒 安全机制

### 1. 杠杆限制
```solidity
uint256 public constant MAX_LEVERAGE = 50; // 最大 50x 杠杆
```

如果开仓时杠杆超过 50x，交易会被拒绝。

### 2. 滑点保护
```solidity
uint256 public constant MAX_SLIPPAGE_BPS = 200; // 最大 2% 滑点
```

通过设置 `acceptablePrice` 参数控制滑点。

### 3. 执行费用
```solidity
uint256 public constant MIN_EXECUTION_FEE = 0.0001 ether;
```

所有订单必须提供至少 0.0001 ETH 的执行费用。

### 4. 重入保护
```solidity
nonReentrant modifier
```

所有状态变更函数都使用 ReentrancyGuard 防止重入攻击。

### 5. 暂停机制
```solidity
function pause() external onlyRole(DEFAULT_ADMIN_ROLE)
function unpause() external onlyRole(DEFAULT_ADMIN_ROLE)
```

紧急情况下可以暂停合约。

## 📊 事件日志

### PositionOpened
```solidity
event PositionOpened(
    address indexed user,
    bytes32 indexed orderKey,
    address market,
    address collateralToken,
    bool isLong,
    uint256 sizeInUsd,
    uint256 collateralAmount,
    uint256 leverage,
    bool isHedge
);
```

### PositionClosed
```solidity
event PositionClosed(
    address indexed user,
    bytes32 indexed orderKey,
    address market,
    uint256 sizeInUsd,
    int256 pnl
);
```

### EmergencyHedgeExecuted
```solidity
event EmergencyHedgeExecuted(
    address indexed user,
    address indexed market,
    uint256 hedgeSize,
    string reason,
    bytes32 orderKey
);
```

## 🎯 实际使用场景

### 场景 1: Delta 对冲

**问题**: 用户在 Aave 存入 10 ETH 作抵押，担心 ETH 价格下跌导致清算。

**解决方案**: 在 GMX 开相同规模的空单对冲。

```javascript
// 1. 计算对冲规模
const ethAmount = ethers.parseEther("10");  // 10 ETH
const ethPrice = ethers.parseEther("2000"); // 假设 ETH = 2000 USD
const hedgeSize = (ethAmount * ethPrice) / ethers.parseEther("1"); // 20,000 USD

// 2. 开空单对冲
const collateral = ethers.parseUnits("2000", 6); // 2000 USDC (10x 杠杆)

await adapter.openPosition(
  ETH_USD_MARKET,
  USDC,
  collateral,
  hedgeSize,
  false,  // 空单
  ethers.parseEther("2100"),
  executionFee,
  { value: executionFee }
);

// 结果:
// ETH 跌到 1800: Aave 抵押品贬值 -2000 USD, GMX 空单盈利 +2000 USD
// ETH 涨到 2200: Aave 抵押品升值 +2000 USD, GMX 空单亏损 -2000 USD
// 总风险敞口 = 0
```

### 场景 2: 清算保护

**问题**: Health Factor 降至 1.15，即将被清算。

**解决方案**: 风控系统自动对冲。

```javascript
// 风控系统监控到 HF < 1.3
const riskSystem = await ethers.getSigner(RISK_MANAGER_ADDRESS);

// 自动对冲 50% 仓位
const hedgeTx = await adapter.connect(riskSystem).emergencyHedge(
  userAddress,
  ETH_USD_MARKET,
  USDC,
  ethers.parseEther("10000"), // 对冲 10k USD
  "Health Factor dropped to 1.15",
  { value: executionFee }
);

// 结果: 风险敞口降低，HF 回升到安全区域
```

### 场景 3: 多空平衡

**问题**: 用户在多个协议有不同方向的仓位，想要平衡风险。

**解决方案**: 计算净敞口，通过 GMX 对冲。

```javascript
// 用户仓位分析
const positions = {
  aave: { long: 20000 },      // Aave 存入 ETH = 20k USD
  uniswap: { long: 10000 },   // Uniswap LP = 10k USD
  compound: { short: 5000 },  // Compound 借款 = 5k USD
};

// 净敞口 = 20k + 10k - 5k = 25k (做多)
const netExposure = 25000;

// 在 GMX 开空单平衡
await adapter.openPosition(
  ETH_USD_MARKET,
  USDC,
  ethers.parseUnits("2500", 6), // 2500 USDC
  ethers.parseEther("25000"),   // 25k USD 空单
  false,
  acceptablePrice,
  executionFee,
  { value: executionFee }
);

// 结果: 净敞口 = 0
```

## 🧪 测试

```bash
# 运行测试
npx hardhat test test/layer2/GMXV2Adapter.test.js

# 查看覆盖率
npx hardhat coverage --testfiles "test/layer2/GMXV2Adapter.test.js"

# Gas 报告
REPORT_GAS=true npx hardhat test test/layer2/GMXV2Adapter.test.js
```

## 🔗 相关资源

- [GMX V2 官方文档](https://docs.gmx.io/)
- [GMX V2 合约代码](https://github.com/gmx-io/gmx-synthetics)
- [Arbitrum 文档](https://docs.arbitrum.io/)

## ⚠️ 注意事项

1. **Gas 费用**: Arbitrum L2 的 Gas 费用很低（~$0.01-0.05），但仍需预留执行费用
2. **价格影响**: 大额订单可能产生滑点，建议分批执行
3. **杠杆风险**: 高杠杆交易风险极高，建议谨慎使用
4. **清算风险**: 永续合约可能被强制平仓，需要密切监控保证金率
5. **测试优先**: 生产环境使用前，务必在测试网充分测试

## 📞 技术支持

遇到问题？
1. 查看部署日志: `deployments/gmx-adapter-*.json`
2. 检查事件日志
3. 运行测试验证
4. 查看 GMX V2 官方文档

---

**Generated**: 2025-11-09
**Version**: 1.0.0
**Status**: ✅ Production Ready
