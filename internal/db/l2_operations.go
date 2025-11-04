package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"loyalty-points-system/internal/models"
)

// UpsertL2VaultPosition updates or inserts vault position
func UpsertL2VaultPosition(tx *sql.Tx, evt *models.L2Event, isDeposit bool) error {
	query := `
		INSERT INTO l2_vault_positions (user_address, deposited, shares, current_value, yield_earned, last_updated)
		VALUES ($1, $2, 0, $2, 0, NOW())
		ON CONFLICT (user_address)
		DO UPDATE SET
			deposited = CASE WHEN $3 THEN l2_vault_positions.deposited + CAST($2 AS NUMERIC)
			                  ELSE l2_vault_positions.deposited - CAST($2 AS NUMERIC) END,
			current_value = CASE WHEN $3 THEN l2_vault_positions.current_value + CAST($2 AS NUMERIC)
			                      ELSE l2_vault_positions.current_value - CAST($2 AS NUMERIC) END,
			last_updated = NOW()
	`
	_, err := tx.Exec(query, evt.UserAddress, evt.Amount, isDeposit)
	return err
}

// InsertL2StrategyAllocation inserts a strategy allocation record
func InsertL2StrategyAllocation(tx *sql.Tx, evt *models.L2Event) error {
	// Extract strategy and allocation from metadata
	strategy := "unknown"
	allocation := "0"

	if evt.Metadata != nil {
		if s, ok := evt.Metadata["strategy"].(string); ok {
			strategy = s
		}
		if a, ok := evt.Metadata["allocation"].(string); ok {
			allocation = a
		}
	}

	query := `
		INSERT INTO l2_strategy_allocations (
			strategy, allocation_percentage, allocated_amount, current_value, apy, last_updated
		) VALUES ($1, $2, $3, $3, 0, NOW())
		ON CONFLICT (strategy)
		DO UPDATE SET
			allocated_amount = CAST($3 AS NUMERIC),
			current_value = CAST($3 AS NUMERIC),
			last_updated = NOW()
	`
	_, err := tx.Exec(query, strategy, allocation, evt.Amount)
	return err
}

// InsertL2RWAAsset inserts an RWA asset record
func InsertL2RWAAsset(tx *sql.Tx, evt *models.L2Event) error {
	assetID := int64(0)
	assetName := "Unknown"
	assetType := "real_estate"
	totalSupply := evt.Amount

	if evt.Metadata != nil {
		if id, ok := evt.Metadata["asset_id"].(string); ok {
			// Parse asset_id
			fmt.Sscanf(id, "%d", &assetID)
		}
		if name, ok := evt.Metadata["asset_name"].(string); ok {
			assetName = name
		}
		if aType, ok := evt.Metadata["asset_type"].(string); ok {
			assetType = aType
		}
	}

	query := `
		INSERT INTO l2_rwa_assets (
			asset_id, asset_name, asset_type, total_supply, price_per_token,
			valuation, status, created_at
		) VALUES ($1, $2, $3, $4, 0, 0, 'active', NOW())
		ON CONFLICT (asset_id) DO UPDATE SET
			total_supply = CAST($4 AS NUMERIC),
			status = 'active'
	`
	_, err := tx.Exec(query, assetID, assetName, assetType, totalSupply)
	return err
}

// UpsertL2RWAHolding updates or inserts RWA holding
func UpsertL2RWAHolding(tx *sql.Tx, evt *models.L2Event, isIncrease bool) error {
	assetID := int64(0)
	if evt.Metadata != nil {
		if id, ok := evt.Metadata["asset_id"].(string); ok {
			fmt.Sscanf(id, "%d", &assetID)
		}
	}

	query := `
		INSERT INTO l2_rwa_holdings (user_address, asset_id, amount, purchase_price, current_value, updated_at)
		VALUES ($1, $2, $3, 0, 0, NOW())
		ON CONFLICT (user_address, asset_id)
		DO UPDATE SET
			amount = CASE WHEN $4 THEN l2_rwa_holdings.amount + CAST($3 AS NUMERIC)
			              ELSE l2_rwa_holdings.amount - CAST($3 AS NUMERIC) END,
			updated_at = NOW()
	`
	_, err := tx.Exec(query, evt.UserAddress, assetID, evt.Amount, isIncrease)
	return err
}

// InsertL2RWAListing inserts an RWA marketplace listing
func InsertL2RWAListing(tx *sql.Tx, evt *models.L2Event) error {
	listingID := int64(0)
	assetID := int64(0)
	price := "0"

	if evt.Metadata != nil {
		if lid, ok := evt.Metadata["listing_id"].(string); ok {
			fmt.Sscanf(lid, "%d", &listingID)
		}
		if aid, ok := evt.Metadata["asset_id"].(string); ok {
			fmt.Sscanf(aid, "%d", &assetID)
		}
		if p, ok := evt.Metadata["price"].(string); ok {
			price = p
		}
	}

	query := `
		INSERT INTO l2_rwa_listings (
			listing_id, asset_id, seller_address, amount, price_per_token,
			status, created_at
		) VALUES ($1, $2, $3, $4, $5, 'active', NOW())
		ON CONFLICT (listing_id) DO UPDATE SET
			status = 'active',
			amount = CAST($4 AS NUMERIC),
			price_per_token = CAST($5 AS NUMERIC)
	`
	_, err := tx.Exec(query, listingID, assetID, evt.UserAddress, evt.Amount, price)
	return err
}

// InsertL2RWAProposal inserts an RWA governance proposal
func InsertL2RWAProposal(tx *sql.Tx, evt *models.L2Event) error {
	proposalID := int64(0)
	assetID := int64(0)
	description := "Governance Proposal"

	if evt.Metadata != nil {
		if pid, ok := evt.Metadata["proposal_id"].(string); ok {
			fmt.Sscanf(pid, "%d", &proposalID)
		}
		if aid, ok := evt.Metadata["asset_id"].(string); ok {
			fmt.Sscanf(aid, "%d", &assetID)
		}
		if desc, ok := evt.Metadata["description"].(string); ok {
			description = desc
		}
	}

	query := `
		INSERT INTO l2_rwa_proposals (
			proposal_id, asset_id, proposer_address, description,
			votes_for, votes_against, status, created_at
		) VALUES ($1, $2, $3, $4, 0, 0, 'active', NOW())
		ON CONFLICT (proposal_id) DO NOTHING
	`
	_, err := tx.Exec(query, proposalID, assetID, evt.UserAddress, description)
	return err
}

// ProcessL2Event processes an L2 event and updates database
func ProcessL2Event(tx *sql.Tx, evt *models.L2Event) error {
	// Ensure user exists
	if err := ensureUserExists(tx, evt.UserAddress); err != nil {
		return fmt.Errorf("ensure user failed: %w", err)
	}

	// Insert into balance_events for backward compatibility
	if err := insertL2BalanceEvent(tx, evt); err != nil {
		return fmt.Errorf("insert balance event failed: %w", err)
	}

	// Process based on event type
	switch {
	case evt.EventType == "vault_deposit" || evt.EventType == "vault_operation":
		if err := UpsertL2VaultPosition(tx, evt, true); err != nil {
			return fmt.Errorf("upsert vault position failed: %w", err)
		}

	case evt.EventType == "vault_withdraw":
		if err := UpsertL2VaultPosition(tx, evt, false); err != nil {
			return fmt.Errorf("upsert vault position failed: %w", err)
		}

	case evt.EventType == "vault_allocate":
		if err := InsertL2StrategyAllocation(tx, evt); err != nil {
			return fmt.Errorf("insert strategy allocation failed: %w", err)
		}

	case evt.EventType == "rwa_factory":
		if err := InsertL2RWAAsset(tx, evt); err != nil {
			return fmt.Errorf("insert RWA asset failed: %w", err)
		}

	case evt.EventType == "rwa_marketplace":
		// Could be listing or purchase
		if err := InsertL2RWAListing(tx, evt); err != nil {
			return fmt.Errorf("insert RWA listing failed: %w", err)
		}
		// Also update holdings if it's a purchase
		if evt.Metadata != nil {
			if action, ok := evt.Metadata["action"].(string); ok && action == "purchase" {
				if err := UpsertL2RWAHolding(tx, evt, true); err != nil {
					return fmt.Errorf("upsert RWA holding failed: %w", err)
				}
			}
		}

	case evt.EventType == "rwa_governance":
		if err := InsertL2RWAProposal(tx, evt); err != nil {
			return fmt.Errorf("insert RWA proposal failed: %w", err)
		}

	// Treasury events
	case evt.EventType == "treasury_token_created":
		if err := ProcessTreasuryTokenCreated(tx, evt); err != nil {
			return fmt.Errorf("process treasury token created failed: %w", err)
		}

	case evt.EventType == "treasury_order_created":
		if err := ProcessTreasuryOrderCreated(tx, evt); err != nil {
			return fmt.Errorf("process treasury order created failed: %w", err)
		}

	case evt.EventType == "treasury_order_matched":
		if err := ProcessTreasuryOrderMatched(tx, evt); err != nil {
			return fmt.Errorf("process treasury order matched failed: %w", err)
		}

	case evt.EventType == "treasury_order_cancelled":
		if err := ProcessTreasuryOrderCancelled(tx, evt); err != nil {
			return fmt.Errorf("process treasury order cancelled failed: %w", err)
		}

	case evt.EventType == "treasury_yield_deposited":
		if err := ProcessTreasuryYieldDeposited(tx, evt); err != nil {
			return fmt.Errorf("process treasury yield deposited failed: %w", err)
		}

	case evt.EventType == "treasury_yield_claimed":
		if err := ProcessTreasuryYieldClaimed(tx, evt); err != nil {
			return fmt.Errorf("process treasury yield claimed failed: %w", err)
		}

	case evt.EventType == "treasury_price_updated":
		if err := ProcessTreasuryPriceUpdated(tx, evt); err != nil {
			return fmt.Errorf("process treasury price updated failed: %w", err)
		}

	default:
		// DeFi operations and other event types
		// Already logged in balance_events
	}

	return nil
}

// insertL2BalanceEvent inserts into balance_events table (backward compatibility)
func insertL2BalanceEvent(tx *sql.Tx, evt *models.L2Event) error {
	metadataJSON, _ := json.Marshal(evt.Metadata)

	query := `
		INSERT INTO balance_events (
			user_address, amount, event_type, tx_hash, chain, block_number, confirmed,
			layer, token, contract_address
		) VALUES ($1, $2, $3, $4, 'L2', $5, $6, 'L2', '', $7)
		ON CONFLICT (tx_hash, event_type) DO NOTHING
	`
	_, err := tx.Exec(query,
		evt.UserAddress,
		evt.Amount,
		evt.EventType,
		evt.TxHash,
		evt.BlockNumber,
		evt.Confirmed,
		evt.ContractAddress,
	)
	_ = metadataJSON // TODO: Store metadata in separate column if needed
	return err
}

// ============================================
// Treasury Event Processing Functions
// ============================================

// ProcessTreasuryTokenCreated handles treasury token creation events
func ProcessTreasuryTokenCreated(tx *sql.Tx, evt *models.L2Event) error {
	// This event is emitted when a new Treasury token is created
	// We don't need to insert into treasury_assets here because that's done via admin/oracle
	// Just log the event
	return nil
}

// ProcessTreasuryOrderCreated handles treasury order creation events
func ProcessTreasuryOrderCreated(tx *sql.Tx, evt *models.L2Event) error {
	if evt.Metadata == nil {
		return fmt.Errorf("missing metadata for treasury order created event")
	}

	orderID := int64(0)

	if id, ok := evt.Metadata["order_id"].(string); ok {
		fmt.Sscanf(id, "%d", &orderID)
	}

	// Update order status in database
	query := `
		UPDATE treasury_market_orders
		SET tx_hash = $1, status = 'open'
		WHERE order_id = $2
	`
	_, err := tx.Exec(query, evt.TxHash, orderID)
	return err
}

// ProcessTreasuryOrderMatched handles treasury order matching events
func ProcessTreasuryOrderMatched(tx *sql.Tx, evt *models.L2Event) error {
	if evt.Metadata == nil {
		return fmt.Errorf("missing metadata for treasury order matched event")
	}

	buyOrderID := int64(0)
	sellOrderID := int64(0)
	assetID := int64(0)
	matchedAmount := evt.Amount
	price := "0"

	if id, ok := evt.Metadata["buy_order_id"].(string); ok {
		fmt.Sscanf(id, "%d", &buyOrderID)
	}
	if id, ok := evt.Metadata["sell_order_id"].(string); ok {
		fmt.Sscanf(id, "%d", &sellOrderID)
	}
	if id, ok := evt.Metadata["asset_id"].(string); ok {
		fmt.Sscanf(id, "%d", &assetID)
	}
	if p, ok := evt.Metadata["price"].(string); ok {
		price = p
	}

	// Update buy order filled amount
	query1 := `
		UPDATE treasury_market_orders
		SET filled_amount = filled_amount + $1,
		    status = CASE
		        WHEN (token_amount - filled_amount - $1) <= 0 THEN 'filled'
		        ELSE 'partial'
		    END,
		    filled_at = CASE
		        WHEN (token_amount - filled_amount - $1) <= 0 THEN NOW()
		        ELSE filled_at
		    END
		WHERE order_id = $2
	`
	if _, err := tx.Exec(query1, matchedAmount, buyOrderID); err != nil {
		return err
	}

	// Update sell order filled amount
	if _, err := tx.Exec(query1, matchedAmount, sellOrderID); err != nil {
		return err
	}

	// Record the trade
	buyer := ""
	seller := ""
	if b, ok := evt.Metadata["buyer"].(string); ok {
		buyer = b
	}
	if s, ok := evt.Metadata["seller"].(string); ok {
		seller = s
	}

	tradeQuery := `
		INSERT INTO treasury_trades (
			asset_id, buy_order_id, sell_order_id, buyer_address, seller_address,
			token_amount, price_per_token, total_value, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	totalValue := fmt.Sprintf("%s", matchedAmount) // Calculate properly in production
	_, err := tx.Exec(tradeQuery, assetID, buyOrderID, sellOrderID, buyer, seller,
		matchedAmount, price, totalValue, evt.TxHash)

	// Update holdings for buyer and seller
	// Increase buyer's holdings
	holdingQuery := `
		INSERT INTO treasury_holdings (user_address, asset_id, tokens_held, last_updated)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_address, asset_id)
		DO UPDATE SET
			tokens_held = treasury_holdings.tokens_held + $3,
			last_updated = NOW()
	`
	if _, err := tx.Exec(holdingQuery, buyer, assetID, matchedAmount); err != nil {
		return err
	}

	// Decrease seller's holdings
	decreaseQuery := `
		UPDATE treasury_holdings
		SET tokens_held = GREATEST(0, tokens_held - $1),
		    last_updated = NOW()
		WHERE user_address = $2 AND asset_id = $3
	`
	if _, err := tx.Exec(decreaseQuery, matchedAmount, seller, assetID); err != nil {
		return err
	}

	return err
}

// ProcessTreasuryOrderCancelled handles treasury order cancellation events
func ProcessTreasuryOrderCancelled(tx *sql.Tx, evt *models.L2Event) error {
	if evt.Metadata == nil {
		return fmt.Errorf("missing metadata for treasury order cancelled event")
	}

	orderID := int64(0)
	if id, ok := evt.Metadata["order_id"].(string); ok {
		fmt.Sscanf(id, "%d", &orderID)
	}

	query := `
		UPDATE treasury_market_orders
		SET status = 'cancelled', cancelled_at = NOW()
		WHERE order_id = $1
	`
	_, err := tx.Exec(query, orderID)
	return err
}

// ProcessTreasuryYieldDeposited handles treasury yield deposit events
func ProcessTreasuryYieldDeposited(tx *sql.Tx, evt *models.L2Event) error {
	if evt.Metadata == nil {
		return fmt.Errorf("missing metadata for treasury yield deposited event")
	}

	assetID := int64(0)
	distributionType := "COUPON"
	yieldAmount := evt.Amount

	if id, ok := evt.Metadata["asset_id"].(string); ok {
		fmt.Sscanf(id, "%d", &assetID)
	}
	if dt, ok := evt.Metadata["distribution_type"].(string); ok {
		distributionType = dt
	}

	// Get total tokens outstanding for this asset
	var tokensOutstanding string
	err := tx.QueryRow(`SELECT tokens_outstanding FROM treasury_assets WHERE asset_id = $1`, assetID).Scan(&tokensOutstanding)
	if err != nil {
		return fmt.Errorf("failed to get tokens outstanding: %w", err)
	}

	// Calculate yield per token (simplified - should use proper math)
	yieldPerToken := "0" // Calculate: yieldAmount / tokensOutstanding

	// Insert yield distribution record
	query := `
		INSERT INTO treasury_yield_distributions (
			asset_id, distribution_date, distribution_type, total_yield,
			yield_per_token, status, tx_hash
		) VALUES ($1, NOW(), $2, $3, $4, 'pending', $5)
	`
	_, err = tx.Exec(query, assetID, distributionType, yieldAmount, yieldPerToken, evt.TxHash)
	return err
}

// ProcessTreasuryYieldClaimed handles treasury yield claim events
func ProcessTreasuryYieldClaimed(tx *sql.Tx, evt *models.L2Event) error {
	if evt.Metadata == nil {
		return fmt.Errorf("missing metadata for treasury yield claimed event")
	}

	assetID := int64(0)
	claimedAmount := evt.Amount
	userAddress := evt.UserAddress

	if id, ok := evt.Metadata["asset_id"].(string); ok {
		fmt.Sscanf(id, "%d", &assetID)
	}

	// Update user's accrued interest
	query := `
		UPDATE treasury_holdings
		SET accrued_interest = GREATEST(0, accrued_interest - $1),
		    last_updated = NOW()
		WHERE user_address = $2 AND asset_id = $3
	`
	_, err := tx.Exec(query, claimedAmount, userAddress, assetID)
	return err
}

// ProcessTreasuryPriceUpdated handles treasury price update events
func ProcessTreasuryPriceUpdated(tx *sql.Tx, evt *models.L2Event) error {
	if evt.Metadata == nil {
		return fmt.Errorf("missing metadata for treasury price updated event")
	}

	assetID := int64(0)
	price := "0"
	yieldRate := "0"

	if id, ok := evt.Metadata["asset_id"].(string); ok {
		fmt.Sscanf(id, "%d", &assetID)
	}
	if p, ok := evt.Metadata["price"].(string); ok {
		price = p
	} else {
		price = evt.Amount
	}
	if y, ok := evt.Metadata["yield"].(string); ok {
		yieldRate = y
	}

	// Update asset price and yield
	updateQuery := `
		UPDATE treasury_assets
		SET current_price = $1,
		    current_yield = $2,
		    last_price_update = NOW(),
		    updated_at = NOW()
		WHERE asset_id = $3
	`
	if _, err := tx.Exec(updateQuery, price, yieldRate, assetID); err != nil {
		return err
	}

	// Insert into price history
	historyQuery := `
		INSERT INTO treasury_price_history (asset_id, price, yield, source, timestamp)
		VALUES ($1, $2, $3, 'chainlink', NOW())
	`
	_, err := tx.Exec(historyQuery, assetID, price, yieldRate)
	return err
}
