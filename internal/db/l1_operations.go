package db

import (
	"database/sql"
	"fmt"

	"loyalty-points-system/internal/models"
)

// InsertL1CollateralDeposit inserts a collateral deposit record
func InsertL1CollateralDeposit(tx *sql.Tx, evt *models.L1Event) error {
	query := `
		INSERT INTO l1_collateral_deposits (
			user_address, token, amount, tx_hash, block_number, confirmed, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (tx_hash) DO NOTHING
	`
	_, err := tx.Exec(query, evt.UserAddress, evt.Token, evt.Amount, evt.TxHash, evt.BlockNumber, evt.Confirmed)
	return err
}

// UpsertL1CollateralBalance updates or inserts collateral balance
func UpsertL1CollateralBalance(tx *sql.Tx, evt *models.L1Event, isDeposit bool) error {
	query := `
		INSERT INTO l1_collateral_balances (user_address, token, amount, usd_value, updated_at)
		VALUES ($1, $2, $3, 0, NOW())
		ON CONFLICT (user_address, token)
		DO UPDATE SET
			amount = CASE WHEN $4 THEN l1_collateral_balances.amount + CAST($3 AS NUMERIC)
			              ELSE l1_collateral_balances.amount - CAST($3 AS NUMERIC) END,
			updated_at = NOW()
	`
	_, err := tx.Exec(query, evt.UserAddress, evt.Token, evt.Amount, isDeposit)
	return err
}

// InsertL1StateSnapshot inserts a state snapshot
func InsertL1StateSnapshot(tx *sql.Tx, evt *models.L1Event) error {
	query := `
		INSERT INTO l1_state_snapshots (
			l2_block_number, state_root, tx_hash, block_number, created_at
		) VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (l2_block_number) DO NOTHING
	`
	// Extract L2 block number from metadata if available
	l2BlockNumber := int64(0)
	// TODO: Parse from event data when available

	_, err := tx.Exec(query, l2BlockNumber, "", evt.TxHash, evt.BlockNumber)
	return err
}

// ProcessL1Event processes an L1 event and updates database
func ProcessL1Event(tx *sql.Tx, evt *models.L1Event) error {
	// Ensure user exists
	if err := ensureUserExists(tx, evt.UserAddress); err != nil {
		return fmt.Errorf("ensure user failed: %w", err)
	}

	// Insert into balance_events for backward compatibility
	if err := insertBalanceEvent(tx, evt); err != nil {
		return fmt.Errorf("insert balance event failed: %w", err)
	}

	// Process based on event type
	switch evt.EventType {
	case "collateral_deposit", "collateral_operation":
		// Insert deposit record
		if err := InsertL1CollateralDeposit(tx, evt); err != nil {
			return fmt.Errorf("insert deposit failed: %w", err)
		}
		// Update balance
		if err := UpsertL1CollateralBalance(tx, evt, true); err != nil {
			return fmt.Errorf("update balance failed: %w", err)
		}

	case "collateral_withdraw":
		// Insert withdrawal record (can reuse deposit table with negative amount)
		if err := InsertL1CollateralDeposit(tx, evt); err != nil {
			return fmt.Errorf("insert withdrawal failed: %w", err)
		}
		// Update balance
		if err := UpsertL1CollateralBalance(tx, evt, false); err != nil {
			return fmt.Errorf("update balance failed: %w", err)
		}

	case "state_update":
		// Insert state snapshot
		if err := InsertL1StateSnapshot(tx, evt); err != nil {
			return fmt.Errorf("insert state snapshot failed: %w", err)
		}

	case "loyaltyusd_transfer", "gateway_message":
		// These are logged in balance_events only
		// No specific table updates needed

	default:
		// Unknown event type, log but don't fail
		// Already inserted into balance_events
	}

	return nil
}

// insertBalanceEvent inserts into balance_events table (backward compatibility)
func insertBalanceEvent(tx *sql.Tx, evt *models.L1Event) error {
	query := `
		INSERT INTO balance_events (
			user_address, amount, event_type, tx_hash, chain, block_number, confirmed,
			layer, token, contract_address
		) VALUES ($1, $2, $3, $4, 'L1', $5, $6, 'L1', $7, $8)
		ON CONFLICT (tx_hash, event_type) DO NOTHING
	`
	_, err := tx.Exec(query,
		evt.UserAddress,
		evt.Amount,
		evt.EventType,
		evt.TxHash,
		evt.BlockNumber,
		evt.Confirmed,
		evt.Token,
		evt.ContractAddress,
	)
	return err
}

// ensureUserExists ensures user record exists
func ensureUserExists(tx *sql.Tx, address string) error {
	query := `INSERT INTO users(address) VALUES($1) ON CONFLICT (address) DO NOTHING`
	_, err := tx.Exec(query, address)
	return err
}
