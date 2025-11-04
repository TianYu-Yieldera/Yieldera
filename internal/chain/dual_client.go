package chain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// DualClientManager manages connections to both L1 (Ethereum) and L2 (Arbitrum) chains
type DualClientManager struct {
	L1Client  *ethclient.Client
	L2Client  *ethclient.Client
	L1ChainID int64
	L2ChainID int64
	L1RPCURL  string
	L2RPCURL  string
}

// NewDualClientManager creates a new manager with connections to both L1 and L2
func NewDualClientManager(ctx context.Context, l1RPC, l2RPC string, l1ChainID, l2ChainID int64) (*DualClientManager, error) {
	log.Printf("🔌 Connecting to L1 (chain ID: %d) at %s", l1ChainID, l1RPC)
	l1Client, err := ethclient.DialContext(ctx, l1RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to L1: %w", err)
	}

	log.Printf("🔌 Connecting to L2 (chain ID: %d) at %s", l2ChainID, l2RPC)
	l2Client, err := ethclient.DialContext(ctx, l2RPC)
	if err != nil {
		l1Client.Close()
		return nil, fmt.Errorf("failed to connect to L2: %w", err)
	}

	// Verify L1 chain ID
	l1ID, err := l1Client.ChainID(ctx)
	if err != nil {
		l1Client.Close()
		l2Client.Close()
		return nil, fmt.Errorf("failed to get L1 chain ID: %w", err)
	}
	if l1ID.Int64() != l1ChainID {
		l1Client.Close()
		l2Client.Close()
		return nil, fmt.Errorf("L1 chain ID mismatch: expected %d, got %d", l1ChainID, l1ID.Int64())
	}

	// Verify L2 chain ID
	l2ID, err := l2Client.ChainID(ctx)
	if err != nil {
		l1Client.Close()
		l2Client.Close()
		return nil, fmt.Errorf("failed to get L2 chain ID: %w", err)
	}
	if l2ID.Int64() != l2ChainID {
		l1Client.Close()
		l2Client.Close()
		return nil, fmt.Errorf("L2 chain ID mismatch: expected %d, got %d", l2ChainID, l2ID.Int64())
	}

	log.Printf("✅ Connected to L1 (chain ID: %d) and L2 (chain ID: %d)", l1ID.Int64(), l2ID.Int64())

	return &DualClientManager{
		L1Client:  l1Client,
		L2Client:  l2Client,
		L1ChainID: l1ChainID,
		L2ChainID: l2ChainID,
		L1RPCURL:  l1RPC,
		L2RPCURL:  l2RPC,
	}, nil
}

// Close closes both L1 and L2 client connections
func (m *DualClientManager) Close() {
	if m.L1Client != nil {
		m.L1Client.Close()
		log.Println("🔌 L1 client closed")
	}
	if m.L2Client != nil {
		m.L2Client.Close()
		log.Println("🔌 L2 client closed")
	}
}

// HealthCheck checks the health of both L1 and L2 connections
func (m *DualClientManager) HealthCheck(ctx context.Context) error {
	// Check L1
	l1Block, err := m.L1Client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("L1 health check failed: %w", err)
	}
	log.Printf("🏥 L1 Health: OK (block: %d)", l1Block)

	// Check L2
	l2Block, err := m.L2Client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("L2 health check failed: %w", err)
	}
	log.Printf("🏥 L2 Health: OK (block: %d)", l2Block)

	return nil
}

// GetL1BlockNumber returns the latest L1 block number
func (m *DualClientManager) GetL1BlockNumber(ctx context.Context) (uint64, error) {
	return m.L1Client.BlockNumber(ctx)
}

// GetL2BlockNumber returns the latest L2 block number
func (m *DualClientManager) GetL2BlockNumber(ctx context.Context) (uint64, error) {
	return m.L2Client.BlockNumber(ctx)
}

// ReconnectL1 attempts to reconnect to L1 with exponential backoff
func (m *DualClientManager) ReconnectL1(ctx context.Context) error {
	log.Println("🔄 Reconnecting to L1...")

	if m.L1Client != nil {
		m.L1Client.Close()
	}

	maxRetries := 5
	baseDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		l1Client, err := ethclient.DialContext(ctx, m.L1RPCURL)
		if err != nil {
			delay := baseDelay * time.Duration(1<<uint(i)) // Exponential backoff: 2s, 4s, 8s, 16s, 32s
			log.Printf("❌ L1 reconnect attempt %d/%d failed: %v. Retrying in %v...", i+1, maxRetries, err, delay)
			time.Sleep(delay)
			continue
		}

		// Verify chain ID
		l1ID, err := l1Client.ChainID(ctx)
		if err != nil {
			l1Client.Close()
			log.Printf("❌ L1 chain ID verification failed: %v", err)
			continue
		}
		if l1ID.Int64() != m.L1ChainID {
			l1Client.Close()
			return fmt.Errorf("L1 chain ID mismatch after reconnect: expected %d, got %d", m.L1ChainID, l1ID.Int64())
		}

		m.L1Client = l1Client
		log.Println("✅ L1 reconnected successfully")
		return nil
	}

	return fmt.Errorf("failed to reconnect to L1 after %d attempts", maxRetries)
}

// ReconnectL2 attempts to reconnect to L2 with exponential backoff
func (m *DualClientManager) ReconnectL2(ctx context.Context) error {
	log.Println("🔄 Reconnecting to L2...")

	if m.L2Client != nil {
		m.L2Client.Close()
	}

	maxRetries := 5
	baseDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		l2Client, err := ethclient.DialContext(ctx, m.L2RPCURL)
		if err != nil {
			delay := baseDelay * time.Duration(1<<uint(i)) // Exponential backoff
			log.Printf("❌ L2 reconnect attempt %d/%d failed: %v. Retrying in %v...", i+1, maxRetries, err, delay)
			time.Sleep(delay)
			continue
		}

		// Verify chain ID
		l2ID, err := l2Client.ChainID(ctx)
		if err != nil {
			l2Client.Close()
			log.Printf("❌ L2 chain ID verification failed: %v", err)
			continue
		}
		if l2ID.Int64() != m.L2ChainID {
			l2Client.Close()
			return fmt.Errorf("L2 chain ID mismatch after reconnect: expected %d, got %d", m.L2ChainID, l2ID.Int64())
		}

		m.L2Client = l2Client
		log.Println("✅ L2 reconnected successfully")
		return nil
	}

	return fmt.Errorf("failed to reconnect to L2 after %d attempts", maxRetries)
}

// GetL1GasPrice returns the current L1 gas price
func (m *DualClientManager) GetL1GasPrice(ctx context.Context) (*big.Int, error) {
	return m.L1Client.SuggestGasPrice(ctx)
}

// GetL2GasPrice returns the current L2 gas price
func (m *DualClientManager) GetL2GasPrice(ctx context.Context) (*big.Int, error) {
	return m.L2Client.SuggestGasPrice(ctx)
}

// StartHealthCheckLoop starts a background goroutine that periodically checks connection health
func (m *DualClientManager) StartHealthCheckLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 Health check loop stopped")
				return
			case <-ticker.C:
				if err := m.HealthCheck(ctx); err != nil {
					log.Printf("⚠️ Health check failed: %v. Attempting reconnection...", err)

					// Try to reconnect
					if _, err := m.L1Client.BlockNumber(ctx); err != nil {
						log.Println("🔄 L1 connection lost, reconnecting...")
						if err := m.ReconnectL1(ctx); err != nil {
							log.Printf("❌ Failed to reconnect L1: %v", err)
						}
					}

					if _, err := m.L2Client.BlockNumber(ctx); err != nil {
						log.Println("🔄 L2 connection lost, reconnecting...")
						if err := m.ReconnectL2(ctx); err != nil {
							log.Printf("❌ Failed to reconnect L2: %v", err)
						}
					}
				}
			}
		}
	}()
	log.Printf("🏥 Health check loop started (interval: %v)", interval)
}
