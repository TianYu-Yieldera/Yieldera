# Database Migration 008: Multi-Chain Support Guide

## Overview
This migration adds comprehensive multi-chain support to enable Yieldera's dual-chain architecture:
- **Arbitrum**: Aggressive DeFi investments (GMX, Aave, etc.)
- **Base**: Conservative RWA investments (US Treasury bonds)

## Architecture Principles
✅ **Chain Isolation**: No cross-chain bridge, each chain operates independently
✅ **AI Orchestration**: AI analyzes both chains but doesn't execute trades
✅ **User Control**: Users manage their own wallets on each chain

## What This Migration Does

### 1. Adds `chain_id` to Existing Tables
- `vault_positions` - DeFi positions on both chains
- `treasury_holdings` - Treasury bonds (Base only)
- `price_history` - Market data from all chains
- `liquidation_events` - Liquidation tracking per chain
- `risk_assessments` - Risk calculations per chain
- All AI/ML tables

### 2. Creates New Tables
- `gmx_performance_comparison` - GMX Arbitrum vs Base metrics
- `liquidation_predictions` - AI-powered early warning system
- `aerodrome_swaps` - Base DEX activity tracking
- `base_pay_transactions` - Credit card onramp tracking
- `supported_chains` - Chain configuration reference

### 3. Creates Cross-Chain Views
- `cross_chain_user_portfolio` - User positions breakdown by chain
- `cross_chain_total_value` - Asset allocation percentages
- `gmx_performance_summary` - Performance comparison data

### 4. Adds Helper Functions
- `get_cross_chain_risk_profile()` - Calculate user risk across chains

## Chain IDs Reference

| Network | Chain ID | Environment | Purpose |
|---------|----------|-------------|---------|
| Arbitrum Sepolia | 421614 | Testnet | Aggressive DeFi testing |
| Arbitrum One | 42161 | Mainnet | Aggressive DeFi production |
| Base Sepolia | 84532 | Testnet | Conservative RWA testing |
| Base Mainnet | 8453 | Mainnet | Conservative RWA production |

## Pre-Migration Checklist

- [ ] Backup your database
- [ ] Ensure TimescaleDB extension is installed
- [ ] Verify all previous migrations (001-007) are applied
- [ ] Check that `vault_positions` and `treasury_holdings` tables exist
- [ ] Confirm no active connections writing to the database

## How to Run the Migration

### Option 1: Using psql (Recommended for Testing)

```bash
# Backup first!
pg_dump -U your_user -d yieldera > backup_before_008_$(date +%Y%m%d_%H%M%S).sql

# Apply migration
psql -U your_user -d yieldera -f db/migrations/008_multi_chain_support.sql

# Verify
psql -U your_user -d yieldera -c "SELECT * FROM supported_chains;"
psql -U your_user -d yieldera -c "SELECT chain_id, COUNT(*) FROM vault_positions GROUP BY chain_id;"
```

### Option 2: Using Migration Tool

```bash
# If you have a migration runner
npm run migrate:up 008
# or
go run services/api/cmd/migrate.go up
```

## Post-Migration Verification

### 1. Check Chain Data Distribution
```sql
-- Verify vault positions have chain_id
SELECT chain_id, COUNT(*) as count, SUM(amount) as total_amount
FROM vault_positions
GROUP BY chain_id;

-- Expected: All existing positions on Arbitrum (421614)
```

### 2. Check Treasury Migration
```sql
-- Verify treasury is on Base
SELECT chain_id, COUNT(*) as count, SUM(principal_usd) as total_principal
FROM treasury_holdings
GROUP BY chain_id;

-- Expected: All treasury on Base Sepolia (84532)
```

### 3. Test Cross-Chain Views
```sql
-- Test user portfolio view
SELECT * FROM cross_chain_user_portfolio LIMIT 5;

-- Test total value view
SELECT * FROM cross_chain_total_value LIMIT 5;
```

### 4. Test Helper Function
```sql
-- Replace with actual user address
SELECT * FROM get_cross_chain_risk_profile('0x1234...');
```

## Data Migration Impact

### Existing Data
- All existing `vault_positions` will have `chain_id = 421614` (Arbitrum Sepolia)
- All existing `treasury_holdings` will have `chain_id = 84532` (Base Sepolia)
- No data loss or modification to existing records

### New Fields Added
- `vault_positions`: `chain_id`, `collateral_value_usd`, `debt_value_usd`, `health_factor`, `active`, `position_id`
- `treasury_holdings`: `chain_id`, `asset_type`
- All TimescaleDB tables: `chain_id`

## Performance Considerations

### New Indexes
The migration creates 15+ new indexes for efficient multi-chain queries:
- `idx_vault_positions_chain_user`
- `idx_treasury_chain_user`
- `idx_gmx_perf_chain_time`
- And more...

### Query Performance
- Single-chain queries: Same as before (uses chain_id index)
- Cross-chain queries: Uses new materialized views
- Expected query time: <100ms for most operations

## Rollback Procedure

If you need to rollback this migration:

```sql
-- ROLLBACK SCRIPT (Run with caution!)
BEGIN;

-- Drop new tables
DROP TABLE IF EXISTS gmx_performance_comparison CASCADE;
DROP TABLE IF EXISTS liquidation_predictions CASCADE;
DROP TABLE IF EXISTS aerodrome_swaps CASCADE;
DROP TABLE IF EXISTS base_pay_transactions CASCADE;
DROP TABLE IF EXISTS supported_chains CASCADE;

-- Drop views
DROP VIEW IF EXISTS cross_chain_user_portfolio CASCADE;
DROP VIEW IF EXISTS cross_chain_total_value CASCADE;
DROP VIEW IF EXISTS gmx_performance_summary CASCADE;

-- Drop function
DROP FUNCTION IF EXISTS get_cross_chain_risk_profile CASCADE;

-- Remove chain_id columns (WARNING: This may cause data loss!)
ALTER TABLE vault_positions DROP COLUMN IF EXISTS chain_id CASCADE;
ALTER TABLE treasury_holdings DROP COLUMN IF EXISTS chain_id CASCADE;
ALTER TABLE price_history DROP COLUMN IF EXISTS chain_id CASCADE;
-- ... (repeat for all tables)

COMMIT;
```

## Troubleshooting

### Issue: Migration fails with "relation does not exist"
**Solution**: Ensure migrations 001-007 are applied first
```bash
psql -U your_user -d yieldera -c "\dt"
```

### Issue: "cannot alter type of column referenced in a view"
**Solution**: Drop dependent views first, then re-run migration
```sql
DROP VIEW IF EXISTS cross_chain_user_portfolio CASCADE;
-- Then re-run migration
```

### Issue: TimescaleDB hypertable error
**Solution**: Check TimescaleDB is installed
```sql
SELECT extname, extversion FROM pg_extension WHERE extname = 'timescaledb';
```

## Next Steps After Migration

### 1. Update Backend Config
```go
// internal/config/config.go
type Config struct {
    Chains map[int64]*ChainConfig
}
```

### 2. Deploy Contracts to Base
```bash
# Deploy treasury system to Base Sepolia
npx hardhat run scripts/deploy-treasury-base.js --network baseSepolia
```

### 3. Start Multi-Chain Listeners
```typescript
// Start Base event listeners
const baseListener = new BaseTreasuryListener(baseProvider, db);
await baseListener.start();
```

### 4. Update Frontend
```javascript
// Add chain switcher
<ChainSwitcher currentChain={chainId} onChange={setChainId} />
```

## Support & Questions

- Check migration logs: `tail -f /var/log/postgresql/postgresql.log`
- Verify migration applied: `SELECT * FROM schema_migrations WHERE version = '008';`
- Report issues: https://github.com/TianYu-Yieldera/Yieldera/issues

## References

- [MULTI_CHAIN_STRATEGY_PLAN.md](../../MULTI_CHAIN_STRATEGY_PLAN.md)
- [BASE_ECOSYSTEM_STRATEGY.md](../../BASE_ECOSYSTEM_STRATEGY.md)
- TimescaleDB docs: https://docs.timescale.com
- PostgreSQL docs: https://www.postgresql.org/docs/

---

**Migration Status**: ✅ Ready for testing
**Tested On**: PostgreSQL 14+ with TimescaleDB 2.x
**Breaking Changes**: None (backward compatible)
**Data Loss Risk**: None (adds columns only)
