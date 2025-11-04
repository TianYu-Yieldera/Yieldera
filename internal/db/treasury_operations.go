package db

import (
	"database/sql"
	"fmt"

	"loyalty-points-system/internal/models"
)

// InsertTreasuryAsset inserts a new treasury asset
func InsertTreasuryAsset(db *sql.DB, asset *models.TreasuryAsset) error {
	query := `
		INSERT INTO treasury_assets (
			treasury_type, maturity_term, cusip, issue_date, maturity_date,
			face_value, coupon_rate, tokens_issued, tokens_outstanding,
			token_address, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (cusip) DO UPDATE SET
			current_price = EXCLUDED.current_price,
			current_yield = EXCLUDED.current_yield,
			tokens_issued = EXCLUDED.tokens_issued,
			tokens_outstanding = EXCLUDED.tokens_outstanding,
			updated_at = NOW()
		RETURNING asset_id, created_at, updated_at
	`

	return db.QueryRow(
		query,
		asset.TreasuryType,
		asset.MaturityTerm,
		asset.CUSIP,
		asset.IssueDate,
		asset.MaturityDate,
		asset.FaceValue,
		asset.CouponRate,
		asset.TokensIssued,
		asset.TokensOutstanding,
		asset.TokenAddress,
		asset.Status,
	).Scan(&asset.AssetID, &asset.CreatedAt, &asset.UpdatedAt)
}

// UpdateTreasuryPrice updates the current price and yield for an asset
func UpdateTreasuryPrice(db *sql.DB, assetID int64, price, yield string, source string) error {
	// Update the asset table
	query := `
		UPDATE treasury_assets
		SET current_price = $1,
		    current_yield = $2,
		    last_price_update = NOW(),
		    updated_at = NOW()
		WHERE asset_id = $3
	`

	_, err := db.Exec(query, price, yield, assetID)
	if err != nil {
		return err
	}

	// Insert price history
	historyQuery := `
		INSERT INTO treasury_price_history (asset_id, price, yield, source, timestamp)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err = db.Exec(historyQuery, assetID, price, yield, source)
	return err
}

// UpsertTreasuryHolding inserts or updates a user's treasury holding
func UpsertTreasuryHolding(db *sql.DB, holding *models.TreasuryHolding) error {
	query := `
		INSERT INTO treasury_holdings (
			user_address, asset_id, tokens_held, avg_purchase_price,
			total_invested, current_value, unrealized_gain, accrued_interest
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_address, asset_id) DO UPDATE SET
			tokens_held = EXCLUDED.tokens_held,
			avg_purchase_price = EXCLUDED.avg_purchase_price,
			total_invested = EXCLUDED.total_invested,
			current_value = EXCLUDED.current_value,
			unrealized_gain = EXCLUDED.unrealized_gain,
			accrued_interest = EXCLUDED.accrued_interest,
			last_updated = NOW()
		RETURNING id, last_updated, created_at
	`

	return db.QueryRow(
		query,
		holding.UserAddress,
		holding.AssetID,
		holding.TokensHeld,
		holding.AvgPurchasePrice,
		holding.TotalInvested,
		holding.CurrentValue,
		holding.UnrealizedGain,
		holding.AccruedInterest,
	).Scan(&holding.ID, &holding.LastUpdated, &holding.CreatedAt)
}

// InsertTreasuryMarketOrder inserts a new market order
func InsertTreasuryMarketOrder(db *sql.DB, order *models.TreasuryMarketOrder) error {
	query := `
		INSERT INTO treasury_market_orders (
			asset_id, order_type, user_address, token_amount,
			price_per_token, total_value, filled_amount, status, tx_hash, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING order_id, created_at
	`

	return db.QueryRow(
		query,
		order.AssetID,
		order.OrderType,
		order.UserAddress,
		order.TokenAmount,
		order.PricePerToken,
		order.TotalValue,
		order.FilledAmount,
		order.Status,
		order.TxHash,
		order.ExpiresAt,
	).Scan(&order.OrderID, &order.CreatedAt)
}

// UpdateTreasuryOrderStatus updates order status and filled amount
func UpdateTreasuryOrderStatus(db *sql.DB, orderID int64, status string, filledAmount string) error {
	query := `
		UPDATE treasury_market_orders
		SET status = $1,
		    filled_amount = $2,
		    filled_at = CASE WHEN $1 = 'filled' THEN NOW() ELSE filled_at END
		WHERE order_id = $3
	`

	_, err := db.Exec(query, status, filledAmount, orderID)
	return err
}

// CancelTreasuryOrder cancels an order
func CancelTreasuryOrder(db *sql.DB, orderID int64) error {
	query := `
		UPDATE treasury_market_orders
		SET status = 'cancelled',
		    cancelled_at = NOW()
		WHERE order_id = $1 AND status IN ('open', 'partial')
	`

	result, err := db.Exec(query, orderID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found or not cancellable")
	}

	return nil
}

// InsertTreasuryTrade records a trade execution
func InsertTreasuryTrade(db *sql.DB, trade *models.TreasuryTrade) error {
	query := `
		INSERT INTO treasury_trades (
			asset_id, buy_order_id, sell_order_id, buyer_address, seller_address,
			token_amount, price_per_token, total_value, fee_amount, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING trade_id, executed_at
	`

	return db.QueryRow(
		query,
		trade.AssetID,
		trade.BuyOrderID,
		trade.SellOrderID,
		trade.BuyerAddress,
		trade.SellerAddress,
		trade.TokenAmount,
		trade.PricePerToken,
		trade.TotalValue,
		trade.FeeAmount,
		trade.TxHash,
	).Scan(&trade.TradeID, &trade.ExecutedAt)
}

// InsertYieldDistribution records a yield distribution
func InsertYieldDistribution(db *sql.DB, dist *models.TreasuryYieldDistribution) error {
	query := `
		INSERT INTO treasury_yield_distributions (
			asset_id, distribution_date, distribution_type, total_yield,
			yield_per_token, recipients_count, total_distributed, status, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	return db.QueryRow(
		query,
		dist.AssetID,
		dist.DistributionDate,
		dist.DistributionType,
		dist.TotalYield,
		dist.YieldPerToken,
		dist.RecipientsCount,
		dist.TotalDistributed,
		dist.Status,
		dist.TxHash,
	).Scan(&dist.ID, &dist.CreatedAt)
}

// UpdateYieldDistributionStatus updates distribution status
func UpdateYieldDistributionStatus(db *sql.DB, distID int64, status string) error {
	query := `
		UPDATE treasury_yield_distributions
		SET status = $1,
		    distributed_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE distributed_at END
		WHERE id = $2
	`

	_, err := db.Exec(query, status, distID)
	return err
}

// ============================================
// QUERY FUNCTIONS
// ============================================

// GetTreasuryAsset retrieves a treasury asset by ID
func GetTreasuryAsset(db *sql.DB, assetID int64) (*models.TreasuryAsset, error) {
	query := `
		SELECT asset_id, treasury_type, maturity_term, cusip, issue_date, maturity_date,
		       face_value, coupon_rate, current_price, current_yield,
		       tokens_issued, tokens_outstanding, token_address, status,
		       last_price_update, created_at, updated_at
		FROM treasury_assets
		WHERE asset_id = $1
	`

	var asset models.TreasuryAsset
	err := db.QueryRow(query, assetID).Scan(
		&asset.AssetID,
		&asset.TreasuryType,
		&asset.MaturityTerm,
		&asset.CUSIP,
		&asset.IssueDate,
		&asset.MaturityDate,
		&asset.FaceValue,
		&asset.CouponRate,
		&asset.CurrentPrice,
		&asset.CurrentYield,
		&asset.TokensIssued,
		&asset.TokensOutstanding,
		&asset.TokenAddress,
		&asset.Status,
		&asset.LastPriceUpdate,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAllTreasuryAssets retrieves all active treasury assets
func GetAllTreasuryAssets(db *sql.DB, treasuryType *string) ([]*models.TreasuryAsset, error) {
	query := `
		SELECT asset_id, treasury_type, maturity_term, cusip, issue_date, maturity_date,
		       face_value, coupon_rate, current_price, current_yield,
		       tokens_issued, tokens_outstanding, token_address, status,
		       last_price_update, created_at, updated_at
		FROM treasury_assets
		WHERE status = 'active'
	`

	args := []interface{}{}
	argCount := 1

	if treasuryType != nil && *treasuryType != "" {
		query += fmt.Sprintf(" AND treasury_type = $%d", argCount)
		args = append(args, *treasuryType)
		argCount++
	}

	query += " ORDER BY maturity_date ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*models.TreasuryAsset
	for rows.Next() {
		var asset models.TreasuryAsset
		err := rows.Scan(
			&asset.AssetID,
			&asset.TreasuryType,
			&asset.MaturityTerm,
			&asset.CUSIP,
			&asset.IssueDate,
			&asset.MaturityDate,
			&asset.FaceValue,
			&asset.CouponRate,
			&asset.CurrentPrice,
			&asset.CurrentYield,
			&asset.TokensIssued,
			&asset.TokensOutstanding,
			&asset.TokenAddress,
			&asset.Status,
			&asset.LastPriceUpdate,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetUserTreasuryHoldings retrieves all holdings for a user
func GetUserTreasuryHoldings(db *sql.DB, userAddress string) ([]*models.UserTreasuryPortfolio, error) {
	query := `
		SELECT h.user_address, h.asset_id, a.treasury_type, a.maturity_term, a.cusip,
		       h.tokens_held, h.avg_purchase_price, h.total_invested, a.current_price,
		       h.current_value, h.unrealized_gain, h.accrued_interest, a.maturity_date
		FROM treasury_holdings h
		JOIN treasury_assets a ON h.asset_id = a.asset_id
		WHERE h.user_address = $1 AND h.tokens_held > 0
		ORDER BY a.maturity_date ASC
	`

	rows, err := db.Query(query, userAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolio []*models.UserTreasuryPortfolio
	for rows.Next() {
		var item models.UserTreasuryPortfolio
		err := rows.Scan(
			&item.UserAddress,
			&item.AssetID,
			&item.TreasuryType,
			&item.MaturityTerm,
			&item.CUSIP,
			&item.TokensHeld,
			&item.AvgPurchasePrice,
			&item.TotalInvested,
			&item.CurrentPrice,
			&item.CurrentValue,
			&item.UnrealizedGainLoss,
			&item.AccruedInterest,
			&item.MaturityDate,
		)
		if err != nil {
			return nil, err
		}
		portfolio = append(portfolio, &item)
	}

	return portfolio, nil
}

// GetTreasuryMarketOrders retrieves open orders for an asset
func GetTreasuryMarketOrders(db *sql.DB, assetID int64, orderType *string) ([]*models.MarketOrderBook, error) {
	query := `
		SELECT o.order_id, o.asset_id, a.cusip, a.treasury_type, o.order_type,
		       o.user_address, o.token_amount, o.filled_amount,
		       (o.token_amount - o.filled_amount) AS remaining_amount,
		       o.price_per_token, o.total_value, o.status, o.created_at, o.expires_at
		FROM treasury_market_orders o
		JOIN treasury_assets a ON o.asset_id = a.asset_id
		WHERE o.asset_id = $1 AND o.status IN ('open', 'partial') AND o.expires_at > NOW()
	`

	args := []interface{}{assetID}
	argCount := 2

	if orderType != nil && *orderType != "" {
		query += fmt.Sprintf(" AND o.order_type = $%d", argCount)
		args = append(args, *orderType)
		argCount++
	}

	query += ` ORDER BY
		CASE WHEN o.order_type = 'BUY' THEN o.price_per_token END DESC,
		CASE WHEN o.order_type = 'SELL' THEN o.price_per_token END ASC,
		o.created_at ASC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.MarketOrderBook
	for rows.Next() {
		var order models.MarketOrderBook
		err := rows.Scan(
			&order.OrderID,
			&order.AssetID,
			&order.CUSIP,
			&order.TreasuryType,
			&order.OrderType,
			&order.UserAddress,
			&order.TokenAmount,
			&order.FilledAmount,
			&order.RemainingAmount,
			&order.PricePerToken,
			&order.TotalValue,
			&order.Status,
			&order.CreatedAt,
			&order.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// GetTreasuryPriceHistory retrieves price history for an asset
func GetTreasuryPriceHistory(db *sql.DB, assetID int64, limit int) ([]*models.TreasuryPriceHistory, error) {
	query := `
		SELECT id, asset_id, price, yield, source, timestamp
		FROM treasury_price_history
		WHERE asset_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := db.Query(query, assetID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.TreasuryPriceHistory
	for rows.Next() {
		var item models.TreasuryPriceHistory
		err := rows.Scan(
			&item.ID,
			&item.AssetID,
			&item.Price,
			&item.Yield,
			&item.Source,
			&item.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, &item)
	}

	return history, nil
}

// GetTreasuryTrades retrieves recent trades for an asset
func GetTreasuryTrades(db *sql.DB, assetID int64, limit int) ([]*models.TreasuryTrade, error) {
	query := `
		SELECT trade_id, asset_id, buy_order_id, sell_order_id, buyer_address, seller_address,
		       token_amount, price_per_token, total_value, fee_amount, tx_hash, executed_at
		FROM treasury_trades
		WHERE asset_id = $1
		ORDER BY executed_at DESC
		LIMIT $2
	`

	rows, err := db.Query(query, assetID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*models.TreasuryTrade
	for rows.Next() {
		var trade models.TreasuryTrade
		err := rows.Scan(
			&trade.TradeID,
			&trade.AssetID,
			&trade.BuyOrderID,
			&trade.SellOrderID,
			&trade.BuyerAddress,
			&trade.SellerAddress,
			&trade.TokenAmount,
			&trade.PricePerToken,
			&trade.TotalValue,
			&trade.FeeAmount,
			&trade.TxHash,
			&trade.ExecutedAt,
		)
		if err != nil {
			return nil, err
		}
		trades = append(trades, &trade)
	}

	return trades, nil
}
