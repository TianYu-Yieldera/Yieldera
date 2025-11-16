-- Phase 1: Data Layer - TimescaleDB Schema
-- Created: 2025-11-09
-- Purpose: Time-series optimized schema for AI risk engine

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- ============================================================
-- 1. PRICE HISTORY
-- ============================================================

CREATE TABLE IF NOT EXISTS price_history (
  time TIMESTAMPTZ NOT NULL,
  asset VARCHAR(42) NOT NULL,
  price NUMERIC(78, 18) NOT NULL,
  source VARCHAR(50) NOT NULL,  -- 'chainlink', 'uniswap_twap', 'coingecko'
  volume NUMERIC(78, 18),
  market_cap NUMERIC(78, 18),
  PRIMARY KEY (time, asset, source)
);

-- Convert to hypertable
SELECT create_hypertable('price_history', 'time', if_not_exists => TRUE);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_price_history_asset
  ON price_history(asset, time DESC);
CREATE INDEX IF NOT EXISTS idx_price_history_source
  ON price_history(source, time DESC);

-- Retention policy: keep 1 year
SELECT add_retention_policy('price_history', INTERVAL '365 days', if_not_exists => TRUE);

-- Comments
COMMENT ON TABLE price_history IS 'Historical price data from multiple sources for AI model training';
COMMENT ON COLUMN price_history.source IS 'Data source: chainlink (most reliable), uniswap_twap (DEX price), coingecko (centralized)';

-- ============================================================
-- 2. LIQUIDATION HISTORY
-- ============================================================

CREATE TABLE IF NOT EXISTS liquidation_history (
  id SERIAL,
  time TIMESTAMPTZ NOT NULL,
  protocol VARCHAR(50) NOT NULL,  -- 'aave', 'compound', 'gmx'
  user_address VARCHAR(42) NOT NULL,
  collateral_asset VARCHAR(42) NOT NULL,
  debt_asset VARCHAR(42),
  collateral_amount NUMERIC(78, 18) NOT NULL,
  debt_amount NUMERIC(78, 18),
  liquidation_price NUMERIC(78, 18),
  health_factor NUMERIC(10, 4),
  gas_price NUMERIC(20, 0),
  tx_hash VARCHAR(66) NOT NULL UNIQUE,
  block_number BIGINT NOT NULL,
  PRIMARY KEY (time, id)
);

-- Convert to hypertable
SELECT create_hypertable('liquidation_history', 'time', if_not_exists => TRUE);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_liquidation_user
  ON liquidation_history(user_address, time DESC);
CREATE INDEX IF NOT EXISTS idx_liquidation_protocol
  ON liquidation_history(protocol, time DESC);
CREATE INDEX IF NOT EXISTS idx_liquidation_asset
  ON liquidation_history(collateral_asset, time DESC);
CREATE INDEX IF NOT EXISTS idx_liquidation_health_factor
  ON liquidation_history(health_factor)
  WHERE health_factor IS NOT NULL;

-- Retention policy: keep 2 years (valuable ML training data)
SELECT add_retention_policy('liquidation_history', INTERVAL '730 days', if_not_exists => TRUE);

-- Comments
COMMENT ON TABLE liquidation_history IS 'Historical liquidation events for training liquidation prediction models';
COMMENT ON COLUMN liquidation_history.health_factor IS 'Health factor at time of liquidation (if available)';

-- ============================================================
-- 3. USER RISK PROFILES
-- ============================================================

CREATE TABLE IF NOT EXISTS user_risk_profiles (
  user_address VARCHAR(42) PRIMARY KEY,
  risk_score NUMERIC(5, 2) NOT NULL DEFAULT 50.0,  -- 0-100 scale
  avg_leverage NUMERIC(10, 4),
  liquidation_count INTEGER DEFAULT 0,
  total_positions INTEGER DEFAULT 0,
  max_position_size NUMERIC(78, 18),
  preferred_assets JSONB,
  preferred_protocols JSONB,
  avg_health_factor NUMERIC(10, 4),
  total_volume_usd NUMERIC(78, 18) DEFAULT 0,
  last_activity TIMESTAMPTZ,
  last_updated TIMESTAMPTZ DEFAULT NOW(),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_user_risk_score
  ON user_risk_profiles(risk_score DESC);
CREATE INDEX IF NOT EXISTS idx_user_last_activity
  ON user_risk_profiles(last_activity DESC);
CREATE INDEX IF NOT EXISTS idx_user_preferred_assets
  ON user_risk_profiles USING GIN(preferred_assets);

-- Comments
COMMENT ON TABLE user_risk_profiles IS 'User behavior profiles for personalized risk assessment';
COMMENT ON COLUMN user_risk_profiles.risk_score IS 'Composite risk score: 0=safest, 100=riskiest';

-- ============================================================
-- 4. MARKET DEPTH SNAPSHOTS
-- ============================================================

CREATE TABLE IF NOT EXISTS market_depth_snapshots (
  time TIMESTAMPTZ NOT NULL,
  protocol VARCHAR(50) NOT NULL,  -- 'aave', 'compound', 'uniswap', 'gmx'
  market VARCHAR(100) NOT NULL,  -- 'ETH-USDC', 'WBTC-USDC', etc.
  total_supply NUMERIC(78, 18),
  total_borrow NUMERIC(78, 18),
  utilization_rate NUMERIC(5, 4),  -- 0-1
  supply_apy NUMERIC(10, 6),  -- Annual percentage yield
  borrow_apy NUMERIC(10, 6),
  available_liquidity NUMERIC(78, 18),
  total_reserves NUMERIC(78, 18),
  PRIMARY KEY (time, protocol, market)
);

-- Convert to hypertable
SELECT create_hypertable('market_depth_snapshots', 'time', if_not_exists => TRUE);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_market_depth_protocol_market
  ON market_depth_snapshots(protocol, market, time DESC);
CREATE INDEX IF NOT EXISTS idx_market_depth_utilization
  ON market_depth_snapshots(utilization_rate DESC, time DESC);

-- Retention policy: keep 6 months
SELECT add_retention_policy('market_depth_snapshots', INTERVAL '180 days', if_not_exists => TRUE);

-- Comments
COMMENT ON TABLE market_depth_snapshots IS 'Regular snapshots of market liquidity and rates for trend analysis';

-- ============================================================
-- 5. PROTOCOL EVENTS (RAW)
-- ============================================================

CREATE TABLE IF NOT EXISTS protocol_events (
  id SERIAL,
  time TIMESTAMPTZ NOT NULL,
  protocol VARCHAR(50) NOT NULL,
  event_name VARCHAR(100) NOT NULL,
  contract_address VARCHAR(42) NOT NULL,
  tx_hash VARCHAR(66) NOT NULL,
  block_number BIGINT NOT NULL,
  log_index INTEGER NOT NULL,
  event_data JSONB NOT NULL,
  processed BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (time, id)
);

-- Convert to hypertable
SELECT create_hypertable('protocol_events', 'time', if_not_exists => TRUE);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_protocol_events_name
  ON protocol_events(event_name, time DESC);
CREATE INDEX IF NOT EXISTS idx_protocol_events_contract
  ON protocol_events(contract_address, time DESC);
CREATE INDEX IF NOT EXISTS idx_protocol_events_processed
  ON protocol_events(processed, time DESC)
  WHERE processed = FALSE;
CREATE INDEX IF NOT EXISTS idx_protocol_events_tx
  ON protocol_events(tx_hash);

-- Retention policy: keep 90 days (raw events processed into other tables)
SELECT add_retention_policy('protocol_events', INTERVAL '90 days', if_not_exists => TRUE);

-- Comments
COMMENT ON TABLE protocol_events IS 'Raw blockchain events before processing into structured tables';

-- ============================================================
-- 6. POSITION SNAPSHOTS (For ML Training)
-- ============================================================

CREATE TABLE IF NOT EXISTS position_snapshots (
  id SERIAL,
  time TIMESTAMPTZ NOT NULL,
  user_address VARCHAR(42) NOT NULL,
  protocol VARCHAR(50) NOT NULL,
  position_id VARCHAR(100),
  collateral_asset VARCHAR(42) NOT NULL,
  collateral_amount NUMERIC(78, 18) NOT NULL,
  debt_asset VARCHAR(42),
  debt_amount NUMERIC(78, 18),
  health_factor NUMERIC(10, 4),
  ltv NUMERIC(5, 4),  -- Loan-to-value ratio
  leverage NUMERIC(10, 4),
  liquidation_price NUMERIC(78, 18),
  collateral_value_usd NUMERIC(78, 18),
  debt_value_usd NUMERIC(78, 18),
  was_liquidated BOOLEAN DEFAULT FALSE,
  days_until_liquidation INTEGER,  -- For ML labels
  PRIMARY KEY (time, id)
);

-- Convert to hypertable
SELECT create_hypertable('position_snapshots', 'time', if_not_exists => TRUE);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_position_snapshots_user
  ON position_snapshots(user_address, time DESC);
CREATE INDEX IF NOT EXISTS idx_position_snapshots_protocol
  ON position_snapshots(protocol, time DESC);
CREATE INDEX IF NOT EXISTS idx_position_snapshots_health_factor
  ON position_snapshots(health_factor, time DESC)
  WHERE health_factor IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_position_snapshots_liquidated
  ON position_snapshots(was_liquidated, time DESC)
  WHERE was_liquidated = TRUE;

-- Retention policy: keep 1 year (critical for ML)
SELECT add_retention_policy('position_snapshots', INTERVAL '365 days', if_not_exists => TRUE);

-- Comments
COMMENT ON TABLE position_snapshots IS 'Regular snapshots of user positions for training liquidation prediction models';
COMMENT ON COLUMN position_snapshots.was_liquidated IS 'Label for ML: was this position liquidated?';
COMMENT ON COLUMN position_snapshots.days_until_liquidation IS 'ML label: days until liquidation (NULL if not liquidated)';

-- ============================================================
-- CONTINUOUS AGGREGATES (Pre-computed Analytics)
-- ============================================================

-- Daily price statistics
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_price_stats
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 day', time) AS day,
  asset,
  source,
  FIRST(price, time) AS open,
  MAX(price) AS high,
  MIN(price) AS low,
  LAST(price, time) AS close,
  AVG(price) AS avg_price,
  STDDEV(price) AS volatility,
  COUNT(*) AS sample_count
FROM price_history
GROUP BY day, asset, source
WITH NO DATA;

-- Refresh policy: update hourly
SELECT add_continuous_aggregate_policy('daily_price_stats',
  start_offset => INTERVAL '3 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '1 hour',
  if_not_exists => TRUE);

-- Hourly market metrics
CREATE MATERIALIZED VIEW IF NOT EXISTS hourly_market_metrics
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 hour', time) AS hour,
  protocol,
  market,
  AVG(utilization_rate) AS avg_utilization,
  AVG(supply_apy) AS avg_supply_apy,
  AVG(borrow_apy) AS avg_borrow_apy,
  AVG(available_liquidity) AS avg_liquidity,
  LAST(total_supply, time) AS last_total_supply,
  LAST(total_borrow, time) AS last_total_borrow
FROM market_depth_snapshots
GROUP BY hour, protocol, market
WITH NO DATA;

-- Refresh policy: update every 10 minutes
SELECT add_continuous_aggregate_policy('hourly_market_metrics',
  start_offset => INTERVAL '1 day',
  end_offset => INTERVAL '10 minutes',
  schedule_interval => INTERVAL '10 minutes',
  if_not_exists => TRUE);

-- Daily liquidation statistics
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_liquidation_stats
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 day', time) AS day,
  protocol,
  collateral_asset,
  COUNT(*) AS liquidation_count,
  SUM(collateral_amount) AS total_collateral_liquidated,
  SUM(debt_amount) AS total_debt_liquidated,
  AVG(health_factor) AS avg_health_factor_at_liquidation,
  PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY health_factor) AS median_health_factor,
  MIN(health_factor) AS min_health_factor,
  MAX(health_factor) AS max_health_factor
FROM liquidation_history
WHERE health_factor IS NOT NULL
GROUP BY day, protocol, collateral_asset
WITH NO DATA;

-- Refresh policy: update hourly
SELECT add_continuous_aggregate_policy('daily_liquidation_stats',
  start_offset => INTERVAL '7 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '1 hour',
  if_not_exists => TRUE);

-- User activity summary (last 30 days)
CREATE MATERIALIZED VIEW IF NOT EXISTS user_activity_summary
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 day', time) AS day,
  user_address,
  COUNT(DISTINCT protocol) AS protocols_used,
  COUNT(*) AS position_count,
  AVG(health_factor) AS avg_health_factor,
  MIN(health_factor) AS min_health_factor,
  MAX(leverage) AS max_leverage,
  SUM(collateral_value_usd) AS total_collateral_usd
FROM position_snapshots
WHERE health_factor IS NOT NULL
GROUP BY day, user_address
WITH NO DATA;

-- Refresh policy: update every 6 hours
SELECT add_continuous_aggregate_policy('user_activity_summary',
  start_offset => INTERVAL '30 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '6 hours',
  if_not_exists => TRUE);

-- ============================================================
-- UTILITY VIEWS
-- ============================================================

-- Latest market rates
CREATE OR REPLACE VIEW latest_market_rates AS
SELECT DISTINCT ON (protocol, market)
  protocol,
  market,
  time,
  supply_apy,
  borrow_apy,
  utilization_rate,
  available_liquidity
FROM market_depth_snapshots
ORDER BY protocol, market, time DESC;

-- User risk rankings
CREATE OR REPLACE VIEW user_risk_rankings AS
SELECT
  user_address,
  risk_score,
  liquidation_count,
  avg_health_factor,
  total_positions,
  RANK() OVER (ORDER BY risk_score DESC) AS risk_rank,
  PERCENT_RANK() OVER (ORDER BY risk_score DESC) AS risk_percentile
FROM user_risk_profiles
WHERE total_positions > 0;

-- High risk positions (real-time)
CREATE OR REPLACE VIEW high_risk_positions AS
SELECT
  ps.user_address,
  ps.protocol,
  ps.collateral_asset,
  ps.collateral_amount,
  ps.debt_asset,
  ps.debt_amount,
  ps.health_factor,
  ps.liquidation_price,
  ps.time AS snapshot_time,
  urp.risk_score AS user_risk_score
FROM position_snapshots ps
JOIN user_risk_profiles urp ON ps.user_address = urp.user_address
WHERE ps.time > NOW() - INTERVAL '1 hour'
  AND (ps.health_factor < 1.2 OR urp.risk_score > 70)
ORDER BY ps.health_factor ASC NULLS LAST;

-- ============================================================
-- HELPER FUNCTIONS
-- ============================================================

-- Calculate price volatility
CREATE OR REPLACE FUNCTION calculate_volatility(
  p_asset VARCHAR,
  p_days INTEGER DEFAULT 30
)
RETURNS NUMERIC AS $$
DECLARE
  volatility NUMERIC;
BEGIN
  SELECT STDDEV(price) / AVG(price) INTO volatility
  FROM price_history
  WHERE asset = p_asset
    AND time > NOW() - (p_days || ' days')::INTERVAL
    AND source = 'chainlink';  -- Use most reliable source

  RETURN COALESCE(volatility, 0);
END;
$$ LANGUAGE plpgsql;

-- Get liquidation probability (simple heuristic, will be replaced by ML)
CREATE OR REPLACE FUNCTION get_liquidation_probability(
  p_protocol VARCHAR,
  p_health_factor NUMERIC
)
RETURNS NUMERIC AS $$
DECLARE
  probability NUMERIC;
  historical_rate NUMERIC;
BEGIN
  -- Get historical liquidation rate for similar health factors
  SELECT COUNT(*) FILTER (WHERE was_liquidated = TRUE)::NUMERIC /
         NULLIF(COUNT(*), 0)
  INTO historical_rate
  FROM position_snapshots
  WHERE protocol = p_protocol
    AND health_factor BETWEEN p_health_factor - 0.1 AND p_health_factor + 0.1
    AND time > NOW() - INTERVAL '90 days';

  -- Simple probability model (will be replaced by ML)
  IF p_health_factor < 1.0 THEN
    probability := 1.0;
  ELSIF p_health_factor > 2.0 THEN
    probability := COALESCE(historical_rate, 0.01);
  ELSE
    -- Linear interpolation between 1.0 and 2.0
    probability := 1.0 - (p_health_factor - 1.0);
    probability := GREATEST(probability, COALESCE(historical_rate, 0.01));
  END IF;

  RETURN LEAST(probability, 1.0);
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- COMMENTS
-- ============================================================

COMMENT ON SCHEMA public IS 'Phase 1 Data Layer: Time-series optimized schema for AI risk engine';
COMMENT ON FUNCTION calculate_volatility IS 'Calculate historical volatility for an asset';
COMMENT ON FUNCTION get_liquidation_probability IS 'Estimate liquidation probability (placeholder for ML model)';

-- ============================================================
-- SAMPLE QUERIES (For Testing)
-- ============================================================

/*
-- Get price volatility for ETH in last 30 days
SELECT calculate_volatility('0xEthAddress', 30);

-- Get liquidation probability for a position
SELECT get_liquidation_probability('aave', 1.5);

-- Find similar historical liquidations
SELECT *
FROM liquidation_history
WHERE protocol = 'aave'
  AND collateral_asset = '0xEthAddress'
  AND health_factor BETWEEN 1.3 AND 1.7
ORDER BY time DESC
LIMIT 100;

-- Get daily price stats for last week
SELECT *
FROM daily_price_stats
WHERE asset = '0xEthAddress'
  AND day > NOW() - INTERVAL '7 days'
ORDER BY day DESC;

-- Get market utilization trends
SELECT
  hour,
  protocol,
  market,
  avg_utilization,
  avg_supply_apy,
  avg_borrow_apy
FROM hourly_market_metrics
WHERE protocol = 'aave'
  AND hour > NOW() - INTERVAL '24 hours'
ORDER BY hour DESC;
*/

-- End of schema
