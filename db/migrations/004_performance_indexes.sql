-- Migration 004: Performance Indexes
-- Created: 2025-10-16
-- Purpose: Add missing indexes for query performance optimization

-- Indexes for balances table
CREATE INDEX IF NOT EXISTS idx_balances_user ON balances(user_address);

-- Indexes for points table
CREATE INDEX IF NOT EXISTS idx_points_user ON points(user_address);
CREATE INDEX IF NOT EXISTS idx_points_value ON points(CAST(points AS NUMERIC) DESC);

-- Indexes for badges table
CREATE INDEX IF NOT EXISTS idx_badges_user ON badges(user_address);
CREATE INDEX IF NOT EXISTS idx_badges_created ON badges(created_at DESC);

-- Composite index for user DeFi positions
CREATE INDEX IF NOT EXISTS idx_user_defi_pos_composite ON user_defi_positions(user_address, pool_id);
CREATE INDEX IF NOT EXISTS idx_user_defi_pos_updated ON user_defi_positions(last_updated DESC);

-- Index for stablecoin positions
CREATE INDEX IF NOT EXISTS idx_stablecoin_user ON stablecoin_positions(user_address);
CREATE INDEX IF NOT EXISTS idx_stablecoin_health ON stablecoin_positions(health_status);

-- Indexes for balance_events (if table exists)
CREATE INDEX IF NOT EXISTS idx_balance_events_txhash ON balance_events(tx_hash);
CREATE INDEX IF NOT EXISTS idx_balance_events_user ON balance_events(user_address);
CREATE INDEX IF NOT EXISTS idx_balance_events_time ON balance_events(timestamp DESC);

-- Verify indexes
SELECT
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

-- Performance test queries
EXPLAIN ANALYZE SELECT * FROM balances WHERE user_address = '0x3c07226a3F1488320426eb5fE9976f72E5712346';
EXPLAIN ANALYZE SELECT * FROM points ORDER BY CAST(points AS NUMERIC) DESC LIMIT 20;
EXPLAIN ANALYZE SELECT * FROM user_defi_positions WHERE user_address = '0x3c07226a3F1488320426eb5fE9976f72E5712346';
