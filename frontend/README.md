# PointFi Frontend

完整的 PointFi Protocol 前端，包含钱包集成、路由、Dashboard 等功能。

## 快速开始

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

## 功能

- **钱包集成**: MetaMask 钱包连接，支持 Sepolia 测试网
- **路由**: 多页面应用，包括 Landing、Dashboard 等
- **Dashboard**: 显示用户余额、积分、排行榜、徽章等数据
- **响应式设计**: 深色主题，渐变样式

## 项目结构

```
frontend_back/
├── public/
│   └── pointfi-logo-mark.svg    # Logo
├── src/
│   ├── components/
│   │   └── Header.jsx            # 导航栏
│   ├── views/
│   │   ├── Landing.jsx           # 首页
│   │   └── HomeView.jsx          # Dashboard
│   ├── web3/
│   │   ├── provider.js           # 钱包 Provider
│   │   └── WalletContext.jsx     # 钱包上下文
│   └── main.jsx                  # 入口文件
├── index.html                    # HTML 模板
├── vite.config.js                # Vite 配置
└── package.json
```

## API 集成

Dashboard 默认从 `http://localhost:8080` 读取数据：
- `/users/:addr/balance` - 用户余额
- `/users/:addr/points` - 用户积分
- `/leaderboard` - 排行榜
- `/users/:addr/badges` - 用户徽章

确保后端 API 服务正在运行。

## 构建

```bash
npm run build      # 生产构建
npm run preview    # 预览生产构建
```
