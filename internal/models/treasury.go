package models

import "time"

// TreasuryAsset represents a tokenized US Treasury security
type TreasuryAsset struct {
	AssetID           int64     `json:"asset_id" db:"asset_id"`
	TreasuryType      string    `json:"treasury_type" db:"treasury_type"`       // T-BILL, T-NOTE, T-BOND
	MaturityTerm      string    `json:"maturity_term" db:"maturity_term"`       // 13W, 2Y, 10Y, 30Y
	CUSIP             string    `json:"cusip" db:"cusip"`                       // Unique identifier
	IssueDate         time.Time `json:"issue_date" db:"issue_date"`
	MaturityDate      time.Time `json:"maturity_date" db:"maturity_date"`
	FaceValue         string    `json:"face_value" db:"face_value"`
	CouponRate        string    `json:"coupon_rate" db:"coupon_rate"`
	CurrentPrice      *string   `json:"current_price,omitempty" db:"current_price"`
	CurrentYield      *string   `json:"current_yield,omitempty" db:"current_yield"`
	TokensIssued      string    `json:"tokens_issued" db:"tokens_issued"`
	TokensOutstanding string    `json:"tokens_outstanding" db:"tokens_outstanding"`
	TokenAddress      *string   `json:"token_address,omitempty" db:"token_address"`
	Status            string    `json:"status" db:"status"` // active, matured, suspended
	LastPriceUpdate   *time.Time `json:"last_price_update,omitempty" db:"last_price_update"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// TreasuryHolding represents a user's holdings of treasury tokens
type TreasuryHolding struct {
	ID                int64     `json:"id" db:"id"`
	UserAddress       string    `json:"user_address" db:"user_address"`
	AssetID           int64     `json:"asset_id" db:"asset_id"`
	TokensHeld        string    `json:"tokens_held" db:"tokens_held"`
	AvgPurchasePrice  *string   `json:"avg_purchase_price,omitempty" db:"avg_purchase_price"`
	TotalInvested     *string   `json:"total_invested,omitempty" db:"total_invested"`
	CurrentValue      *string   `json:"current_value,omitempty" db:"current_value"`
	UnrealizedGain    *string   `json:"unrealized_gain,omitempty" db:"unrealized_gain"`
	AccruedInterest   string    `json:"accrued_interest" db:"accrued_interest"`
	LastUpdated       time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// TreasuryMarketOrder represents a buy/sell order in the secondary market
type TreasuryMarketOrder struct {
	OrderID       int64      `json:"order_id" db:"order_id"`
	AssetID       int64      `json:"asset_id" db:"asset_id"`
	OrderType     string     `json:"order_type" db:"order_type"` // BUY, SELL
	UserAddress   string     `json:"user_address" db:"user_address"`
	TokenAmount   string     `json:"token_amount" db:"token_amount"`
	PricePerToken string     `json:"price_per_token" db:"price_per_token"`
	TotalValue    string     `json:"total_value" db:"total_value"`
	FilledAmount  string     `json:"filled_amount" db:"filled_amount"`
	Status        string     `json:"status" db:"status"` // open, partial, filled, cancelled
	TxHash        *string    `json:"tx_hash,omitempty" db:"tx_hash"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	FilledAt      *time.Time `json:"filled_at,omitempty" db:"filled_at"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
}

// TreasuryYieldDistribution represents a yield/coupon payment distribution
type TreasuryYieldDistribution struct {
	ID                int64      `json:"id" db:"id"`
	AssetID           int64      `json:"asset_id" db:"asset_id"`
	DistributionDate  time.Time  `json:"distribution_date" db:"distribution_date"`
	DistributionType  string     `json:"distribution_type" db:"distribution_type"` // COUPON, MATURITY
	TotalYield        string     `json:"total_yield" db:"total_yield"`
	YieldPerToken     string     `json:"yield_per_token" db:"yield_per_token"`
	RecipientsCount   int        `json:"recipients_count" db:"recipients_count"`
	TotalDistributed  string     `json:"total_distributed" db:"total_distributed"`
	Status            string     `json:"status" db:"status"` // pending, completed
	TxHash            *string    `json:"tx_hash,omitempty" db:"tx_hash"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	DistributedAt     *time.Time `json:"distributed_at,omitempty" db:"distributed_at"`
}

// TreasuryPriceHistory represents historical price data
type TreasuryPriceHistory struct {
	ID        int64     `json:"id" db:"id"`
	AssetID   int64     `json:"asset_id" db:"asset_id"`
	Price     string    `json:"price" db:"price"`
	Yield     string    `json:"yield" db:"yield"`
	Source    *string   `json:"source,omitempty" db:"source"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// TreasuryTrade represents an executed trade
type TreasuryTrade struct {
	TradeID       int64     `json:"trade_id" db:"trade_id"`
	AssetID       int64     `json:"asset_id" db:"asset_id"`
	BuyOrderID    *int64    `json:"buy_order_id,omitempty" db:"buy_order_id"`
	SellOrderID   *int64    `json:"sell_order_id,omitempty" db:"sell_order_id"`
	BuyerAddress  string    `json:"buyer_address" db:"buyer_address"`
	SellerAddress string    `json:"seller_address" db:"seller_address"`
	TokenAmount   string    `json:"token_amount" db:"token_amount"`
	PricePerToken string    `json:"price_per_token" db:"price_per_token"`
	TotalValue    string    `json:"total_value" db:"total_value"`
	FeeAmount     string    `json:"fee_amount" db:"fee_amount"`
	TxHash        *string   `json:"tx_hash,omitempty" db:"tx_hash"`
	ExecutedAt    time.Time `json:"executed_at" db:"executed_at"`
}

// TreasuryAssetSummary is a view combining asset and market data
type TreasuryAssetSummary struct {
	AssetID           int64      `json:"asset_id" db:"asset_id"`
	TreasuryType      string     `json:"treasury_type" db:"treasury_type"`
	MaturityTerm      string     `json:"maturity_term" db:"maturity_term"`
	CUSIP             string     `json:"cusip" db:"cusip"`
	CurrentPrice      *string    `json:"current_price,omitempty" db:"current_price"`
	CurrentYield      *string    `json:"current_yield,omitempty" db:"current_yield"`
	TokensOutstanding string     `json:"tokens_outstanding" db:"tokens_outstanding"`
	TotalMarketValue  *string    `json:"total_market_value,omitempty" db:"total_market_value"`
	MaturityDate      time.Time  `json:"maturity_date" db:"maturity_date"`
	Status            string     `json:"status" db:"status"`
	UniqueHolders     *int       `json:"unique_holders,omitempty" db:"unique_holders"`
	ActiveOrders      *int       `json:"active_orders,omitempty" db:"active_orders"`
}

// UserTreasuryPortfolio represents a user's complete portfolio
type UserTreasuryPortfolio struct {
	UserAddress       string     `json:"user_address" db:"user_address"`
	AssetID           int64      `json:"asset_id" db:"asset_id"`
	TreasuryType      string     `json:"treasury_type" db:"treasury_type"`
	MaturityTerm      string     `json:"maturity_term" db:"maturity_term"`
	CUSIP             string     `json:"cusip" db:"cusip"`
	TokensHeld        string     `json:"tokens_held" db:"tokens_held"`
	AvgPurchasePrice  *string    `json:"avg_purchase_price,omitempty" db:"avg_purchase_price"`
	TotalInvested     *string    `json:"total_invested,omitempty" db:"total_invested"`
	CurrentPrice      *string    `json:"current_price,omitempty" db:"current_price"`
	CurrentValue      *string    `json:"current_value,omitempty" db:"current_value"`
	UnrealizedGainLoss *string   `json:"unrealized_gain_loss,omitempty" db:"unrealized_gain_loss"`
	AccruedInterest   string     `json:"accrued_interest" db:"accrued_interest"`
	MaturityDate      time.Time  `json:"maturity_date" db:"maturity_date"`
}

// MarketOrderBook represents the order book for an asset
type MarketOrderBook struct {
	OrderID         int64      `json:"order_id" db:"order_id"`
	AssetID         int64      `json:"asset_id" db:"asset_id"`
	CUSIP           string     `json:"cusip" db:"cusip"`
	TreasuryType    string     `json:"treasury_type" db:"treasury_type"`
	OrderType       string     `json:"order_type" db:"order_type"`
	UserAddress     string     `json:"user_address" db:"user_address"`
	TokenAmount     string     `json:"token_amount" db:"token_amount"`
	FilledAmount    string     `json:"filled_amount" db:"filled_amount"`
	RemainingAmount string     `json:"remaining_amount" db:"remaining_amount"`
	PricePerToken   string     `json:"price_per_token" db:"price_per_token"`
	TotalValue      string     `json:"total_value" db:"total_value"`
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}
