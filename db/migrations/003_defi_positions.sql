-- DeFi Position Tracking Tables
-- Critical for delivering institutional risk management features

-- Table for tracking individual DeFi positions
CREATE TABLE IF NOT EXISTS defi_positions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(42) NOT NULL,
    protocol VARCHAR(20) NOT NULL CHECK (protocol IN ('aave', 'compound', 'uniswap', 'gmx')),
    position_id VARCHAR(100) UNIQUE NOT NULL,

    -- Position details
    collateral_asset VARCHAR(42),
    collateral_amount NUMERIC(78, 0) DEFAULT 0,
    debt_asset VARCHAR(42),
    debt_amount NUMERIC(78, 0) DEFAULT 0,

    -- Risk metrics
    health_factor DECIMAL(10, 4) DEFAULT 999,
    liquidation_price DECIMAL(20, 6),
    liquidation_threshold DECIMAL(5, 2),
    current_ltv DECIMAL(5, 2),

    -- Value and returns
    value_usd DECIMAL(20, 2) NOT NULL DEFAULT 0,
    pnl_usd DECIMAL(20, 2) DEFAULT 0,
    apy DECIMAL(10, 2),

    -- Risk assessment
    risk_score INTEGER CHECK (risk_score >= 0 AND risk_score <= 100),
    ai_risk_assessment JSONB,

    -- Timestamps
    opened_at TIMESTAMP DEFAULT NOW(),
    last_updated TIMESTAMP DEFAULT NOW(),
    closed_at TIMESTAMP,

    -- Indexes
    INDEX idx_user_positions (user_id),
    INDEX idx_protocol (protocol),
    INDEX idx_health_factor (health_factor),
    INDEX idx_risk_score (risk_score),
    INDEX idx_last_updated (last_updated)
);

-- Portfolio snapshots for historical tracking
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(42) NOT NULL,

    -- Portfolio metrics
    total_value DECIMAL(20, 2) NOT NULL,
    total_collateral DECIMAL(20, 2) DEFAULT 0,
    total_debt DECIMAL(20, 2) DEFAULT 0,
    net_worth DECIMAL(20, 2) GENERATED ALWAYS AS (total_value - total_debt) STORED,

    -- Risk metrics
    average_health_factor DECIMAL(10, 4) DEFAULT 999,
    overall_risk INTEGER CHECK (overall_risk >= 0 AND overall_risk <= 100),
    liquidation_risk DECIMAL(5, 2),

    -- Protocol distribution
    protocol_distribution JSONB, -- {"aave": 0.4, "compound": 0.3, ...}
    asset_distribution JSONB, -- {"ETH": 0.5, "USDC": 0.3, ...}

    -- AI insights
    recommendations TEXT[],
    ai_analysis JSONB,

    -- Timestamp
    timestamp TIMESTAMP DEFAULT NOW(),

    -- Indexes
    INDEX idx_user_snapshots (user_id, timestamp DESC),
    INDEX idx_snapshot_timestamp (timestamp DESC)
);

-- Liquidation alerts table
CREATE TABLE IF NOT EXISTS liquidation_alerts (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(42) NOT NULL,
    position_id VARCHAR(100) NOT NULL,

    -- Alert details
    alert_type VARCHAR(20) CHECK (alert_type IN ('warning', 'critical', 'liquidated')),
    health_factor DECIMAL(10, 4),
    risk_score INTEGER,

    -- Prediction
    predicted_liquidation_time TIMESTAMP,
    confidence_score DECIMAL(5, 2),

    -- Recommendations
    recommended_action TEXT,
    required_collateral DECIMAL(20, 2),

    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'resolved', 'ignored', 'expired')),
    created_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP,

    -- Indexes
    INDEX idx_user_alerts (user_id, status),
    INDEX idx_alert_type (alert_type, status)
);

-- Risk events for audit trail
CREATE TABLE IF NOT EXISTS risk_events (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(42) NOT NULL,
    position_id VARCHAR(100),

    -- Event details
    event_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) CHECK (severity IN ('info', 'warning', 'critical')),

    -- Metrics at time of event
    health_factor_before DECIMAL(10, 4),
    health_factor_after DECIMAL(10, 4),
    risk_score_before INTEGER,
    risk_score_after INTEGER,

    -- Event data
    description TEXT,
    metadata JSONB,

    -- Response
    action_taken TEXT,
    auto_hedged BOOLEAN DEFAULT FALSE,

    timestamp TIMESTAMP DEFAULT NOW(),

    -- Indexes
    INDEX idx_user_events (user_id, timestamp DESC),
    INDEX idx_event_type (event_type, timestamp DESC)
);

-- Automated hedging transactions
CREATE TABLE IF NOT EXISTS hedging_transactions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(42) NOT NULL,
    position_id VARCHAR(100) NOT NULL,

    -- Hedging details
    hedge_type VARCHAR(30) CHECK (hedge_type IN ('add_collateral', 'reduce_debt', 'close_position', 'rebalance')),

    -- Transaction data
    tx_hash VARCHAR(66) UNIQUE,
    from_protocol VARCHAR(20),
    to_protocol VARCHAR(20),

    -- Amounts
    amount NUMERIC(78, 0),
    asset VARCHAR(42),
    value_usd DECIMAL(20, 2),

    -- Risk metrics
    health_factor_before DECIMAL(10, 4),
    health_factor_after DECIMAL(10, 4),
    risk_score_improvement INTEGER,

    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'failed', 'reverted')),

    -- Timestamps
    initiated_at TIMESTAMP DEFAULT NOW(),
    confirmed_at TIMESTAMP,

    -- Indexes
    INDEX idx_user_hedges (user_id, initiated_at DESC),
    INDEX idx_tx_hash (tx_hash),
    INDEX idx_status (status)
);

-- Protocol TVL tracking
CREATE TABLE IF NOT EXISTS protocol_tvl (
    id SERIAL PRIMARY KEY,
    protocol VARCHAR(20) NOT NULL,

    -- TVL metrics
    total_value_locked DECIMAL(30, 2),
    user_count INTEGER,

    -- Protocol specific
    total_borrowed DECIMAL(30, 2),
    total_supplied DECIMAL(30, 2),
    utilization_rate DECIMAL(5, 2),

    -- Risk metrics
    protocol_risk_score INTEGER,
    systemic_risk_contribution DECIMAL(5, 2),

    timestamp TIMESTAMP DEFAULT NOW(),

    -- Indexes
    INDEX idx_protocol_tvl (protocol, timestamp DESC),
    UNIQUE (protocol, timestamp)
);

-- Create hypertable for time-series data (TimescaleDB)
-- SELECT create_hypertable('portfolio_snapshots', 'timestamp', if_not_exists => TRUE);
-- SELECT create_hypertable('protocol_tvl', 'timestamp', if_not_exists => TRUE);
-- SELECT create_hypertable('risk_events', 'timestamp', if_not_exists => TRUE);

-- Create continuous aggregates for fast queries
-- CREATE MATERIALIZED VIEW portfolio_hourly
-- WITH (timescaledb.continuous) AS
-- SELECT
--     user_id,
--     time_bucket('1 hour', timestamp) AS hour,
--     AVG(total_value) as avg_value,
--     AVG(overall_risk) as avg_risk,
--     MIN(average_health_factor) as min_health_factor
-- FROM portfolio_snapshots
-- GROUP BY user_id, hour
-- WITH NO DATA;

-- Add comments for documentation
COMMENT ON TABLE defi_positions IS 'Real-time tracking of user positions across all DeFi protocols';
COMMENT ON TABLE portfolio_snapshots IS 'Historical snapshots of user portfolios for trend analysis';
COMMENT ON TABLE liquidation_alerts IS 'AI-generated alerts for positions at risk of liquidation';
COMMENT ON TABLE risk_events IS 'Audit trail of all risk-related events';
COMMENT ON TABLE hedging_transactions IS 'Automated hedging actions taken to protect user positions';
COMMENT ON TABLE protocol_tvl IS 'Protocol-level TVL and risk metrics for systemic risk analysis';