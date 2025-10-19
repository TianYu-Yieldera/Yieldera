#!/bin/bash

# Test script for Vault and RWA services
# This script tests the new microservices functionality

set -e

echo "🚀 Testing Vault and RWA Services"
echo "================================="

# Configuration
API_BASE="http://localhost:8080"
TEST_ADDRESS="0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
    fi
}

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4

    echo -e "\n${YELLOW}Testing:${NC} $description"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_BASE$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE$endpoint")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ] || [ "$http_code" = "201" ]; then
        print_status 0 "$description (HTTP $http_code)"
        echo "Response: $body" | head -n 1
    else
        print_status 1 "$description (HTTP $http_code)"
        echo "Response: $body"
    fi
}

echo -e "\n📊 Testing Oracle Service"
echo "------------------------"

test_endpoint "GET" "/api/oracle/apys" "" "Get all protocol APYs"
test_endpoint "GET" "/api/oracle/price/bAAPL" "" "Get price for bAAPL"
test_endpoint "GET" "/api/oracle/stats" "" "Get market statistics"

echo -e "\n💰 Testing Vault Service"
echo "------------------------"

test_endpoint "GET" "/api/vault/protocols" "" "Get available DeFi protocols"
test_endpoint "GET" "/api/vault/strategies" "" "Get vault strategies"
test_endpoint "GET" "/api/vault/balance/$TEST_ADDRESS" "" "Get user vault balance"

# Test deposit (will fail if user has no points, but tests the endpoint)
test_endpoint "POST" "/api/vault/deposit" \
    '{"user_address":"'$TEST_ADDRESS'","amount":1000,"mode":"smart","strategy":"balanced"}' \
    "Test vault deposit"

echo -e "\n💎 Testing RWA Service"
echo "----------------------"

test_endpoint "GET" "/api/rwa/assets?type=stock" "" "Get stock assets"
test_endpoint "GET" "/api/rwa/assets?type=bond" "" "Get bond assets"
test_endpoint "GET" "/api/rwa/assets?type=commodity" "" "Get commodity assets"
test_endpoint "GET" "/api/rwa/assets/bAAPL" "" "Get bAAPL asset details"
test_endpoint "GET" "/api/rwa/holdings/$TEST_ADDRESS" "" "Get user RWA holdings"
test_endpoint "GET" "/api/rwa/prices/bAAPL?period=7d" "" "Get bAAPL price history"

echo -e "\n🔄 Testing API Gateway Proxy"
echo "----------------------------"

# These should work through the main API gateway
test_endpoint "GET" "/health" "" "Main API health check"

echo -e "\n📈 Testing Price Updates"
echo "------------------------"

echo "Waiting for price update cycle (60 seconds)..."
echo "You should see prices changing in the database"

# Optional: Monitor price changes
if command -v watch &> /dev/null; then
    echo -e "\n${YELLOW}Tip:${NC} Run this command in another terminal to watch price updates:"
    echo "watch -n 5 'curl -s http://localhost:8080/api/oracle/price/bAAPL | jq .'"
fi

echo -e "\n✅ Basic tests completed!"
echo "=========================="
echo ""
echo "Next steps:"
echo "1. Check Docker logs: docker-compose logs -f vault rwa oracle"
echo "2. Connect to database: docker exec -it loyalty-postgres psql -U postgres -d loyalty_points"
echo "3. Run SQL queries to verify data:"
echo "   - SELECT * FROM defi_protocols;"
echo "   - SELECT * FROM rwa_assets;"
echo "   - SELECT * FROM vault_positions WHERE user_address='$TEST_ADDRESS';"
echo ""
echo "To test with frontend:"
echo "1. Open http://localhost:5173"
echo "2. Enable demo mode in Tutorial section"
echo "3. Navigate to Vault and RWA Market"
echo "4. Test operations with demo data"
echo "5. Disable demo mode to use real API"

# Make script executable
chmod +x "$0"