#!/bin/bash

# Load environment variables
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC_RAW=events.raw
export DATABASE_URL="postgres://loyalty_user:loyalty_pass@localhost:5432/loyalty_db?sslmode=disable"
export POINTS_RATE=0.05
export SCHEDULER_INTERVAL_SEC=60
export API_PORT=8080
export API_ALLOW_ORIGIN=*
export LISTENER_MODE=real
export CHAINS_JSON='[{"name":"sepolia","wss_url":"wss://eth-sepolia.g.alchemy.com/v2/FP6JOVxZoc4lDScODskcP","token_address":"0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B","staking_address":"0x3dBF997f45b7AF6EA274A7c4e9BaaBAC79989eed","confirmations":6}]'

echo "🚀 Starting all services..."
echo ""

cd ~/loyalty-points-system-final

echo "Starting Listener..."
go run services/listener/cmd/main.go > logs/listener.log 2>&1 &
LISTENER_PID=$!

sleep 2

echo "Starting Consumer..."
go run services/consumer/cmd/main.go > logs/consumer.log 2>&1 &
CONSUMER_PID=$!

sleep 2

echo "Starting Scheduler..."
go run services/scheduler/cmd/main.go > logs/scheduler.log 2>&1 &
SCHEDULER_PID=$!

sleep 2

echo "Starting API..."
go run services/api/cmd/main.go > logs/api.log 2>&1 &
API_PID=$!

sleep 3

echo ""
echo "✅ All services started!"
echo ""
echo "📊 Service PIDs:"
echo "  Listener:  $LISTENER_PID"
echo "  Consumer:  $CONSUMER_PID"
echo "  Scheduler: $SCHEDULER_PID"
echo "  API:       $API_PID"
echo ""
echo "📝 Logs available in logs/ directory"
echo "🌐 API running on: http://localhost:8080"
echo ""
echo "⚡ Test API: curl http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop all services"

# Create stop function
cleanup() {
    echo ""
    echo "🛑 Stopping all services..."
    kill $LISTENER_PID $CONSUMER_PID $SCHEDULER_PID $API_PID 2>/dev/null
    echo "✅ All services stopped"
    exit 0
}

trap cleanup INT TERM

# Wait for user to press Ctrl+C
wait
