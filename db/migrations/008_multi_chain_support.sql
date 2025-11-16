-- ============================================================================
-- Migration 008: Multi-Chain Support for Dual-Chain Architecture
-- ============================================================================
-- Purpose: Add chain_id support to enable Arbitrum + Base dual-chain system
-- Architecture: Chain isolation (no cross-chain bridge needed)
--   - Arbitrum (421614 testnet, 42161 mainnet): Aggressive DeFi (GMX, Aave, etc.)
--   - Base (84532 testnet, 8453 mainnet): Conservative RWA (US Treasury bonds)
-- ============================================================================

BEGIN;

-- ============================================================================
-- STEP 1: Add chain_id to existing tables
-- ============================================================================

-- Vault positions (DeFi positions on both chains)
-- Existing fields: id, user_address, protocol, amount, earned, apy, last_harvest, created_at, updated_at
ALTER TABLE vault_positions
ADD COLUMN IF NOT EXISTS chain_id BIGINT NOT NULL DEFAULT 421614;

COMMENT ON COLUMN vault_positions.chain_id IS 'Chain ID: 421614=Arbitrum Sepolia, 84532=Base Sepolia, 42161=Arbitrum One, 8453=Base';

-- Add additional fields for risk management if not exist
ALTER TABLE vault_positions
ADD COLUMN IF NOT EXISTS collateral_value_usd DECIMAL(30, 2),
ADD COLUMN IF NOT EXISTS debt_value_usd DECIMAL(30, 2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS health_factor DECIMAL(10, 4),
ADD COLUMN IF NOT EXISTS active BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS position_id TEXT;

-- Treasury holdings (will be ONLY on Base after migration)
-- Existing fields: user_id, bond_type, token_amount, principal_usd, purchase_date, last_yield_date, total_yield_earned, compounding_enabled, created_at, updated_at
ALTER TABLE treasury_holdings
ADD COLUMN IF NOT EXISTS chain_id BIGINT;

COMMENT ON COLUMN treasury_holdings.chain_id IS 'Chain ID for treasury (Base only): 84532=Base Sepolia, 8453=Base Mainnet';

-- Update existing treasury data to Base Sepolia
UPDATE treasury_holdings
SET chain_id = 84532
WHERE chain_id IS NULL;

ALTER TABLE treasury_holdings
ALTER COLUMN chain_id SET NOT NULL,
ALTER COLUMN chain_id SET DEFAULT 84532;

-- Add user_address alias for consistency (treasury_holdings uses user_id)
ALTER TABLE treasury_holdings
ADD COLUMN IF NOT EXISTS asset_type TEXT DEFAULT 'US_TREASURY';

-- Price history (market data from all chains)
ALTER TABLE price_history
ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 421614;

-- Liquidation events (track on both chains)
ALTER TABLE liquidation_events
ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 421614;

-- User risk profiles (aggregate across chains)
ALTER TABLE user_risk_profiles
ADD COLUMN chain_id BIGINT;

COMMENT ON COLUMN user_risk_profiles.chain_id IS 'NULL means cross-chain aggregate, specific value means chain-specific';

-- High risk positions monitoring
ALTER TABLE high_risk_positions
ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 421614;

-- Risk assessments
ALTER TABLE risk_assessments
ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 421614;

-- ML predictions
ALTER TABLE ml_predictions
ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 421614;

-- Market snapshots
ALTER TABLE market_snapshots
ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 421614;

-- Monitoring alerts
ALTER TABLE monitoring_alerts
ADD COLUMN chain_id BIGINT;

COMMENT ON COLUMN monitoring_alerts.chain_id IS 'NULL means cross-chain alert, specific value means chain-specific';


-- ============================================================================
-- STEP 2: Create GMX performance comparison table (for Coinbase contribution)
-- ============================================================================

CREATE TABLE gmx_performance_comparison (
    id SERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Transaction type
    action VARCHAR(50) NOT NULL,  -- 'open_long', 'open_short', 'close_position', 'increase', 'decrease'

    -- Performance metrics
    gas_used BIGINT NOT NULL,
    gas_price_gwei DECIMAL(20, 9) NOT NULL,
    total_cost_usd DECIMAL(20, 2) NOT NULL,
    execution_time_ms INTEGER,

    -- Trading metrics
    position_size_usd DECIMAL(30, 2),
    leverage DECIMAL(5, 2),
    slippage_bps INTEGER,  -- basis points
    price_impact_bps INTEGER,

    -- Result
    success BOOLEAN NOT NULL DEFAULT true,
    error_message TEXT,

    -- Market conditions at time of tx
    eth_price_usd DECIMAL(20, 2),
    btc_price_usd DECIMAL(20, 2),
    gas_network_congestion VARCHAR(20),  -- 'low', 'medium', 'high'

    CONSTRAINT chk_chain_id CHECK (chain_id IN (421614, 42161, 84532, 8453)),
    CONSTRAINT chk_action CHECK (action IN ('open_long', 'open_short', 'close_position', 'increase_collateral', 'decrease_collateral', 'emergency_hedge'))
);

CREATE INDEX idx_gmx_perf_chain_time ON gmx_performance_comparison(chain_id, timestamp DESC);
CREATE INDEX idx_gmx_perf_user ON gmx_performance_comparison(user_address, chain_id);
CREATE INDEX idx_gmx_perf_action ON gmx_performance_comparison(action, chain_id);

COMMENT ON TABLE gmx_performance_comparison IS 'GMX V2 performance comparison: Arbitrum vs Base. Data contributed to Coinbase Base team.';


-- ============================================================================
-- STEP 3: Create cross-chain aggregation views
-- ============================================================================

-- View 1: Cross-chain user portfolio
CREATE OR REPLACE VIEW cross_chain_user_portfolio AS
SELECT
    user_address,
    chain_id,
    COUNT(*) as position_count,
    SUM(COALESCE(collateral_value_usd, amount::DECIMAL(30,2))) as total_collateral_usd,
    SUM(COALESCE(debt_value_usd, 0)) as total_debt_usd,
    AVG(health_factor) as avg_health_factor,
    MIN(health_factor) as min_health_factor,
    SUM(CASE WHEN health_factor IS NOT NULL AND health_factor < 1.2 THEN 1 ELSE 0 END) as high_risk_position_count
FROM vault_positions
WHERE COALESCE(active, true) = true
GROUP BY user_address, chain_id;

COMMENT ON VIEW cross_chain_user_portfolio IS 'User portfolio breakdown by chain';

-- View 2: Cross-chain total value
CREATE OR REPLACE VIEW cross_chain_total_value AS
SELECT
    user_address,
    SUM(CASE WHEN chain_id IN (421614, 42161) THEN total_collateral_usd ELSE 0 END) as arbitrum_value_usd,
    SUM(CASE WHEN chain_id IN (84532, 8453) THEN total_collateral_usd ELSE 0 END) as base_value_usd,
    SUM(total_collateral_usd) as total_value_usd,
    ROUND(100.0 * SUM(CASE WHEN chain_id IN (421614, 42161) THEN total_collateral_usd ELSE 0 END) / NULLIF(SUM(total_collateral_usd), 0), 2) as arbitrum_allocation_pct,
    ROUND(100.0 * SUM(CASE WHEN chain_id IN (84532, 8453) THEN total_collateral_usd ELSE 0 END) / NULLIF(SUM(total_collateral_usd), 0), 2) as base_allocation_pct
FROM cross_chain_user_portfolio
GROUP BY user_address;

COMMENT ON VIEW cross_chain_total_value IS 'User total value and allocation across Arbitrum and Base';

-- View 3: GMX performance comparison summary
CREATE OR REPLACE VIEW gmx_performance_summary AS
SELECT
    chain_id,
    CASE
        WHEN chain_id IN (421614, 42161) THEN 'Arbitrum'
        WHEN chain_id IN (84532, 8453) THEN 'Base'
    END as chain_name,
    COUNT(*) as total_transactions,
    SUM(CASE WHEN success = true THEN 1 ELSE 0 END) as successful_transactions,
    ROUND(100.0 * SUM(CASE WHEN success = true THEN 1 ELSE 0 END) / COUNT(*), 2) as success_rate_pct,
    AVG(gas_used) as avg_gas_used,
    AVG(gas_price_gwei) as avg_gas_price_gwei,
    AVG(total_cost_usd) as avg_cost_usd,
    AVG(execution_time_ms) as avg_execution_time_ms,
    AVG(slippage_bps) as avg_slippage_bps,
    MIN(timestamp) as first_transaction,
    MAX(timestamp) as last_transaction
FROM gmx_performance_comparison
GROUP BY chain_id;

COMMENT ON VIEW gmx_performance_summary IS 'GMX performance comparison: Arbitrum vs Base aggregated metrics';


-- ============================================================================
-- STEP 4: Create indexes for multi-chain queries
-- ============================================================================

-- Vault positions
CREATE INDEX idx_vault_positions_chain_user ON vault_positions(chain_id, user_address);
CREATE INDEX idx_vault_positions_chain_protocol ON vault_positions(chain_id, protocol);
CREATE INDEX idx_vault_positions_chain_active ON vault_positions(chain_id, active) WHERE active = true;

-- Treasury holdings
CREATE INDEX idx_treasury_chain_user ON treasury_holdings(chain_id, user_address);
CREATE INDEX idx_treasury_chain_asset ON treasury_holdings(chain_id, asset_type);

-- Price history
CREATE INDEX idx_price_history_chain_asset ON price_history(chain_id, asset_symbol, timestamp DESC);

-- Liquidation events
CREATE INDEX idx_liquidation_chain_time ON liquidation_events(chain_id, timestamp DESC);

-- Risk assessments
CREATE INDEX idx_risk_assessment_chain_user ON risk_assessments(chain_id, user_address, timestamp DESC);

-- High risk positions
CREATE INDEX idx_high_risk_chain ON high_risk_positions(chain_id, risk_score DESC);


-- ============================================================================
-- STEP 5: Update primary keys to include chain_id (where appropriate)
-- ============================================================================

-- Note: Some tables may need to drop and recreate constraints
-- Only do this for tables where chain_id is part of natural key

-- Example for vault_positions (if needed - be careful with existing data)
-- ALTER TABLE vault_positions DROP CONSTRAINT IF EXISTS vault_positions_pkey;
-- ALTER TABLE vault_positions ADD PRIMARY KEY (chain_id, user_address, protocol, position_id);


-- ============================================================================
-- STEP 6: Create function to get user's cross-chain risk profile
-- ============================================================================

CREATE OR REPLACE FUNCTION get_cross_chain_risk_profile(p_user_address VARCHAR(42))
RETURNS TABLE (
    user_address VARCHAR(42),
    arbitrum_risk_score DECIMAL(5,2),
    base_risk_score DECIMAL(5,2),
    total_risk_score DECIMAL(5,2),
    arbitrum_value_usd DECIMAL(30,2),
    base_value_usd DECIMAL(30,2),
    diversification_score DECIMAL(5,2)
) AS $$
BEGIN
    RETURN QUERY
    WITH chain_risks AS (
        SELECT
            vp.user_address,
            vp.chain_id,
            AVG(ra.risk_score) as avg_risk_score,
            SUM(vp.collateral_value_usd) as total_value
        FROM vault_positions vp
        LEFT JOIN risk_assessments ra ON ra.user_address = vp.user_address AND ra.chain_id = vp.chain_id
        WHERE vp.user_address = p_user_address
        AND vp.active = true
        GROUP BY vp.user_address, vp.chain_id
    )
    SELECT
        p_user_address,
        MAX(CASE WHEN chain_id IN (421614, 42161) THEN avg_risk_score ELSE 0 END) as arbitrum_risk_score,
        MAX(CASE WHEN chain_id IN (84532, 8453) THEN avg_risk_score ELSE 0 END) as base_risk_score,
        AVG(avg_risk_score) as total_risk_score,
        SUM(CASE WHEN chain_id IN (421614, 42161) THEN total_value ELSE 0 END) as arbitrum_value_usd,
        SUM(CASE WHEN chain_id IN (84532, 8453) THEN total_value ELSE 0 END) as base_value_usd,
        -- Diversification score: 100 = perfectly balanced, 0 = all on one chain
        (100.0 - ABS(
            SUM(CASE WHEN chain_id IN (421614, 42161) THEN total_value ELSE 0 END) -
            SUM(CASE WHEN chain_id IN (84532, 8453) THEN total_value ELSE 0 END)
        ) / NULLIF(SUM(total_value), 0) * 100.0)::DECIMAL(5,2) as diversification_score
    FROM chain_risks;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_cross_chain_risk_profile IS 'Get user risk profile across all chains';


-- ============================================================================
-- STEP 7: Create liquidation prediction table (for AI early warning)
-- ============================================================================

CREATE TABLE liquidation_predictions (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    chain_id BIGINT NOT NULL,
    position_id VARCHAR(100) NOT NULL,

    -- Prediction details
    prediction_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    predicted_liquidation_time TIMESTAMPTZ,
    liquidation_probability_1h DECIMAL(5,4),
    liquidation_probability_6h DECIMAL(5,4),
    liquidation_probability_24h DECIMAL(5,4),
    liquidation_probability_48h DECIMAL(5,4),

    -- Current state
    current_health_factor DECIMAL(10,4),
    current_collateral_value_usd DECIMAL(30,2),
    current_debt_value_usd DECIMAL(30,2),

    -- Risk factors
    price_volatility_score DECIMAL(5,2),
    leverage_score DECIMAL(5,2),
    market_risk_score DECIMAL(5,2),

    -- AI model info
    model_version VARCHAR(50),
    confidence_score DECIMAL(5,4),

    -- Alert status
    alert_sent BOOLEAN DEFAULT false,
    alert_sent_at TIMESTAMPTZ,

    CONSTRAINT chk_liq_chain_id CHECK (chain_id IN (421614, 42161, 84532, 8453))
);

CREATE INDEX idx_liq_pred_user_chain ON liquidation_predictions(user_address, chain_id, prediction_timestamp DESC);
CREATE INDEX idx_liq_pred_high_risk ON liquidation_predictions(chain_id, liquidation_probability_24h DESC)
    WHERE liquidation_probability_24h > 0.15;
CREATE INDEX idx_liq_pred_alert_pending ON liquidation_predictions(user_address, chain_id)
    WHERE alert_sent = false AND liquidation_probability_6h > 0.1;

COMMENT ON TABLE liquidation_predictions IS 'AI-powered liquidation predictions for early warning system';


-- ============================================================================
-- STEP 8: Create Base ecosystem tracking tables
-- ============================================================================

-- Aerodrome DEX swaps
CREATE TABLE aerodrome_swaps (
    id SERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL DEFAULT 84532,  -- Base only
    transaction_hash VARCHAR(66) NOT NULL UNIQUE,
    user_address VARCHAR(42) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    token_in VARCHAR(42) NOT NULL,
    token_out VARCHAR(42) NOT NULL,
    amount_in DECIMAL(30, 18) NOT NULL,
    amount_out DECIMAL(30, 18) NOT NULL,

    price_impact_bps INTEGER,
    fee_amount_usd DECIMAL(20, 2),

    CONSTRAINT chk_aero_chain CHECK (chain_id IN (84532, 8453))
);

CREATE INDEX idx_aero_user ON aerodrome_swaps(user_address, timestamp DESC);
CREATE INDEX idx_aero_pair ON aerodrome_swaps(token_in, token_out);

-- Base Pay onramp transactions
CREATE TABLE base_pay_transactions (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    transaction_id VARCHAR(100) NOT NULL UNIQUE,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    fiat_amount DECIMAL(20, 2) NOT NULL,
    fiat_currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    crypto_amount DECIMAL(30, 18) NOT NULL,
    crypto_asset VARCHAR(20) NOT NULL,

    payment_method VARCHAR(50),  -- 'credit_card', 'debit_card', 'bank_transfer'
    status VARCHAR(20) NOT NULL,  -- 'pending', 'completed', 'failed'

    -- Auto-purchase of treasury if configured
    auto_treasury_purchase BOOLEAN DEFAULT false,
    treasury_purchase_tx_hash VARCHAR(66)
);

CREATE INDEX idx_base_pay_user ON base_pay_transactions(user_address, timestamp DESC);
CREATE INDEX idx_base_pay_status ON base_pay_transactions(status, timestamp DESC);

COMMENT ON TABLE base_pay_transactions IS 'Base Pay credit card onramp transactions';


-- ============================================================================
-- STEP 9: Data validation and constraints
-- ============================================================================

-- Add check constraint to ensure treasury is only on Base
ALTER TABLE treasury_holdings
ADD CONSTRAINT chk_treasury_only_base
CHECK (chain_id IN (84532, 8453));

-- Add check to prevent negative risk scores
ALTER TABLE risk_assessments
ADD CONSTRAINT chk_risk_score_range
CHECK (risk_score >= 0 AND risk_score <= 100);


-- ============================================================================
-- STEP 10: Insert reference data for chains
-- ============================================================================

CREATE TABLE IF NOT EXISTS supported_chains (
    chain_id BIGINT PRIMARY KEY,
    chain_name VARCHAR(100) NOT NULL,
    network_type VARCHAR(20) NOT NULL,  -- 'mainnet', 'testnet'
    rpc_url TEXT,
    wss_url TEXT,
    block_explorer_url TEXT,

    -- Chain characteristics
    features JSONB,  -- ["defi", "rwa", "aggressive", "conservative"]
    is_active BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO supported_chains (chain_id, chain_name, network_type, features) VALUES
(421614, 'Arbitrum Sepolia', 'testnet', '["defi", "gmx", "aggressive", "high_leverage"]'),
(42161, 'Arbitrum One', 'mainnet', '["defi", "gmx", "aggressive", "high_leverage"]'),
(84532, 'Base Sepolia', 'testnet', '["rwa", "treasury", "conservative", "base_ecosystem"]'),
(8453, 'Base Mainnet', 'mainnet', '["rwa", "treasury", "conservative", "base_ecosystem"]')
ON CONFLICT (chain_id) DO NOTHING;

COMMENT ON TABLE supported_chains IS 'Supported blockchain networks configuration';


-- ============================================================================
-- COMMIT TRANSACTION
-- ============================================================================

COMMIT;

-- ============================================================================
-- Verification queries (to run after migration)
-- ============================================================================

-- Check chain distribution
-- SELECT chain_id, COUNT(*) FROM vault_positions GROUP BY chain_id;
-- SELECT chain_id, COUNT(*) FROM treasury_holdings GROUP BY chain_id;

-- Test cross-chain views
-- SELECT * FROM cross_chain_total_value LIMIT 5;
-- SELECT * FROM gmx_performance_summary;

-- Test function
-- SELECT * FROM get_cross_chain_risk_profile('0x1234...');
