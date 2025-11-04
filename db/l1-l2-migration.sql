-- =====================================================
-- L1/L2 Migration - Extend Existing Schema
-- =====================================================

-- Add L1/L2 distinction to balance_events
ALTER TABLE balance_events
ADD COLUMN IF NOT EXISTS layer TEXT CHECK (layer IN ('L1', 'L2')) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS token TEXT DEFAULT NULL,  -- For L1 collateral (USDC, USDT, DAI)
ADD COLUMN IF NOT EXISTS contract_address TEXT DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_balance_events_layer ON balance_events(layer);

-- =====================================================
-- L1-specific Tables
-- =====================================================

-- L1 collateral deposits
CREATE TABLE IF NOT EXISTS l1_collateral_deposits (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  token TEXT NOT NULL CHECK (token IN ('USDC', 'USDT', 'DAI')),
  amount NUMERIC(78, 18) NOT NULL,
  tx_hash TEXT UNIQUE NOT NULL,
  block_number BIGINT NOT NULL,
  confirmed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_l1_deposits_user ON l1_collateral_deposits(user_address);
CREATE INDEX IF NOT EXISTS idx_l1_deposits_token ON l1_collateral_deposits(token);

-- L1 collateral balances (per token)
CREATE TABLE IF NOT EXISTS l1_collateral_balances (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  usdc_balance NUMERIC(78, 18) NOT NULL DEFAULT 0,
  usdt_balance NUMERIC(78, 18) NOT NULL DEFAULT 0,
  dai_balance NUMERIC(78, 18) NOT NULL DEFAULT 0,
  total_usd_value NUMERIC(78, 18) NOT NULL DEFAULT 0,
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address)
);
CREATE INDEX IF NOT EXISTS idx_l1_balances_user ON l1_collateral_balances(user_address);

-- L1 state snapshots (from L2StateAggregator)
CREATE TABLE IF NOT EXISTS l1_state_snapshots (
  id SERIAL PRIMARY KEY,
  merkle_root TEXT NOT NULL,
  total_deposited NUMERIC(78, 18) NOT NULL,
  user_count INT NOT NULL,
  l2_block_number BIGINT NOT NULL,
  l1_tx_hash TEXT UNIQUE NOT NULL,
  l1_block_number BIGINT NOT NULL,
  submitted_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_l1_snapshots_l2block ON l1_state_snapshots(l2_block_number);

-- =====================================================
-- L2-specific Tables
-- =====================================================

-- L2 vault positions (NEW - replacing old vault_positions)
CREATE TABLE IF NOT EXISTS l2_vault_positions (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  deposited NUMERIC(78, 18) NOT NULL DEFAULT 0,        -- Total deposited to IntegratedVault
  shares NUMERIC(78, 18) NOT NULL DEFAULT 0,            -- Vault share tokens
  current_value NUMERIC(78, 18) NOT NULL DEFAULT 0,     -- Current USD value
  yield_earned NUMERIC(78, 18) NOT NULL DEFAULT 0,      -- Lifetime yield
  last_updated TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address)
);
CREATE INDEX IF NOT EXISTS idx_l2_positions_user ON l2_vault_positions(user_address);

-- L2 vault strategy allocations (current state)
CREATE TABLE IF NOT EXISTS l2_strategy_allocations (
  id SERIAL PRIMARY KEY,
  protocol TEXT NOT NULL, -- 'Aave', 'Compound', 'Uniswap'
  allocated_amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  percentage NUMERIC(5, 2) NOT NULL DEFAULT 0,
  apy NUMERIC(10, 4) NOT NULL DEFAULT 0,
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(protocol)
);

-- L2 RWA assets (NEW - on-chain tracking)
CREATE TABLE IF NOT EXISTS l2_rwa_assets (
  id SERIAL PRIMARY KEY,
  asset_id BIGINT UNIQUE NOT NULL,          -- On-chain asset ID from RWAAssetFactory
  asset_type TEXT NOT NULL,                  -- RealEstate, Bonds, Equity, etc.
  total_value NUMERIC(78, 18) NOT NULL,
  token_address TEXT NOT NULL,               -- FractionalRWAToken address
  total_tokens NUMERIC(78, 18) NOT NULL,
  status TEXT NOT NULL,                      -- Pending, Active, Matured, Defaulted
  yield_rate NUMERIC(10, 4) DEFAULT 0,       -- APY
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_l2_rwa_assets_status ON l2_rwa_assets(status);

-- L2 RWA token holdings
CREATE TABLE IF NOT EXISTS l2_rwa_holdings (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  asset_id BIGINT NOT NULL,
  token_amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  average_price NUMERIC(78, 18) NOT NULL DEFAULT 0,
  current_value NUMERIC(78, 18) NOT NULL DEFAULT 0,
  yield_claimed NUMERIC(78, 18) NOT NULL DEFAULT 0,
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_address, asset_id)
);
CREATE INDEX IF NOT EXISTS idx_l2_rwa_holdings_user ON l2_rwa_holdings(user_address);

-- L2 RWA marketplace listings
CREATE TABLE IF NOT EXISTS l2_rwa_listings (
  id SERIAL PRIMARY KEY,
  listing_id BIGINT UNIQUE NOT NULL,         -- On-chain listing ID
  asset_id BIGINT NOT NULL,
  seller TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  price_per_token NUMERIC(78, 18) NOT NULL,
  listing_type TEXT NOT NULL,                -- Fixed, Auction
  status TEXT NOT NULL,                      -- Active, PartiallyFilled, Filled, Cancelled
  filled_amount NUMERIC(78, 18) DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_l2_listings_asset ON l2_rwa_listings(asset_id);
CREATE INDEX IF NOT EXISTS idx_l2_listings_status ON l2_rwa_listings(status);

-- L2 RWA governance proposals
CREATE TABLE IF NOT EXISTS l2_rwa_proposals (
  id SERIAL PRIMARY KEY,
  proposal_id BIGINT UNIQUE NOT NULL,
  asset_id BIGINT NOT NULL,
  proposer TEXT NOT NULL,
  proposal_type TEXT NOT NULL,               -- ParameterChange, AssetSale, etc.
  description TEXT,
  votes_for NUMERIC(78, 18) DEFAULT 0,
  votes_against NUMERIC(78, 18) DEFAULT 0,
  status TEXT NOT NULL,                      -- Active, Passed, Rejected, Executed, Cancelled
  voting_ends_at TIMESTAMP,
  executed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_l2_proposals_asset ON l2_rwa_proposals(asset_id);
CREATE INDEX IF NOT EXISTS idx_l2_proposals_status ON l2_rwa_proposals(status);

-- =====================================================
-- Bridge Tables
-- =====================================================

-- Bridge messages (L1 ↔ L2)
CREATE TABLE IF NOT EXISTS bridge_messages (
  id SERIAL PRIMARY KEY,
  message_hash TEXT UNIQUE NOT NULL,
  direction TEXT NOT NULL CHECK (direction IN ('L1_TO_L2', 'L2_TO_L1')),
  user_address TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('initiated', 'pending', 'confirmed', 'failed')),

  l1_tx_hash TEXT,
  l2_tx_hash TEXT,
  l1_block_number BIGINT,
  l2_block_number BIGINT,

  retry_count INT DEFAULT 0,
  error_message TEXT,

  initiated_at TIMESTAMP DEFAULT NOW(),
  confirmed_at TIMESTAMP,

  -- For L2→L1 withdrawals (needs proof)
  merkle_proof JSONB,
  proof_submitted BOOLEAN DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_bridge_user ON bridge_messages(user_address);
CREATE INDEX IF NOT EXISTS idx_bridge_status ON bridge_messages(status);
CREATE INDEX IF NOT EXISTS idx_bridge_direction ON bridge_messages(direction);

-- =====================================================
-- Views for Easy Querying
-- =====================================================

-- User's total balance across L1 and L2
CREATE OR REPLACE VIEW user_total_balances AS
SELECT
  COALESCE(l1.user_address, l2.user_address) AS user_address,
  COALESCE(l1.total_usd_value, 0) AS l1_collateral,
  COALESCE(l2.current_value, 0) AS l2_vault_value,
  COALESCE(l1.total_usd_value, 0) + COALESCE(l2.current_value, 0) AS total_value
FROM l1_collateral_balances l1
FULL OUTER JOIN l2_vault_positions l2 ON l1.user_address = l2.user_address;

-- System-wide statistics
CREATE OR REPLACE VIEW system_stats AS
SELECT
  (SELECT SUM(total_usd_value) FROM l1_collateral_balances) AS total_l1_tvl,
  (SELECT SUM(current_value) FROM l2_vault_positions) AS total_l2_vault_tvl,
  (SELECT SUM(current_value) FROM l2_rwa_holdings) AS total_l2_rwa_tvl,
  (SELECT COUNT(DISTINCT user_address) FROM l1_collateral_balances) AS l1_users,
  (SELECT COUNT(DISTINCT user_address) FROM l2_vault_positions) AS l2_vault_users,
  (SELECT COUNT(DISTINCT user_address) FROM l2_rwa_holdings) AS l2_rwa_users,
  (SELECT COUNT(*) FROM bridge_messages WHERE status = 'pending') AS pending_bridge_txs;
