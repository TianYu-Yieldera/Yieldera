# Monitoring System Setup Guide

## Prerequisites

- Node.js 18+
- TypeScript
- Deployed Treasury contracts on Arbitrum Sepolia

## Step 1: Get WebSocket RPC Provider API Key

### Option A: Alchemy (Recommended)

1. Visit [https://www.alchemy.com/](https://www.alchemy.com/)
2. Sign up for a free account
3. Create a new app:
   - **Chain**: Arbitrum
   - **Network**: Arbitrum Sepolia
   - **Name**: loyalty-points-monitoring
4. Copy the WebSocket URL from your dashboard
   - Format: `wss://arb-sepolia.g.alchemy.com/v2/YOUR_API_KEY`

### Option B: Infura

1. Visit [https://www.infura.io/](https://www.infura.io/)
2. Sign up for a free account
3. Create a new project
4. Enable Arbitrum network
5. Copy the WebSocket endpoint
   - Format: `wss://arbitrum-sepolia.infura.io/ws/v3/YOUR_PROJECT_ID`

### Option C: QuickNode

1. Visit [https://www.quicknode.com/](https://www.quicknode.com/)
2. Sign up for a free account
3. Create endpoint:
   - **Chain**: Arbitrum
   - **Network**: Sepolia Testnet
4. Copy the WebSocket URL

## Step 2: Configure Environment Variables

Edit `backend/.env`:

```env
# WebSocket Provider (choose one)
ARBITRUM_SEPOLIA_WS=wss://arb-sepolia.g.alchemy.com/v2/YOUR_API_KEY

# HTTP Provider (for queries)
ARBITRUM_SEPOLIA_RPC=https://sepolia-rollup.arbitrum.io/rpc

# Treasury Contracts (already deployed)
TREASURY_MARKETPLACE_ADDRESS=0x90708d3663C3BE0DF3002dC293Bb06c45b67a334
TREASURY_ASSET_FACTORY_ADDRESS=0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B
TREASURY_YIELD_DISTRIBUTOR_ADDRESS=0x0BE14D40188FCB5924c36af46630faBD76698A80
TREASURY_PRICE_ORACLE_ADDRESS=0xB478ca7F5f03f2700BfC56613bb22546D6D10681
USDC_ADDRESS=0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d

# DeFi Adapters (to be deployed)
UNISWAP_ADAPTER_ADDRESS=
AAVE_ADAPTER_ADDRESS=
COMPOUND_ADAPTER_ADDRESS=

# Alerts (optional)
SLACK_ENABLED=false
SLACK_WEBHOOK_URL=

# Performance
STATS_REPORT_INTERVAL_MS=300000
RISK_CALCULATION_INTERVAL_MS=60000
```

## Step 3: Install Dependencies

```bash
cd backend
npm install
```

## Step 4: Run Test Suite

Test the monitoring system configuration:

```bash
npx ts-node test-monitoring.ts
```

Expected output:
```
🧪 Starting Monitoring System Tests
============================================================

📋 Test 1: Environment Configuration
------------------------------------------------------------
  ✅ ARBITRUM_SEPOLIA_WS: wss://arb-sepolia.g.alchemy.com/v2/...
  ✅ All required environment variables are set

🔌 Test 2: WebSocket Connection
------------------------------------------------------------
  ✅ Connected to Arbitrum Sepolia
  📦 Current block: 12345678

...

📊 Test Summary
============================================================
Total Tests: 5
✅ Passed: 5
❌ Failed: 0
📈 Pass Rate: 100.0%

🎉 All tests passed! Monitoring system is ready.
```

## Step 5: Start Monitoring System

```bash
npm run start
```

Or with ts-node:
```bash
npx ts-node index.ts
```

## Monitoring Output

The system will display real-time events:

### Treasury Events
```
📝 Marketplace orderCreated: { orderId: '1', seller: '0x1234...', ... }
✅ Marketplace orderFilled: { orderId: '1', buyer: '0x5678...', totalPrice: '$1000.00' }
🆕 Asset created: { assetId: '1', symbol: 'T-BILL-2025', value: '$100000' }
✔️ Asset verified: { assetId: '1', verifier: '0xabcd...' }
```

### Alerts
```
⚠️ ALERT [WARNING]: LARGE_TRADE
   Large trade executed: $50,000.00

🚨 ALERT [CRITICAL]: ASSET_ISSUE
   Asset #1 status changed to: DEFAULTED
```

### Statistics (every 5 minutes)
```
📊 === Monitoring Stats ===
TreasuryMarketplace: {
  totalOrders: 10,
  totalTrades: 8,
  totalVolume: '$125,000.00',
  activeOrders: 2,
  averageTradeSize: '$15,625.00'
}
TreasuryAssetFactory: {
  totalAssetsCreated: 5,
  verifiedAssets: 4,
  totalValue: '$500,000.00',
  verificationRate: '80.0%'
}
========================
```

## Troubleshooting

### Error: "Unexpected server response: 400"
- **Cause**: Invalid or demo API key
- **Solution**: Replace with your actual API key from Alchemy/Infura

### Error: "Connection timeout"
- **Cause**: Network issues or rate limiting
- **Solution**:
  1. Check internet connection
  2. Verify API key is valid
  3. Check rate limits on your RPC provider plan

### Error: "No contract found at address"
- **Cause**: Contract not deployed or wrong network
- **Solution**: Verify contract addresses in deployment logs

### Error: "Missing environment variables"
- **Cause**: .env file not properly configured
- **Solution**: Copy .env.example to .env and fill in values

## Next Steps

1. **Deploy DeFi Adapters** (if not already deployed):
   ```bash
   npx hardhat run scripts/layer2/deploy-adapters.js --network arbitrumSepolia
   ```

2. **Update DeFi Adapter Addresses** in `.env`:
   ```env
   UNISWAP_ADAPTER_ADDRESS=0x...
   AAVE_ADAPTER_ADDRESS=0x...
   COMPOUND_ADAPTER_ADDRESS=0x...
   ```

3. **Enable Slack Notifications** (optional):
   - Create Slack webhook URL
   - Set `SLACK_ENABLED=true`
   - Add `SLACK_WEBHOOK_URL=https://hooks.slack.com/...`

4. **Set up Database** (optional, for persistent storage):
   - Install PostgreSQL
   - Create database: `createdb loyalty_monitoring`
   - Update DB credentials in `.env`

## Architecture Overview

```
MonitoringSystem
├── DeFi Adapter Listeners
│   ├── UniswapListener (swaps, slippage)
│   ├── AaveListener (supply, borrow, flash loans)
│   └── CompoundListener (supply, rates)
└── Treasury Listeners
    ├── MarketplaceListener (orders, trades)
    ├── AssetFactoryListener (assets, verification)
    └── YieldDistributorListener (yields, distributions)
```

Each listener:
- Connects via WebSocket for real-time events
- Auto-reconnects on disconnection
- Emits alerts for risk conditions
- Tracks statistics and metrics
- Reports to central monitoring system

## Performance Considerations

- **WebSocket Connection**: Persistent connection for real-time updates
- **Event Processing**: Asynchronous with 5s delay for batching
- **Stats Reporting**: Every 5 minutes (configurable)
- **Risk Calculation**: Every 1 minute (configurable)

## Security

- Never commit `.env` file to version control
- Rotate API keys regularly
- Use read-only RPC endpoints when possible
- Monitor rate limits to avoid service interruption
