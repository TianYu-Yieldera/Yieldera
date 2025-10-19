# 快速启动指南 - PointFi Frontend

## 前置准备

确保你已经安装：
- Node.js (v16 或更高版本)
- npm
- MetaMask 浏览器扩展

## Step 1: 启动后端服务

### 选项 A: 使用 Docker Compose (推荐)

```bash
# 在项目根目录
cd /home/tianyu/loyalty-points-system-final

# 启动所有服务（Kafka, PostgreSQL, API, Listener, Consumer, Scheduler）
docker-compose up --build

# 或者在后台运行
docker-compose up -d --build
```

### 选项 B: 本地开发模式

```bash
# 在项目根目录
cd /home/tianyu/loyalty-points-system-final

# 1. 先启动基础设施（Kafka, Zookeeper, PostgreSQL）
docker-compose up zookeeper kafka postgres -d

# 2. 等待几秒让服务启动
sleep 5

# 3. 启动所有 Go 服务
./start-services.sh
```

### 验证后端是否运行

```bash
# 测试 API 健康检查
curl http://localhost:8080/health

# 应该返回：{"ok":true}
```

## Step 2: 安装前端依赖

```bash
# 进入前端目录
cd /home/tianyu/loyalty-points-system-final/frontend_back

# 安装依赖
npm install

# 如果在国内，可以使用镜像加速
# npm install --registry=https://registry.npmmirror.com
```

## Step 3: 启动前端开发服务器

```bash
# 在 frontend_back 目录中
npm run dev
```

你应该看到类似输出：
```
  VITE v5.4.10  ready in 500 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: http://0.0.0.0:5173/
  ➜  press h + enter to show help
```

## Step 4: 在浏览器中打开

1. **打开浏览器**访问：`http://localhost:5173`

2. **你会看到首页** (Landing Page)，展示所有功能模块

## Step 5: 连接 MetaMask 钱包

1. **点击右上角的 "Connect Wallet" 按钮**

2. **MetaMask 弹窗会出现**，选择要连接的账户

3. **点击"连接"**

4. **如果你不在 Sepolia 测试网**，会出现橙色按钮"切到 Sepolia"
   - 点击按钮
   - MetaMask 会提示切换网络
   - 确认切换

## Step 6: 查看 Dashboard

1. **点击导航栏的"概览"** 或访问 `http://localhost:5173/dashboard`

2. **Dashboard 会显示**：
   - Net Worth (余额)
   - Points (积分)
   - Staked TVL (质押总值)
   - Badges (徽章数量)

3. **数据从后端 API 获取**：
   - 如果你的地址有活动，会显示实际数据
   - 如果是新地址，显示为 0

## Step 7: 测试其他页面

导航栏中的其他链接：
- **资产** - Portfolio (Coming Soon)
- **DeFi 池** - Staking (Coming Soon)
- **排行榜** - Leaderboard (Coming Soon)
- **NFT 徽章** - Badges (Coming Soon)
- **空投** - Airdrops (Coming Soon)
- **指数** - Subgraph (Coming Soon)
- **系统状态** - System Status (Coming Soon)

## 常见问题排查

### 问题 1: 前端启动失败

```bash
# 清除 node_modules 重新安装
rm -rf node_modules package-lock.json
npm install
npm run dev
```

### 问题 2: API 连接失败

检查后端是否运行：
```bash
curl http://localhost:8080/health
```

如果没有响应：
```bash
# 检查 Docker 容器状态
docker-compose ps

# 查看 API 日志
docker-compose logs api

# 或者如果使用本地服务
tail -f logs/api.log
```

### 问题 3: MetaMask 无法连接

- 确保已安装 MetaMask 扩展
- 刷新页面
- 检查浏览器控制台是否有错误
- 尝试在 MetaMask 中手动断开并重新连接

### 问题 4: 看不到数据

这是正常的，如果：
- 你的地址是新的，还没有任何链上活动
- 后端 Listener 还没有监听到事件

**测试地址**：可以使用默认地址 `0x3C07226A3f1488320426eB5FE9976f72E5712346`

### 问题 5: 端口被占用

如果 5173 端口被占用：
```bash
# 修改 vite.config.js 中的端口
# server: { host: '0.0.0.0', port: 5174 }

# 或者关闭占用端口的进程
lsof -ti:5173 | xargs kill -9
```

## 开发模式特性

### 热重载
修改代码后，页面会自动刷新

### 查看网络请求
打开浏览器开发者工具 (F12) → Network 标签，可以看到所有 API 请求

### 查看 Console 日志
开发者工具 → Console 标签，可以看到任何错误或警告

## 生产构建

```bash
# 构建生产版本
npm run build

# 输出在 dist/ 目录

# 预览生产版本
npm run preview
```

## 下一步

1. **获取测试 ETH**：
   - 访问 Sepolia 水龙头获取测试 ETH
   - https://sepoliafaucet.com/

2. **进行测试交易**：
   - 使用配置的 ERC-20 代币合约进行转账或质押
   - Listener 会监听事件并更新数据库

3. **查看积分累积**：
   - Scheduler 每 60 秒根据余额计算积分
   - 刷新 Dashboard 查看积分变化

## 完整的本地开发流程

```bash
# Terminal 1 - 启动后端
cd /home/tianyu/loyalty-points-system-final
docker-compose up

# Terminal 2 - 启动前端
cd /home/tianyu/loyalty-points-system-final/frontend_back
npm run dev

# Terminal 3 - 查看日志（可选）
cd /home/tianyu/loyalty-points-system-final
docker-compose logs -f api listener consumer scheduler
```

## 停止服务

```bash
# 停止前端：在运行 npm run dev 的终端按 Ctrl+C

# 停止 Docker 服务
docker-compose down

# 停止并清除数据
docker-compose down -v
```

---

享受你的 PointFi Protocol！🚀
