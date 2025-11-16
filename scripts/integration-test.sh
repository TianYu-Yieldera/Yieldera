#!/bin/bash

# Simple integration test script using curl
# Tests all new API endpoints

API_BASE="http://localhost:8080"
USER_ADDRESS="0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

# Test function
test_endpoint() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_code="${5:-200}"

    echo -n "Testing: $name... "

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi

    # Extract HTTP code (last line)
    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -eq "$expected_code" ]; then
        echo -e "${GREEN}PASS${NC} (HTTP $http_code)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}FAIL${NC} (HTTP $http_code, expected $expected_code)"
        echo "Response: $body"
        ((FAILED++))
        return 1
    fi
}

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Integration Test Suite${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Check if API is running
echo -n "Checking API health... "
if curl -s "$API_BASE/health" | grep -q '"ok":true'; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}FAIL - API not responding${NC}"
    exit 1
fi
echo ""

echo -e "${YELLOW}Yield Calculation Service Tests:${NC}"
test_endpoint "Get Treasury Rates" "GET" "$API_BASE/api/v1/yields/rates"
test_endpoint "Get User Total Yield" "GET" "$API_BASE/api/v1/yields/total/$USER_ADDRESS"
test_endpoint "Project Yield" "POST" "$API_BASE/api/v1/yields/project" \
    '{"bond_type":"TBILL_3M","principal_usd":1000,"duration_days":90,"compounding":true}'

echo ""
echo -e "${YELLOW}Notification Service Tests:${NC}"
test_endpoint "Get User Notifications" "GET" "$API_BASE/api/v1/notifications/$USER_ADDRESS"
test_endpoint "Get Notification Preferences" "GET" "$API_BASE/api/v1/notifications/$USER_ADDRESS/preferences"
test_endpoint "Update Notification Preferences" "PUT" "$API_BASE/api/v1/notifications/$USER_ADDRESS/preferences" \
    '{"channels":["email"],"min_priority":"medium","enabled_types":["liquidation_warning"],"frequency":"realtime"}'

echo ""
echo -e "${YELLOW}Auto Hedge Executor Tests:${NC}"
test_endpoint "Get Hedge Settings" "GET" "$API_BASE/api/v1/hedge/settings/$USER_ADDRESS"
test_endpoint "Get Hedge History" "GET" "$API_BASE/api/v1/hedge/history/$USER_ADDRESS"
test_endpoint "Get Hedge Stats" "GET" "$API_BASE/api/v1/hedge/stats"
test_endpoint "Update Hedge Settings" "PUT" "$API_BASE/api/v1/hedge/settings/$USER_ADDRESS" \
    '{"auto_hedge_enabled":true,"max_hedge_amount":5000,"min_health_factor":1.5,"target_health_factor":2.0}'

echo ""
echo -e "${YELLOW}Yield Distribution Service Tests:${NC}"
test_endpoint "Get Distribution Stats" "GET" "$API_BASE/api/v1/distribution/stats"
test_endpoint "Get Distribution Stats (7 days)" "GET" "$API_BASE/api/v1/distribution/stats?days=7"

echo ""
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Test Results${NC}"
echo -e "${YELLOW}========================================${NC}"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "\n${RED}Some tests failed! ✗${NC}"
    exit 1
fi
