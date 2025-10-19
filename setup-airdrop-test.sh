#!/bin/bash

echo "=========================================="
echo "Setting up Airdrop Test Environment"
echo "=========================================="

echo ""
echo "Step 1: Stopping existing containers..."
docker-compose down

echo ""
echo "Step 2: Rebuilding containers..."
docker-compose build --no-cache api scheduler frontend

echo ""
echo "Step 3: Starting services..."
docker-compose up -d

echo ""
echo "Waiting for database to be ready..."
sleep 10

echo ""
echo "Step 4: Loading test data..."
docker exec -i loyalty-postgres psql -U postgres -d loyalty_points < db/test-airdrop-data.sql

echo ""
echo "=========================================="
echo "✅ Setup Complete!"
echo "=========================================="
echo ""
echo "🌐 Frontend: http://localhost:5173"
echo "🔧 API: http://localhost:8080"
echo ""
echo "📋 Test Admin Addresses:"
echo "   - 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
echo "   - 0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
echo "   - 0x742d35cc6634c0532925a3b844bc9e7595f0beb1"
echo ""
echo "📖 Next steps:"
echo "   1. Visit http://localhost:5173"
echo "   2. Connect with one of the test admin addresses"
echo "   3. Go to '空投管理' in the navigation"
echo "   4. You can create campaigns, upload CSV, and manage airdrops"
echo "   5. Go to '空投' to claim as a user"
echo ""
