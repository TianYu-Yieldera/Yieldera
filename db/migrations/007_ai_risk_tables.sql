-- ============================================================
-- AI Risk Engine Database Schema
-- Migration 007: AI Risk Management Tables
-- ============================================================
-- This migration creates all tables required by the AI Risk Engine:
-- - Time-series price history
-- - User risk profiles
-- - High-risk position monitoring
-- - Market data aggregations
-- - ML prediction storage
-- ============================================================

-- Enable TimescaleDB extension (if not already enabled)
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- ============================================================
-- PRICE HISTORY - Time-series price data
-- ============================================================
CREATE TABLE IF NOT EXISTS price_history (
    time TIMESTAMPTZ NOT NULL,
    asset TEXT NOT NULL,
    price NUMERIC(30, 18) NOT NULL,
    source TEXT NOT NULL,
    volume NUMERIC(30, 18),
    market_cap NUMERIC(30, 18),
    PRIMARY KEY (time, asset, source)
);

-- Convert to TimescaleDB hypertable for efficient time-series queries
SELECT create_hypertable('price_history', 'time', if_not_exists => TRUE);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_price_history_asset ON price_history(asset, time DESC);
CREATE INDEX IF NOT EXISTS idx_price_history_source ON price_history(source, time DESC);

-- ============================================================
-- USER RISK PROFILES - Aggregated user risk metrics
-- ============================================================
CREATE TABLE IF NOT EXISTS user_risk_profiles (
    user_address TEXT PRIMARY KEY,
    risk_score NUMERIC(10, 6) DEFAULT 0,
    total_collateral NUMERIC(30, 18) DEFAULT 0,
    total_debt NUMERIC(30, 18) DEFAULT 0,
    health_factor NUMERIC(10, 6),
    liquidation_probability NUMERIC(10, 6),
    position_count INT DEFAULT 0,
    protocols TEXT[], -- Array of protocols user is active in
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_risk_profiles_risk_score ON user_risk_profiles(risk_score DESC);
CREATE INDEX IF NOT EXISTS idx_user_risk_profiles_updated ON user_risk_profiles(last_updated DESC);

-- ============================================================
-- HIGH RISK POSITIONS - Real-time monitoring view
-- ============================================================
CREATE TABLE IF NOT EXISTS high_risk_positions (
    position_id TEXT PRIMARY KEY,
    user_address TEXT NOT NULL,
    protocol TEXT NOT NULL,
    collateral_asset TEXT NOT NULL,
    collateral_amount NUMERIC(30, 18) NOT NULL,
    debt_asset TEXT,
    debt_amount NUMERIC(30, 18) DEFAULT 0,
    health_factor NUMERIC(10, 6),
    risk_score NUMERIC(10, 6) NOT NULL,
    liquidation_probability NUMERIC(10, 6),
    alert_level TEXT CHECK (alert_level IN ('safe', 'warning', 'danger', 'critical')),
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_high_risk_user ON high_risk_positions(user_address);
CREATE INDEX IF NOT EXISTS idx_high_risk_protocol ON high_risk_positions(protocol);
CREATE INDEX IF NOT EXISTS idx_high_risk_score ON high_risk_positions(risk_score DESC);
CREATE INDEX IF NOT EXISTS idx_high_risk_alert ON high_risk_positions(alert_level, last_updated DESC);

-- ============================================================
-- RISK ASSESSMENTS - Historical risk calculation results
-- ============================================================
CREATE TABLE IF NOT EXISTS risk_assessments (
    id BIGSERIAL PRIMARY KEY,
    assessment_id TEXT UNIQUE NOT NULL,
    position_id TEXT,
    user_address TEXT NOT NULL,
    protocol TEXT NOT NULL,
    risk_score NUMERIC(10, 6) NOT NULL,
    liquidation_probability NUMERIC(10, 6),
    alert_level TEXT,
    confidence_level NUMERIC(10, 6),
    market_risk NUMERIC(10, 6),
    liquidity_risk NUMERIC(10, 6),
    value_at_risk NUMERIC(30, 18),
    recommendations JSONB,
    risk_metrics JSONB,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Convert to hypertable for time-series efficiency
SELECT create_hypertable('risk_assessments', 'timestamp', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_risk_assessments_user ON risk_assessments(user_address, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_risk_assessments_protocol ON risk_assessments(protocol, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_risk_assessments_score ON risk_assessments(risk_score DESC, timestamp DESC);

-- ============================================================
-- ML PREDICTIONS - AI model prediction results
-- ============================================================
CREATE TABLE IF NOT EXISTS ml_predictions (
    id BIGSERIAL PRIMARY KEY,
    prediction_id TEXT UNIQUE NOT NULL,
    user_address TEXT,
    position_id TEXT,
    model_type TEXT NOT NULL, -- 'xgboost', 'lstm', 'hybrid'
    prediction_type TEXT NOT NULL, -- 'liquidation', 'risk_score', 'market_crash'
    prediction_value NUMERIC(10, 6) NOT NULL,
    confidence NUMERIC(10, 6),
    features JSONB, -- Input features used
    model_version TEXT,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Convert to hypertable
SELECT create_hypertable('ml_predictions', 'timestamp', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_ml_predictions_user ON ml_predictions(user_address, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_ml_predictions_type ON ml_predictions(prediction_type, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_ml_predictions_model ON ml_predictions(model_type, timestamp DESC);

-- ============================================================
-- MARKET SNAPSHOTS - Aggregated market data
-- ============================================================
CREATE TABLE IF NOT EXISTS market_snapshots (
    time TIMESTAMPTZ NOT NULL,
    snapshot_id TEXT NOT NULL,
    price_data JSONB,
    volatility_metrics JSONB,
    market_depth JSONB,
    liquidity_data JSONB,
    PRIMARY KEY (time, snapshot_id)
);

-- Convert to hypertable
SELECT create_hypertable('market_snapshots', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_market_snapshots_id ON market_snapshots(snapshot_id, time DESC);

-- ============================================================
-- LIQUIDATION EVENTS - Historical liquidation data
-- ============================================================
CREATE TABLE IF NOT EXISTS liquidation_events (
    id BIGSERIAL PRIMARY KEY,
    event_id TEXT UNIQUE NOT NULL,
    user_address TEXT NOT NULL,
    protocol TEXT NOT NULL,
    collateral_asset TEXT NOT NULL,
    collateral_amount NUMERIC(30, 18),
    debt_asset TEXT,
    debt_amount NUMERIC(30, 18),
    liquidation_price NUMERIC(30, 18),
    liquidation_value NUMERIC(30, 18),
    tx_hash TEXT,
    block_number BIGINT,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Convert to hypertable
SELECT create_hypertable('liquidation_events', 'timestamp', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_liquidation_user ON liquidation_events(user_address, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_liquidation_protocol ON liquidation_events(protocol, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_liquidation_asset ON liquidation_events(collateral_asset, timestamp DESC);

-- ============================================================
-- MONITORING ALERTS - Alert history
-- ============================================================
CREATE TABLE IF NOT EXISTS monitoring_alerts (
    id BIGSERIAL PRIMARY KEY,
    alert_id TEXT UNIQUE NOT NULL,
    user_address TEXT,
    position_id TEXT,
    alert_type TEXT NOT NULL, -- 'liquidation_risk', 'systemic_risk', etc.
    alert_level TEXT NOT NULL, -- 'info', 'warning', 'danger', 'critical'
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    action_required TEXT,
    auto_response TEXT,
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_at TIMESTAMPTZ,
    data JSONB,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Convert to hypertable
SELECT create_hypertable('monitoring_alerts', 'timestamp', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_monitoring_alerts_user ON monitoring_alerts(user_address, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_monitoring_alerts_level ON monitoring_alerts(alert_level, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_monitoring_alerts_type ON monitoring_alerts(alert_type, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_monitoring_alerts_acked ON monitoring_alerts(acknowledged, timestamp DESC);

-- ============================================================
-- MARKET DATA AGGREGATIONS - 1h continuous aggregates
-- ============================================================
CREATE MATERIALIZED VIEW IF NOT EXISTS market_data_1h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', time) AS bucket,
    asset,
    first(price, time) AS open,
    max(price) AS high,
    min(price) AS low,
    last(price, time) AS close,
    avg(price) AS avg_price,
    sum(volume) AS total_volume,
    count(*) AS data_points
FROM price_history
GROUP BY bucket, asset
WITH NO DATA;

-- Add refresh policy for continuous aggregate
SELECT add_continuous_aggregate_policy('market_data_1h',
    start_offset => INTERVAL '3 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => TRUE
);

-- ============================================================
-- HELPER FUNCTIONS
-- ============================================================

-- Function to get liquidation probability based on health factor
CREATE OR REPLACE FUNCTION get_liquidation_probability(
    p_protocol TEXT,
    p_health_factor NUMERIC
) RETURNS NUMERIC AS $$
BEGIN
    -- Simple liquidation probability model
    -- Can be enhanced with protocol-specific logic
    IF p_health_factor IS NULL OR p_health_factor <= 0 THEN
        RETURN 1.0; -- 100% probability
    ELSIF p_health_factor < 1.0 THEN
        RETURN 1.0; -- Already liquidatable
    ELSIF p_health_factor < 1.1 THEN
        RETURN 0.95; -- Very high risk
    ELSIF p_health_factor < 1.3 THEN
        RETURN 0.75; -- High risk
    ELSIF p_health_factor < 1.5 THEN
        RETURN 0.40; -- Moderate risk
    ELSIF p_health_factor < 2.0 THEN
        RETURN 0.15; -- Low risk
    ELSE
        RETURN 0.05; -- Very low risk
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Function to update user risk profile
CREATE OR REPLACE FUNCTION update_user_risk_profile(
    p_user_address TEXT,
    p_risk_score NUMERIC,
    p_total_collateral NUMERIC,
    p_total_debt NUMERIC,
    p_health_factor NUMERIC,
    p_liquidation_probability NUMERIC,
    p_position_count INT,
    p_protocols TEXT[]
) RETURNS VOID AS $$
BEGIN
    INSERT INTO user_risk_profiles (
        user_address,
        risk_score,
        total_collateral,
        total_debt,
        health_factor,
        liquidation_probability,
        position_count,
        protocols,
        last_updated
    ) VALUES (
        p_user_address,
        p_risk_score,
        p_total_collateral,
        p_total_debt,
        p_health_factor,
        p_liquidation_probability,
        p_position_count,
        p_protocols,
        NOW()
    )
    ON CONFLICT (user_address) DO UPDATE SET
        risk_score = EXCLUDED.risk_score,
        total_collateral = EXCLUDED.total_collateral,
        total_debt = EXCLUDED.total_debt,
        health_factor = EXCLUDED.health_factor,
        liquidation_probability = EXCLUDED.liquidation_probability,
        position_count = EXCLUDED.position_count,
        protocols = EXCLUDED.protocols,
        last_updated = NOW();
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- DATA RETENTION POLICIES
-- ============================================================

-- Retain raw price data for 90 days
SELECT add_retention_policy('price_history', INTERVAL '90 days', if_not_exists => TRUE);

-- Retain risk assessments for 180 days
SELECT add_retention_policy('risk_assessments', INTERVAL '180 days', if_not_exists => TRUE);

-- Retain ML predictions for 180 days
SELECT add_retention_policy('ml_predictions', INTERVAL '180 days', if_not_exists => TRUE);

-- Retain liquidation events for 1 year
SELECT add_retention_policy('liquidation_events', INTERVAL '1 year', if_not_exists => TRUE);

-- Retain monitoring alerts for 90 days
SELECT add_retention_policy('monitoring_alerts', INTERVAL '90 days', if_not_exists => TRUE);

-- Retain market snapshots for 30 days
SELECT add_retention_policy('market_snapshots', INTERVAL '30 days', if_not_exists => TRUE);

-- ============================================================
-- INITIAL DEMO DATA (Optional)
-- ============================================================

-- Insert some sample price data for testing
INSERT INTO price_history (time, asset, price, source, volume, market_cap)
VALUES
    (NOW(), 'ETH', 2000.00, 'chainlink', 1000000, 240000000000),
    (NOW(), 'BTC', 35000.00, 'chainlink', 500000, 680000000000),
    (NOW(), 'USDC', 1.00, 'chainlink', 10000000, 25000000000),
    (NOW(), 'USDT', 1.00, 'chainlink', 15000000, 90000000000)
ON CONFLICT (time, asset, source) DO NOTHING;

-- ============================================================
-- COMMENTS AND DOCUMENTATION
-- ============================================================

COMMENT ON TABLE price_history IS 'Time-series price data for assets from various oracles';
COMMENT ON TABLE user_risk_profiles IS 'Aggregated risk metrics per user address';
COMMENT ON TABLE high_risk_positions IS 'Real-time view of positions requiring monitoring';
COMMENT ON TABLE risk_assessments IS 'Historical risk calculation results';
COMMENT ON TABLE ml_predictions IS 'AI/ML model prediction outputs';
COMMENT ON TABLE market_snapshots IS 'Periodic snapshots of market conditions';
COMMENT ON TABLE liquidation_events IS 'Historical liquidation event records';
COMMENT ON TABLE monitoring_alerts IS 'Alert history for risk monitoring';

-- Migration complete
SELECT 'AI Risk Engine database schema created successfully' AS status;
