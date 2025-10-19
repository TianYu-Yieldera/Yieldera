-- Migration: Add demo mode support for hackathon users
-- This allows new users to try the system with test tokens without real blockchain transactions

-- Add demo mode columns to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_demo BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS demo_expires_at TIMESTAMP;

-- Create index for demo users query performance
CREATE INDEX IF NOT EXISTS idx_users_demo ON users(is_demo, demo_expires_at) WHERE is_demo = TRUE;

-- Create demo transaction flag for stablecoin positions
ALTER TABLE stablecoin_positions
ADD COLUMN IF NOT EXISTS is_demo_transaction BOOLEAN DEFAULT FALSE;

-- Add comments for documentation
COMMENT ON COLUMN users.is_demo IS 'Flag indicating if user is in demo mode with test tokens';
COMMENT ON COLUMN users.demo_expires_at IS 'Expiration timestamp for demo mode (typically 24 hours)';
COMMENT ON COLUMN stablecoin_positions.is_demo_transaction IS 'Flag indicating if position was created in demo mode';

-- Create view for active demo users
CREATE OR REPLACE VIEW active_demo_users AS
SELECT
    user_id,
    wallet_address,
    username,
    points,
    demo_expires_at,
    created_at
FROM users
WHERE is_demo = TRUE
  AND demo_expires_at > NOW()
ORDER BY created_at DESC;

-- Create function to clean up expired demo users
CREATE OR REPLACE FUNCTION cleanup_expired_demo_users()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Delete expired demo users (older than 7 days after expiration)
    DELETE FROM users
    WHERE is_demo = TRUE
      AND demo_expires_at < NOW() - INTERVAL '7 days';

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Grant permissions
GRANT SELECT ON active_demo_users TO PUBLIC;
