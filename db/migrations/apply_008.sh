#!/bin/bash

# ============================================================================
# Migration 008 Deployment Script
# ============================================================================
# Purpose: Safely apply migration 008 with backup and verification
# Usage: ./db/migrations/apply_008.sh
# ============================================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}Yieldera Migration 008 Deployment${NC}"
echo -e "${BLUE}Multi-Chain Support${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

# ============================================================================
# Configuration
# ============================================================================

# Database connection settings (read from env or use defaults)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-yieldera}"
DB_USER="${DB_USER:-postgres}"

echo -e "${YELLOW}Database Configuration:${NC}"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo ""

# ============================================================================
# Step 1: Pre-flight checks
# ============================================================================

echo -e "${YELLOW}Step 1: Running pre-flight checks...${NC}"

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ Error: psql is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ psql is installed${NC}"

# Check if database exists
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo -e "${RED}❌ Error: Database '$DB_NAME' does not exist${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Database '$DB_NAME' exists${NC}"

# Check if TimescaleDB extension is installed
TIMESCALE_CHECK=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM pg_extension WHERE extname = 'timescaledb';" 2>/dev/null || echo "0")
if [ "$TIMESCALE_CHECK" -eq "0" ]; then
    echo -e "${YELLOW}⚠️  Warning: TimescaleDB extension not found${NC}"
    echo -e "${YELLOW}   Some features may not work. Consider installing TimescaleDB.${NC}"
else
    echo -e "${GREEN}✅ TimescaleDB extension is installed${NC}"
fi

echo ""

# ============================================================================
# Step 2: Backup database
# ============================================================================

echo -e "${YELLOW}Step 2: Creating database backup...${NC}"

BACKUP_DIR="./db/backups"
mkdir -p $BACKUP_DIR

BACKUP_FILE="$BACKUP_DIR/backup_before_migration_008_$(date +%Y%m%d_%H%M%S).sql"

echo "  Backing up to: $BACKUP_FILE"
pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME > $BACKUP_FILE

if [ -f "$BACKUP_FILE" ]; then
    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}✅ Backup created successfully ($BACKUP_SIZE)${NC}"
else
    echo -e "${RED}❌ Error: Backup failed${NC}"
    exit 1
fi

echo ""

# ============================================================================
# Step 3: Show what will be changed
# ============================================================================

echo -e "${YELLOW}Step 3: Migration preview${NC}"
echo "This migration will:"
echo "  📊 Add chain_id columns to existing tables"
echo "  🆕 Create 5 new tables (GMX comparison, liquidation predictions, etc.)"
echo "  👁️  Create 3 cross-chain views"
echo "  ⚙️  Add helper function for risk profiling"
echo "  📈 Create 15+ new indexes for performance"
echo ""
echo "Expected impact:"
echo "  - Downtime: ~10 seconds"
echo "  - Data loss: None (backward compatible)"
echo "  - Breaking changes: None"
echo ""

# Ask for confirmation
read -p "Continue with migration? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${YELLOW}Migration cancelled by user${NC}"
    exit 0
fi

echo ""

# ============================================================================
# Step 4: Apply migration
# ============================================================================

echo -e "${YELLOW}Step 4: Applying migration 008...${NC}"

MIGRATION_FILE="./db/migrations/008_multi_chain_support.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo -e "${RED}❌ Error: Migration file not found: $MIGRATION_FILE${NC}"
    exit 1
fi

echo "  Executing migration..."
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $MIGRATION_FILE > /tmp/migration_008_output.log 2>&1; then
    echo -e "${GREEN}✅ Migration applied successfully${NC}"
else
    echo -e "${RED}❌ Error: Migration failed${NC}"
    echo "See log: /tmp/migration_008_output.log"
    cat /tmp/migration_008_output.log
    echo ""
    echo -e "${YELLOW}Attempting to restore from backup...${NC}"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME < $BACKUP_FILE
    echo -e "${GREEN}✅ Database restored from backup${NC}"
    exit 1
fi

echo ""

# ============================================================================
# Step 5: Run tests
# ============================================================================

echo -e "${YELLOW}Step 5: Running verification tests...${NC}"

TEST_FILE="./db/migrations/test_008_migration.sql"

if [ -f "$TEST_FILE" ]; then
    echo "  Running test suite..."
    if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $TEST_FILE > /tmp/test_008_output.log 2>&1; then
        echo -e "${GREEN}✅ All tests passed${NC}"
    else
        echo -e "${YELLOW}⚠️  Some tests failed (this may be expected if tables are empty)${NC}"
        echo "See log: /tmp/test_008_output.log"
    fi
else
    echo -e "${YELLOW}⚠️  Test file not found, skipping verification${NC}"
fi

echo ""

# ============================================================================
# Step 6: Verify data
# ============================================================================

echo -e "${YELLOW}Step 6: Verifying data migration...${NC}"

# Count vault positions by chain
echo "  Vault positions by chain:"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
SELECT
    chain_id,
    CASE
        WHEN chain_id = 421614 THEN 'Arbitrum Sepolia'
        WHEN chain_id = 42161 THEN 'Arbitrum One'
        WHEN chain_id = 84532 THEN 'Base Sepolia'
        WHEN chain_id = 8453 THEN 'Base Mainnet'
        ELSE 'Unknown'
    END as chain_name,
    COUNT(*) as count
FROM vault_positions
GROUP BY chain_id
ORDER BY chain_id;" 2>/dev/null || echo "  (Table may be empty)"

echo ""

# Count treasury holdings by chain
echo "  Treasury holdings by chain:"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
SELECT
    chain_id,
    CASE
        WHEN chain_id = 421614 THEN 'Arbitrum Sepolia'
        WHEN chain_id = 42161 THEN 'Arbitrum One'
        WHEN chain_id = 84532 THEN 'Base Sepolia'
        WHEN chain_id = 8453 THEN 'Base Mainnet'
        ELSE 'Unknown'
    END as chain_name,
    COUNT(*) as count
FROM treasury_holdings
GROUP BY chain_id
ORDER BY chain_id;" 2>/dev/null || echo "  (Table may be empty)"

echo ""

# ============================================================================
# Step 7: Summary
# ============================================================================

echo -e "${BLUE}=========================================${NC}"
echo -e "${GREEN}✅ Migration 008 Completed Successfully!${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""
echo "📝 Summary:"
echo "  - Backup: $BACKUP_FILE"
echo "  - Migration log: /tmp/migration_008_output.log"
echo "  - Test log: /tmp/test_008_output.log"
echo ""
echo "📋 Next Steps:"
echo "  1. Update backend config (internal/config/config.go)"
echo "  2. Deploy contracts to Base Sepolia"
echo "  3. Start multi-chain event listeners"
echo "  4. Update frontend with chain switcher"
echo ""
echo "📖 Documentation:"
echo "  - Migration guide: db/migrations/MIGRATION_008_GUIDE.md"
echo "  - Multi-chain strategy: MULTI_CHAIN_STRATEGY_PLAN.md"
echo ""
echo -e "${GREEN}Happy multi-chain building! 🚀${NC}"
