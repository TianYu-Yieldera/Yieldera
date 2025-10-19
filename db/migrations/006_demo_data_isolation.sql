-- Migration 006: Demo Data Isolation Flags
-- Purpose: add is_demo flags to core tables so demo data stays separated from production records

ALTER TABLE balances
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE balance_events
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE points
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE points_events
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE badges
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE user_defi_positions
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE defi_transactions
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE stablecoin_positions
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

ALTER TABLE stablecoin_transactions
    ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE;

-- Helpful indexes for filtering demo data
CREATE INDEX IF NOT EXISTS idx_balances_is_demo ON balances(is_demo);
CREATE INDEX IF NOT EXISTS idx_balance_events_is_demo ON balance_events(is_demo);
CREATE INDEX IF NOT EXISTS idx_points_is_demo ON points(is_demo);
CREATE INDEX IF NOT EXISTS idx_points_events_is_demo ON points_events(is_demo);
CREATE INDEX IF NOT EXISTS idx_badges_is_demo ON badges(is_demo);
CREATE INDEX IF NOT EXISTS idx_user_defi_positions_is_demo ON user_defi_positions(is_demo);
CREATE INDEX IF NOT EXISTS idx_defi_transactions_is_demo ON defi_transactions(is_demo);
CREATE INDEX IF NOT EXISTS idx_stablecoin_positions_is_demo ON stablecoin_positions(is_demo);
CREATE INDEX IF NOT EXISTS idx_stablecoin_transactions_is_demo ON stablecoin_transactions(is_demo);

