package models

// ====== Legacy Event (keep for backward compatibility) ======
type BalanceEvent struct {
	UserAddress string `json:"user_address"`
	Amount      string `json:"amount"`
	EventType   string `json:"event_type"`
	TxHash      string `json:"tx_hash"`
	Chain       string `json:"chain"`
	BlockNumber int64  `json:"block_number"`
	Confirmed   bool   `json:"confirmed"`
	Timestamp   int64  `json:"timestamp"`
}

// ====== NEW: L1 Event ======
type L1Event struct {
	UserAddress     string `json:"user_address"`
	Amount          string `json:"amount"`
	EventType       string `json:"event_type"` // "collateral_deposit", "collateral_withdraw", "state_update"
	Token           string `json:"token"`      // "USDC", "USDT", "DAI"
	TxHash          string `json:"tx_hash"`
	BlockNumber     int64  `json:"block_number"`
	Confirmed       bool   `json:"confirmed"`
	Timestamp       int64  `json:"timestamp"`
	ContractAddress string `json:"contract_address"`     // CollateralVaultL1 address
	L2TxHash        string `json:"l2_tx_hash,omitempty"` // For bridge events
}

// ====== NEW: L2 Event ======
type L2Event struct {
	UserAddress     string                 `json:"user_address"`
	Amount          string                 `json:"amount"`
	EventType       string                 `json:"event_type"` // "vault_deposit", "vault_withdraw", "rwa_trade", etc.
	TxHash          string                 `json:"tx_hash"`
	BlockNumber     int64                  `json:"block_number"`
	Confirmed       bool                   `json:"confirmed"`
	Timestamp       int64                  `json:"timestamp"`
	ContractAddress string                 `json:"contract_address"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"` // Protocol, asset_id, etc.
}

// ====== NEW: Bridge Event ======
type BridgeEvent struct {
	UserAddress   string `json:"user_address"`
	Amount        string `json:"amount"`
	Direction     string `json:"direction"` // "L1_TO_L2" or "L2_TO_L1"
	Status        string `json:"status"`    // "initiated", "pending", "confirmed", "failed"
	L1TxHash      string `json:"l1_tx_hash"`
	L2TxHash      string `json:"l2_tx_hash,omitempty"`
	MessageHash   string `json:"message_hash"`
	L1BlockNumber int64  `json:"l1_block_number"`
	L2BlockNumber int64  `json:"l2_block_number,omitempty"`
	InitiatedAt   int64  `json:"initiated_at"`
	ConfirmedAt   int64  `json:"confirmed_at,omitempty"`
	RetryCount    int    `json:"retry_count"`
	ErrorMsg      string `json:"error_msg,omitempty"`
}
