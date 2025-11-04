package listener

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	k "github.com/segmentio/kafka-go"

	"loyalty-points-system/internal/models"
)

// BridgeListenerConfig holds configuration for bridge listener
type BridgeListenerConfig struct {
	L1RPCURL        string
	L1WSSURL        string
	L1ChainID       int64
	L1Gateway       string
	L1Confirmations int

	L2RPCURL        string
	L2WSSURL        string
	L2ChainID       int64
	L2Gateway       string
	L2Confirmations int
}

// BridgeListener listens to bridge events on both L1 and L2
type BridgeListener struct {
	cfg        BridgeListenerConfig
	l1Client   *ethclient.Client
	l2Client   *ethclient.Client
	writer     Writer
	pending    map[string]*models.BridgeEvent // keyed by message hash
	pendingMu  sync.Mutex
}

// NewBridgeListener creates a new bridge listener
func NewBridgeListener(cfg BridgeListenerConfig, writer Writer) *BridgeListener {
	return &BridgeListener{
		cfg:     cfg,
		writer:  writer,
		pending: make(map[string]*models.BridgeEvent),
	}
}

// Start begins listening to bridge events on both chains
func (l *BridgeListener) Start(ctx context.Context) error {
	var err error

	// Connect to L1
	l.l1Client, err = ethclient.DialContext(ctx, l.cfg.L1WSSURL)
	if err != nil {
		return err
	}
	log.Printf("🔌 [Bridge] Connected to L1 at %s", l.cfg.L1WSSURL)

	// Connect to L2
	l.l2Client, err = ethclient.DialContext(ctx, l.cfg.L2WSSURL)
	if err != nil {
		return err
	}
	log.Printf("🔌 [Bridge] Connected to L2 at %s", l.cfg.L2WSSURL)

	// Start L1 bridge listener
	go l.listenL1Bridge(ctx)

	// Start L2 bridge listener
	go l.listenL2Bridge(ctx)

	// Start status updater (checks pending messages)
	go l.updateBridgeStatus(ctx)

	return nil
}

// listenL1Bridge listens to L1 gateway events
func (l *BridgeListener) listenL1Bridge(ctx context.Context) {
	headers := make(chan *types.Header, 32)
	subHead, err := l.l1Client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Printf("❌ [Bridge-L1] Failed to subscribe to headers: %v", err)
		return
	}

	logsCh := make(chan types.Log, 256)
	addr := common.HexToAddress(l.cfg.L1Gateway)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	subLogs, err := l.l1Client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		log.Printf("❌ [Bridge-L1] Failed to subscribe to logs: %v", err)
		return
	}

	log.Printf("👂 [Bridge-L1] Subscribed to Gateway events at %s", addr.Hex())

	conf := l.cfg.L1Confirmations
	if conf <= 0 {
		conf = 12
	}

	for {
		select {
		case err := <-subHead.Err():
			log.Printf("❌ [Bridge-L1] header subscription error: %v", err)
			return
		case err := <-subLogs.Err():
			log.Printf("❌ [Bridge-L1] logs subscription error: %v", err)
			return
		case <-headers:
			// Block confirmation handled separately
		case lg := <-logsCh:
			l.onL1BridgeLog(lg)
		case <-ctx.Done():
			log.Println("🛑 [Bridge-L1] Listener stopped")
			return
		}
	}
}

// listenL2Bridge listens to L2 gateway/receiver events
func (l *BridgeListener) listenL2Bridge(ctx context.Context) {
	headers := make(chan *types.Header, 32)
	subHead, err := l.l2Client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Printf("❌ [Bridge-L2] Failed to subscribe to headers: %v", err)
		return
	}

	logsCh := make(chan types.Log, 256)

	// L2 might have different contract for receiving messages
	// For now, use the same gateway address pattern
	addr := common.HexToAddress(l.cfg.L2Gateway)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	subLogs, err := l.l2Client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		log.Printf("❌ [Bridge-L2] Failed to subscribe to logs: %v", err)
		return
	}

	log.Printf("👂 [Bridge-L2] Subscribed to Gateway events at %s", addr.Hex())

	conf := l.cfg.L2Confirmations
	if conf <= 0 {
		conf = 1
	}

	for {
		select {
		case err := <-subHead.Err():
			log.Printf("❌ [Bridge-L2] header subscription error: %v", err)
			return
		case err := <-subLogs.Err():
			log.Printf("❌ [Bridge-L2] logs subscription error: %v", err)
			return
		case <-headers:
			// Block confirmation handled separately
		case lg := <-logsCh:
			l.onL2BridgeLog(lg)
		case <-ctx.Done():
			log.Println("🛑 [Bridge-L2] Listener stopped")
			return
		}
	}
}

// onL1BridgeLog processes L1 bridge events
func (l *BridgeListener) onL1BridgeLog(lg types.Log) {
	if len(lg.Topics) < 2 {
		return
	}

	// Parse MessageSent event from L1 Gateway
	// MessageSent(bytes32 indexed messageHash, address indexed sender, uint256 amount, bytes data)

	messageHash := lg.Topics[1].Hex()
	var userAddress string

	if len(lg.Topics) >= 3 {
		userAddress = common.BytesToAddress(lg.Topics[2].Bytes()).Hex()
	}

	var amount string
	if len(lg.Data) >= 32 {
		amt := new(big.Int).SetBytes(lg.Data[:32])
		amount = amt.String()
	} else {
		amount = "0"
	}

	l.pendingMu.Lock()
	defer l.pendingMu.Unlock()

	// Create or update bridge event
	if existingEvent, exists := l.pending[messageHash]; exists {
		// Update existing event with L1 confirmation
		existingEvent.L1TxHash = lg.TxHash.Hex()
		existingEvent.L1BlockNumber = int64(lg.BlockNumber)
		existingEvent.Status = "pending"
		log.Printf("🔄 [Bridge] Updated L1→L2 message %s", messageHash)
	} else {
		// Create new bridge event
		event := &models.BridgeEvent{
			UserAddress:   userAddress,
			Amount:        amount,
			Direction:     "L1_TO_L2",
			Status:        "initiated",
			L1TxHash:      lg.TxHash.Hex(),
			L2TxHash:      "",
			MessageHash:   messageHash,
			L1BlockNumber: int64(lg.BlockNumber),
			L2BlockNumber: 0,
			InitiatedAt:   time.Now().Unix(),
			ConfirmedAt:   0,
			RetryCount:    0,
			ErrorMsg:      "",
		}
		l.pending[messageHash] = event
		log.Printf("🚀 [Bridge] L1→L2 message initiated %s user=%s amt=%s", messageHash, userAddress, amount)
	}
}

// onL2BridgeLog processes L2 bridge events
func (l *BridgeListener) onL2BridgeLog(lg types.Log) {
	if len(lg.Topics) < 2 {
		return
	}

	// Parse MessageReceived event from L2
	// MessageReceived(bytes32 indexed messageHash, address indexed recipient, uint256 amount)

	messageHash := lg.Topics[1].Hex()

	l.pendingMu.Lock()
	defer l.pendingMu.Unlock()

	// Update existing bridge event
	if existingEvent, exists := l.pending[messageHash]; exists {
		existingEvent.L2TxHash = lg.TxHash.Hex()
		existingEvent.L2BlockNumber = int64(lg.BlockNumber)
		existingEvent.Status = "confirmed"
		existingEvent.ConfirmedAt = time.Now().Unix()

		// Publish to Kafka
		go l.publishBridgeEvent(existingEvent)

		// Remove from pending
		delete(l.pending, messageHash)

		log.Printf("✅ [Bridge] L1→L2 message confirmed %s", messageHash)
	} else {
		// Received L2 event without L1 event (might be L2→L1)
		// Handle reverse direction
		var userAddress string
		if len(lg.Topics) >= 3 {
			userAddress = common.BytesToAddress(lg.Topics[2].Bytes()).Hex()
		}

		var amount string
		if len(lg.Data) >= 32 {
			amt := new(big.Int).SetBytes(lg.Data[:32])
			amount = amt.String()
		} else {
			amount = "0"
		}

		event := &models.BridgeEvent{
			UserAddress:   userAddress,
			Amount:        amount,
			Direction:     "L2_TO_L1",
			Status:        "initiated",
			L1TxHash:      "",
			L2TxHash:      lg.TxHash.Hex(),
			MessageHash:   messageHash,
			L1BlockNumber: 0,
			L2BlockNumber: int64(lg.BlockNumber),
			InitiatedAt:   time.Now().Unix(),
			ConfirmedAt:   0,
			RetryCount:    0,
			ErrorMsg:      "",
		}
		l.pending[messageHash] = event
		log.Printf("🚀 [Bridge] L2→L1 message initiated %s user=%s amt=%s", messageHash, userAddress, amount)
	}
}

// updateBridgeStatus periodically checks and updates pending bridge messages
func (l *BridgeListener) updateBridgeStatus(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.checkPendingMessages()
		case <-ctx.Done():
			return
		}
	}
}

// checkPendingMessages checks status of pending messages and retries if needed
func (l *BridgeListener) checkPendingMessages() {
	l.pendingMu.Lock()
	defer l.pendingMu.Unlock()

	now := time.Now().Unix()

	for messageHash, event := range l.pending {
		// Check if message is stuck (initiated >5 minutes ago)
		if now-event.InitiatedAt > 300 {
			event.RetryCount++

			if event.RetryCount > 10 {
				// Mark as failed after 10 retries
				event.Status = "failed"
				event.ErrorMsg = "timeout after 10 retries"
				go l.publishBridgeEvent(event)
				delete(l.pending, messageHash)
				log.Printf("❌ [Bridge] Message %s failed after timeout", messageHash)
			} else {
				log.Printf("⚠️  [Bridge] Message %s pending for %d seconds (retry %d)",
					messageHash, now-event.InitiatedAt, event.RetryCount)
			}
		}
	}
}

// publishBridgeEvent publishes a bridge event to Kafka
func (l *BridgeListener) publishBridgeEvent(event *models.BridgeEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, _ := json.Marshal(event)

	msg := k.Message{
		Key:   []byte(event.MessageHash),
		Value: data,
		Headers: []k.Header{
			{Key: "direction", Value: []byte(event.Direction)},
			{Key: "status", Value: []byte(event.Status)},
		},
	}

	if err := l.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("❌ [Bridge] Kafka write error: %v", err)
	} else {
		log.Printf("✅ [Bridge] Published %s message %s status=%s",
			event.Direction, event.MessageHash, event.Status)
	}
}

// Close closes the bridge listener
func (l *BridgeListener) Close() error {
	if l.l1Client != nil {
		l.l1Client.Close()
	}
	if l.l2Client != nil {
		l.l2Client.Close()
	}
	return nil
}
