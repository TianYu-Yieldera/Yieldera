#!/bin/bash

# Script to run all tests for Yieldera platform
# Usage: ./scripts/run-tests.sh [unit|integration|all]

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test type (default: all)
TEST_TYPE=${1:-all}

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Yieldera Platform Test Suite${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Function to run Go tests
run_go_tests() {
    echo -e "${YELLOW}Running Go API tests...${NC}"

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Go is not installed. Skipping Go tests.${NC}"
        return 1
    fi

    cd services/api/handlers || exit 1

    # Install testify if needed
    if ! go list -m github.com/stretchr/testify &> /dev/null; then
        echo "Installing testify..."
        go get github.com/stretchr/testify/assert
    fi

    # Run tests
    go test -v -cover ./...

    cd ../../.. || exit 1
}

# Function to run TypeScript unit tests
run_ts_unit_tests() {
    echo -e "${YELLOW}Running TypeScript unit tests...${NC}"

    cd backend || exit 1

    # Check if node_modules exist
    if [ ! -d "node_modules" ]; then
        echo "Installing dependencies..."
        npm install
    fi

    # Run Jest tests
    npm test -- --testPathPattern=__tests__ --coverage

    cd .. || exit 1
}

# Function to run integration tests
run_integration_tests() {
    echo -e "${YELLOW}Running integration tests...${NC}"

    # Check if services are running
    if ! curl -s http://localhost:8080/health > /dev/null; then
        echo -e "${RED}API service is not running at localhost:8080${NC}"
        echo "Please start services with: docker-compose up -d"
        return 1
    fi

    cd backend || exit 1

    # Run integration tests
    npm test -- --testPathPattern=integration --runInBand

    cd .. || exit 1
}

# Main test execution
case $TEST_TYPE in
    unit)
        echo "Running unit tests only..."
        run_go_tests || true
        run_ts_unit_tests || true
        ;;

    integration)
        echo "Running integration tests only..."
        run_integration_tests
        ;;

    all)
        echo "Running all tests..."
        run_go_tests || true
        run_ts_unit_tests || true
        run_integration_tests
        ;;

    *)
        echo -e "${RED}Invalid test type: $TEST_TYPE${NC}"
        echo "Usage: $0 [unit|integration|all]"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Test Suite Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
