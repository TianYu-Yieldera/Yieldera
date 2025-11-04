package bridge

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"loyalty-points-system/internal/db"
	"loyalty-points-system/internal/models"
)

// Monitor tracks cross-chain bridge messages and handles retries
type Monitor struct {
	database     *sql.DB
	checkInterval time.Duration
	messageTimeout time.Duration
	maxRetries    int
	mu           sync.RWMutex
	running      bool
}

// NewMonitor creates a new bridge monitor
func NewMonitor(database *sql.DB) *Monitor {
	return &Monitor{
		database:      database,
		checkInterval: 30 * time.Second,      // Check every 30 seconds
		messageTimeout: 15 * time.Minute,     // Timeout after 15 minutes
		maxRetries:    10,                    // Max 10 retry attempts
		running:       false,
	}
}

// Start begins monitoring bridge messages
func (m *Monitor) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return nil
	}
	m.running = true
	m.mu.Unlock()

	log.Println("🌉 [Bridge Monitor] Started")

	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 [Bridge Monitor] Stopped")
			return nil
		case <-ticker.C:
			if err := m.processPendingMessages(ctx); err != nil {
				log.Printf("❌ [Bridge Monitor] Error processing pending messages: %v", err)
			}
		}
	}
}

// Stop halts the bridge monitor
func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = false
}

// processPendingMessages checks and retries pending bridge messages
func (m *Monitor) processPendingMessages(ctx context.Context) error {
	messages, err := db.GetPendingBridgeMessages(m.database)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	log.Printf("🔍 [Bridge Monitor] Found %d pending messages", len(messages))

	for _, msg := range messages {
		if err := m.processMessage(ctx, msg); err != nil {
			log.Printf("❌ [Bridge Monitor] Error processing message %s: %v", msg.MessageHash, err)
		}
	}

	return nil
}

// processMessage handles a single pending bridge message
func (m *Monitor) processMessage(ctx context.Context, msg *models.BridgeEvent) error {
	// Check if message has timed out
	if m.isTimedOut(msg) {
		log.Printf("⏰ [Bridge Monitor] Message %s timed out (initiated %s ago)",
			msg.MessageHash, time.Since(time.Unix(msg.InitiatedAt, 0)))
		return db.MarkBridgeMessageFailed(m.database, msg.MessageHash, "Message timed out after 15 minutes")
	}

	// Check if max retries exceeded
	if msg.RetryCount >= m.maxRetries {
		log.Printf("🔄 [Bridge Monitor] Message %s exceeded max retries (%d/%d)",
			msg.MessageHash, msg.RetryCount, m.maxRetries)
		return db.MarkBridgeMessageFailed(m.database, msg.MessageHash, "Max retry attempts exceeded")
	}

	// Attempt to check message status
	// In production, this would query L1/L2 chains for transaction receipts
	log.Printf("🔄 [Bridge Monitor] Checking message %s (retry %d/%d, direction: %s)",
		msg.MessageHash, msg.RetryCount, m.maxRetries, msg.Direction)

	// TODO: Implement actual chain query logic
	// For now, we just increment retry count
	// In production, you would:
	// 1. Query L1 for L1→L2 messages using Arbitrum Inbox
	// 2. Query L2 for L2→L1 messages using Arbitrum Outbox
	// 3. Update status based on transaction receipts

	return nil
}

// isTimedOut checks if a message has exceeded the timeout period
func (m *Monitor) isTimedOut(msg *models.BridgeEvent) bool {
	initiatedTime := time.Unix(msg.InitiatedAt, 0)
	return time.Since(initiatedTime) > m.messageTimeout
}

// GetMessageStatus retrieves the current status of a bridge message
func (m *Monitor) GetMessageStatus(messageHash string) (*MessageStatus, error) {
	// Query database for message
	var status MessageStatus
	err := m.database.QueryRow(`
		SELECT message_hash, direction, status, l1_tx_hash, l2_tx_hash,
		       initiated_at, confirmed_at, retry_count, error_msg
		FROM bridge_messages
		WHERE message_hash = $1
	`, messageHash).Scan(
		&status.MessageHash,
		&status.Direction,
		&status.Status,
		&status.L1TxHash,
		&status.L2TxHash,
		&status.InitiatedAt,
		&status.ConfirmedAt,
		&status.RetryCount,
		&status.ErrorMsg,
	)

	if err != nil {
		return nil, err
	}

	return &status, nil
}

// MessageStatus represents the current state of a bridge message
type MessageStatus struct {
	MessageHash string
	Direction   string
	Status      string
	L1TxHash    string
	L2TxHash    sql.NullString
	InitiatedAt time.Time
	ConfirmedAt sql.NullTime
	RetryCount  int
	ErrorMsg    sql.NullString
}

// GetPendingMessagesCount returns the count of pending messages
func (m *Monitor) GetPendingMessagesCount() (int, error) {
	var count int
	err := m.database.QueryRow(`
		SELECT COUNT(*)
		FROM bridge_messages
		WHERE status IN ('initiated', 'pending')
		  AND retry_count < $1
	`, m.maxRetries).Scan(&count)
	return count, err
}

// GetFailedMessagesCount returns the count of failed messages
func (m *Monitor) GetFailedMessagesCount() (int, error) {
	var count int
	err := m.database.QueryRow(`
		SELECT COUNT(*)
		FROM bridge_messages
		WHERE status = 'failed'
	`).Scan(&count)
	return count, err
}

// GetStats returns bridge statistics
func (m *Monitor) GetStats() (*BridgeStats, error) {
	var stats BridgeStats

	// Get pending count
	pendingCount, err := m.GetPendingMessagesCount()
	if err != nil {
		return nil, err
	}
	stats.PendingCount = pendingCount

	// Get failed count
	failedCount, err := m.GetFailedMessagesCount()
	if err != nil {
		return nil, err
	}
	stats.FailedCount = failedCount

	// Get confirmed count
	err = m.database.QueryRow(`
		SELECT COUNT(*)
		FROM bridge_messages
		WHERE status = 'confirmed'
	`).Scan(&stats.ConfirmedCount)
	if err != nil {
		return nil, err
	}

	// Get average confirmation time
	err = m.database.QueryRow(`
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (confirmed_at - initiated_at))), 0)
		FROM bridge_messages
		WHERE status = 'confirmed'
		  AND confirmed_at IS NOT NULL
	`).Scan(&stats.AvgConfirmationTime)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// BridgeStats contains bridge operation statistics
type BridgeStats struct {
	PendingCount        int
	ConfirmedCount      int
	FailedCount         int
	AvgConfirmationTime float64 // in seconds
}

// RetryMessage manually retries a failed or stuck message
func (m *Monitor) RetryMessage(messageHash string) error {
	// Get current message
	status, err := m.GetMessageStatus(messageHash)
	if err != nil {
		return err
	}

	// Check if retry is allowed
	if status.Status == "confirmed" {
		return nil // Already confirmed, no retry needed
	}

	if status.RetryCount >= m.maxRetries {
		return db.MarkBridgeMessageFailed(m.database, messageHash, "Max retries exceeded")
	}

	// Reset status to pending for retry
	_, err = m.database.Exec(`
		UPDATE bridge_messages
		SET status = 'pending', retry_count = retry_count + 1, updated_at = NOW()
		WHERE message_hash = $1
	`, messageHash)

	if err == nil {
		log.Printf("🔄 [Bridge Monitor] Message %s queued for retry (attempt %d)",
			messageHash, status.RetryCount+1)
	}

	return err
}

// GetRecentMessages returns recent bridge messages for monitoring
func (m *Monitor) GetRecentMessages(limit int) ([]*models.BridgeEvent, error) {
	query := `
		SELECT message_hash, direction, user_address, amount, status,
		       l1_tx_hash, l2_tx_hash, l1_block_number, l2_block_number,
		       EXTRACT(EPOCH FROM initiated_at)::bigint,
		       COALESCE(EXTRACT(EPOCH FROM confirmed_at)::bigint, 0),
		       retry_count, COALESCE(error_msg, '')
		FROM bridge_messages
		ORDER BY initiated_at DESC
		LIMIT $1
	`

	rows, err := m.database.Query(query, limit)
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
