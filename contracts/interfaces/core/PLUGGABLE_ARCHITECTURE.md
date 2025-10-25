# 可插拔架构核心接口文档 (Phase 1a)

## 📋 概述

本目录包含可插拔架构的所有核心接口定义。这些接口是整个系统的基础，定义了模块间通信和交互的标准。

## 🏗️ 接口架构

```
contracts/interfaces/core/
├── IModule.sol                  # 基础模块接口（所有模块必须实现）
├── IModuleRegistry.sol          # 模块注册中心接口
├── IAccessController.sol        # 统一访问控制接口
├── IEventHub.sol               # 事件总线接口
├── IVaultModule.sol            # 金库模块接口
├── IRWAModule.sol              # RWA交易模块接口
├── IPriceOracleModule.sol      # 价格预言机模块接口
└── IAuditModule.sol            # 审计日志模块接口
```

## 📚 接口详解

### 1️⃣ IModule.sol - 基础模块接口

**作用**: 所有可插拔模块的基础接口，定义模块生命周期和基本功能。

**核心功能**:
- 模块标识和元数据管理
- 生命周期管理（初始化、暂停、恢复）
- 依赖关系声明
- 健康检查机制

**模块状态**:
```solidity
enum ModuleState {
    UNINITIALIZED,  // 已部署但未初始化
    ACTIVE,         // 活跃运行中
    PAUSED,         // 暂停（可恢复）
    DEPRECATED,     // 已废弃（永久）
    UPGRADED        // 已升级到新版本
}
```

**关键方法**:
- `getModuleId()`: 获取模块唯一标识符
- `getDependencies()`: 声明模块依赖
- `initialize()`: 初始化模块
- `healthCheck()`: 健康状态检查

### 2️⃣ IModuleRegistry.sol - 模块注册中心

**作用**: 管理所有模块的注册、发现和生命周期。

**核心功能**:
- 模块注册/注销
- 依赖关系验证
- 模块发现和查询
- 升级管理
- 系统健康检查

**关键方法**:
- `registerModule()`: 注册新模块
- `enableModule()`: 启用模块
- `validateDependencies()`: 验证依赖关系
- `systemHealthCheck()`: 系统级健康检查

**使用示例**:
```solidity
// 注册模块
bytes32 moduleId = registry.registerModule(vaultModuleAddress);

// 验证依赖
bool satisfied = registry.validateDependencies(moduleId);

// 启用模块
registry.enableModule(moduleId);
```

### 3️⃣ IAccessController.sol - 统一访问控制

**作用**: 提供全局统一的权限管理系统。

**核心功能**:
- 基于角色的访问控制（RBAC）
- 细粒度权限管理
- 时间锁操作
- 模块级别权限隔离
- 紧急暂停机制

**权限层级**:
```
DEFAULT_ADMIN_ROLE
    ├── Role 1
    │   ├── Permission A
    │   └── Permission B
    └── Role 2
        └── Permission C
```

**关键方法**:
- `grantRole()`: 授予角色
- `hasPermission()`: 检查权限
- `scheduleOperation()`: 安排时间锁操作
- `activateEmergencyPause()`: 激活紧急暂停

**时间锁示例**:
```solidity
// 安排需要时间锁的操作
bytes32 opId = accessController.scheduleOperation(
    CRITICAL_PERMISSION,
    targetContract,
    callData
);

// 等待时间锁期满后执行
accessController.executeOperation(opId);
```

### 4️⃣ IEventHub.sol - 事件总线

**作用**: 提供模块间松耦合的事件通信机制。

**核心功能**:
- 事件发布/订阅
- 事件路由
- 事件历史查询
- 跨模块通信
- 回调机制

**事件分类**:
```solidity
enum EventCategory {
    SYSTEM,         // 系统事件
    MODULE,         // 模块事件
    TRANSACTION,    // 交易事件
    GOVERNANCE,     // 治理事件
    ORACLE,         // 预言机事件
    AUDIT,          // 审计事件
    CUSTOM          // 自定义事件
}
```

**使用示例**:
```solidity
// 模块A发布事件
bytes32 eventId = eventHub.publishEvent(
    EventCategory.TRANSACTION,
    EventSeverity.INFO,
    "DEPOSIT",
    abi.encode(user, amount)
);

// 模块B订阅事件
bytes32 subId = eventHub.subscribe(
    moduleA_ID,
    EventCategory.TRANSACTION,
    "DEPOSIT",
    callbackAddress,
    callbackSelector
);
```

### 5️⃣ IVaultModule.sol - 金库模块接口

**作用**: 定义抵押品管理和债务追踪的标准接口。

**核心功能**:
- 抵押品存取
- 债务管理（铸造/销毁稳定币）
- 抵押率计算
- 清算机制
- 利息累积

**关键数据结构**:
```solidity
struct Position {
    uint256 collateralAmount;   // 抵押品数量
    uint256 debtAmount;          // 债务数量
    uint256 lastInterestUpdate;  // 上次利息更新时间
    uint256 accruedInterest;     // 累积利息
    bool isActive;               // 仓位状态
}
```

**关键方法**:
- `depositCollateral()`: 存入抵押品
- `increaseDebt()`: 增加债务（铸造）
- `getCollateralRatio()`: 获取抵押率
- `liquidate()`: 清算

### 6️⃣ IRWAModule.sol - RWA交易模块接口

**作用**: 定义RWA资产交易的标准接口。

**核心功能**:
- 限价单管理
- 市价单执行
- 订单簿维护
- 交易撮合
- 市场统计

**订单类型**:
```solidity
enum OrderType { BUY, SELL }
enum OrderStatus {
    OPEN,               // 开放
    PARTIALLY_FILLED,   // 部分成交
    FILLED,             // 完全成交
    CANCELLED           // 已取消
}
```

**关键方法**:
- `placeOrder()`: 下限价单
- `placeMarketOrder()`: 下市价单
- `cancelOrder()`: 取消订单
- `matchOrders()`: 撮合订单
- `getOrderBookDepth()`: 获取订单簿深度

### 7️⃣ IPriceOracleModule.sol - 价格预言机接口

**作用**: 提供可靠的价格数据服务。

**核心功能**:
- 多源价格聚合
- 价格验证
- 熔断机制
- TWAP计算
- 历史价格查询

**价格源类型**:
```solidity
enum PriceSource {
    CHAINLINK,      // Chainlink 预言机
    UNISWAP_V3,    // Uniswap V3 TWAP
    CUSTOM,        // 自定义预言机
    MANUAL,        // 手动设置
    AGGREGATED     // 多源聚合
}
```

**关键方法**:
- `getLatestPrice()`: 获取最新价格
- `getPriceWithConfidence()`: 获取带置信度的价格
- `addPriceFeed()`: 添加价格源
- `getTWAP()`: 获取时间加权平均价格

### 8️⃣ IAuditModule.sol - 审计模块接口

**作用**: 提供全面的审计日志和合规报告功能。

**核心功能**:
- 事件记录
- 审计追踪
- 合规报告
- 日志查询
- 归档管理

**事件类型**:
```solidity
enum AuditEventType {
    DEPOSIT,
    WITHDRAWAL,
    MINT,
    BURN,
    LIQUIDATION,
    PRICE_UPDATE,
    CONFIG_CHANGE,
    EMERGENCY_ACTION,
    // ... 更多
}
```

**关键方法**:
- `logEvent()`: 记录事件
- `getAuditLogs()`: 查询日志
- `generateComplianceReport()`: 生成合规报告
- `exportAuditLogs()`: 导出审计日志

## 🔗 接口依赖关系

```
IModule (基础)
    ├── IVaultModule (继承)
    ├── IRWAModule (继承)
    ├── IPriceOracleModule (继承)
    └── IAuditModule (继承)

IModuleRegistry (管理所有实现 IModule 的合约)
    └── 依赖: IAccessController (权限检查)

IAccessController (独立)
    └── 可选: IEventHub (发布权限变更事件)

IEventHub (独立)
    └── 依赖: IModuleRegistry (验证模块身份)
```

## 📊 模块通信流程示例

### 场景: 用户存入抵押品

```
1. 用户 → VaultModule.depositCollateral()
   ↓
2. VaultModule → AccessController.hasPermission()
   ↓
3. VaultModule → EventHub.publishEvent(DEPOSIT)
   ↓
4. VaultModule → AuditModule.logDeposit()
   ↓
5. EventHub → 通知订阅模块
```

### 场景: 价格更新触发清算检查

```
1. PriceOracle → EventHub.publishEvent(PRICE_UPDATE)
   ↓
2. EventHub → VaultModule (订阅回调)
   ↓
3. VaultModule → 检查所有仓位
   ↓
4. VaultModule → liquidate() (如果需要)
   ↓
5. VaultModule → AuditModule.logLiquidation()
```

## 🎯 设计原则

1. **单一职责**: 每个接口只关注一个核心功能
2. **松耦合**: 模块间通过事件总线通信，避免直接依赖
3. **可扩展**: 接口设计支持未来功能扩展
4. **向后兼容**: 保持接口稳定，使用版本管理
5. **安全优先**: 所有关键操作都需要权限验证

## 📝 下一步 (Phase 1b)

完成接口定义后，下一步将实现这些接口：

1. ✅ **ModuleRegistry** - 模块注册中心实现
2. ✅ **AccessController** - 访问控制实现
3. ✅ **EventHub** - 事件总线实现
4. ✅ **PriceOracleModule** - 升级现有 PriceOracle
5. ✅ **AuditModule** - 升级现有 AuditLogger

## 🔍 使用说明

### 如何实现一个新模块

```solidity
// 1. 继承 IModule 接口
contract MyModule is IModule {
    bytes32 public constant MODULE_ID = keccak256("MY_MODULE");

    // 2. 实现必要的函数
    function getModuleId() external pure returns (bytes32) {
        return MODULE_ID;
    }

    function getDependencies() external pure returns (bytes32[] memory) {
        bytes32[] memory deps = new bytes32[](1);
        deps[0] = keccak256("PRICE_ORACLE_MODULE");
        return deps;
    }

    // 3. 实现业务逻辑
    // ...
}
```

### 如何注册模块

```solidity
// 1. 部署模块
MyModule module = new MyModule();

// 2. 在注册中心注册
bytes32 moduleId = moduleRegistry.registerModule(address(module));

// 3. 启用模块
moduleRegistry.enableModule(moduleId);
```

## 📖 参考资料

- [EIP-2535 Diamond Standard](https://eips.ethereum.org/EIPS/eip-2535)
- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Upgradeable Contracts](https://docs.openzeppelin.com/upgrades-plugins/)

## 🤝 贡献指南

如果需要扩展接口：

1. 确保向后兼容
2. 添加详细的 NatSpec 注释
3. 更新此文档
4. 创建对应的测试用例

---

**Phase 1a 完成时间**: 2025-10-25
**接口总数**: 8 个
**代码行数**: ~2000 行
**下一阶段**: Phase 1b - 核心合约实现
