-- ============================================================================
-- Migration 004: Support for New Services
-- Date: 2025-11-10
-- Description: Add database tables for:
--   1. Yield Calculation Service
--   2. Notification Service
--   3. Auto Hedge Executor
--   4. Yield Distribution Service
-- ============================================================================

-- ============================================================================
-- 1. Yield Calculation Service Tables
-- ============================================================================

-- Treasury yield rates (historical rates from US Treasury API)
CREATE TABLE IF NOT EXISTS treasury_yield_rates (
    id SERIAL PRIMARY KEY,
    bond_type VARCHAR(20) NOT NULL,
    annual_yield NUMERIC(8, 6) NOT NULL,
    effective_date DATE NOT NULL,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(bond_type, effective_date)
);

CREATE INDEX idx_yield_rates_bond_type ON treasury_yield_rates(bond_type);
CREATE INDEX idx_yield_rates_date ON treasury_yield_rates(effective_date DESC);

-- Treasury holdings (user bond positions)
CREATE TABLE IF NOT EXISTS treasury_holdings (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    token_amount NUMERIC(78, 0) NOT NULL,
    principal_usd NUMERIC(18, 2) NOT NULL,
    purchase_date TIMESTAMP NOT NULL,
    last_yield_date TIMESTAMP DEFAULT NOW(),
    total_yield_earned NUMERIC(18, 2) DEFAULT 0,
    compounding_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, bond_type)
);

CREATE INDEX idx_holdings_user ON treasury_holdings(user_id);
CREATE INDEX idx_holdings_bond_type ON treasury_holdings(bond_type);
CREATE INDEX idx_holdings_updated ON treasury_holdings(updated_at DESC);

-- Daily yield accruals (time-series data)
CREATE TABLE IF NOT EXISTS treasury_yield_accruals (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    principal_amount NUMERIC(18, 2) NOT NULL,
    daily_yield_rate NUMERIC(12, 9) NOT NULL,
    yield_amount_usd NUMERIC(18, 2) NOT NULL,
    cumulative_yield NUMERIC(18, 2) NOT NULL,
    compounded BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Convert to TimescaleDB hypertable for efficient time-series queries
SELECT create_hypertable(
    'treasury_yield_accruals',
    'date',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 month'
);

CREATE INDEX idx_accruals_user_date ON treasury_yield_accruals(user_id, date DESC);
CREATE INDEX idx_accruals_bond_type ON treasury_yield_accruals(bond_type);

-- Continuous aggregate for monthly yield summaries
CREATE MATERIALIZED VIEW IF NOT EXISTS monthly_yield_summary
WITH (timescaledb.continuous) AS
SELECT
    user_id,
    bond_type,
    time_bucket('1 month', date) AS month,
    SUM(yield_amount_usd) as total_yield,
    AVG(daily_yield_rate) as avg_rate,
    COUNT(*) as days_count
FROM treasury_yield_accruals
GROUP BY user_id, bond_type, month
WITH NO DATA;

-- Refresh policy for continuous aggregate
SELECT add_continuous_aggregate_policy(
    'monthly_yield_summary',
    start_offset => INTERVAL '3 months',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Tax reports
CREATE TABLE IF NOT EXISTS treasury_tax_reports (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    tax_year INTEGER NOT NULL,
    total_interest NUMERIC(18, 2) NOT NULL,
    report_data JSONB NOT NULL,
    generated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, tax_year)
);

CREATE INDEX idx_tax_reports_user ON treasury_tax_reports(user_id);
CREATE INDEX idx_tax_reports_year ON treasury_tax_reports(tax_year DESC);

-- ============================================================================
-- 2. Notification Service Tables
-- ============================================================================

-- User notification preferences
CREATE TABLE IF NOT EXISTS notification_preferences (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL UNIQUE,
    channels TEXT[] NOT NULL DEFAULT '{"email"}',
    min_priority VARCHAR(20) DEFAULT 'medium',
    quiet_hours_start INTEGER,
    quiet_hours_end INTEGER,
    enabled_types TEXT[] NOT NULL DEFAULT '{}',
    frequency VARCHAR(20) DEFAULT 'realtime',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notif_prefs_user ON notification_preferences(user_id);

-- Notification history
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    type VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    channels TEXT[] NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW(),
    read_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'sent',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable(
    'notifications',
    'sent_at',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 week'
);

CREATE INDEX idx_notifications_user_sent ON notifications(user_id, sent_at DESC);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_priority ON notifications(priority);
CREATE INDEX idx_notifications_unread ON notifications(user_id, read_at) WHERE read_at IS NULL;

-- Notification queue (for delayed/scheduled notifications)
CREATE TABLE IF NOT EXISTS notification_queue (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    notification_data JSONB NOT NULL,
    scheduled_for TIMESTAMP NOT NULL,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notif_queue_scheduled ON notification_queue(scheduled_for) WHERE NOT processed;
CREATE INDEX idx_notif_queue_user ON notification_queue(user_id);

-- User contact information
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(66) PRIMARY KEY,
    email VARCHAR(255),
    phone_number VARCHAR(20),
    telegram_chat_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
) ON CONFLICT (id) DO NOTHING;

-- User devices (for push notifications)
CREATE TABLE IF NOT EXISTS user_devices (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    fcm_token TEXT NOT NULL,
    device_type VARCHAR(20),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    last_seen TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_devices_user ON user_devices(user_id);
CREATE INDEX idx_devices_active ON user_devices(user_id, active) WHERE active = true;

-- ============================================================================
-- 3. Auto Hedge Executor Tables
-- ============================================================================

-- User settings for auto-hedge
CREATE TABLE IF NOT EXISTS user_settings (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL UNIQUE,
    auto_hedge_enabled BOOLEAN DEFAULT false,
    max_hedge_amount NUMERIC(18, 2) DEFAULT 10000,
    min_health_factor NUMERIC(8, 4) DEFAULT 1.5,
    target_health_factor NUMERIC(8, 4) DEFAULT 2.0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_settings_user ON user_settings(user_id);
CREATE INDEX idx_user_settings_enabled ON user_settings(auto_hedge_enabled) WHERE auto_hedge_enabled = true;

-- Hedge execution records
CREATE TABLE IF NOT EXISTS hedge_executions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    risk_signal JSONB NOT NULL,
    strategy VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    simulated_outcome JSONB,
    tx_hash VARCHAR(66),
    gas_used NUMERIC(78, 0),
    executed_at TIMESTAMP,
    completed_at TIMESTAMP,
    error TEXT,
    resulting_health_factor NUMERIC(8, 4),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable(
    'hedge_executions',
    'created_at',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 month'
);

CREATE INDEX idx_hedge_exec_user ON hedge_executions(user_id, created_at DESC);
CREATE INDEX idx_hedge_exec_status ON hedge_executions(status);
CREATE INDEX idx_hedge_exec_strategy ON hedge_executions(strategy);

-- Continuous aggregate for hedge statistics
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_hedge_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', created_at) AS day,
    COUNT(*) as total_executions,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
    AVG(CASE WHEN resulting_health_factor IS NOT NULL THEN resulting_health_factor END) as avg_resulting_hf,
    SUM(CASE WHEN gas_used IS NOT NULL THEN gas_used::numeric END) as total_gas_used
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

-- Individual yield distributions
CREATE TABLE IF NOT EXISTS yield_distributions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(66) NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    yield_amount NUMERIC(18, 2) NOT NULL,
    token_amount NUMERIC(78, 0) NOT NULL,
    distribution_date DATE NOT NULL,
    tx_hash VARCHAR(66),
    status VARCHAR(20) DEFAULT 'pending',
    retry_count INTEGER DEFAULT 0,
    error TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable(
    'yield_distributions',
    'distribution_date',
    if_not_exists => TRUE,
    chunk_time_interval => INTERVAL '1 month'
);

CREATE INDEX idx_yield_dist_user ON yield_distributions(user_id, distribution_date DESC);
CREATE INDEX idx_yield_dist_status ON yield_distributions(status);
CREATE INDEX idx_yield_dist_pending ON yield_distributions(status, retry_count) WHERE status = 'pending' OR status = 'failed';

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
    executed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_batch_dist_date ON batch_distributions(distribution_date DESC);
CREATE INDEX idx_batch_dist_status ON batch_distributions(status);

-- Distribution reports (daily summaries)
CREATE TABLE IF NOT EXISTS distribution_reports (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    total_distributed NUMERIC(18, 2) NOT NULL,
    recipient_count INTEGER NOT NULL,
    success_rate NUMERIC(5, 4) NOT NULL,
    average_yield NUMERIC(18, 2) NOT NULL,
    total_gas_cost NUMERIC(18, 8) NOT NULL,
    generated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_dist_reports_date ON distribution_reports(date DESC);

-- Distribution errors (for debugging)
CREATE TABLE IF NOT EXISTS distribution_errors (
    id SERIAL PRIMARY KEY,
    error_message TEXT NOT NULL,
    stack_trace TEXT,
    occurred_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_dist_errors_time ON distribution_errors(occurred_at DESC);

-- ============================================================================
-- 5. Data Retention Policies
-- ============================================================================

-- Auto-delete old notifications after 90 days
SELECT add_retention_policy(
    'notifications',
    INTERVAL '90 days',
    if_not_exists => TRUE
);

-- Auto-delete old yield accruals after 2 years (keep aggregates)
SELECT add_retention_policy(
    'treasury_yield_accruals',
    INTERVAL '2 years',
    if_not_exists => TRUE
);

-- Auto-delete old hedge executions after 1 year
SELECT add_retention_policy(
    'hedge_executions',
    INTERVAL '1 year',
    if_not_exists => TRUE
);

-- Auto-delete old distribution records after 2 years
SELECT add_retention_policy(
    'yield_distributions',
    INTERVAL '2 years',
    if_not_exists => TRUE
);

-- ============================================================================
-- 6. Helper Functions
-- ============================================================================

-- Function to get user's total portfolio value
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

-- Function to check if user should receive notification
CREATE OR REPLACE FUNCTION should_send_notification(
    p_user_id VARCHAR,
    p_type VARCHAR,
    p_priority VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    prefs RECORD;
    current_hour INTEGER;
BEGIN
    -- Get user preferences
    SELECT * INTO prefs
    FROM notification_preferences
    WHERE user_id = p_user_id;

    -- If no preferences, use defaults
    IF NOT FOUND THEN
        RETURN TRUE;
    END IF;

    -- Check if notification type is enabled
    IF NOT (p_type = ANY(prefs.enabled_types)) AND p_priority != 'critical' THEN
        RETURN FALSE;
    END IF;

    -- Check quiet hours
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
    id,
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

-- Insert current treasury rates (fallback values)
INSERT INTO treasury_yield_rates (bond_type, annual_yield, effective_date, source)
VALUES
    ('TBILL_3M', 0.0450, CURRENT_DATE, 'Initial'),
    ('TBILL_6M', 0.0470, CURRENT_DATE, 'Initial'),
    ('TNOTE_2Y', 0.0480, CURRENT_DATE, 'Initial'),
    ('TNOTE_5Y', 0.0500, CURRENT_DATE, 'Initial'),
    ('TNOTE_10Y', 0.0520, CURRENT_DATE, 'Initial'),
    ('TBOND_20Y', 0.0530, CURRENT_DATE, 'Initial'),
    ('TBOND_30Y', 0.0550, CURRENT_DATE, 'Initial')
ON CONFLICT (bond_type, effective_date) DO NOTHING;

-- ============================================================================
-- 8. Grants and Permissions
-- ============================================================================

-- Grant permissions to application user
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO yieldera_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO yieldera_app;

-- ============================================================================
-- Migration Complete
-- ============================================================================

-- Log migration completion
DO $$
BEGIN
    RAISE NOTICE 'Migration 004 completed successfully!';
    RAISE NOTICE 'Created tables for:';
    RAISE NOTICE '  - Yield Calculation Service (4 tables + 1 view)';
    RAISE NOTICE '  - Notification Service (5 tables)';
    RAISE NOTICE '  - Auto Hedge Executor (3 tables + 1 view)';
    RAISE NOTICE '  - Yield Distribution Service (4 tables)';
    RAISE NOTICE 'Total: 16 new tables, 2 continuous aggregates, 2 helper functions';
END $$;
