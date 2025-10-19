# Yieldera - DeFi Yield Aggregation & RWA Platform

<div align="center">

![Yieldera Logo](frontend/public/pointfi-logo-mark.svg)

**Enter the Yieldera**

[![Version](https://img.shields.io/badge/version-2.0.0-purple)](https://github.com/yieldera/loyalty-points-system)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-compose-blue.svg)](docker-compose.yml)

*将 DeFi 收益优化与真实世界资产（RWA）投资完美结合的下一代金融平台*

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [架构说明](#-系统架构) • [演示指南](#-演示指南) • [开发文档](#-开发指南)

</div>

---

## 📖 项目概述

Yieldera 是一个完整的 Web3 金融平台，提供：

1. **DeFi 收益聚合** - 自动优化资金分配到 Aave、Compound、Uniswap V3、GMX 等协议
2. **RWA 资产购买** - 投资代币化的真实世界资产（股票、债券、黄金）
3. **积分奖励系统** - 基于链上活动自动累积积分
4. **空投分发** - 支持链下积分空投 + 链上代币空投（备用升级方案）

### 🎯 核心价值主张

- ✅ **收益最大化**：智能算法自动选择最高收益协议组合
- ✅ **真实资产**：通过区块链投资传统金融资产
- ✅ **双模式操作**：智能自动化 + 手动精细控制
- ✅ **透明可验证**：所有资金流动链上可查
- ✅ **商业化就绪**：链上空投系统符合 Web3 行业标准

---

## ✨ 功能特性

### 前端功能（React + Vite）

#### ✅ 已完成并可演示

| 功能模块 | 状态 | 说明 | 路由 |
|---------|------|------|------|
| **品牌首页** | ✅ 完成 | Yieldera 品牌展示，核心数据展示 | `/` |
| **用户概览** | ✅ 完成 | 钱包余额、积分、活动历史 | `/dashboard` |
| **理财金库** | ✅ 完成 | 智能模式 + 手动模式双模式收益优化 | `/vault` |
| **RWA 商城** | ✅ 完成 | 代币化股票/债券/黄金购买 | `/rwa-market` |
| **排行榜** | ✅ 完成 | 积分排名、社区竞争 | `/leaderboard` |
| **空投系统** | ✅ 完成 | 链下积分空投（演示用） | `/airdrop` |
| **空投管理** | ✅ 完成 | 管理员创建和管理空投活动 | `/admin/airdrop` |
| **系统状态** | ✅ 完成 | 微服务健康监控 | `/status` |
| **教程指南** | ✅ 完成 | 用户引导和帮助 | `/tutorial` |

#### 🔧 备用升级方案

| 功能模块 | 状态 | 说明 | 路由 |
|---------|------|------|------|
| **链上空投** | ✅ 已开发 | Merkle Tree + Subgraph 链上空投系统 | `/airdrop/onchain` |

---

### 后端微服务（Go + PostgreSQL + Kafka）

#### ✅ 核心服务（已完成）

| 服务 | 端口 | 状态 | 功能说明 |
|------|------|------|---------|
| **Listener** | 8090 | ✅ 运行中 | 监听链上事件（Transfer、Staking 等） |
| **Consumer** | - | ✅ 运行中 | 处理 Kafka 消息，更新用户余额 |
| **Scheduler** | - | ✅ 运行中 | 定时发放积分（每 60 秒） |
| **API** | 8080 | ✅ 运行中 | REST + GraphQL 接口 |
| **PostgreSQL** | 5432 | ✅ 运行中 | 主数据库 |
| **Kafka** | 9092 | ✅ 运行中 | 消息队列 |
| **Zookeeper** | 2181 | ✅ 运行中 | Kafka 协调器 |
| **Frontend** | 5173 | ✅ 运行中 | Nginx 静态文件服务 |

#### 📊 数据库表结构

**用户相关**
- `users` - 用户基础信息
- `balances` - 代币余额
- `balance_events` - 余额变动事件
- `points` - 积分余额
- `points_events` - 积分变动记录

**空投相关**
- `admin_whitelist` - 管理员白名单
- `airdrop_campaigns` - 空投活动
- `airdrop_allocations` - 分配白名单
- `airdrop_claims` - 领取记录

**其他**
- `badges` - 徽章数据（功能保留但不展示）

---

### 链上空投系统（备用方案）

完整的去中心化空投解决方案，位于 `onchain-airdrop/` 目录。

#### ✅ 已开发完成

| 组件 | 状态 | 说明 |
|------|------|------|
| **智能合约** | ✅ 完成 | YielderaAirdrop.sol - 支持 Merkle Proof 验证 |
| **测试代币** | ✅ 完成 | MockERC20.sol - 用于开发测试 |
| **部署脚本** | ✅ 完成 | Hardhat 自动化部署 |
| **Merkle 工具** | ✅ 完成 | CSV → Merkle Tree 生成器 |
| **Subgraph Schema** | ✅ 完成 | GraphQL 数据模型 |
| **Subgraph Mapping** | ✅ 完成 | 事件索引处理器 |
| **前端组件** | ✅ 完成 | OnchainAirdropView.jsx |
| **迁移文档** | ✅ 完成 | 详细的升级指南 |

**部署状态**: ⏳ 未部署到链上（作为备用升级方案）

**优势**:
- ✅ 完全透明可验证
- ✅ 去中心化执行
- ✅ Gas 优化（Merkle Tree）
- ✅ 支持批量领取
- ✅ Subgraph 实时索引

**文档**: 参见 `onchain-airdrop/README.md` 和 `onchain-airdrop/MIGRATION_GUIDE.md`

---

## 🏗️ 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
│  Landing │ Dashboard │ Vault │ RWA │ Leaderboard │ Airdrop      │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ↓
              ┌──────────────────────┐
              │   API Service (Go)    │
              │   Port: 8080          │
              │   REST + GraphQL      │
              └──────────┬───────────┘
                         │
         ┌───────────────┼───────────────┐
         ↓               ↓               ↓
┌────────────┐  ┌────────────┐  ┌────────────┐
│ PostgreSQL │  │   Kafka    │  │ Blockchain │
│  (数据库)   │  │  (消息队列) │  │  (事件源)   │
└────────────┘  └────────────┘  └────────────┘
         ↑               ↑               │
         │               │               ↓
         │       ┌───────────────┐  ┌──────────┐
         │       │   Consumer    │  │ Listener │
         │       │   (处理事件)   │←─┤ (监听链)  │
         │       └───────────────┘  └──────────┘
         │
         ↓
  ┌──────────────┐
  │  Scheduler   │
  │  (发放积分)   │
  └──────────────┘
```

### 数据流

1. **链上事件** → Listener → Kafka → Consumer → PostgreSQL
2. **积分发放** → Scheduler → PostgreSQL
3. **用户查询** → Frontend → API → PostgreSQL
4. **空投管理** → Admin UI → API → PostgreSQL

---

## 🚀 快速开始

### 环境要求

- Docker & Docker Compose
- Node.js 20+ (用于本地开发)
- Go 1.21+ (用于后端开发)
- MetaMask 或其他 Web3 钱包

### 一键启动（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd loyalty-points-system-final

# 启动所有服务
docker-compose up --build
```

服务启动后访问：
- **前端**: http://localhost:5173
- **API**: http://localhost:8080
- **API 文档**: http://localhost:8080/health

### 验证部署

```bash
# 检查所有容器状态
docker-compose ps

# 应该看到 8 个服务都在运行:
# - zookeeper
# - kafka
# - postgres
# - listener
# - consumer
# - scheduler
# - api
# - frontend
```

### 快速测试

```bash
# 健康检查
curl http://localhost:8080/health

# 获取排行榜
curl http://localhost:8080/leaderboard

# 查看用户积分（替换为实际地址）
curl http://localhost:8080/users/0xYourAddress/points
```

---


## 🛠️ 开发指南

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

**主要技术栈**:
- React 18
- Vite 5
- React Router 6
- ethers.js 6
- lucide-react (图标)

**关键文件**:
- `src/main.jsx` - 路由配置
- `src/components/Header.jsx` - 导航栏
- `src/web3/WalletContext.jsx` - 钱包连接
- `src/views/VaultView.jsx` - 理财金库（双模式）
- `src/views/RWAMarketView.jsx` - RWA 商城
- `src/views/OnchainAirdropView.jsx` - 链上空投（备用）

---

### 后端开发

#### 本地运行 Go 服务

```bash
# 启动基础设施（仅 Kafka、Zookeeper、PostgreSQL）
docker-compose up -d zookeeper kafka postgres

# 运行所有 Go 服务
./start-services.sh

# 或单独运行
go run services/api/cmd/main.go
go run services/listener/cmd/main.go
go run services/consumer/cmd/main.go
go run services/scheduler/cmd/main.go
```

**日志位置**: `logs/` 目录
- `logs/api.log`
- `logs/listener.log`
- `logs/consumer.log`
- `logs/scheduler.log`

#### API 端点

**用户端点**:
```
GET  /health
GET  /users/:address/balance
GET  /users/:address/points
GET  /leaderboard
POST /graphql
```

**空投端点**:
```
GET  /api/airdrop/campaigns
GET  /api/airdrop/campaigns/:id
GET  /api/airdrop/campaigns/:id/eligibility?address=
POST /api/airdrop/campaigns/:id/claim
```

**管理员端点** (需要 Bearer Token):
```
POST /api/admin/airdrop/campaigns
POST /api/admin/airdrop/campaigns/:id/allocations/import
POST /api/admin/airdrop/campaigns/:id/activate
POST /api/admin/airdrop/campaigns/:id/close
```

---

### 链上空投开发

详细文档请参考：
- `onchain-airdrop/README.md` - 使用指南
- `onchain-airdrop/MIGRATION_GUIDE.md` - 迁移指南

```bash
cd onchain-airdrop/contracts

# 安装依赖
npm install

# 编译合约
npm run compile

# 生成 Merkle Tree
npm run merkle:generate data/whitelist.csv

# 部署到 Sepolia
npm run deploy:sepolia

# 部署 Subgraph
cd ../subgraph
npm install
npm run codegen
npm run build
npm run deploy
```

---

## 📊 监控和日志

### 容器状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看特定服务日志
docker-compose logs -f frontend
docker-compose logs -f api
docker-compose logs -f listener

# 查看所有日志
docker-compose logs -f
```

### 数据库管理

```bash
# 连接到 PostgreSQL
docker exec -it loyalty-postgres psql -U postgres -d loyalty_points

# 常用查询
SELECT COUNT(*) FROM users;
SELECT * FROM points ORDER BY points DESC LIMIT 10;
SELECT * FROM airdrop_campaigns;
```

### 健康检查

访问 http://localhost:5173/status 查看所有微服务健康状态。

---

## 🔧 配置说明

### 环境变量

主要环境变量（`.env` 文件）：

```env
# 数据库
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=loyalty_points
DATABASE_URL=postgres://postgres:postgres@postgres:5432/loyalty_points?sslmode=disable

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC_RAW=events.raw

# API
API_PORT=8080
API_ALLOW_ORIGIN=*

# 积分系统
POINTS_RATE=0.05
SCHEDULER_INTERVAL_SEC=60

# 区块链监听
LISTENER_MODE=real
CHAINS_JSON=[{"name":"sepolia","wss_url":"wss://...","token_address":"0x...","staking_address":"0x...","confirmations":6}]
```

### 端口映射

| 服务 | 容器端口 | 宿主机端口 | 说明 |
|------|---------|-----------|------|
| Frontend | 80 | 5173 | 前端应用 |
| API | 8080 | 8080 | 后端 API |
| PostgreSQL | 5432 | 5432 | 数据库 |
| Kafka | 9092 | 9092 | 消息队列 |
| Zookeeper | 2181 | 2181 | 协调器 |
| Listener | 8090 | 8090 | 监听服务 |

---

---

## 🧪 测试

### 手动测试

```bash
# API 健康检查
curl http://localhost:8080/health

# 获取用户余额
curl http://localhost:8080/users/0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1/balance

# 获取用户积分
curl http://localhost:8080/users/0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1/points

# 获取排行榜
curl http://localhost:8080/leaderboard

# GraphQL 查询
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ balance(address: \"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1\") { address balance } }"}'
```

### 前端测试

```bash
cd frontend
npm run test              # 运行测试
npm run test:watch        # 监听模式
npm run test:coverage     # 生成覆盖率报告
```

---

---

## 🤝 贡献

欢迎贡献代码、提出问题或建议！

### 开发流程

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

---

## 📄 许可证

MIT License

---

---

## 🙏 致谢

本项目使用了以下优秀的开源项目：

**前端**:
- React - UI 框架
- Vite - 构建工具
- ethers.js - 以太坊库
- lucide-react - 图标库

**后端**:
- Go - 编程语言
- Gin - Web 框架
- PostgreSQL - 数据库
- Apache Kafka - 消息队列

**区块链**:
- Hardhat - 智能合约开发
- OpenZeppelin - 安全合约库
- The Graph - 链上数据索引

---

<div align="center">

**Built with ❤️ by Yieldera Team**

Enter the Yieldera 🚀

</div>
