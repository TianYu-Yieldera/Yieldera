-- ============================================================================
-- Test Script for Migration 008: Multi-Chain Support
-- ============================================================================
-- Purpose: Verify that migration 008 was applied correctly
-- Usage: psql -U your_user -d yieldera -f db/migrations/test_008_migration.sql
-- ============================================================================

\echo '========================================='
\echo 'Testing Migration 008: Multi-Chain Support'
\echo '========================================='
\echo ''

-- ============================================================================
-- TEST 1: Check if chain_id columns exist
-- ============================================================================
\echo 'TEST 1: Checking if chain_id columns were added...'

DO $$
DECLARE
    column_count INTEGER;
BEGIN
    -- Check vault_positions
    SELECT COUNT(*) INTO column_count
    FROM information_schema.columns
    WHERE table_name = 'vault_positions' AND column_name = 'chain_id';

    IF column_count = 1 THEN
        RAISE NOTICE '✅ vault_positions.chain_id exists';
    ELSE
        RAISE EXCEPTION '❌ vault_positions.chain_id NOT FOUND!';
    END IF;

    -- Check treasury_holdings
    SELECT COUNT(*) INTO column_count
    FROM information_schema.columns
    WHERE table_name = 'treasury_holdings' AND column_name = 'chain_id';

    IF column_count = 1 THEN
        RAISE NOTICE '✅ treasury_holdings.chain_id exists';
    ELSE
        RAISE EXCEPTION '❌ treasury_holdings.chain_id NOT FOUND!';
    END IF;

    -- Check price_history
    SELECT COUNT(*) INTO column_count
    FROM information_schema.columns
    WHERE table_name = 'price_history' AND column_name = 'chain_id';

    IF column_count = 1 THEN
        RAISE NOTICE '✅ price_history.chain_id exists';
    ELSE
        RAISE WARNING '⚠️ price_history.chain_id NOT FOUND (may not exist yet)';
    END IF;
END $$;

\echo ''

-- ============================================================================
-- TEST 2: Check if new tables were created
-- ============================================================================
\echo 'TEST 2: Checking if new tables were created...'

DO $$
DECLARE
    table_count INTEGER;
BEGIN
    -- Check gmx_performance_comparison
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_name = 'gmx_performance_comparison';

    IF table_count = 1 THEN
        RAISE NOTICE '✅ gmx_performance_comparison table exists';
    ELSE
        RAISE EXCEPTION '❌ gmx_performance_comparison table NOT FOUND!';
    END IF;

    -- Check liquidation_predictions
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_name = 'liquidation_predictions';

    IF table_count = 1 THEN
        RAISE NOTICE '✅ liquidation_predictions table exists';
    ELSE
        RAISE EXCEPTION '❌ liquidation_predictions table NOT FOUND!';
    END IF;

    -- Check supported_chains
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_name = 'supported_chains';

    IF table_count = 1 THEN
        RAISE NOTICE '✅ supported_chains table exists';
    ELSE
        RAISE EXCEPTION '❌ supported_chains table NOT FOUND!';
    END IF;

    -- Check aerodrome_swaps
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_name = 'aerodrome_swaps';

    IF table_count = 1 THEN
        RAISE NOTICE '✅ aerodrome_swaps table exists';
    ELSE
        RAISE EXCEPTION '❌ aerodrome_swaps table NOT FOUND!';
    END IF;

    -- Check base_pay_transactions
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_name = 'base_pay_transactions';

    IF table_count = 1 THEN
        RAISE NOTICE '✅ base_pay_transactions table exists';
    ELSE
        RAISE EXCEPTION '❌ base_pay_transactions table NOT FOUND!';
    END IF;
END $$;

\echo ''

-- ============================================================================
-- TEST 3: Check if views were created
-- ============================================================================
\echo 'TEST 3: Checking if cross-chain views were created...'

DO $$
DECLARE
    view_count INTEGER;
BEGIN
    -- Check cross_chain_user_portfolio
    SELECT COUNT(*) INTO view_count
    FROM information_schema.views
    WHERE table_name = 'cross_chain_user_portfolio';

    IF view_count = 1 THEN
        RAISE NOTICE '✅ cross_chain_user_portfolio view exists';
    ELSE
        RAISE EXCEPTION '❌ cross_chain_user_portfolio view NOT FOUND!';
    END IF;

    -- Check cross_chain_total_value
    SELECT COUNT(*) INTO view_count
    FROM information_schema.views
    WHERE table_name = 'cross_chain_total_value';

    IF view_count = 1 THEN
        RAISE NOTICE '✅ cross_chain_total_value view exists';
    ELSE
        RAISE EXCEPTION '❌ cross_chain_total_value view NOT FOUND!';
    END IF;

    -- Check gmx_performance_summary
    SELECT COUNT(*) INTO view_count
    FROM information_schema.views
    WHERE table_name = 'gmx_performance_summary';

    IF view_count = 1 THEN
        RAISE NOTICE '✅ gmx_performance_summary view exists';
    ELSE
        RAISE EXCEPTION '❌ gmx_performance_summary view NOT FOUND!';
    END IF;
END $$;

\echo ''

-- ============================================================================
-- TEST 4: Check if helper function exists
-- ============================================================================
\echo 'TEST 4: Checking if helper function was created...'

DO $$
DECLARE
    function_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO function_count
    FROM pg_proc p
    JOIN pg_namespace n ON p.pronamespace = n.oid
    WHERE n.nspname = 'public' AND p.proname = 'get_cross_chain_risk_profile';

    IF function_count >= 1 THEN
        RAISE NOTICE '✅ get_cross_chain_risk_profile function exists';
    ELSE
        RAISE EXCEPTION '❌ get_cross_chain_risk_profile function NOT FOUND!';
    END IF;
END $$;

\echo ''

-- ============================================================================
-- TEST 5: Check supported chains reference data
-- ============================================================================
\echo 'TEST 5: Checking supported chains reference data...'

SELECT
    chain_id,
    chain_name,
    network_type,
    features,
    is_active
FROM supported_chains
ORDER BY chain_id;

\echo ''
\echo 'Expected: 4 rows (Arbitrum Sepolia, Arbitrum One, Base Sepolia, Base Mainnet)'
\echo ''

-- ============================================================================
-- TEST 6: Check data distribution
-- ============================================================================
\echo 'TEST 6: Checking chain_id data distribution...'

\echo 'Vault positions by chain:'
SELECT
    chain_id,
    CASE
        WHEN chain_id = 421614 THEN 'Arbitrum Sepolia'
        WHEN chain_id = 42161 THEN 'Arbitrum One'
        WHEN chain_id = 84532 THEN 'Base Sepolia'
        WHEN chain_id = 8453 THEN 'Base Mainnet'
        ELSE 'Unknown'
    END as chain_name,
    COUNT(*) as position_count
FROM vault_positions
GROUP BY chain_id
ORDER BY chain_id;

\echo ''
\echo 'Treasury holdings by chain:'
SELECT
    chain_id,
    CASE
        WHEN chain_id = 421614 THEN 'Arbitrum Sepolia'
        WHEN chain_id = 42161 THEN 'Arbitrum One'
        WHEN chain_id = 84532 THEN 'Base Sepolia'
        WHEN chain_id = 8453 THEN 'Base Mainnet'
        ELSE 'Unknown'
    END as chain_name,
    COUNT(*) as holding_count
FROM treasury_holdings
GROUP BY chain_id
ORDER BY chain_id;

\echo ''

-- ============================================================================
-- TEST 7: Check indexes were created
-- ============================================================================
\echo 'TEST 7: Checking if indexes were created...'

DO $$
DECLARE
    index_count INTEGER;
BEGIN
    -- Check vault positions indexes
    SELECT COUNT(*) INTO index_count
    FROM pg_indexes
    WHERE tablename = 'vault_positions' AND indexname LIKE '%chain%';

    IF index_count >= 1 THEN
        RAISE NOTICE '✅ Chain-related indexes on vault_positions exist (% indexes)', index_count;
    ELSE
        RAISE WARNING '⚠️ No chain-related indexes on vault_positions found';
    END IF;

    -- Check GMX performance indexes
    SELECT COUNT(*) INTO index_count
    FROM pg_indexes
    WHERE tablename = 'gmx_performance_comparison';

    IF index_count >= 1 THEN
        RAISE NOTICE '✅ Indexes on gmx_performance_comparison exist (% indexes)', index_count;
    ELSE
        RAISE WARNING '⚠️ No indexes on gmx_performance_comparison found';
    END IF;
END $$;

\echo ''

-- ============================================================================
-- TEST 8: Test cross-chain views with sample data
-- ============================================================================
\echo 'TEST 8: Testing cross-chain views (showing first 5 rows)...'

\echo 'Cross-chain user portfolio view:'
SELECT * FROM cross_chain_user_portfolio LIMIT 5;

\echo ''
\echo 'Cross-chain total value view:'
SELECT * FROM cross_chain_total_value LIMIT 5;

\echo ''

-- ============================================================================
-- TEST 9: Test GMX performance comparison structure
-- ============================================================================
\echo 'TEST 9: Checking GMX performance comparison table structure...'

SELECT
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'gmx_performance_comparison'
ORDER BY ordinal_position;

\echo ''

-- ============================================================================
-- TEST 10: Performance test - check query speed
-- ============================================================================
\echo 'TEST 10: Testing query performance...'

EXPLAIN ANALYZE
SELECT
    chain_id,
    COUNT(*) as count,
    SUM(amount) as total_amount
FROM vault_positions
WHERE chain_id = 421614
GROUP BY chain_id;

\echo ''

-- ============================================================================
-- SUMMARY
-- ============================================================================
\echo '========================================='
\echo 'Migration 008 Test Summary'
\echo '========================================='
\echo ''
\echo 'If all tests passed, migration 008 is successfully applied!'
\echo ''
\echo 'Next steps:'
\echo '  1. Update backend configuration (internal/config/config.go)'
\echo '  2. Deploy contracts to Base Sepolia'
\echo '  3. Start multi-chain event listeners'
\echo '  4. Update frontend with chain switcher'
\echo ''
\echo 'See MIGRATION_008_GUIDE.md for detailed instructions.'
\echo '========================================='
