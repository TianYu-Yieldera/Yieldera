# 启用 Arbitrum Sepolia 网络

## 问题

当前的 Alchemy API key (`FP6JOVxZoc4lDScODskcP`) 只启用了 Ethereum Sepolia，需要添加 Arbitrum Sepolia 支持。

## 解决方案

### 方法 1：在现有 App 中启用 Arbitrum Sepolia（推荐）

1. 访问 Alchemy Dashboard：
   ```
   https://dashboard.alchemy.com/apps/a58zedmkn5zq3cid/networks
   ```

2. 在 Networks 页面，找到 "Add Network" 或类似按钮

3. 选择添加 **Arbitrum Sepolia** 网络

4. 保存设置

5. API key 将自动支持 Arbitrum Sepolia

### 方法 2：创建新的 Arbitrum Sepolia App

如果上述方法不行，可以创建一个新的 App：

1. 访问 Alchemy Dashboard：
   ```
   https://dashboard.alchemy.com/
   ```

2. 点击 "Create new app"

3. 配置：
   - **Name**: loyalty-points-arbitrum
   - **Chain**: Arbitrum
   - **Network**: Arbitrum Sepolia
   - **Plan**: Free

4. 创建后，复制新的 API key

5. 更新 `backend/.env`：
   ```env
   ARBITRUM_SEPOLIA_WS=wss://arb-sepolia.g.alchemy.com/v2/YOUR_NEW_KEY
   ARBITRUM_SEPOLIA_RPC=https://arb-sepolia.g.alchemy.com/v2/YOUR_NEW_KEY
   ```

## 测试连接

完成配置后，运行测试验证：

```bash
cd backend
npx ts-node test-monitoring.ts
```

成功的输出应该类似：
```
✅ Connected to Arbitrum Sepolia
📦 Current block: 213140xxx
```

## 备用方案：使用公共 RPC

如果 Alchemy 遇到问题，可以使用 Arbitrum 官方公共 RPC（但不推荐用于生产环境）：

```env
ARBITRUM_SEPOLIA_WS=wss://sepolia-rollup.arbitrum.io/rpc
ARBITRUM_SEPOLIA_RPC=https://sepolia-rollup.arbitrum.io/rpc
```

**注意**：公共 RPC 有严格的速率限制，不适合监控系统使用。

## 当前状态

- ✅ Ethereum Sepolia: 已启用
- ❌ Arbitrum Sepolia: 需要启用
- 📝 API Key ID: `a58zedmkn5zq3cid`

## 下一步

1. 按照上述方法启用 Arbitrum Sepolia
2. 运行测试验证连接
3. 启动监控系统
