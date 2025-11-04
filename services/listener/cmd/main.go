package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"loyalty-points-system/internal/config"
	"loyalty-points-system/internal/kafka"
	"loyalty-points-system/internal/listener"
)

func main() {
	log.Println("🚀 Starting L1/L2 Listener Service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Ensure Kafka topics exist
	ctx := context.Background()
	broker := cfg.KafkaBrokers
	if err := kafka.EnsureTopic(ctx, broker, cfg.KafkaTopicL1, 3, 1); err != nil {
		log.Printf("⚠️  Failed to create L1 topic: %v", err)
	}
	if err := kafka.EnsureTopic(ctx, broker, cfg.KafkaTopicL2, 3, 1); err != nil {
		log.Printf("⚠️  Failed to create L2 topic: %v", err)
	}
	if err := kafka.EnsureTopic(ctx, broker, cfg.KafkaTopicBridge, 3, 1); err != nil {
		log.Printf("⚠️  Failed to create Bridge topic: %v", err)
	}

	// Create Kafka writers
	brokers := []string{cfg.KafkaBrokers}
	l1Writer := kafka.NewWriter(brokers, cfg.KafkaTopicL1)
	l2Writer := kafka.NewWriter(brokers, cfg.KafkaTopicL2)
	bridgeWriter := kafka.NewWriter(brokers, cfg.KafkaTopicBridge)

	defer l1Writer.Close()
	defer l2Writer.Close()
	defer bridgeWriter.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and start L1 listener
	l1Config := listener.L1ListenerConfig{
		RPCURL:          cfg.L1RPCURL,
		WSSURL:          cfg.L1WSSURL,
		ChainID:         cfg.L1ChainID,
		Confirmations:   cfg.L1Confirmations,
		CollateralVault: cfg.L1CollateralVault,
		StateRegistry:   cfg.L1StateRegistry,
		LoyaltyUSD:      cfg.L1LoyaltyUSD,
		Gateway:         cfg.L1Gateway,
	}

	l1Listener := listener.NewL1Listener(l1Config, l1Writer)
	if err := l1Listener.Start(ctx); err != nil {
		log.Fatalf("❌ Failed to start L1 listener: %v", err)
	}
	log.Println("✅ L1 Listener started")

	// Create and start L2 listener
	l2Config := listener.L2ListenerConfig{
		RPCURL:          cfg.L2RPCURL,
		WSSURL:          cfg.L2WSSURL,
		ChainID:         cfg.L2ChainID,
		Confirmations:   cfg.L2Confirmations,
		IntegratedVault: cfg.L2IntegratedVault,
		StateAggregator: cfg.L2StateAggregator,
		AaveAdapter:     cfg.L2AaveAdapter,
		CompoundAdapter: cfg.L2CompoundAdapter,
		UniswapAdapter:  cfg.L2UniswapAdapter,
		RWAFactory:      cfg.L2RWAFactory,
		RWAMarketplace:  cfg.L2RWAMarketplace,
		RWAYield:        cfg.L2RWAYieldDistributor,
		RWACompliance:   cfg.L2RWACompliance,
		RWAValuation:    cfg.L2RWAValuation,
		RWAGovernance:   cfg.L2RWAGovernance,
	}

	l2Listener := listener.NewL2Listener(l2Config, l2Writer)
	if err := l2Listener.Start(ctx); err != nil {
		log.Fatalf("❌ Failed to start L2 listener: %v", err)
	}
	log.Println("✅ L2 Listener started")

	// Create and start Bridge listener
	bridgeConfig := listener.BridgeListenerConfig{
		L1RPCURL:        cfg.L1RPCURL,
		L1WSSURL:        cfg.L1WSSURL,
		L1ChainID:       cfg.L1ChainID,
		L1Gateway:       cfg.L1Gateway,
		L1Confirmations: cfg.L1Confirmations,
		L2RPCURL:        cfg.L2RPCURL,
		L2WSSURL:        cfg.L2WSSURL,
		L2ChainID:       cfg.L2ChainID,
		L2Gateway:       cfg.L2IntegratedVault, // TODO: Use actual L2 Gateway address
		L2Confirmations: cfg.L2Confirmations,
	}

	bridgeListener := listener.NewBridgeListener(bridgeConfig, bridgeWriter)
	if err := bridgeListener.Start(ctx); err != nil {
		log.Fatalf("❌ Failed to start Bridge listener: %v", err)
	}
	log.Println("✅ Bridge Listener started")

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	// Status endpoint with listener info
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"status": "running",
			"listeners": {
				"l1": "active",
				"l2": "active",
				"bridge": "active"
			},
			"topics": {
				"l1": "` + cfg.KafkaTopicL1 + `",
				"l2": "` + cfg.KafkaTopicL2 + `",
				"bridge": "` + cfg.KafkaTopicBridge + `"
			}
		}`))
	})

	// Start HTTP server in goroutine
	go func() {
		log.Println("🏥 Health check server on :8090")
		if err := http.ListenAndServe(":8090", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("✅ All listeners running. Press Ctrl+C to stop.")

	<-sigChan
	log.Println("\n🛑 Shutdown signal received, stopping listeners...")

	// Cancel context to stop all listeners
	cancel()

	// Close listeners
	l1Listener.Close()
	l2Listener.Close()
	bridgeListener.Close()

	log.Println("✅ All listeners stopped gracefully")
}
