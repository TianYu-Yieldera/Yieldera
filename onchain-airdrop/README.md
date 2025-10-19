# Yieldera 链上空投系统

完整的链上空投解决方案，包含智能合约、Subgraph 索引和前端界面。

---

## 🌟 特性

- ✅ **基于 Merkle Tree**：Gas 高效的白名单验证
- ✅ **去中心化**：完全链上执行，透明可验证
- ✅ **多活动支持**：同一合约管理多个空投活动
- ✅ **批量领取**：用户可一次性领取多个活动
- ✅ **实时索引**：使用 The Graph 提供快速查询
- ✅ **安全性**：使用 OpenZeppelin 库，防止重入攻击
- ✅ **紧急提取**：管理员可在活动结束后回收剩余代币

---

## 📂 项目结构

```
onchain-airdrop/
├── contracts/                    # 智能合约
│   ├── YielderaAirdrop.sol      # 主合约
│   ├── MockERC20.sol            # 测试 ERC20 代币
│   ├── hardhat.config.js        # Hardhat 配置
│   ├── scripts/
│   │   ├── deploy.js            # 部署脚本
│   │   └── generate-merkle-tree.js  # Merkle Tree 生成工具
│   ├── data/                    # 白名单和 Merkle Tree 数据
│   └── package.json
│
├── subgraph/                    # The Graph 索引服务
│   ├── schema.graphql           # GraphQL Schema
│   ├── subgraph.yaml            # Subgraph 配置
│   ├── src/
│   │   └── mapping.ts           # 事件处理器
│   └── package.json
│
├── MIGRATION_GUIDE.md           # 迁移指南
└── README.md                    # 本文件
```

---

## 🚀 快速开始

### 1. 安装依赖

```bash
# 智能合约
cd contracts
npm install

# Subgraph
cd ../subgraph
npm install
```

### 2. 生成 Merkle Tree

准备白名单 CSV 文件（`contracts/data/whitelist.csv`）：

```csv
address,amount
0x70997970C51812dc3A010C7d01b50e0d17dc79C8,1000000000000000000000
0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC,2000000000000000000000
```

生成 Merkle Tree：

```bash
cd contracts
npm run merkle:generate data/whitelist.csv
```

保存输出的 Merkle Root！

### 3. 部署智能合约

配置环境变量（`.env`）：

```env
PRIVATE_KEY=your_private_key
SEPOLIA_RPC_URL=https://rpc.sepolia.org
ETHERSCAN_API_KEY=your_api_key
```

部署到 Sepolia 测试网：

```bash
npm run deploy:sepolia
```

记录输出的合约地址！

### 4. 创建空投活动

使用 Hardhat Console：

```bash
npx hardhat console --network sepolia
```

```javascript
const airdrop = await ethers.getContractAt("YielderaAirdrop", "0xYourAirdropAddress");
const token = await ethers.getContractAt("IERC20", "0xYourTokenAddress");

// 授权
await token.approve("0xYourAirdropAddress", ethers.parseEther("100000"));

// 创建活动
await airdrop.createCampaign(
  "Genesis Airdrop",
  "Thank you for your support!",
  "0xYourMerkleRoot",
  ethers.parseEther("100000"),
  Math.floor(Date.now() / 1000),
  Math.floor(Date.now() / 1000) + 86400 * 30
);
```

### 5. 部署 Subgraph

更新 `subgraph/subgraph.yaml` 中的合约地址和起始区块：

```yaml
source:
  address: "0xYourAirdropAddress"
  startBlock: 5123456
```

部署：

```bash
cd subgraph
npm run codegen
npm run build
npm run deploy
```

---

## 🔧 开发指南

### 编译合约

```bash
cd contracts
npm run compile
```

### 运行测试

```bash
npm run test
```

### 本地部署

启动本地节点：

```bash
npx hardhat node
```

部署到本地网络：

```bash
npm run deploy:local
```

---

## 📖 智能合约接口

### 核心函数

#### `createCampaign`

创建新的空投活动（仅管理员）

```solidity
function createCampaign(
    string memory name,
    string memory description,
    bytes32 merkleRoot,
    uint256 totalBudget,
    uint256 startTime,
    uint256 endTime
) external onlyOwner returns (uint256)
```

#### `claim`

领取单个活动的空投

```solidity
function claim(
    uint256 campaignId,
    uint256 amount,
    bytes32[] calldata merkleProof
) external
```

#### `claimMultiple`

批量领取多个活动

```solidity
function claimMultiple(
    uint256[] calldata campaignIds,
    uint256[] calldata amounts,
    bytes32[][] calldata merkleProofs
) external
```

#### `updateCampaignStatus`

更新活动状态（仅管理员）

```solidity
function updateCampaignStatus(
    uint256 campaignId,
    bool isActive
) external onlyOwner
```

#### `emergencyWithdraw`

紧急提取剩余代币（仅管理员）

```solidity
function emergencyWithdraw(uint256 campaignId) external onlyOwner
```

### 查询函数

```solidity
// 检查用户是否已领取
function hasClaimed(uint256 campaignId, address user) external view returns (bool)

// 获取活动信息
function getCampaign(uint256 campaignId) external view returns (Campaign memory)

// 获取剩余预算
function getRemainingBudget(uint256 campaignId) external view returns (uint256)

// 获取用户已领取金额
function getUserClaimedAmount(uint256 campaignId, address user) external view returns (uint256)
```

---

## 📊 Subgraph 查询

### 获取所有活动

```graphql
query {
  airdropCampaigns(first: 10, orderBy: createdAt, orderDirection: desc) {
    id
    name
    description
    totalBudget
    claimedAmount
    remainingBudget
    startTime
    endTime
    isActive
    participantCount
    claimCount
  }
}
```

### 获取用户领取记录

```graphql
query($user: String!) {
  user(id: $user) {
    id
    totalClaimed
    claimCount
    campaignCount
    claims {
      id
      campaign {
        name
      }
      amount
      timestamp
    }
  }
}
```

### 获取每日统计

```graphql
query($campaignId: String!) {
  dailySnapshots(
    where: { campaign: $campaignId }
    orderBy: date
    orderDirection: desc
  ) {
    date
    claimCount
    totalAmount
    uniqueClaimers
    cumulativeClaimCount
    cumulativeAmount
  }
}
```

### 全局统计

```graphql
query {
  globalStats(id: "global") {
    totalCampaigns
    activeCampaigns
    totalDistributed
    totalClaims
    totalUsers
    updatedAt
  }
}
```

---

## 🔐 安全考虑

### 已实施的安全措施

1. **ReentrancyGuard**：防止重入攻击
2. **Ownable**：关键函数仅管理员可调用
3. **Merkle Proof 验证**：防止未授权领取
4. **状态检查**：活动时间、预算、领取状态验证
5. **OpenZeppelin 库**：使用经过审计的标准库

### 建议

- ✅ 主网部署前进行专业审计
- ✅ 使用多签钱包管理合约
- ✅ 设置合理的活动时间和预算上限
- ✅ 定期监控合约事件和余额

---

## 💡 使用案例

### 1. 早期用户奖励

```javascript
// 根据用户在旧系统中的积分生成白名单
const users = await db.query("SELECT address, points FROM points WHERE points > 100");
const whitelist = users.map(u => ({
  address: u.address,
  amount: ethers.parseEther(u.points.toString())
}));
```

### 2. DeFi 用户激励

```javascript
// 奖励在 Vault 中质押的用户
const stakers = await vault.getStakers();
const whitelist = stakers.map(s => ({
  address: s.address,
  amount: s.stakedAmount * 0.1 // 10% 奖励
}));
```

### 3. RWA 投资者空投

```javascript
// 奖励购买过 RWA 资产的用户
const buyers = await rwaMarket.getBuyers();
const whitelist = buyers.map(b => ({
  address: b.address,
  amount: b.totalPurchased * 0.05 // 5% 返现
}));
```

---

## 🛠️ 工具和脚本

### Merkle Tree 生成器

位置：`contracts/scripts/generate-merkle-tree.js`

功能：
- 读取 CSV 白名单
- 生成 Merkle Tree
- 输出 Merkle Root 和所有 Proof
- 保存为 JSON 文件

用法：

```bash
npm run merkle:generate data/whitelist.csv
```

### 部署脚本

位置：`contracts/scripts/deploy.js`

功能：
- 自动部署合约
- 可选择部署测试代币
- 保存部署信息
- 提供验证命令

用法：

```bash
npm run deploy:sepolia
```

---

## 📈 Gas 优化

### 合约层面

- ✅ 使用 `immutable` 存储 token 地址
- ✅ 使用 Merkle Tree 而非链上白名单
- ✅ 批量操作支持（`claimMultiple`）
- ✅ 事件参数使用 `indexed` 优化查询

### 用户层面

- 在 Gas 价格较低时领取
- 使用 `claimMultiple` 批量领取多个活动
- 考虑实现 Gasless（Meta Transaction）方案

---

## 🔄 升级路径

### 当前系统：不可升级

合约部署后无法修改逻辑。如需升级：

1. 部署新版本合约
2. 迁移数据到新合约
3. 更新前端和 Subgraph 配置

### 未来可实现的升级方案

- **代理模式**：使用 UUPS 或 Transparent Proxy
- **模块化**：将功能拆分为多个可替换的模块
- **DAO 治理**：由社区投票决定升级

---

## 📚 相关资源

- [Merkle Tree 详解](https://en.wikipedia.org/wiki/Merkle_tree)
- [The Graph 官方文档](https://thegraph.com/docs/)
- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Hardhat 使用指南](https://hardhat.org/tutorial)
- [EIP-2612: Permit](https://eips.ethereum.org/EIPS/eip-2612)

---

## 🤝 贡献

欢迎贡献代码、提出问题或建议！

---

## 📄 许可证

MIT License

---

## 联系方式

- 📧 Email: dev@yieldera.com
- 💬 Discord: https://discord.gg/yieldera
- 🐦 Twitter: @yieldera
- 📖 文档: https://docs.yieldera.com

---

**Built with ❤️ by Yieldera Team**
