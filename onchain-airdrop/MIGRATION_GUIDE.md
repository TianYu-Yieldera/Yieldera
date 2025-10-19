# Yieldera 空投系统迁移指南

## 📋 目录

1. [为什么要迁移](#为什么要迁移)
2. [架构对比](#架构对比)
3. [部署步骤](#部署步骤)
4. [用户迁移](#用户迁移)
5. [常见问题](#常见问题)

---

## 为什么要迁移

### 当前系统（链下积分空投）的局限性

❌ **不透明**：用户无法验证积分分配的公平性
❌ **中心化**：完全依赖后端数据库，存在单点故障风险
❌ **可信度低**：企业可以随时修改积分规则或余额
❌ **不可交易**：积分无法在二级市场流通
❌ **监管风险**：链下系统难以满足合规要求

### 新系统（链上代币空投）的优势

✅ **完全透明**：所有分配记录永久存储在区块链上，可公开验证
✅ **去中心化**：智能合约自动执行，无需信任中间方
✅ **可验证性**：用户可以通过 Merkle Proof 独立验证自己的资格
✅ **可交易**：ERC-20 代币可在 DEX 和 CEX 自由交易
✅ **成熟商业系统必备**：符合 Web3 行业标准，提升品牌信誉度

---

## 架构对比

### 旧架构（链下积分）

```
┌─────────────┐
│   Frontend  │
└──────┬──────┘
       │ REST API
       ↓
┌─────────────┐
│   API 服务   │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│  PostgreSQL │ ← 积分存储在数据库中
└─────────────┘
```

**问题**：
- 用户必须信任 API 服务器
- 数据可被管理员随意修改
- 无法独立验证

---

### 新架构（链上代币）

```
┌─────────────┐
│   Frontend  │
└──────┬──────┘
       │
       ├─────────────────────┐
       │                     │
       ↓                     ↓
┌─────────────┐      ┌─────────────┐
│  Smart      │      │  Subgraph   │
│  Contract   │      │  (索引服务)  │
└──────┬──────┘      └──────┬──────┘
       │                     │
       └──────────┬──────────┘
                  ↓
           ┌─────────────┐
           │  Ethereum   │ ← 数据存储在区块链上
           └─────────────┘
```

**优势**：
- 去中心化：智能合约公开透明
- 可验证：Merkle Tree 加密证明
- 不可篡改：区块链永久记录
- 高性能：Subgraph 快速查询

---

## 部署步骤

### 1️⃣ 环境准备

#### 安装依赖

```bash
# 进入合约目录
cd onchain-airdrop/contracts

# 安装依赖
npm install

# 配置环境变量
cp .env.example .env
```

#### 配置 `.env` 文件

```env
# 部署者私钥（重要：不要泄露！）
PRIVATE_KEY=your_private_key_here

# RPC 节点
SEPOLIA_RPC_URL=https://rpc.sepolia.org
MAINNET_RPC_URL=https://eth.llamarpc.com

# Etherscan API Key（用于验证合约）
ETHERSCAN_API_KEY=your_api_key

# 代币地址（如果已部署）
TOKEN_ADDRESS=0x0000000000000000000000000000000000000000
```

---

### 2️⃣ 生成 Merkle Tree

#### 准备白名单 CSV

创建 `onchain-airdrop/contracts/data/whitelist.csv`：

```csv
address,amount
0x70997970C51812dc3A010C7d01b50e0d17dc79C8,1000000000000000000000
0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC,2000000000000000000000
0x90F79bf6EB2c4f870365E785982E1f101E93b906,500000000000000000000
```

> **注意**：`amount` 使用 wei 单位（18 位小数）
> 例如：`1000000000000000000000` = 1000 YLD

#### 生成 Merkle Tree

```bash
npm run merkle:generate data/whitelist.csv
```

**输出示例**：

```
🌳 开始生成 Merkle Tree...
📄 读取 CSV 文件: data/whitelist.csv
✅ 成功解析 3 个地址

🌲 Merkle Tree 信息:
   根哈希 (Merkle Root): 0x1234567890abcdef...
   叶子节点数量: 3
   树的深度: 2

📦 完整数据已保存: data/merkle-tree-1234567890.json
📄 摘要数据已保存: data/merkle-summary.json
```

> **重要**：保存 Merkle Root，创建活动时需要使用！

---

### 3️⃣ 部署智能合约

#### 测试网部署（Sepolia）

```bash
# 编译合约
npm run compile

# 部署到 Sepolia
npm run deploy:sepolia
```

**输出示例**：

```
🚀 开始部署 YielderaAirdrop 合约...
📝 部署者地址: 0xYourAddress
💰 部署者余额: 0.5 ETH

✅ 测试代币部署成功: 0xTokenAddress
   代币名称: Yieldera Token (YLD)
   总供应量: 1,000,000,000 YLD

✅ YielderaAirdrop 部署成功: 0xAirdropAddress

============================================================
🎉 部署完成！
============================================================
网络: sepolia
链ID: 11155111
区块高度: 5123456
代币地址: 0xTokenAddress
空投合约: 0xAirdropAddress
============================================================
```

#### 验证合约（可选）

```bash
npx hardhat verify --network sepolia 0xAirdropAddress 0xTokenAddress
```

---

### 4️⃣ 创建空投活动

使用 Hardhat Console 或编写脚本：

```javascript
// scripts/create-campaign.js
const hre = require("hardhat");

async function main() {
  const airdropAddress = "0xYourAirdropContract";
  const tokenAddress = "0xYourTokenAddress";

  const airdrop = await hre.ethers.getContractAt("YielderaAirdrop", airdropAddress);
  const token = await hre.ethers.getContractAt("IERC20", tokenAddress);

  // 授权合约使用代币
  console.log("授权代币...");
  const totalBudget = hre.ethers.parseEther("100000");
  await token.approve(airdropAddress, totalBudget);

  // 创建活动
  console.log("创建空投活动...");
  const tx = await airdrop.createCampaign(
    "Yieldera Genesis Airdrop",                    // 活动名称
    "感谢早期用户的支持！",                          // 活动描述
    "0x1234567890abcdef...",                       // Merkle Root（从步骤2获取）
    totalBudget,                                   // 总预算
    Math.floor(Date.now() / 1000),                // 开始时间（现在）
    Math.floor(Date.now() / 1000) + 86400 * 30    // 结束时间（30天后）
  );

  const receipt = await tx.wait();
  console.log("活动创建成功！Campaign ID:", receipt.logs[0].args.campaignId);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
```

运行脚本：

```bash
npx hardhat run scripts/create-campaign.js --network sepolia
```

---

### 5️⃣ 部署 Subgraph

#### 配置 Subgraph

编辑 `onchain-airdrop/subgraph/subgraph.yaml`：

```yaml
dataSources:
  - kind: ethereum
    name: YielderaAirdrop
    network: sepolia
    source:
      address: "0xYourAirdropContract"  # 替换为实际合约地址
      abi: YielderaAirdrop
      startBlock: 5123456               # 替换为部署区块号
```

#### 部署到 The Graph Studio

```bash
cd onchain-airdrop/subgraph

# 安装依赖
npm install

# 生成类型
npm run codegen

# 构建
npm run build

# 部署到 The Graph Studio
npm run deploy
```

> **提示**：首次使用需要在 [The Graph Studio](https://thegraph.com/studio/) 创建 Subgraph

---

### 6️⃣ 配置前端

编辑 `frontend/src/views/OnchainAirdropView.jsx`：

```javascript
// 更新合约地址
const AIRDROP_CONTRACT_ADDRESS = "0xYourAirdropContract";

// 更新 Subgraph 查询地址（如果使用）
const SUBGRAPH_URL = "https://api.studio.thegraph.com/query/...";
```

#### 部署 Merkle Proof 数据

有两种方式：

**方式 1：上传到 IPFS（推荐）**

```bash
# 使用 Pinata 或其他 IPFS 服务
ipfs add data/merkle-tree-1234567890.json

# 获取 CID，在前端通过 IPFS Gateway 读取
# https://gateway.pinata.cloud/ipfs/<CID>
```

**方式 2：后端 API**

```javascript
// 在后端提供 API 接口
app.get('/api/merkle-proof/:address', (req, res) => {
  const address = req.params.address;
  const data = require('./data/merkle-tree-1234567890.json');
  const proof = data.proofs[address];
  res.json(proof);
});
```

---

### 7️⃣ 重新部署前端

```bash
# 构建前端
docker-compose build frontend

# 重启服务
docker-compose up -d frontend
```

---

## 用户迁移

### 迁移策略

#### 选项 A：快照空投（推荐）

1. **导出现有积分数据**

```sql
-- 从 PostgreSQL 导出积分余额
SELECT address, points
FROM points
WHERE points > 0
ORDER BY points DESC;
```

2. **生成白名单 CSV**

```python
# convert_points_to_csv.py
import csv

# 从数据库读取
points_data = [
  ("0x1234...", 5000),
  ("0x5678...", 3000),
  # ...
]

# 写入 CSV（积分转换为代币，例如 1:1）
with open('whitelist.csv', 'w') as f:
  writer = csv.writer(f)
  writer.writerow(['address', 'amount'])
  for addr, points in points_data:
    # 转换为 wei（18位小数）
    amount = str(points * 10**18)
    writer.writerow([addr, amount])
```

3. **按照上述步骤部署空投活动**

#### 选项 B：双系统并行运行

- 保留旧系统用于历史查询
- 新用户使用链上空投
- 逐步引导用户迁移

---

### 用户通知

**邮件/公告模板**：

```markdown
## 🎉 Yieldera 空投系统升级公告

亲爱的 Yieldera 用户：

我们很高兴地宣布，Yieldera 空投系统已全面升级至链上版本！

### ✨ 新系统亮点

- ✅ **完全透明**：所有分配记录永久存储在区块链上
- ✅ **可独立验证**：通过 Merkle Proof 验证您的资格
- ✅ **真正的所有权**：ERC-20 代币可自由交易

### 📋 如何领取

1. 访问 https://app.yieldera.com/airdrop
2. 连接您的钱包
3. 查看您的空投资格
4. 点击"立即领取"按钮

### 💰 您的空投份额

根据您在旧系统中的积分余额，您有资格领取：

**5,000 YLD 代币**

感谢您对 Yieldera 的支持！

---

Yieldera 团队
```

---

## 常见问题

### Q1: 我的旧积分会怎么样？

**A**: 旧积分将按 1:1 比例转换为新代币空投资格。您可以在新系统中领取等值的 YLD 代币。

---

### Q2: 我可以一次性领取多个活动的空投吗？

**A**: 可以！智能合约提供了 `claimMultiple` 函数，支持批量领取。

```javascript
await contract.claimMultiple(
  [0, 1, 2],           // campaignIds
  [amount0, amount1, amount2],
  [proof0, proof1, proof2]
);
```

---

### Q3: Gas 费用由谁承担？

**A**: 用户需要支付 Gas 费用。在 Sepolia 测试网上，可以从 [水龙头](https://sepoliafaucet.com/) 获取免费 ETH。

---

### Q4: 如果我错过了领取时间怎么办？

**A**: 活动结束后，管理员可以通过 `emergencyWithdraw` 回收剩余代币。建议在活动期内及时领取。

---

### Q5: 如何验证 Merkle Proof 的真实性？

**A**: 您可以使用以下工具：

- [Merkle Tree Verifier](https://lab.miguelmota.com/merkletreejs/example/)
- 或运行本地验证脚本：

```javascript
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');

// 验证您的 proof
const leaf = keccak256(abi.encodePacked(yourAddress, yourAmount));
const isValid = MerkleTree.verify(proof, leaf, merkleRoot);
console.log("Proof valid:", isValid);
```

---

### Q6: Subgraph 查询慢怎么办？

**A**: Subgraph 索引需要时间。如果刚部署，等待 5-10 分钟让区块完全索引。您也可以直接调用智能合约查询：

```javascript
const campaign = await contract.campaigns(0);
const hasClaimed = await contract.hasClaimed(0, yourAddress);
```

---

### Q7: 如何处理大规模空投（10万+ 用户）？

**A**:

1. **分批创建活动**：每个活动 1-2 万用户
2. **优化 Gas**：使用 `claimMultiple` 批量领取
3. **提供 Gasless 方案**：实现 EIP-2612 Permit + Meta Transaction

---

### Q8: 合约安全性如何保证？

**A**:

- ✅ 使用 OpenZeppelin 经过审计的库
- ✅ 实现 ReentrancyGuard 防止重入攻击
- ✅ 使用 Merkle Proof 防止伪造
- ✅ 建议在主网部署前进行专业审计

---

## 📚 参考资料

- [Merkle Tree 原理](https://en.wikipedia.org/wiki/Merkle_tree)
- [The Graph 文档](https://thegraph.com/docs/)
- [OpenZeppelin 合约](https://docs.openzeppelin.com/contracts/)
- [Hardhat 文档](https://hardhat.org/docs)

---

## 🆘 技术支持

如有问题，请联系：

- 📧 Email: support@yieldera.com
- 💬 Discord: https://discord.gg/yieldera
- 📖 文档: https://docs.yieldera.com

---

**祝您顺利完成迁移！🚀**
