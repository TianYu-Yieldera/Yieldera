-- ============================================================================
-- Migration 004: Support for New Services (Fixed for TimescaleDB)
-- Date: 2025-11-10
-- Description: Add database tables for new services with proper TimescaleDB support
-- ============================================================================

-- ============================================================================
-- 1. Yield Calculation Service Tables
-- ============================================================================

-- Treasury yield rates
CREATE TABLE IF NOT EXISTS treasury_yield_rates (
    bond_type VARCHAR(20) NOT NULL,
    effective_date DATE NOT NULL,
    annual_yield NUMERIC(8, 6) NOT NULL,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (bond_type, effective_date)
);

CREATE INDEX IF NOT EXISTS idx_yield_rates_date ON treasury_yield_rates(effective_date DESC);

-- Treasury holdings
CREATE TABLE IF NOT EXISTS treasury_holdings (
    user_id VARCHAR(66) NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    token_amount NUMERIC(78, 0) NOT NULL,
    principal_usd NUMERIC(18, 2) NOT NULL,
    purchase_date TIMESTAMPTZ NOT NULL,
    last_yield_date TIMESTAMPTZ DEFAULT NOW(),
    total_yield_earned NUMERIC(18, 2) DEFAULT 0,
    compounding_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, bond_type)
);

CREATE INDEX IF NOT EXISTS idx_holdings_user ON treasury_holdings(user_id);
CREATE INDEX IF NOT EXISTS idx_holdings_updated ON treasury_holdings(updated_at DESC);

-- Daily yield accruals (hypertable)
CREATE TABLE IF NOT EXISTS treasury_yield_accruals (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(66) NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    principal_amount NUMERIC(18, 2) NOT NULL,
    daily_yield_rate NUMERIC(12, 9) NOT NULL,
    yield_amount_usd NUMERIC(18, 2) NOT NULL,
    cumulative_yield NUMERIC(18, 2) NOT NULL,
    compounded BOOLEAN DEFAULT false
);

SELECT create_hypertable(
    'treasury_yield_accruals',
    'time',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 month'
);

CREATE INDEX IF NOT EXISTS idx_accruals_user ON treasury_yield_accruals(user_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_accruals_bond_type ON treasury_yield_accruals(bond_type, time DESC);

-- Monthly yield summary (continuous aggregate)
CREATE MATERIALIZED VIEW IF NOT EXISTS monthly_yield_summary
WITH (timescaledb.continuous) AS
SELECT
    user_id,
    bond_type,
    time_bucket('1 month', time) AS month,
    SUM(yield_amount_usd) as total_yield,
    AVG(daily_yield_rate) as avg_rate,
    COUNT(*) as days_count
FROM treasury_yield_accruals
GROUP BY user_id, bond_type, month
WITH NO DATA;

SELECT add_continuous_aggregate_policy(
    'monthly_yield_summary',
    start_offset => INTERVAL '3 months',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Tax reports
CREATE TABLE IF NOT EXISTS treasury_tax_reports (
    user_id VARCHAR(66) NOT NULL,
    tax_year INTEGER NOT NULL,
    total_interest NUMERIC(18, 2) NOT NULL,
    report_data JSONB NOT NULL,
    generated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, tax_year)
);

CREATE INDEX IF NOT EXISTS idx_tax_reports_year ON treasury_tax_reports(tax_year DESC);

-- ============================================================================
-- 2. Notification Service Tables
-- ============================================================================

-- User notification preferences
CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id VARCHAR(66) PRIMARY KEY,
    channels TEXT[] NOT NULL DEFAULT '{"email"}',
    min_priority VARCHAR(20) DEFAULT 'medium',
    quiet_hours_start INTEGER,
    quiet_hours_end INTEGER,
    enabled_types TEXT[] NOT NULL DEFAULT '{}',
    frequency VARCHAR(20) DEFAULT 'realtime',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Notification history (hypertable)
CREATE TABLE IF NOT EXISTS notifications (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(66) NOT NULL,
    type VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    channels TEXT[] NOT NULL,
    read_at TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'sent'
);

SELECT create_hypertable(
    'notifications',
    'time',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 week'
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type, time DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(user_id, time DESC) WHERE read_at IS NULL;

-- Notification queue
CREATE TABLE IF NOT EXISTS notification_queue (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    notification_data JSONB NOT NULL,
    scheduled_for TIMESTAMPTZ NOT NULL,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notif_queue_scheduled ON notification_queue(scheduled_for) WHERE NOT processed;

-- User devices (for push notifications)
CREATE TABLE IF NOT EXISTS user_devices (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    fcm_token TEXT NOT NULL,
    device_type VARCHAR(20),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_seen TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_devices_user_active ON user_devices(user_id) WHERE active = true;

-- ============================================================================
-- 3. Auto Hedge Executor Tables
-- ============================================================================

-- User settings for auto-hedge
CREATE TABLE IF NOT EXISTS user_settings (
    user_id VARCHAR(66) PRIMARY KEY,
    auto_hedge_enabled BOOLEAN DEFAULT false,
    max_hedge_amount NUMERIC(18, 2) DEFAULT 10000,
    min_health_factor NUMERIC(8, 4) DEFAULT 1.5,
    target_health_factor NUMERIC(8, 4) DEFAULT 2.0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_settings_enabled ON user_settings(user_id) WHERE auto_hedge_enabled = true;

-- Hedge execution records (hypertable)
CREATE TABLE IF NOT EXISTS hedge_executions (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(66) NOT NULL,
    risk_signal JSONB NOT NULL,
    strategy VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    simulated_outcome JSONB,
    tx_hash VARCHAR(66),
    gas_used NUMERIC(78, 0),
    executed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error TEXT,
    resulting_health_factor NUMERIC(8, 4)
);

SELECT create_hypertable(
    'hedge_executions',
    'time',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 month'
);

CREATE INDEX IF NOT EXISTS idx_hedge_exec_user ON hedge_executions(user_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_hedge_exec_status ON hedge_executions(status, time DESC);

-- Daily hedge statistics (continuous aggregate)
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_hedge_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', time) AS day,
    COUNT(*) as total_executions,
    COUNT(*) FILTER (WHERE status = 'completed') as successful,
    COUNT(*) FILTER (WHERE status = 'failed') as failed,
    AVG(resulting_health_factor) FILTER (WHERE resulting_health_factor IS NOT NULL) as avg_resulting_hf,
    SUM(gas_used) FILTER (WHERE gas_used IS NOT NULL) as total_gas_used
FROM hedge_executions
GROUP BY day
WITH NO DATA;

SELECT add_continuous_aggregate_policy(
    'daily_hedge_stats',
    start_offset => INTERVAL '3 months',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- ============================================================================
-- 4. Yield Distribution Service Tables
-- ============================================================================

-- Individual yield distributions (hypertable)
CREATE TABLE IF NOT EXISTS yield_distributions (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(66) NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    yield_amount NUMERIC(18, 2) NOT NULL,
    token_amount NUMERIC(78, 0) NOT NULL,
    tx_hash VARCHAR(66),
    status VARCHAR(20) DEFAULT 'pending',
    retry_count INTEGER DEFAULT 0,
    error TEXT
);

SELECT create_hypertable(
    'yield_distributions',
    'time',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 month'
);

CREATE INDEX IF NOT EXISTS idx_yield_dist_user ON yield_distributions(user_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_yield_dist_status ON yield_distributions(status, time DESC);
CREATE INDEX IF NOT EXISTS idx_yield_dist_pending ON yield_distributions(time DESC) WHERE status IN ('pending', 'failed');

-- Batch distribution records
CREATE TABLE IF NOT EXISTS batch_distributions (
    id SERIAL PRIMARY KEY,
    batch_number INTEGER NOT NULL,
    distribution_date DATE NOT NULL,
    total_recipients INTEGER NOT NULL,
    total_yield_usd NUMERIC(18, 2) NOT NULL,
    total_token_amount NUMERIC(78, 0) NOT NULL,
    merkle_root VARCHAR(66) NOT NULL,
    tx_hash VARCHAR(66),
    gas_used NUMERIC(78, 0),
    gas_price NUMERIC(78, 0),
    status VARCHAR(20) DEFAULT 'pending',
    executed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_batch_dist_date ON batch_distributions(distribution_date DESC);
CREATE INDEX IF NOT EXISTS idx_batch_dist_status ON batch_distributions(status);

-- Distribution reports
CREATE TABLE IF NOT EXISTS distribution_reports (
    date DATE PRIMARY KEY,
    total_distributed NUMERIC(18, 2) NOT NULL,
    recipient_count INTEGER NOT NULL,
    success_rate NUMERIC(5, 4) NOT NULL,
    average_yield NUMERIC(18, 2) NOT NULL,
    total_gas_cost NUMERIC(18, 8) NOT NULL,
    generated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Distribution errors
CREATE TABLE IF NOT EXISTS distribution_errors (
    id SERIAL PRIMARY KEY,
    error_message TEXT NOT NULL,
    stack_trace TEXT,
    occurred_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dist_errors_time ON distribution_errors(occurred_at DESC);

-- ============================================================================
-- 5. Data Retention Policies
-- ============================================================================

SELECT add_retention_policy(
    'notifications',
    INTERVAL '90 days',
    if_not_exists => TRUE
);

SELECT add_retention_policy(
    'treasury_yield_accruals',
    INTERVAL '2 years',
    if_not_exists => TRUE
);

SELECT add_retention_policy(
    'hedge_executions',
    INTERVAL '1 year',
    if_not_exists => TRUE
);

SELECT add_retention_policy(
    'yield_distributions',
    INTERVAL '2 years',
    if_not_exists => TRUE
);

-- ============================================================================
-- 6. Helper Functions
-- ============================================================================

CREATE OR REPLACE FUNCTION get_user_portfolio_value(p_user_id VARCHAR)
RETURNS NUMERIC AS $$
DECLARE
    total_value NUMERIC;
BEGIN
    SELECT
        COALESCE(SUM(principal_usd + total_yield_earned), 0)
    INTO total_value
    FROM treasury_holdings
    WHERE user_id = p_user_id;

    RETURN total_value;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION should_send_notification(
    p_user_id VARCHAR,
    p_type VARCHAR,
    p_priority VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    prefs RECORD;
    current_hour INTEGER;
BEGIN
    SELECT * INTO prefs
    FROM notification_preferences
    WHERE user_id = p_user_id;

    IF NOT FOUND THEN
        RETURN TRUE;
    END IF;

    IF NOT (p_type = ANY(prefs.enabled_types)) AND p_priority != 'critical' THEN
        RETURN FALSE;
    END IF;

    current_hour := EXTRACT(HOUR FROM NOW() AT TIME ZONE 'UTC');
    IF prefs.quiet_hours_start IS NOT NULL AND prefs.quiet_hours_end IS NOT NULL THEN
        IF prefs.quiet_hours_start < prefs.quiet_hours_end THEN
            IF current_hour >= prefs.quiet_hours_start AND current_hour < prefs.quiet_hours_end THEN
                RETURN FALSE;
            END IF;
        ELSE
            IF current_hour >= prefs.quiet_hours_start OR current_hour < prefs.quiet_hours_end THEN
                RETURN FALSE;
            END IF;
        END IF;
    END IF;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- 7. Initial Data
-- ============================================================================

-- Insert default notification preferences for existing users
INSERT INTO notification_preferences (user_id, channels, min_priority, enabled_types)
SELECT
    address,
    ARRAY['email']::TEXT[],
    'medium',
    ARRAY[
        'liquidation_warning',
        'high_risk_position',
        'daily_yield_report',
        'weekly_summary'
    ]::TEXT[]
FROM users
ON CONFLICT (user_id) DO NOTHING;

-- Insert current treasury rates
INSERT INTO treasury_yield_rates (bond_type, effective_date, annual_yield, source)
VALUES
    ('TBILL_3M', CURRENT_DATE, 0.0450, 'Initial'),
    ('TBILL_6M', CURRENT_DATE, 0.0470, 'Initial'),
    ('TNOTE_2Y', CURRENT_DATE, 0.0480, 'Initial'),
    ('TNOTE_5Y', CURRENT_DATE, 0.0500, 'Initial'),
    ('TNOTE_10Y', CURRENT_DATE, 0.0520, 'Initial'),
    ('TBOND_20Y', CURRENT_DATE, 0.0530, 'Initial'),
    ('TBOND_30Y', CURRENT_DATE, 0.0550, 'Initial')
ON CONFLICT (bond_type, effective_date) DO NOTHING;

-- ============================================================================
-- Migration Complete
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '✅ Migration 004 completed successfully!';
    RAISE NOTICE 'Created tables for:';
    RAISE NOTICE '  - Yield Calculation Service (4 tables + 1 view)';
    RAISE NOTICE '  - Notification Service (5 tables)';
    RAISE NOTICE '  - Auto Hedge Executor (3 tables + 1 view)';
    RAISE NOTICE '  - Yield Distribution Service (4 tables)';
    RAISE NOTICE 'Total: 16 new tables, 2 continuous aggregates, 2 helper functions';
    RAISE NOTICE '🚀 TimescaleDB hypertables and retention policies configured!';
END $$;
