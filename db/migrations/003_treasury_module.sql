-- Migration 003: Treasury Module
-- Description: Replace generic RWA module with US Treasury tokenization system
-- Author: System
-- Date: 2025-11-04

-- =============================================================
-- STEP 1: DROP OLD RWA TABLES
-- =============================================================

-- Drop old RWA tables (if they exist)
DROP TABLE IF EXISTS l2_rwa_proposals CASCADE;
DROP TABLE IF EXISTS l2_rwa_listings CASCADE;
DROP TABLE IF EXISTS l2_rwa_holdings CASCADE;
DROP TABLE IF EXISTS l2_rwa_assets CASCADE;

-- Drop points and rewards tables (removed from scope)
DROP TABLE IF EXISTS points CASCADE;
DROP TABLE IF EXISTS badges CASCADE;
DROP TABLE IF EXISTS airdrop_campaigns CASCADE;
DROP TABLE IF EXISTS airdrop_allocations CASCADE;

-- =============================================================
-- STEP 2: CREATE TREASURY TABLES
-- =============================================================

-- Table 1: Treasury Assets
-- Stores information about tokenized US Treasury securities
CREATE TABLE treasury_assets (
    asset_id            BIGSERIAL PRIMARY KEY,
    treasury_type       TEXT NOT NULL,              -- 'T-BILL', 'T-NOTE', 'T-BOND'
    maturity_term       TEXT NOT NULL,              -- '4W', '13W', '26W', '2Y', '10Y', '30Y', etc.
    cusip               TEXT UNIQUE NOT NULL,       -- CUSIP identifier (unique)
    issue_date          DATE NOT NULL,
    maturity_date       DATE NOT NULL,
    face_value          NUMERIC(78, 18) NOT NULL,   -- Face value in USD
    coupon_rate         NUMERIC(5, 4) NOT NULL,     -- Annual coupon rate (0.0425 = 4.25%)
    current_price       NUMERIC(78, 18),            -- Current market price
    current_yield       NUMERIC(5, 4),              -- Current yield
    tokens_issued       NUMERIC(78, 18) DEFAULT 0,  -- Total tokens issued
    tokens_outstanding  NUMERIC(78, 18) DEFAULT 0,  -- Tokens in circulation
    token_address       TEXT,                       -- L2 ERC20 token contract address
    status              TEXT DEFAULT 'active',      -- 'active', 'matured', 'suspended'
    last_price_update   TIMESTAMP,                  -- Last price update timestamp
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT valid_treasury_type CHECK (treasury_type IN ('T-BILL', 'T-NOTE', 'T-BOND')),
    CONSTRAINT valid_status CHECK (status IN ('active', 'matured', 'suspended')),
    CONSTRAINT valid_face_value CHECK (face_value > 0),
    CONSTRAINT valid_coupon_rate CHECK (coupon_rate >= 0 AND coupon_rate <= 1),
    CONSTRAINT valid_dates CHECK (maturity_date > issue_date)
);

-- Indexes for treasury_assets
CREATE INDEX idx_treasury_type_status ON treasury_assets(treasury_type, status);
CREATE INDEX idx_treasury_maturity ON treasury_assets(maturity_date);
CREATE INDEX idx_treasury_cusip ON treasury_assets(cusip);
CREATE INDEX idx_treasury_token_address ON treasury_assets(token_address);

-- Table 2: Treasury Holdings
-- Tracks user holdings of treasury tokens
CREATE TABLE treasury_holdings (
    id                  BIGSERIAL PRIMARY KEY,
    user_address        TEXT NOT NULL,
    asset_id            BIGINT REFERENCES treasury_assets(asset_id) ON DELETE CASCADE,
    tokens_held         NUMERIC(78, 18) NOT NULL DEFAULT 0,
    avg_purchase_price  NUMERIC(78, 18),            -- Average purchase price
    total_invested      NUMERIC(78, 18),            -- Total USD invested
    current_value       NUMERIC(78, 18),            -- Current market value
    unrealized_gain     NUMERIC(78, 18),            -- Unrealized profit/loss
    accrued_interest    NUMERIC(78, 18) DEFAULT 0,  -- Accrued interest
    last_updated        TIMESTAMP DEFAULT NOW(),
    created_at          TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT valid_tokens_held CHECK (tokens_held >= 0),
    UNIQUE(user_address, asset_id)
);

-- Indexes for treasury_holdings
CREATE INDEX idx_treasury_holdings_user ON treasury_holdings(user_address);
CREATE INDEX idx_treasury_holdings_asset ON treasury_holdings(asset_id);
CREATE INDEX idx_treasury_holdings_user_asset ON treasury_holdings(user_address, asset_id);

-- Table 3: Treasury Market Orders
-- Secondary market buy/sell orders
CREATE TABLE treasury_market_orders (
    order_id            BIGSERIAL PRIMARY KEY,
    asset_id            BIGINT REFERENCES treasury_assets(asset_id) ON DELETE CASCADE,
    order_type          TEXT NOT NULL,              -- 'BUY', 'SELL'
    user_address        TEXT NOT NULL,
    token_amount        NUMERIC(78, 18) NOT NULL,   -- Number of tokens
    price_per_token     NUMERIC(78, 18) NOT NULL,   -- Price per token in USD
    total_value         NUMERIC(78, 18) NOT NULL,   -- Total order value
    filled_amount       NUMERIC(78, 18) DEFAULT 0,  -- Amount filled
    status              TEXT DEFAULT 'open',        -- 'open', 'partial', 'filled', 'cancelled'
    tx_hash             TEXT,                       -- Transaction hash
    created_at          TIMESTAMP DEFAULT NOW(),
    expires_at          TIMESTAMP,
    filled_at           TIMESTAMP,
    cancelled_at        TIMESTAMP,

    -- Constraints
    CONSTRAINT valid_order_type CHECK (order_type IN ('BUY', 'SELL')),
    CONSTRAINT valid_order_status CHECK (status IN ('open', 'partial', 'filled', 'cancelled')),
    CONSTRAINT valid_token_amount CHECK (token_amount > 0),
    CONSTRAINT valid_price CHECK (price_per_token > 0),
    CONSTRAINT valid_filled CHECK (filled_amount >= 0 AND filled_amount <= token_amount)
);

-- Indexes for treasury_market_orders
CREATE INDEX idx_treasury_orders_asset ON treasury_market_orders(asset_id);
CREATE INDEX idx_treasury_orders_user ON treasury_market_orders(user_address);
CREATE INDEX idx_treasury_orders_status ON treasury_market_orders(status);
CREATE INDEX idx_treasury_orders_asset_status ON treasury_market_orders(asset_id, status, order_type);
CREATE INDEX idx_treasury_orders_created ON treasury_market_orders(created_at DESC);

-- Table 4: Treasury Yield Distributions
-- Records of interest/coupon payments
CREATE TABLE treasury_yield_distributions (
    id                  BIGSERIAL PRIMARY KEY,
    asset_id            BIGINT REFERENCES treasury_assets(asset_id) ON DELETE CASCADE,
    distribution_date   DATE NOT NULL,
    distribution_type   TEXT NOT NULL,              -- 'COUPON', 'MATURITY'
    total_yield         NUMERIC(78, 18) NOT NULL,   -- Total yield amount (USD)
    yield_per_token     NUMERIC(78, 18) NOT NULL,   -- Yield per token
    recipients_count    INT DEFAULT 0,
    total_distributed   NUMERIC(78, 18) DEFAULT 0,
    status              TEXT DEFAULT 'pending',     -- 'pending', 'completed'
    tx_hash             TEXT,
    created_at          TIMESTAMP DEFAULT NOW(),
    distributed_at      TIMESTAMP,

    -- Constraints
    CONSTRAINT valid_distribution_type CHECK (distribution_type IN ('COUPON', 'MATURITY')),
    CONSTRAINT valid_distribution_status CHECK (status IN ('pending', 'completed')),
    CONSTRAINT valid_yield_amount CHECK (total_yield > 0)
);

-- Indexes for treasury_yield_distributions
CREATE INDEX idx_treasury_yield_asset ON treasury_yield_distributions(asset_id);
CREATE INDEX idx_treasury_yield_date ON treasury_yield_distributions(distribution_date DESC);
CREATE INDEX idx_treasury_yield_status ON treasury_yield_distributions(status);

-- Table 5: Treasury Price History
-- Historical price and yield data
CREATE TABLE treasury_price_history (
    id                  BIGSERIAL PRIMARY KEY,
    asset_id            BIGINT REFERENCES treasury_assets(asset_id) ON DELETE CASCADE,
    price               NUMERIC(78, 18) NOT NULL,   -- Price in USD
    yield               NUMERIC(5, 4) NOT NULL,     -- Yield rate
    source              TEXT,                        -- 'chainlink', 'manual', 'api', etc.
    timestamp           TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT valid_price CHECK (price > 0),
    CONSTRAINT valid_yield CHECK (yield >= 0 AND yield <= 1)
);

-- Indexes for treasury_price_history
CREATE INDEX idx_treasury_price_asset ON treasury_price_history(asset_id);
CREATE INDEX idx_treasury_price_timestamp ON treasury_price_history(timestamp DESC);
CREATE INDEX idx_treasury_price_asset_time ON treasury_price_history(asset_id, timestamp DESC);

-- Table 6: Treasury Trades
-- Record of executed trades
CREATE TABLE treasury_trades (
    trade_id            BIGSERIAL PRIMARY KEY,
    asset_id            BIGINT REFERENCES treasury_assets(asset_id) ON DELETE CASCADE,
    buy_order_id        BIGINT REFERENCES treasury_market_orders(order_id),
    sell_order_id       BIGINT REFERENCES treasury_market_orders(order_id),
    buyer_address       TEXT NOT NULL,
    seller_address      TEXT NOT NULL,
    token_amount        NUMERIC(78, 18) NOT NULL,
    price_per_token     NUMERIC(78, 18) NOT NULL,
    total_value         NUMERIC(78, 18) NOT NULL,
    fee_amount          NUMERIC(78, 18) DEFAULT 0,
    tx_hash             TEXT,
    executed_at         TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT valid_trade_amount CHECK (token_amount > 0),
    CONSTRAINT valid_trade_price CHECK (price_per_token > 0)
);

-- Indexes for treasury_trades
CREATE INDEX idx_treasury_trades_asset ON treasury_trades(asset_id);
CREATE INDEX idx_treasury_trades_buyer ON treasury_trades(buyer_address);
CREATE INDEX idx_treasury_trades_seller ON treasury_trades(seller_address);
CREATE INDEX idx_treasury_trades_time ON treasury_trades(executed_at DESC);

-- =============================================================
-- STEP 3: CREATE VIEWS FOR ANALYTICS
-- =============================================================

-- View: Treasury Asset Summary
CREATE OR REPLACE VIEW v_treasury_asset_summary AS
SELECT
    a.asset_id,
    a.treasury_type,
    a.maturity_term,
    a.cusip,
    a.current_price,
    a.current_yield,
    a.tokens_outstanding,
    (a.tokens_outstanding * COALESCE(a.current_price, a.face_value)) AS total_market_value,
    a.maturity_date,
    a.status,
    COUNT(DISTINCT h.user_address) AS unique_holders,
    COUNT(DISTINCT o.order_id) AS active_orders
FROM treasury_assets a
LEFT JOIN treasury_holdings h ON a.asset_id = h.asset_id AND h.tokens_held > 0
LEFT JOIN treasury_market_orders o ON a.asset_id = o.asset_id AND o.status IN ('open', 'partial')
GROUP BY a.asset_id;

-- View: User Treasury Portfolio
CREATE OR REPLACE VIEW v_user_treasury_portfolio AS
SELECT
    h.user_address,
    h.asset_id,
    a.treasury_type,
    a.maturity_term,
    a.cusip,
    h.tokens_held,
    h.avg_purchase_price,
    h.total_invested,
    a.current_price,
    (h.tokens_held * COALESCE(a.current_price, a.face_value)) AS current_value,
    ((h.tokens_held * COALESCE(a.current_price, a.face_value)) - COALESCE(h.total_invested, 0)) AS unrealized_gain_loss,
    h.accrued_interest,
    a.maturity_date
FROM treasury_holdings h
JOIN treasury_assets a ON h.asset_id = a.asset_id
WHERE h.tokens_held > 0;

-- View: Market Order Book
CREATE OR REPLACE VIEW v_treasury_order_book AS
SELECT
    o.order_id,
    o.asset_id,
    a.cusip,
    a.treasury_type,
    o.order_type,
    o.user_address,
    o.token_amount,
    o.filled_amount,
    (o.token_amount - o.filled_amount) AS remaining_amount,
    o.price_per_token,
    o.total_value,
    o.status,
    o.created_at,
    o.expires_at
FROM treasury_market_orders o
JOIN treasury_assets a ON o.asset_id = a.asset_id
WHERE o.status IN ('open', 'partial')
  AND o.expires_at > NOW()
ORDER BY o.price_per_token DESC, o.created_at ASC;

-- =============================================================
-- STEP 4: CREATE FUNCTIONS
-- =============================================================

-- Function: Update treasury holding after trade
CREATE OR REPLACE FUNCTION update_treasury_holding()
RETURNS TRIGGER AS $$
BEGIN
    -- This is a placeholder for automatic holding updates
    -- Will be called by triggers when trades are executed
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function: Calculate total portfolio value for user
CREATE OR REPLACE FUNCTION get_user_treasury_total_value(p_user_address TEXT)
RETURNS NUMERIC AS $$
DECLARE
    total_value NUMERIC;
BEGIN
    SELECT COALESCE(SUM(h.tokens_held * COALESCE(a.current_price, a.face_value)), 0)
    INTO total_value
    FROM treasury_holdings h
    JOIN treasury_assets a ON h.asset_id = a.asset_id
    WHERE h.user_address = p_user_address
      AND h.tokens_held > 0;

    RETURN total_value;
END;
$$ LANGUAGE plpgsql;

-- Function: Get active orders for asset
CREATE OR REPLACE FUNCTION get_active_treasury_orders(
    p_asset_id BIGINT,
    p_order_type TEXT DEFAULT NULL
)
RETURNS TABLE (
    order_id BIGINT,
    order_type TEXT,
    user_address TEXT,
    token_amount NUMERIC,
    price_per_token NUMERIC,
    filled_amount NUMERIC,
    remaining_amount NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        o.order_id,
        o.order_type,
        o.user_address,
        o.token_amount,
        o.price_per_token,
        o.filled_amount,
        (o.token_amount - o.filled_amount) AS remaining_amount
    FROM treasury_market_orders o
    WHERE o.asset_id = p_asset_id
      AND o.status IN ('open', 'partial')
      AND o.expires_at > NOW()
      AND (p_order_type IS NULL OR o.order_type = p_order_type)
    ORDER BY
        CASE WHEN o.order_type = 'BUY' THEN o.price_per_token END DESC,
        CASE WHEN o.order_type = 'SELL' THEN o.price_per_token END ASC,
        o.created_at ASC;
END;
$$ LANGUAGE plpgsql;

-- =============================================================
-- STEP 5: INSERT SAMPLE DATA (FOR TESTING)
-- =============================================================

-- Sample Treasury Assets
-- INSERT INTO treasury_assets (
--     treasury_type, maturity_term, cusip, issue_date, maturity_date,
--     face_value, coupon_rate, current_price, current_yield, status
-- ) VALUES
-- ('T-BILL', '13W', '912796YZ1', '2024-10-01', '2025-01-01', 1000, 0.0525, 980, 0.0540, 'active'),
-- ('T-NOTE', '2Y', '91282CHX6', '2024-01-01', '2026-01-01', 1000, 0.0450, 985, 0.0465, 'active'),
-- ('T-NOTE', '10Y', '912828YK4', '2024-01-01', '2034-01-01', 1000, 0.0425, 950, 0.0475, 'active'),
-- ('T-BOND', '30Y', '912810TT4', '2024-01-01', '2054-01-01', 1000, 0.0400, 920, 0.0450, 'active');

-- =============================================================
-- STEP 6: GRANT PERMISSIONS
-- =============================================================

-- Grant permissions to application user (adjust username as needed)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO loyalty_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO loyalty_app;
-- GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO loyalty_app;

-- =============================================================
-- MIGRATION COMPLETE
-- =============================================================

-- Add migration record
-- INSERT INTO schema_migrations (version, description, applied_at)
-- VALUES (3, 'Treasury Module - Replace RWA with US Treasury tokenization', NOW());

COMMENT ON TABLE treasury_assets IS 'US Treasury securities tokenization metadata';
COMMENT ON TABLE treasury_holdings IS 'User holdings of treasury tokens';
COMMENT ON TABLE treasury_market_orders IS 'Secondary market buy/sell orders';
COMMENT ON TABLE treasury_yield_distributions IS 'Coupon payment and yield distribution records';
COMMENT ON TABLE treasury_price_history IS 'Historical price and yield data';
COMMENT ON TABLE treasury_trades IS 'Executed trade records';
