# 🚀 5分钟快速开始指南

**适合人群**: 想快速体验产品的用户

---

## 第1步: 准备钱包（1分钟）

1. **打开 MetaMask**，切换到 **Arbitrum Sepolia** 网络

   没有这个网络？添加网络设置：
   ```
   网络名称: Arbitrum Sepolia
   RPC URL: https://sepolia-rollup.arbitrum.io/rpc
   Chain ID: 421614
   ```

2. **检查余额**：
   - ETH: 需要一些（用于 gas 费）
   - USDC: 需要一些（用于购买国债）

3. **如果没有 USDC**:
   - 访问: https://faucet.circle.com/
   - 选择 Arbitrum Sepolia
   - 领取 10 USDC（免费）

---

## 第2步: 打开网站（30秒）

1. 浏览器访问: **http://localhost:5173**

2. 点击右上角 **"Connect Wallet"**

3. 选择 MetaMask，点击连接

---

## 第3步: 浏览国债（1分钟）

1. 点击菜单中的 **"Treasury Market"**

2. 你会看到 8 个国债产品：

   ```
   📊 3-Month T-Bill  - 5.25% APY
   📊 6-Month T-Bill  - 5.40% APY
   📊 2-Year T-Note   - 4.75% APY
   ... 还有 5 个
   ```

3. 点击任意产品查看详情

---

## 第4步: 购买国债（3分钟）

1. **选择产品**: 点击 "3-Month T-Bill" 的 **"Purchase"** 按钮

2. **输入金额**: 输入 `10`（代表 $10 USDC）

3. **第一笔交易 - 授权**:
   - MetaMask 弹出 → 点击 "确认"
   - 等待 5-10 秒

4. **第二笔交易 - 购买**:
   - MetaMask 再次弹出 → 点击 "确认"
   - 等待 5-10 秒

5. **完成！** 看到 "Purchase Successful!" ✅

---

## 第5步: 查看持仓（30秒）

1. 点击 **"My Holdings"** 或 "我的持仓"

2. 你会看到：
   ```
   持有: 10 TBILL-3M tokens
   价值: $10.00
   APY: 5.25%
   ```

---

## ✅ 测试完成！

恭喜！你已经完成了：
- ✅ 连接钱包
- ✅ 浏览国债市场
- ✅ 购买国债
- ✅ 查看持仓

**用时**: 约 5-7 分钟

---

## 💡 继续探索

想尝试更多功能？

- 购买其他国债产品（6个月、2年、5年等）
- 查看交易历史
- 体验 DeFi Vault 功能
- 查看 AI 风险评估

详细指南：查看 `USER_EXPERIENCE_TEST_GUIDE.md`

---

## 🆘 遇到问题？

**看不到国债产品？**
```bash
docker restart loyalty-rwa
```

**连接钱包失败？**
- 确保 MetaMask 已解锁
- 刷新页面重试

**USDC 不够？**
- 访问 https://faucet.circle.com/ 领取更多

---

**开始测试**: http://localhost:5173 🚀
