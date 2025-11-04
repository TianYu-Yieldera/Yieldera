# 智能合约部署和集成指南

本文档说明如何完成智能合约部署并将其与后端服务集成。

## 当前状态

### ✅ 已完成
1. **前端钱包签名** - 已实现EIP-712签名
   - `frontend/src/components/TradingForm.jsx` - 订单创建签名
   - `frontend/src/views/TreasuryHoldingsView.jsx` - 收益领取签名

2. **区块链客户端工具** - 已创建基础设施
   - `internal/blockchain/client.go` - L1/L2客户端管理
   - `internal/blockchain/signature.go` - EIP-712签名验证

3. **智能合约** - 已编写完成
   - Layer 1: CollateralVault, StateRegistry, LoyaltyUSD, Gateway
   - Layer 2: IntegratedVault, StateAggregator, DeFi Adapters
   - Treasury: TreasuryAssetFactory, TreasuryMarketplace, TreasuryYieldDistributor, TreasuryPriceOracle

4. **部署脚本** - 已准备就绪
   - `scripts/deploy-treasury.js`
   - `scripts/deploy-phase3.js`
   - `scripts/deploy-phase4.js`
   - `scripts/deploy-phase5.js`

### ⏳ 待完成
1. **部署智能合约到测试网**
2. **配置环境变量**
3. **后端集成合约调用**
4. **端到端测试**

---

## 第一步：部署智能合约

### 1.1 准备工作

#### 获取测试币
1. **Sepolia ETH** (L1 测试网)
   - 访问: https://sepoliafaucet.com/
   - 或: https://www.infura.io/faucet/sepolia

2. **Arbitrum Sepolia ETH** (L2 测试网)
   - 访问: https://faucet.quicknode.com/arbitrum/sepolia
   - 或通过Sepolia桥接: https://bridge.arbitrum.io/

#### 配置RPC节点
创建 `.env` 文件（基于 `.env.example`）:

```bash
# 私钥 (确保账户有测试币)
PRIVATE_KEY=your_private_key_here_without_0x_prefix

# L1 (Sepolia)
L1_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_ALCHEMY_KEY
L1_WSS_URL=wss://eth-sepolia.g.alchemy.com/v2/YOUR_ALCHEMY_KEY

# L2 (Arbitrum Sepolia)
L2_RPC_URL=https://sepolia-rollup.arbitrum.io/rpc
ARBITRUM_SEPOLIA_RPC_URL=https://sepolia-rollup.arbitrum.io/rpc

# 可选: Etherscan API Keys (用于合约验证)
ETHERSCAN_API_KEY=your_etherscan_api_key
ARBISCAN_API_KEY=your_arbiscan_api_key
```

#### 获取免费RPC API Key
- **Alchemy**: https://www.alchemy.com/ (推荐)
- **Infura**: https://www.infura.io/
- **QuickNode**: https://www.quicknode.com/

### 1.2 部署步骤

#### Phase 3: 部署L1核心合约
```bash
cd /home/tianyu/loyalty-points-system-final
npx hardhat run scripts/deploy-phase3.js --network sepolia
```

输出示例:
```
L1_COLLATERAL_VAULT=0x1234...
L1_STATE_REGISTRY=0x5678...
L1_LOYALTY_USD=0xabcd...
L1_GATEWAY=0xef01...
```

#### Phase 4: 部署L2核心合约
```bash
npx hardhat run scripts/deploy-phase4.js --network arbitrumSepolia
```

输出示例:
```
L2_INTEGRATED_VAULT=0x2345...
L2_STATE_AGGREGATOR=0x6789...
```

#### Phase 5: 部署L2 DeFi适配器
```bash
npx hardhat run scripts/deploy-phase5.js --network arbitrumSepolia
```

输出示例:
```
L2_AAVE_ADAPTER=0x3456...
L2_COMPOUND_ADAPTER=0x7890...
L2_UNISWAP_ADAPTER=0xbcde...
```

#### Treasury: 部署Treasury模块
```bash
npx hardhat run scripts/deploy-treasury.js --network arbitrumSepolia
```

输出示例:
```
TreasuryAssetFactory: 0x4567...
TreasuryPriceOracle: 0x8901...
TreasuryYieldDistributor: 0xcdef...
TreasuryMarketplace: 0x0123...
```

### 1.3 保存部署信息

每个部署脚本会在 `deployments/` 目录创建JSON文件，记录:
- 合约地址
- 部署时间
- 网络信息
- 部署者地址

**重要**: 备份 `deployments/` 目录!

---

## 第二步：配置环境变量

### 2.1 更新 `.env`

将部署的合约地址添加到 `.env`:

```bash
# L1 Contracts (Sepolia)
L1_COLLATERAL_VAULT=0x...
L1_STATE_REGISTRY=0x...
L1_LOYALTY_USD=0x...
L1_GATEWAY=0x...

# L2 Core Contracts (Arbitrum Sepolia)
L2_INTEGRATED_VAULT=0x...
L2_STATE_AGGREGATOR=0x...

# L2 DeFi Adapters
L2_AAVE_ADAPTER=0x...
L2_COMPOUND_ADAPTER=0x...
L2_UNISWAP_ADAPTER=0x...

# Treasury Contracts (Arbitrum Sepolia)
L2_TREASURY_FACTORY=0x...
L2_TREASURY_MARKETPLACE=0x...
L2_TREASURY_YIELD_DISTRIBUTOR=0x...
L2_TREASURY_ORACLE=0x...
```

### 2.2 更新前端配置

在 `frontend/.env` 或 `frontend/.env.local`:

```bash
VITE_API_URL=http://localhost:8080
VITE_TREASURY_MARKETPLACE_ADDRESS=0x...
VITE_TREASURY_YIELD_DISTRIBUTOR_ADDRESS=0x...
VITE_CHAIN_ID=421614  # Arbitrum Sepolia
```

更新前端中的合约地址:
- `frontend/src/components/TradingForm.jsx` 第60行
- `frontend/src/views/TreasuryHoldingsView.jsx` 第91行

```javascript
// 替换这一行:
verifyingContract: '0x0000000000000000000000000000000000000000',

// 改为:
verifyingContract: import.meta.env.VITE_TREASURY_MARKETPLACE_ADDRESS,
```

---

## 第三步：后端集成合约调用

### 3.1 生成合约Go绑定

需要为每个合约生成Go binding:

```bash
# 安装abigen
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 生成TreasuryMarketplace绑定
abigen --abi artifacts/contracts/layer2/treasury/TreasuryMarketplace.sol/TreasuryMarketplace.json \
       --pkg contracts \
       --type TreasuryMarketplace \
       --out internal/contracts/treasury_marketplace.go

# 生成TreasuryYieldDistributor绑定
abigen --abi artifacts/contracts/layer2/treasury/TreasuryYieldDistributor.sol/TreasuryYieldDistributor.json \
       --pkg contracts \
       --type TreasuryYieldDistributor \
       --out internal/contracts/treasury_yield_distributor.go

# 类似地为其他合约生成绑定...
```

### 3.2 实现Treasury Handler合约调用

修改 `services/api/handlers/treasury_handler.go`:

```go
import (
	"loyalty-points-system/internal/blockchain"
	"loyalty-points-system/internal/contracts"
	"github.com/ethereum/go-ethereum/common"
)

// CreateTreasuryOrder - 完整实现
func CreateTreasuryOrder(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AssetID       int64   `json:"asset_id"`
			UserAddress   string  `json:"user_address"`
			OrderType     string  `json:"order_type"`
			TokenAmount   float64 `json:"token_amount"`
			PricePerToken float64 `json:"price_per_token"`
			Signature     string  `json:"signature"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 1. 验证签名
		domain := blockchain.EIP712Domain{
			Name:              "Treasury Marketplace",
			Version:           "1",
			ChainID:           big.NewInt(421614), // Arbitrum Sepolia
			VerifyingContract: common.HexToAddress(os.Getenv("L2_TREASURY_MARKETPLACE")),
		}

		// ... 构建messageHash并验证签名 ...

		// 2. 调用智能合约
		client, err := blockchain.NewClient()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create blockchain client"})
			return
		}
		defer client.Close()

		// 获取合约实例
		marketplaceAddr := common.HexToAddress(os.Getenv("L2_TREASURY_MARKETPLACE"))
		marketplace, err := contracts.NewTreasuryMarketplace(marketplaceAddr, client.L2Client)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load contract"})
			return
		}

		// 准备交易选项
		auth, err := client.GetL2TransactOpts(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create transaction options"})
			return
		}

		// 调用合约方法
		var tx *types.Transaction
		if req.OrderType == "BUY" {
			tx, err = marketplace.CreateBuyOrder(
				auth,
				big.NewInt(req.AssetID),
				// ... 其他参数
			)
		} else {
			tx, err = marketplace.CreateSellOrder(
				auth,
				big.NewInt(req.AssetID),
				// ... 其他参数
			)
		}

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create order"})
			return
		}

		// 等待交易确认
		receipt, err := client.WaitForL2Transaction(c.Request.Context(), tx)
		if err != nil {
			c.JSON(500, gin.H{"error": "Transaction failed"})
			return
		}

		// 3. 返回结果
		c.JSON(200, gin.H{
			"order_id": "extract_from_event",
			"tx_hash": tx.Hash().Hex(),
			"block_number": receipt.BlockNumber,
		})
	}
}
```

### 3.3 类似地实现其他handlers

需要更新的文件:
- `services/api/handlers/l1_handler.go` - L1存款/提款
- `services/api/handlers/l2_handler.go` - L2存款/提款
- `services/api/handlers/bridge_handler.go` - 跨链桥接

---

## 第四步：测试

### 4.1 编译并启动服务

```bash
# 编译Go服务
go build -o bin/api services/api/cmd/main.go
go build -o bin/listener services/listener/cmd/main.go
go build -o bin/consumer services/consumer/cmd/main.go

# 或使用Docker Compose
docker-compose up -d
```

### 4.2 测试流程

1. **连接钱包** - 在前端连接MetaMask到Arbitrum Sepolia
2. **浏览Treasury市场** - 访问 `/treasury`
3. **创建订单** - 点击资产，输入金额和价格
4. **签名确认** - MetaMask会弹出签名请求
5. **等待确认** - 交易上链后更新UI

### 4.3 调试工具

- **Sepolia Etherscan**: https://sepolia.etherscan.io/
- **Arbiscan Sepolia**: https://sepolia.arbiscan.io/
- **检查交易**: 使用上述网站搜索交易哈希
- **查看日志**: `docker-compose logs -f api`

---

## 常见问题

### Q: 部署失败 "insufficient funds"
A: 确保部署账户有足够的测试币。Sepolia需要约0.5 ETH，Arbitrum Sepolia需要约0.1 ETH。

### Q: 交易revert "execution reverted"
A: 检查合约逻辑，可能是权限不足或参数错误。查看Etherscan的错误信息。

### Q: 前端签名但后端未收到
A: 检查CORS配置和API_ALLOW_ORIGIN环境变量。

### Q: 合约调用超时
A: 增加gas limit或检查RPC节点连接。

---

## 下一步建议

完成基本功能后，建议:

1. **添加单元测试** - 测试合约交互函数
2. **添加集成测试** - 端到端测试流程
3. **监控和告警** - 集成Prometheus指标
4. **审计合约** - 聘请专业审计团队
5. **优化gas** - 批量操作和优化存储
6. **文档完善** - API文档和用户手册

---

## 参考资料

- [Hardhat文档](https://hardhat.org/docs)
- [go-ethereum文档](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)
- [EIP-712规范](https://eips.ethereum.org/EIPS/eip-712)
- [Arbitrum文档](https://docs.arbitrum.io/)

---

**最后更新**: 2025-11-04
**维护者**: TianYu
