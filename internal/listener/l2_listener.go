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

// L2ListenerConfig holds configuration for L2 listener
type L2ListenerConfig struct {
	RPCURL              string
	WSSURL              string
	ChainID             int64
	Confirmations       int
	IntegratedVault     string
	StateAggregator     string
	AaveAdapter         string
	CompoundAdapter     string
	UniswapAdapter      string
	RWAFactory          string
	RWAMarketplace      string
	RWAYield            string
	RWACompliance       string
	RWAValuation        string
	RWAGovernance       string
	// Treasury contracts
	TreasuryFactory         string
	TreasuryMarketplace     string
	TreasuryYieldDistributor string
	TreasuryPriceOracle     string
}

// L2Listener listens to L2 events and publishes to Kafka
type L2Listener struct {
	cfg       L2ListenerConfig
	client    *ethclient.Client
	writer    Writer
	pending   map[uint64][]models.L2Event
	pendingMu sync.Mutex
	head      uint64
}

// NewL2Listener creates a new L2 listener
func NewL2Listener(cfg L2ListenerConfig, writer Writer) *L2Listener {
	return &L2Listener{
		cfg:     cfg,
		writer:  writer,
		pending: make(map[uint64][]models.L2Event),
	}
}

// Start begins listening to L2 events
func (l *L2Listener) Start(ctx context.Context) error {
	var err error
	l.client, err = ethclient.DialContext(ctx, l.cfg.WSSURL)
	if err != nil {
		return err
	}
	log.Printf("🔌 [L2] Connected to %s (chain ID: %d)", l.cfg.WSSURL, l.cfg.ChainID)

	// Subscribe to new block headers
	headers := make(chan *types.Header, 32)
	subHead, err := l.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return err
	}

	// Start listening to each contract
	logsCh := make(chan types.Log, 256)

	// Subscribe to IntegratedVault events
	if l.cfg.IntegratedVault != "" {
		if err := l.subscribeIntegratedVault(ctx, logsCh); err != nil {
			return err
		}
	}

	// Subscribe to StateAggregator events
	if l.cfg.StateAggregator != "" {
		if err := l.subscribeStateAggregator(ctx, logsCh); err != nil {
			return err
		}
	}

	// Subscribe to DeFi adapter events
	if l.cfg.AaveAdapter != "" {
		if err := l.subscribeDeFiAdapter(ctx, logsCh, l.cfg.AaveAdapter, "Aave"); err != nil {
			return err
		}
	}
	if l.cfg.CompoundAdapter != "" {
		if err := l.subscribeDeFiAdapter(ctx, logsCh, l.cfg.CompoundAdapter, "Compound"); err != nil {
			return err
		}
	}
	if l.cfg.UniswapAdapter != "" {
		if err := l.subscribeDeFiAdapter(ctx, logsCh, l.cfg.UniswapAdapter, "Uniswap"); err != nil {
			return err
		}
	}

	// Subscribe to RWA events
	rwaContracts := map[string]string{
		"RWAFactory":      l.cfg.RWAFactory,
		"RWAMarketplace":  l.cfg.RWAMarketplace,
		"RWAYield":        l.cfg.RWAYield,
		"RWACompliance":   l.cfg.RWACompliance,
		"RWAValuation":    l.cfg.RWAValuation,
		"RWAGovernance":   l.cfg.RWAGovernance,
	}

	for name, addr := range rwaContracts {
		if addr != "" {
			if err := l.subscribeRWAContract(ctx, logsCh, addr, name); err != nil {
				return err
			}
		}
	}

	// Subscribe to Treasury events
	treasuryContracts := map[string]string{
		"TreasuryFactory":         l.cfg.TreasuryFactory,
		"TreasuryMarketplace":     l.cfg.TreasuryMarketplace,
		"TreasuryYieldDistributor": l.cfg.TreasuryYieldDistributor,
		"TreasuryPriceOracle":     l.cfg.TreasuryPriceOracle,
	}

	for name, addr := range treasuryContracts {
		if addr != "" {
			if err := l.subscribeTreasuryContract(ctx, logsCh, addr, name); err != nil {
				return err
			}
		}
	}

	conf := l.cfg.Confirmations
	if conf <= 0 {
		conf = 1 // Default for L2 (faster finality)
	}

	// Main event loop
	go func() {
		for {
			select {
			case err := <-subHead.Err():
				log.Printf("❌ [L2] head subscription error: %v", err)
				return
			case h := <-headers:
				if h != nil {
					l.onHead(h.Number.Uint64(), conf)
				}
			case lg := <-logsCh:
				l.onLog(lg)
			case <-ctx.Done():
				log.Println("🛑 [L2] Listener stopped")
				return
			}
		}
	}()

	return nil
}

// subscribeIntegratedVault subscribes to IntegratedVault events
func (l *L2Listener) subscribeIntegratedVault(ctx context.Context, logsCh chan types.Log) error {
	addr := common.HexToAddress(l.cfg.IntegratedVault)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L2] Subscribed to IntegratedVault events at %s", addr.Hex())

	go func() {
		<-sub.Err()
		log.Println("❌ [L2] IntegratedVault subscription error")
	}()

	return nil
}

// subscribeStateAggregator subscribes to StateAggregator events
func (l *L2Listener) subscribeStateAggregator(ctx context.Context, logsCh chan types.Log) error {
	addr := common.HexToAddress(l.cfg.StateAggregator)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L2] Subscribed to StateAggregator events at %s", addr.Hex())

	go func() {
		<-sub.Err()
		log.Println("❌ [L2] StateAggregator subscription error")
	}()

	return nil
}

// subscribeDeFiAdapter subscribes to DeFi adapter events
func (l *L2Listener) subscribeDeFiAdapter(ctx context.Context, logsCh chan types.Log, address string, name string) error {
	addr := common.HexToAddress(address)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L2] Subscribed to %s adapter events at %s", name, addr.Hex())

	go func() {
		<-sub.Err()
		log.Printf("❌ [L2] %s adapter subscription error", name)
	}()

	return nil
}

// subscribeRWAContract subscribes to RWA contract events
func (l *L2Listener) subscribeRWAContract(ctx context.Context, logsCh chan types.Log, address string, name string) error {
	addr := common.HexToAddress(address)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L2] Subscribed to %s events at %s", name, addr.Hex())

	go func() {
		<-sub.Err()
		log.Printf("❌ [L2] %s subscription error", name)
	}()

	return nil
}

// onLog processes incoming log events
func (l *L2Listener) onLog(lg types.Log) {
	var event models.L2Event
	var eventType string

	// Determine event type based on contract address
	contractAddr := lg.Address.Hex()

	switch contractAddr {
	case common.HexToAddress(l.cfg.IntegratedVault).Hex():
		eventType = l.parseIntegratedVaultEvent(lg, &event)
	case common.HexToAddress(l.cfg.StateAggregator).Hex():
		eventType = l.parseStateAggregatorEvent(lg, &event)
	case common.HexToAddress(l.cfg.AaveAdapter).Hex():
		eventType = l.parseDeFiAdapterEvent(lg, &event, "aave")
	case common.HexToAddress(l.cfg.CompoundAdapter).Hex():
		eventType = l.parseDeFiAdapterEvent(lg, &event, "compound")
	case common.HexToAddress(l.cfg.UniswapAdapter).Hex():
		eventType = l.parseDeFiAdapterEvent(lg, &event, "uniswap")
	case common.HexToAddress(l.cfg.RWAFactory).Hex():
		eventType = l.parseRWAEvent(lg, &event, "factory")
	case common.HexToAddress(l.cfg.RWAMarketplace).Hex():
		eventType = l.parseRWAEvent(lg, &event, "marketplace")
	case common.HexToAddress(l.cfg.RWAYield).Hex():
		eventType = l.parseRWAEvent(lg, &event, "yield")
	case common.HexToAddress(l.cfg.RWACompliance).Hex():
		eventType = l.parseRWAEvent(lg, &event, "compliance")
	case common.HexToAddress(l.cfg.RWAValuation).Hex():
		eventType = l.parseRWAEvent(lg, &event, "valuation")
	case common.HexToAddress(l.cfg.RWAGovernance).Hex():
		eventType = l.parseRWAEvent(lg, &event, "governance")
	case common.HexToAddress(l.cfg.TreasuryFactory).Hex():
		eventType = l.parseTreasuryEvent(lg, &event, "factory")
	case common.HexToAddress(l.cfg.TreasuryMarketplace).Hex():
		eventType = l.parseTreasuryEvent(lg, &event, "marketplace")
	case common.HexToAddress(l.cfg.TreasuryYieldDistributor).Hex():
		eventType = l.parseTreasuryEvent(lg, &event, "yield_distributor")
	case common.HexToAddress(l.cfg.TreasuryPriceOracle).Hex():
		eventType = l.parseTreasuryEvent(lg, &event, "price_oracle")
	default:
		log.Printf("⚠️  [L2] Unknown contract address: %s", contractAddr)
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

	log.Printf("🕒 [L2] Pending %s block=%d tx=%s", eventType, lg.BlockNumber, lg.TxHash.Hex())
}

// parseIntegratedVaultEvent parses IntegratedVault events
func (l *L2Listener) parseIntegratedVaultEvent(lg types.Log, event *models.L2Event) string {
	if len(lg.Topics) < 2 {
		return ""
	}

	// Common pattern: user in topics[1], amount in data
	event.UserAddress = common.BytesToAddress(lg.Topics[1].Bytes()).Hex()

	if len(lg.Data) >= 32 {
		amount := new(big.Int).SetBytes(lg.Data[:32])
		event.Amount = amount.String()
	}

	// Determine specific event type based on signature
	eventSig := lg.Topics[0].Hex()
	_ = eventSig // Placeholder until real signatures are added

	// For now, return generic operation type
	// TODO: Update with actual event signatures after contract deployment
	return "vault_operation"
}

// parseStateAggregatorEvent parses StateAggregator events
func (l *L2Listener) parseStateAggregatorEvent(lg types.Log, event *models.L2Event) string {
	event.Amount = "0"
	event.UserAddress = "0x0000000000000000000000000000000000000000"

	return "state_aggregated"
}

// parseDeFiAdapterEvent parses DeFi adapter events
func (l *L2Listener) parseDeFiAdapterEvent(lg types.Log, event *models.L2Event, protocol string) string {
	if len(lg.Topics) < 2 {
		return ""
	}

	event.UserAddress = common.BytesToAddress(lg.Topics[1].Bytes()).Hex()

	if len(lg.Data) >= 32 {
		amount := new(big.Int).SetBytes(lg.Data[:32])
		event.Amount = amount.String()
	}

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}
	event.Metadata["protocol"] = protocol

	return "defi_" + protocol + "_operation"
}

// parseRWAEvent parses RWA contract events
func (l *L2Listener) parseRWAEvent(lg types.Log, event *models.L2Event, contractType string) string {
	if len(lg.Topics) < 2 {
		// Some RWA events might not have indexed user
		event.UserAddress = "0x0000000000000000000000000000000000000000"
	} else {
		event.UserAddress = common.BytesToAddress(lg.Topics[1].Bytes()).Hex()
	}

	// RWA events often have asset IDs or proposal IDs
	if len(lg.Data) >= 32 {
		amount := new(big.Int).SetBytes(lg.Data[:32])
		event.Amount = amount.String()
	} else {
		event.Amount = "0"
	}

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}
	event.Metadata["rwa_contract"] = contractType

	// Extract asset ID if present in topics
	if len(lg.Topics) >= 3 {
		assetID := new(big.Int).SetBytes(lg.Topics[2].Bytes())
		event.Metadata["asset_id"] = assetID.String()
	}

	return "rwa_" + contractType
}

// onHead confirms events after N blocks
func (l *L2Listener) onHead(n uint64, conf int) {
	l.head = n

	var confirmed []models.L2Event
	keep := make(map[uint64][]models.L2Event)

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
			log.Printf("❌ [L2] Kafka write error: %v", err)
		} else {
			log.Printf("✅ [L2] Confirmed %s tx=%s block=%d user=%s", evt.EventType, evt.TxHash, evt.BlockNumber, evt.UserAddress)
		}
	}
}

// subscribeTreasuryContract subscribes to Treasury contract events
func (l *L2Listener) subscribeTreasuryContract(ctx context.Context, logsCh chan types.Log, address string, name string) error {
	addr := common.HexToAddress(address)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}

	log.Printf("👂 [L2] Subscribed to %s events at %s", name, addr.Hex())

	go func() {
		<-sub.Err()
		log.Printf("❌ [L2] %s subscription error", name)
	}()

	return nil
}

// parseTreasuryEvent parses Treasury contract events
func (l *L2Listener) parseTreasuryEvent(lg types.Log, event *models.L2Event, contractType string) string {
	event.ContractAddress = lg.Address.Hex()
	event.TxHash = lg.TxHash.Hex()
	event.BlockNumber = lg.BlockNumber

	// Extract basic metadata
	metadata := map[string]interface{}{
		"contract_type": contractType,
		"log_index":     lg.Index,
	}

	// Parse based on contract type and event signature
	switch contractType {
	case "factory":
		// TreasuryTokenCreated event
		if len(lg.Topics) > 0 {
			eventSig := lg.Topics[0].Hex()
			switch eventSig {
			case "0x...": // TreasuryTokenCreated signature (to be filled with actual signature)
				// Parse asset creation event
				if len(lg.Topics) >= 3 {
					metadata["asset_id"] = new(big.Int).SetBytes(lg.Topics[1][:]).String()
					metadata["token_address"] = common.BytesToAddress(lg.Topics[2][:]).Hex()
				}
				event.Metadata = metadata
				return "treasury_token_created"
			}
		}

	case "marketplace":
		// OrderCreated, OrderMatched, OrderCancelled events
		if len(lg.Topics) > 0 {
			eventSig := lg.Topics[0].Hex()
			switch eventSig {
			case "0x...": // OrderCreated signature
				if len(lg.Topics) >= 3 {
					metadata["order_id"] = new(big.Int).SetBytes(lg.Topics[1][:]).String()
					metadata["user_address"] = common.BytesToAddress(lg.Topics[2][:]).Hex()
					event.UserAddress = metadata["user_address"].(string)
				}
				event.Metadata = metadata
				return "treasury_order_created"
			case "0x...": // OrderMatched signature
				if len(lg.Topics) >= 2 {
					metadata["order_id"] = new(big.Int).SetBytes(lg.Topics[1][:]).String()
				}
				event.Metadata = metadata
				return "treasury_order_matched"
			case "0x...": // OrderCancelled signature
				if len(lg.Topics) >= 2 {
					metadata["order_id"] = new(big.Int).SetBytes(lg.Topics[1][:]).String()
				}
				event.Metadata = metadata
				return "treasury_order_cancelled"
			}
		}

	case "yield_distributor":
		// YieldDeposited, YieldClaimed events
		if len(lg.Topics) > 0 {
			eventSig := lg.Topics[0].Hex()
			switch eventSig {
			case "0x...": // YieldDeposited signature
				if len(lg.Topics) >= 2 {
					metadata["asset_id"] = new(big.Int).SetBytes(lg.Topics[1][:]).String()
				}
				event.Metadata = metadata
				return "treasury_yield_deposited"
			case "0x...": // YieldClaimed signature
				if len(lg.Topics) >= 3 {
					metadata["user_address"] = common.BytesToAddress(lg.Topics[1][:]).Hex()
					metadata["asset_id"] = new(big.Int).SetBytes(lg.Topics[2][:]).String()
					event.UserAddress = metadata["user_address"].(string)
				}
				event.Metadata = metadata
				return "treasury_yield_claimed"
			}
		}

	case "price_oracle":
		// PriceUpdated event
		if len(lg.Topics) > 0 {
			eventSig := lg.Topics[0].Hex()
			switch eventSig {
			case "0x...": // PriceUpdated signature
				if len(lg.Topics) >= 2 {
					metadata["asset_id"] = new(big.Int).SetBytes(lg.Topics[1][:]).String()
				}
				event.Metadata = metadata
				return "treasury_price_updated"
			}
		}
	}

	// Default fallback
	event.Metadata = metadata
	return "treasury_unknown_event"
}

// Close closes the listener
func (l *L2Listener) Close() error {
	if l.client != nil {
		l.client.Close()
	}
	return nil
}
