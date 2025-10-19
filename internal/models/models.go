package models
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
