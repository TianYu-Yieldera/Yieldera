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

	"loyalty-points-system/internal/blockchain/l1"
	"loyalty-points-system/internal/models"
)

// L1ListenerConfig holds configuration for L1 listener
type L1ListenerConfig struct {
	RPCURL              string
	WSSURL              string
	ChainID             int64
	Confirmations       int
	CollateralVault     string
	StateRegistry       string
	LoyaltyUSD          string
	Gateway             string
}

// L1Listener listens to L1 events and publishes to Kafka
type L1Listener struct {
	cfg        L1ListenerConfig
	client     *ethclient.Client
	writer     Writer
	pending    map[uint64][]models.L1Event
	pendingMu  sync.Mutex
	head       uint64
}

// Writer interface for Kafka
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...k.Message) error
}

// NewL1Listener creates a new L1 listener
func NewL1Listener(cfg L1ListenerConfig, writer Writer) *L1Listener {
	return &L1Listener{
		cfg:     cfg,
		writer:  writer,
		pending: make(map[uint64][]models.L1Event),
	}
}

// Start begins listening to L1 events
func (l *L1Listener) Start(ctx context.Context) error {
	var err error
	l.client, err = ethclient.DialContext(ctx, l.cfg.WSSURL)
	if err != nil {
		return err
	}
	log.Printf("🔌 [L1] Connected to %s (chain ID: %d)", l.cfg.WSSURL, l.cfg.ChainID)

	// Subscribe to new block headers
	headers := make(chan *types.Header, 32)
	subHead, err := l.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return err
	}

	// Start listening to each contract
	logsCh := make(chan types.Log, 256)

	// Subscribe to CollateralVault events
	if l.cfg.CollateralVault != "" {
		if err := l.subscribeCollateralVault(ctx, logsCh); err != nil {
			return err
		}
	}

	// Subscribe to StateRegistry events
	if l.cfg.StateRegistry != "" {
		if err := l.subscribeStateRegistry(ctx, logsCh); err != nil {
			return err
		}
	}

	// Subscribe to LoyaltyUSD events
	if l.cfg.LoyaltyUSD != "" {
		if err := l.subscribeLoyaltyUSD(ctx, logsCh); err != nil {
			return err
		}
	}

	// Subscribe to Gateway events
	if l.cfg.Gateway != "" {
		if err := l.subscribeGateway(ctx, logsCh); err != nil {
			return err
		}
	}

	conf := l.cfg.Confirmations
	if conf <= 0 {
		conf = 12 // Default for L1
	}

	// Main event loop
	go func() {
		for {
			select {
			case err := <-subHead.Err():
				log.Printf("❌ [L1] head subscription error: %v", err)
				return
			case h := <-headers:
				if h != nil {
					l.onHead(h.Number.Uint64(), conf)
				}
			case lg := <-logsCh:
				l.onLog(lg)
			case <-ctx.Done():
				log.Println("🛑 [L1] Listener stopped")
				return
			}
		}
	}()

	return nil
}

// subscribeCollateralVault subscribes to CollateralVault events
func (l *L1Listener) subscribeCollateralVault(ctx context.Context, logsCh chan types.Log) error {
	addr := common.HexToAddress(l.cfg.CollateralVault)

	// Get contract to access event signatures
	contract, err := l1.NewCollateralVaultL1(addr, l.client)
	if err != nil {
		return err
	}

	// Subscribe to CollateralLocked and CollateralUnlocked events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L1] Subscribed to CollateralVault events at %s", addr.Hex())

	go func() {
		<-sub.Err()
		log.Println("❌ [L1] CollateralVault subscription error")
	}()

	_ = contract // Use contract to avoid unused warning

	return nil
}

// subscribeStateRegistry subscribes to StateRegistry events
func (l *L1Listener) subscribeStateRegistry(ctx context.Context, logsCh chan types.Log) error {
	addr := common.HexToAddress(l.cfg.StateRegistry)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L1] Subscribed to StateRegistry events at %s", addr.Hex())

	go func() {
		<-sub.Err()
		log.Println("❌ [L1] StateRegistry subscription error")
	}()

	return nil
}

// subscribeLoyaltyUSD subscribes to LoyaltyUSD events
func (l *L1Listener) subscribeLoyaltyUSD(ctx context.Context, logsCh chan types.Log) error {
	addr := common.HexToAddress(l.cfg.LoyaltyUSD)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L1] Subscribed to LoyaltyUSD events at %s", addr.Hex())

	go func() {
		<-sub.Err()
		log.Println("❌ [L1] LoyaltyUSD subscription error")
	}()

	return nil
}

// subscribeGateway subscribes to Gateway events
func (l *L1Listener) subscribeGateway(ctx context.Context, logsCh chan types.Log) error {
	addr := common.HexToAddress(l.cfg.Gateway)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L1] Subscribed to Gateway events at %s", addr.Hex())

	go func() {
		<-sub.Err()
		log.Println("❌ [L1] Gateway subscription error")
	}()

	return nil
}

// onLog processes incoming log events
func (l *L1Listener) onLog(lg types.Log) {
	var event models.L1Event
	var eventType string

	// Determine event type based on contract address
	switch lg.Address.Hex() {
	case common.HexToAddress(l.cfg.CollateralVault).Hex():
		eventType = l.parseCollateralVaultEvent(lg, &event)
	case common.HexToAddress(l.cfg.StateRegistry).Hex():
		eventType = l.parseStateRegistryEvent(lg, &event)
	case common.HexToAddress(l.cfg.LoyaltyUSD).Hex():
		eventType = l.parseLoyaltyUSDEvent(lg, &event)
	case common.HexToAddress(l.cfg.Gateway).Hex():
		eventType = l.parseGatewayEvent(lg, &event)
	default:
		log.Printf("⚠️  [L1] Unknown contract address: %s", lg.Address.Hex())
		return
	}

	if eventType == "" {
		return // Event not recognized
	}

	// Fill in common fields
	event.EventType = eventType
	event.TxHash = lg.TxHash.Hex()
	event.BlockNumber = int64(lg.BlockNumber)
	event.Confirmed = false
	event.Timestamp = time.Now().Unix()
	event.ContractAddress = lg.Address.Hex()

	// Add to pending
	l.pendingMu.Lock()
	l.pending[lg.BlockNumber] = append(l.pending[lg.BlockNumber], event)
	l.pendingMu.Unlock()

	log.Printf("🕒 [L1] Pending %s block=%d tx=%s", eventType, lg.BlockNumber, lg.TxHash.Hex())
}

// parseCollateralVaultEvent parses CollateralVault events
func (l *L1Listener) parseCollateralVaultEvent(lg types.Log, event *models.L1Event) string {
	if len(lg.Topics) < 2 {
		return ""
	}

	// Get event signature
	eventSig := lg.Topics[0].Hex()

	// CollateralLocked(address indexed user, address indexed token, uint256 amount, bytes32 indexed l2TxHash)
	// CollateralUnlocked(address indexed user, address indexed token, uint256 amount, bytes32 indexed l2TxHash)

	if len(lg.Topics) >= 2 {
		event.UserAddress = common.BytesToAddress(lg.Topics[1].Bytes()).Hex()
	}

	if len(lg.Data) >= 32 {
		amount := new(big.Int).SetBytes(lg.Data[:32])
		event.Amount = amount.String()
	}

	if len(lg.Topics) >= 3 {
		event.Token = common.BytesToAddress(lg.Topics[2].Bytes()).Hex()
	}

	if len(lg.Topics) >= 4 {
		event.L2TxHash = lg.Topics[3].Hex()
	}

	// Determine event type based on signature
	// Note: These are placeholder signatures - should be updated with actual event signatures from ABI
	// Use: crypto.Keccak256Hash([]byte("EventName(...)")).Hex()
	_ = eventSig // Placeholder until real signatures are added

	// For now, return generic operation type
	// TODO: Update with actual event signatures after contract deployment
	return "collateral_operation"
}

// parseStateRegistryEvent parses StateRegistry events
func (l *L1Listener) parseStateRegistryEvent(lg types.Log, event *models.L1Event) string {
	if len(lg.Topics) < 1 {
		return ""
	}

	// StateRootReceived(bytes32 indexed stateRoot, uint256 l2Block, uint256 timestamp)
	event.Amount = "0" // State updates don't have amounts
	event.UserAddress = "0x0000000000000000000000000000000000000000" // System event

	return "state_update"
}

// parseLoyaltyUSDEvent parses LoyaltyUSD events
func (l *L1Listener) parseLoyaltyUSDEvent(lg types.Log, event *models.L1Event) string {
	if len(lg.Topics) < 3 {
		return ""
	}

	// Transfer(address indexed from, address indexed to, uint256 value)
	to := common.BytesToAddress(lg.Topics[2].Bytes()).Hex()

	if len(lg.Data) >= 32 {
		amount := new(big.Int).SetBytes(lg.Data)
		event.Amount = amount.String()
	}

	// Create event for recipient (for tracking)
	event.UserAddress = to
	event.Token = "LOYALTY_USD"

	// Could emit two events (from/to) like the old listener
	// For now, just track incoming transfers

	return "loyaltyusd_transfer"
}

// parseGatewayEvent parses Gateway events
func (l *L1Listener) parseGatewayEvent(lg types.Log, event *models.L1Event) string {
	if len(lg.Topics) < 1 {
		return ""
	}

	// MessageSent, MessageReceived, etc.
	event.Amount = "0" // Bridge messages don't have amounts in L1 events
	event.UserAddress = "0x0000000000000000000000000000000000000000"

	return "gateway_message"
}

// onHead confirms events after N blocks
func (l *L1Listener) onHead(n uint64, conf int) {
	l.head = n

	var confirmed []models.L1Event
	keep := make(map[uint64][]models.L1Event)

	l.pendingMu.Lock()
	for bn, list := range l.pending {
		if n >= uint64(conf) && bn <= n-uint64(conf) {
			confirmed = append(confirmed, list...)
		} else {
			keep[bn] = list
		}
	}
	l.pending = keep
	l.pendingMu.Unlock()

	if len(confirmed) == 0 {
		return
	}

	// Publish confirmed events to Kafka
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, evt := range confirmed {
		evt.Confirmed = true
		data, _ := json.Marshal(evt)

		msg := k.Message{
			Key:   []byte(evt.UserAddress),
			Value: data,
			Headers: []k.Header{
				{Key: "event_type", Value: []byte(evt.EventType)},
				{Key: "contract", Value: []byte(evt.ContractAddress)},
			},
		}

		if err := l.writer.WriteMessages(ctx, msg); err != nil {
			log.Printf("❌ [L1] Kafka write error: %v", err)
		} else {
			log.Printf("✅ [L1] Confirmed %s tx=%s block=%d user=%s", evt.EventType, evt.TxHash, evt.BlockNumber, evt.UserAddress)
		}
	}
}

// Close closes the listener
func (l *L1Listener) Close() error {
	if l.client != nil {
		l.client.Close()
	}
	return nil
}
