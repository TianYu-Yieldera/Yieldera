-- =============================================
-- Yieldera Vault & RWA Migration Script
-- Version: 2.0
-- Date: 2025-10-18
-- Description: Add tables for DeFi Vault and RWA Marketplace
-- =============================================

-- =============================================
-- 理财金库 (DeFi Vault) Tables
-- =============================================

-- 金库存款记录
CREATE TABLE IF NOT EXISTS vault_deposits (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  mode TEXT CHECK (mode IN ('smart', 'manual')),
  strategy TEXT DEFAULT 'balanced' CHECK (strategy IN ('conservative', 'balanced', 'aggressive')),
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'failed')),
  tx_hash TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- 金库策略配置
CREATE TABLE IF NOT EXISTS vault_strategies (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  mode TEXT CHECK (mode IN ('conservative', 'balanced', 'aggressive')),
  protocol_allocations JSONB NOT NULL, -- {"aave": 40, "compound": 30, ...}
  min_amount NUMERIC(78, 18) DEFAULT 100,
  max_amount NUMERIC(78, 18),
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- 用户金库仓位
CREATE TABLE IF NOT EXISTS vault_positions (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  protocol TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  earned NUMERIC(78, 18) DEFAULT 0,
  apy NUMERIC(5, 2),
  last_harvest TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address, protocol)
);

-- 收益记录
CREATE TABLE IF NOT EXISTS vault_earnings (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  protocol TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  apy NUMERIC(5, 2),
  tx_hash TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- DeFi 协议信息
CREATE TABLE IF NOT EXISTS defi_protocols (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  protocol_type TEXT CHECK (protocol_type IN ('lending', 'dex', 'derivatives', 'yield')),
  risk_level TEXT CHECK (risk_level IN ('low', 'medium', 'high')),
  current_apy NUMERIC(5, 2),
  tvl NUMERIC(78, 18),
  contract_address TEXT,
  chain TEXT DEFAULT 'ethereum',
  is_active BOOLEAN DEFAULT true,
  metadata JSONB,
  updated_at TIMESTAMP DEFAULT NOW()
);

-- =============================================
-- RWA 商城 (RWA Marketplace) Tables
-- =============================================

-- RWA 资产信息
CREATE TABLE IF NOT EXISTS rwa_assets (
  id SERIAL PRIMARY KEY,
  ticker TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  asset_type TEXT CHECK (asset_type IN ('stock', 'bond', 'commodity', 'realestate', 'other')),
  issuer TEXT NOT NULL, -- 'Backed Finance', 'Ondo', 'Paxos', etc.
  contract_address TEXT,
  current_price NUMERIC(78, 18),
  price_change_24h NUMERIC(5, 2),
  price_change_7d NUMERIC(5, 2),
  market_cap NUMERIC(78, 18),
  total_supply NUMERIC(78, 18),
  circulating_supply NUMERIC(78, 18),
  is_active BOOLEAN DEFAULT true,
  metadata JSONB, -- 存储额外信息如图标URL、描述等
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- RWA 订单
CREATE TABLE IF NOT EXISTS rwa_orders (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  asset_ticker TEXT NOT NULL,
  order_type TEXT CHECK (order_type IN ('buy', 'sell')),
  order_style TEXT CHECK (order_style IN ('market', 'limit', 'stop')),
  amount NUMERIC(78, 18) NOT NULL, -- 买入时为USDC金额，卖出时为代币数量
  price NUMERIC(78, 18), -- 限价单价格
  stop_price NUMERIC(78, 18), -- 止损单价格
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'partial', 'executed', 'cancelled', 'expired')),
  filled_amount NUMERIC(78, 18) DEFAULT 0,
  average_price NUMERIC(78, 18),
  fee_amount NUMERIC(78, 18) DEFAULT 0,
  tx_hash TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  executed_at TIMESTAMP,
  cancelled_at TIMESTAMP,
  expires_at TIMESTAMP
);

-- 用户 RWA 持仓
CREATE TABLE IF NOT EXISTS rwa_holdings (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  asset_ticker TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  average_cost NUMERIC(78, 18),
  current_value NUMERIC(78, 18),
  pnl NUMERIC(78, 18), -- profit and loss
  pnl_percentage NUMERIC(5, 2), -- profit and loss percentage
  last_updated TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address, asset_ticker)
);

-- 价格历史
CREATE TABLE IF NOT EXISTS price_history (
  id SERIAL PRIMARY KEY,
  asset_ticker TEXT NOT NULL,
  price NUMERIC(78, 18) NOT NULL,
  high_24h NUMERIC(78, 18),
  low_24h NUMERIC(78, 18),
  volume_24h NUMERIC(78, 18),
  market_cap NUMERIC(78, 18),
  source TEXT, -- 'chainlink', 'band', 'api3', 'coingecko'
  timestamp TIMESTAMP DEFAULT NOW()
);

-- APY 历史（用于 DeFi 协议）
CREATE TABLE IF NOT EXISTS apy_history (
  id SERIAL PRIMARY KEY,
  protocol TEXT NOT NULL,
  apy NUMERIC(5, 2) NOT NULL,
  tvl NUMERIC(78, 18),
  utilization_rate NUMERIC(5, 2),
  timestamp TIMESTAMP DEFAULT NOW()
);

-- 用户组合总览
CREATE TABLE IF NOT EXISTS user_portfolios (
  id SERIAL PRIMARY KEY,
  user_address TEXT UNIQUE NOT NULL,
  total_deposited NUMERIC(78, 18) DEFAULT 0,
  total_earned NUMERIC(78, 18) DEFAULT 0,
  vault_value NUMERIC(78, 18) DEFAULT 0,
  rwa_value NUMERIC(78, 18) DEFAULT 0,
  total_value NUMERIC(78, 18) DEFAULT 0,
  roi_percentage NUMERIC(5, 2),
  risk_score INTEGER CHECK (risk_score >= 0 AND risk_score <= 100),
  last_rebalance TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- 交易历史（审计日志）
CREATE TABLE IF NOT EXISTS transaction_history (
  id BIGSERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  transaction_type TEXT NOT NULL, -- 'vault_deposit', 'vault_withdraw', 'rwa_buy', 'rwa_sell', etc.
  asset TEXT,
  amount NUMERIC(78, 18),
  price NUMERIC(78, 18),
  fee NUMERIC(78, 18),
  status TEXT,
  tx_hash TEXT,
  metadata JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

-- =============================================
-- Indexes for Performance
-- =============================================

-- Vault indexes
CREATE INDEX IF NOT EXISTS idx_vault_deposits_user ON vault_deposits(user_address);
CREATE INDEX IF NOT EXISTS idx_vault_deposits_status ON vault_deposits(status);
CREATE INDEX IF NOT EXISTS idx_vault_positions_user ON vault_positions(user_address);
CREATE INDEX IF NOT EXISTS idx_vault_positions_protocol ON vault_positions(protocol);
CREATE INDEX IF NOT EXISTS idx_vault_earnings_user ON vault_earnings(user_address);
CREATE INDEX IF NOT EXISTS idx_vault_earnings_created ON vault_earnings(created_at);

-- RWA indexes
CREATE INDEX IF NOT EXISTS idx_rwa_orders_user ON rwa_orders(user_address);
CREATE INDEX IF NOT EXISTS idx_rwa_orders_status ON rwa_orders(status);
CREATE INDEX IF NOT EXISTS idx_rwa_orders_ticker ON rwa_orders(asset_ticker);
CREATE INDEX IF NOT EXISTS idx_rwa_holdings_user ON rwa_holdings(user_address);
CREATE INDEX IF NOT EXISTS idx_rwa_holdings_ticker ON rwa_holdings(asset_ticker);

-- Price history indexes
CREATE INDEX IF NOT EXISTS idx_price_history_ticker_time ON price_history(asset_ticker, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_apy_history_protocol_time ON apy_history(protocol, timestamp DESC);

-- Transaction history index
CREATE INDEX IF NOT EXISTS idx_transaction_history_user ON transaction_history(user_address);
CREATE INDEX IF NOT EXISTS idx_transaction_history_created ON transaction_history(created_at DESC);

-- =============================================
-- Initial Data Seeding
-- =============================================

-- Insert DeFi protocols
INSERT INTO defi_protocols (name, protocol_type, risk_level, current_apy, tvl, chain) VALUES
('Aave V3', 'lending', 'low', 3.52, 5000000000, 'ethereum'),
('Compound V3', 'lending', 'low', 4.18, 3000000000, 'ethereum'),
('Uniswap V3', 'dex', 'medium', 12.85, 7000000000, 'ethereum'),
('GMX', 'derivatives', 'high', 22.30, 500000000, 'arbitrum'),
('Curve', 'dex', 'low', 5.20, 4000000000, 'ethereum'),
('Yearn Finance', 'yield', 'medium', 10.50, 1000000000, 'ethereum'),
('Rocket Pool', 'staking', 'low', 4.80, 2000000000, 'ethereum'),
('Lido', 'staking', 'low', 4.50, 8000000000, 'ethereum')
ON CONFLICT (name) DO UPDATE SET
  current_apy = EXCLUDED.current_apy,
  tvl = EXCLUDED.tvl,
  updated_at = NOW();

-- Insert vault strategies
INSERT INTO vault_strategies (name, mode, protocol_allocations) VALUES
('Conservative Strategy', 'conservative', '{"Aave V3": 40, "Compound V3": 30, "Curve": 20, "Lido": 10}'),
('Balanced Strategy', 'balanced', '{"Aave V3": 25, "Compound V3": 20, "Uniswap V3": 20, "Yearn Finance": 20, "Curve": 15}'),
('Aggressive Strategy', 'aggressive', '{"Uniswap V3": 25, "GMX": 20, "Yearn Finance": 25, "Aave V3": 15, "Rocket Pool": 15}');

-- Insert RWA assets
INSERT INTO rwa_assets (ticker, name, asset_type, issuer, current_price, price_change_24h, market_cap, total_supply, metadata) VALUES
-- Stocks
('bAAPL', 'Apple Inc.', 'stock', 'Backed Finance', 178.52, 2.34, 2800000000000, 15700000000, '{"logo": "🍎", "description": "Technology giant, iPhone maker", "isin": "US0378331005"}'),
('bTSLA', 'Tesla Inc.', 'stock', 'Backed Finance', 242.18, 5.67, 770000000000, 3180000000, '{"logo": "⚡", "description": "Electric vehicle leader", "isin": "US88160R1014"}'),
('bNVDA', 'NVIDIA Corp.', 'stock', 'Backed Finance', 485.93, 8.92, 1200000000000, 2470000000, '{"logo": "🎮", "description": "AI chip leader", "isin": "US67066G1040"}'),
('bGOOGL', 'Alphabet Inc.', 'stock', 'Backed Finance', 139.47, 1.23, 1740000000000, 12480000000, '{"logo": "🔍", "description": "Search engine giant", "isin": "US02079K3059"}'),
('bAMZN', 'Amazon.com Inc.', 'stock', 'Backed Finance', 178.25, 3.45, 1850000000000, 10380000000, '{"logo": "📦", "description": "E-commerce and cloud giant", "isin": "US0231351067"}'),
('bMSFT', 'Microsoft Corp.', 'stock', 'Backed Finance', 398.76, 2.11, 2960000000000, 7430000000, '{"logo": "💻", "description": "Software and cloud services giant", "isin": "US5949181045"}'),

-- Bonds
('OUSG', 'Ondo Short-Term US Government Bond', 'bond', 'Ondo Finance', 105.20, 0.12, 500000000, 4750000, '{"logo": "🏛️", "description": "Short-term US Treasury token", "apy": "4.5%"}'),
('USDY', 'Ondo US Dollar Yield', 'bond', 'Ondo Finance', 100.85, 0.08, 250000000, 2480000, '{"logo": "💵", "description": "US dollar yield token", "apy": "5.2%"}'),
('USTB', 'Matrixdock Short-term Treasury Bill', 'bond', 'Matrixdock', 100.50, 0.05, 150000000, 1490000, '{"logo": "📊", "description": "Tokenized T-Bills", "apy": "4.8%"}'),

-- Commodities
('PAXG', 'Paxos Gold', 'commodity', 'Paxos', 2042.50, 1.05, 600000000, 293706, '{"logo": "🥇", "description": "1 PAXG = 1 troy ounce of gold", "backing": "Physical gold in Brinks vaults"}'),
('XAUT', 'Tether Gold', 'commodity', 'Tether', 2041.80, 1.02, 500000000, 244824, '{"logo": "🪙", "description": "1 XAUT = 1 troy ounce of gold", "backing": "Physical gold in Swiss vaults"}'),
('SLVT', 'Silvertoken', 'commodity', 'Silvertoken', 24.15, 0.85, 50000000, 2070000, '{"logo": "🥈", "description": "Tokenized silver", "backing": "Physical silver"}')
ON CONFLICT (ticker) DO UPDATE SET
  current_price = EXCLUDED.current_price,
  price_change_24h = EXCLUDED.price_change_24h,
  market_cap = EXCLUDED.market_cap,
  updated_at = NOW();

-- =============================================
-- Functions and Triggers
-- =============================================

-- Function to update user portfolio summary
CREATE OR REPLACE FUNCTION update_user_portfolio()
RETURNS TRIGGER AS $$
BEGIN
  -- Update or insert portfolio summary
  INSERT INTO user_portfolios (
    user_address,
    vault_value,
    rwa_value,
    total_value,
    updated_at
  )
  VALUES (
    NEW.user_address,
    COALESCE((SELECT SUM(amount + earned) FROM vault_positions WHERE user_address = NEW.user_address), 0),
    COALESCE((SELECT SUM(current_value) FROM rwa_holdings WHERE user_address = NEW.user_address), 0),
    COALESCE((SELECT SUM(amount + earned) FROM vault_positions WHERE user_address = NEW.user_address), 0) +
    COALESCE((SELECT SUM(current_value) FROM rwa_holdings WHERE user_address = NEW.user_address), 0),
    NOW()
  )
  ON CONFLICT (user_address) DO UPDATE SET
    vault_value = EXCLUDED.vault_value,
    rwa_value = EXCLUDED.rwa_value,
    total_value = EXCLUDED.total_value,
    updated_at = NOW();

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to update portfolio
CREATE TRIGGER update_portfolio_on_vault_change
AFTER INSERT OR UPDATE OR DELETE ON vault_positions
FOR EACH ROW EXECUTE FUNCTION update_user_portfolio();

CREATE TRIGGER update_portfolio_on_rwa_change
AFTER INSERT OR UPDATE OR DELETE ON rwa_holdings
FOR EACH ROW EXECUTE FUNCTION update_user_portfolio();

-- Function to calculate and update PnL
CREATE OR REPLACE FUNCTION update_rwa_pnl()
RETURNS TRIGGER AS $$
BEGIN
  NEW.current_value := NEW.amount * (SELECT current_price FROM rwa_assets WHERE ticker = NEW.asset_ticker);
  NEW.pnl := NEW.current_value - (NEW.amount * NEW.average_cost);
  IF NEW.average_cost > 0 THEN
    NEW.pnl_percentage := ((NEW.current_value / (NEW.amount * NEW.average_cost)) - 1) * 100;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update PnL
CREATE TRIGGER calculate_pnl
BEFORE INSERT OR UPDATE ON rwa_holdings
FOR EACH ROW EXECUTE FUNCTION update_rwa_pnl();

-- =============================================
-- Grants (adjust based on your user setup)
-- =============================================

-- Grant permissions to the application user
-- GRANT ALL ON ALL TABLES IN SCHEMA public TO loyalty_user;
-- GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO loyalty_user;

-- =============================================
-- Migration Complete
-- =============================================