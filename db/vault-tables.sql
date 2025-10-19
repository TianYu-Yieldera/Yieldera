-- =====================================================
-- Vault Service Tables
-- =====================================================

-- Vault deposits with idempotency support
CREATE TABLE IF NOT EXISTS vault_deposits (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  mode TEXT NOT NULL CHECK (mode IN ('smart', 'manual')),
  strategy TEXT,
  status TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'failed')),
  idempotency_key TEXT UNIQUE, -- For preventing duplicate deposits
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_vault_deposits_user ON vault_deposits(user_address);
CREATE INDEX IF NOT EXISTS idx_vault_deposits_idempotency ON vault_deposits(idempotency_key);

-- Vault positions (user's staked amounts per protocol)
CREATE TABLE IF NOT EXISTS vault_positions (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  protocol TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  earned NUMERIC(78, 18) NOT NULL DEFAULT 0,
  apy NUMERIC(10, 2) DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address, protocol)
);
CREATE INDEX IF NOT EXISTS idx_vault_positions_user ON vault_positions(user_address);
CREATE INDEX IF NOT EXISTS idx_vault_positions_protocol ON vault_positions(protocol);

-- Vault earnings history
CREATE TABLE IF NOT EXISTS vault_earnings (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  protocol TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  apy NUMERIC(10, 2),
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_vault_earnings_user ON vault_earnings(user_address);

-- Vault strategies
CREATE TABLE IF NOT EXISTS vault_strategies (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  mode TEXT NOT NULL CHECK (mode IN ('smart', 'manual')),
  protocol_allocations JSONB NOT NULL, -- e.g., {"Aave V3": 40, "Compound V3": 30, ...}
  min_amount NUMERIC(78, 18),
  max_amount NUMERIC(78, 18),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- DeFi protocols
CREATE TABLE IF NOT EXISTS defi_protocols (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  protocol_type TEXT NOT NULL CHECK (protocol_type IN ('lending', 'dex', 'derivatives', 'yield', 'liquid_staking')),
  risk_level TEXT NOT NULL CHECK (risk_level IN ('low', 'medium', 'high')),
  current_apy NUMERIC(10, 2) NOT NULL DEFAULT 0,
  tvl NUMERIC(78, 18) NOT NULL DEFAULT 0,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- RWA Service Tables
-- =====================================================

-- RWA assets
CREATE TABLE IF NOT EXISTS rwa_assets (
  id SERIAL PRIMARY KEY,
  ticker TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  asset_type TEXT NOT NULL CHECK (asset_type IN ('stock', 'bond', 'commodity', 'real_estate')),
  issuer TEXT,
  current_price NUMERIC(78, 18) NOT NULL DEFAULT 0,
  price_change_24h NUMERIC(10, 4) DEFAULT 0,
  price_change_7d NUMERIC(10, 4) DEFAULT 0,
  market_cap NUMERIC(78, 18) DEFAULT 0,
  total_supply NUMERIC(78, 18) DEFAULT 0,
  circulating_supply NUMERIC(78, 18),
  contract_address TEXT,
  metadata JSONB,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_rwa_assets_ticker ON rwa_assets(ticker);
CREATE INDEX IF NOT EXISTS idx_rwa_assets_type ON rwa_assets(asset_type);

-- RWA orders
CREATE TABLE IF NOT EXISTS rwa_orders (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  asset_ticker TEXT NOT NULL,
  order_type TEXT NOT NULL CHECK (order_type IN ('buy', 'sell')),
  order_style TEXT NOT NULL CHECK (order_style IN ('market', 'limit')),
  amount NUMERIC(78, 18) NOT NULL,
  price NUMERIC(78, 18) NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'executed', 'cancelled', 'failed')),
  filled_amount NUMERIC(78, 18) DEFAULT 0,
  average_price NUMERIC(78, 18),
  created_at TIMESTAMP DEFAULT NOW(),
  executed_at TIMESTAMP,
  cancelled_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_rwa_orders_user ON rwa_orders(user_address);
CREATE INDEX IF NOT EXISTS idx_rwa_orders_status ON rwa_orders(status);

-- RWA holdings
CREATE TABLE IF NOT EXISTS rwa_holdings (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  asset_ticker TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  average_cost NUMERIC(78, 18) NOT NULL DEFAULT 0,
  current_value NUMERIC(78, 18) NOT NULL DEFAULT 0,
  pnl NUMERIC(78, 18) DEFAULT 0,
  pnl_percentage NUMERIC(10, 4) DEFAULT 0,
  last_updated TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address, asset_ticker)
);
CREATE INDEX IF NOT EXISTS idx_rwa_holdings_user ON rwa_holdings(user_address);

-- Price history
CREATE TABLE IF NOT EXISTS price_history (
  id SERIAL PRIMARY KEY,
  asset_ticker TEXT NOT NULL,
  price NUMERIC(78, 18) NOT NULL,
  high_24h NUMERIC(78, 18),
  low_24h NUMERIC(78, 18),
  volume_24h NUMERIC(78, 18),
  source TEXT DEFAULT 'internal',
  timestamp TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_price_history_ticker ON price_history(asset_ticker);
CREATE INDEX IF NOT EXISTS idx_price_history_timestamp ON price_history(timestamp);

-- APY history
CREATE TABLE IF NOT EXISTS apy_history (
  id SERIAL PRIMARY KEY,
  protocol TEXT NOT NULL,
  apy NUMERIC(10, 2) NOT NULL,
  tvl NUMERIC(78, 18),
  timestamp TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_apy_history_protocol ON apy_history(protocol);
CREATE INDEX IF NOT EXISTS idx_apy_history_timestamp ON apy_history(timestamp);

-- =====================================================
-- Transaction History (Unified)
-- =====================================================

-- Unified transaction history
CREATE TABLE IF NOT EXISTS transaction_history (
  id BIGSERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  transaction_type TEXT NOT NULL, -- vault_deposit, vault_withdraw, rwa_buy, rwa_sell, etc.
  asset TEXT, -- For RWA transactions
  amount NUMERIC(78, 18) NOT NULL,
  price NUMERIC(78, 18), -- For RWA transactions
  fee NUMERIC(78, 18) DEFAULT 0,
  status TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'failed', 'cancelled')),
  metadata JSONB, -- Additional data
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_transaction_history_user ON transaction_history(user_address);
CREATE INDEX IF NOT EXISTS idx_transaction_history_type ON transaction_history(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transaction_history_created ON transaction_history(created_at);

-- =====================================================
-- Audit Logs (For Security)
-- =====================================================

-- Audit logs for all critical operations
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGSERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  action TEXT NOT NULL, -- DEPOSIT, WITHDRAW, BUY, SELL, MINT, BURN, etc.
  service TEXT NOT NULL, -- vault, rwa, api, etc.
  amount NUMERIC(78, 18),
  status TEXT NOT NULL CHECK (status IN ('success', 'failure')),
  ip_address TEXT,
  user_agent TEXT,
  request_id TEXT, -- For tracing
  error_message TEXT,
  metadata JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_address);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_request ON audit_logs(request_id);

-- =====================================================
-- Initial Data
-- =====================================================

-- Insert default DeFi protocols
INSERT INTO defi_protocols (name, protocol_type, risk_level, current_apy, tvl, is_active) VALUES
  ('Aave V3', 'lending', 'low', 5.5, 10000000000, true),
  ('Compound V3', 'lending', 'low', 4.8, 8000000000, true),
  ('Uniswap V3', 'dex', 'medium', 12.3, 5000000000, true),
  ('Curve', 'dex', 'low', 6.2, 7000000000, true),
  ('Yearn Finance', 'yield', 'medium', 15.7, 3000000000, true),
  ('Lido', 'liquid_staking', 'low', 4.2, 15000000000, true),
  ('Rocket Pool', 'liquid_staking', 'low', 4.5, 2000000000, true),
  ('GMX', 'derivatives', 'high', 25.0, 500000000, true)
ON CONFLICT (name) DO NOTHING;

-- Insert default vault strategies
INSERT INTO vault_strategies (name, mode, protocol_allocations, min_amount, max_amount, is_active) VALUES
  ('conservative', 'smart', '{"Aave V3": 40, "Compound V3": 30, "Curve": 20, "Lido": 10}', 100, NULL, true),
  ('balanced', 'smart', '{"Aave V3": 25, "Compound V3": 20, "Uniswap V3": 20, "Yearn Finance": 20, "Curve": 15}', 100, NULL, true),
  ('aggressive', 'smart', '{"Uniswap V3": 25, "GMX": 20, "Yearn Finance": 25, "Aave V3": 15, "Rocket Pool": 15}', 500, NULL, true)
ON CONFLICT (name) DO NOTHING;

-- Insert sample RWA assets
INSERT INTO rwa_assets (ticker, name, asset_type, issuer, current_price, market_cap, total_supply, is_active, metadata) VALUES
  ('AAPL-RWA', 'Apple Inc. Tokenized Stock', 'stock', 'TokenFi Securities', 175.50, 2800000000000, 16000000000, true, '{"sector": "Technology", "exchange": "NASDAQ", "country": "USA"}'),
  ('TSLA-RWA', 'Tesla Inc. Tokenized Stock', 'stock', 'TokenFi Securities', 245.30, 770000000000, 3100000000, true, '{"sector": "Automotive", "exchange": "NASDAQ", "country": "USA"}'),
  ('GOLD-RWA', 'Tokenized Gold', 'commodity', 'MetalVault', 1950.00, 12000000000000, 6000000000, true, '{"type": "Precious Metal", "backing": "Physical Gold"}'),
  ('UST10Y-RWA', 'US Treasury 10Y Bond', 'bond', 'GovToken', 99.50, 500000000000, 5000000000, true, '{"maturity": "10 years", "coupon": "4.5%", "rating": "AAA"}')
ON CONFLICT (ticker) DO NOTHING;
