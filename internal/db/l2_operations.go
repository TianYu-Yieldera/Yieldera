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
