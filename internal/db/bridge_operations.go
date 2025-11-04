package db

import (
	"database/sql"
	"fmt"

	"loyalty-points-system/internal/models"
)

// UpsertBridgeMessage inserts or updates a bridge message
func UpsertBridgeMessage(tx *sql.Tx, evt *models.BridgeEvent) error {
	query := `
		INSERT INTO bridge_messages (
			message_hash, direction, user_address, amount, status,
			l1_tx_hash, l2_tx_hash, l1_block_number, l2_block_number,
			initiated_at, confirmed_at, retry_count, error_msg
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,
		          to_timestamp($10),
		          CASE WHEN $11 > 0 THEN to_timestamp($11) ELSE NULL END,
		          $12, $13)
		ON CONFLICT (message_hash)
		DO UPDATE SET
			status = $5,
			l1_tx_hash = COALESCE($6, bridge_messages.l1_tx_hash),
			l2_tx_hash = COALESCE($7, bridge_messages.l2_tx_hash),
			l1_block_number = CASE WHEN $8 > 0 THEN $8 ELSE bridge_messages.l1_block_number END,
			l2_block_number = CASE WHEN $9 > 0 THEN $9 ELSE bridge_messages.l2_block_number END,
			confirmed_at = CASE WHEN $11 > 0 THEN to_timestamp($11) ELSE bridge_messages.confirmed_at END,
			retry_count = $12,
			error_msg = $13,
			updated_at = NOW()
	`

	_, err := tx.Exec(query,
		evt.MessageHash,
		evt.Direction,
		evt.UserAddress,
		evt.Amount,
		evt.Status,
		evt.L1TxHash,
		evt.L2TxHash,
		evt.L1BlockNumber,
		evt.L2BlockNumber,
		evt.InitiatedAt,
		evt.ConfirmedAt,
		evt.RetryCount,
		evt.ErrorMsg,
	)

	return err
}

// ProcessBridgeEvent processes a bridge event and updates database
func ProcessBridgeEvent(tx *sql.Tx, evt *models.BridgeEvent) error {
	// Ensure user exists
	if err := ensureUserExists(tx, evt.UserAddress); err != nil {
		return fmt.Errorf("ensure user failed: %w", err)
	}

	// Upsert bridge message
	if err := UpsertBridgeMessage(tx, evt); err != nil {
		return fmt.Errorf("upsert bridge message failed: %w", err)
	}

	// If bridge is confirmed and it's L1→L2, we might want to update L2 balances
	// This depends on your business logic
	if evt.Status == "confirmed" && evt.Direction == "L1_TO_L2" {
		// L1 collateral locked, L2 vault credited
		// This should already be handled by L2 vault events
	}

	return nil
}

// GetPendingBridgeMessages returns pending bridge messages
func GetPendingBridgeMessages(db *sql.DB) ([]*models.BridgeEvent, error) {
	query := `
		SELECT message_hash, direction, user_address, amount, status,
		       l1_tx_hash, l2_tx_hash, l1_block_number, l2_block_number,
		       EXTRACT(EPOCH FROM initiated_at)::bigint,
		       COALESCE(EXTRACT(EPOCH FROM confirmed_at)::bigint, 0),
		       retry_count, COALESCE(error_msg, '')
		FROM bridge_messages
		WHERE status IN ('initiated', 'pending')
		  AND retry_count < 10
		ORDER BY initiated_at ASC
		LIMIT 100
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.BridgeEvent
	for rows.Next() {
		evt := &models.BridgeEvent{}
		err := rows.Scan(
			&evt.MessageHash,
			&evt.Direction,
			&evt.UserAddress,
			&evt.Amount,
			&evt.Status,
			&evt.L1TxHash,
			&evt.L2TxHash,
			&evt.L1BlockNumber,
			&evt.L2BlockNumber,
			&evt.InitiatedAt,
			&evt.ConfirmedAt,
			&evt.RetryCount,
			&evt.ErrorMsg,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, evt)
	}

	return messages, rows.Err()
}

// MarkBridgeMessageFailed marks a bridge message as failed
func MarkBridgeMessageFailed(db *sql.DB, messageHash string, errorMsg string) error {
	query := `
		UPDATE bridge_messages
		SET status = 'failed', error_msg = $2, updated_at = NOW()
		WHERE message_hash = $1
	`
	_, err := db.Exec(query, messageHash, errorMsg)
	return err
}
