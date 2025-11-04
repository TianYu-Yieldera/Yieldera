package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"loyalty-points-system/internal/bridge"
)

// GetBridgeStatus returns the status of a bridge message
func GetBridgeStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageHash := c.Param("messageHash")

		query := `
			SELECT message_hash, direction, user_address, amount, status,
			       l1_tx_hash, l2_tx_hash, l1_block_number, l2_block_number,
			       initiated_at, confirmed_at, retry_count, error_msg
			FROM bridge_messages
			WHERE message_hash = $1
		`

		var msg struct {
			MessageHash   string  `json:"message_hash"`
			Direction     string  `json:"direction"`
			UserAddress   string  `json:"user_address"`
			Amount        string  `json:"amount"`
			Status        string  `json:"status"`
			L1TxHash      string  `json:"l1_tx_hash"`
			L2TxHash      *string `json:"l2_tx_hash,omitempty"`
			L1BlockNumber int64   `json:"l1_block_number"`
			L2BlockNumber *int64  `json:"l2_block_number,omitempty"`
			InitiatedAt   string  `json:"initiated_at"`
			ConfirmedAt   *string `json:"confirmed_at,omitempty"`
			RetryCount    int     `json:"retry_count"`
			ErrorMsg      *string `json:"error_msg,omitempty"`
		}

		var l2TxHash, confirmedAt, errorMsg sql.NullString
		var l2BlockNumber sql.NullInt64

		err := db.QueryRow(query, messageHash).Scan(
			&msg.MessageHash,
			&msg.Direction,
			&msg.UserAddress,
			&msg.Amount,
			&msg.Status,
			&msg.L1TxHash,
			&l2TxHash,
			&msg.L1BlockNumber,
			&l2BlockNumber,
			&msg.InitiatedAt,
			&confirmedAt,
			&msg.RetryCount,
			&errorMsg,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bridge message not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Handle nullable fields
		if l2TxHash.Valid {
			msg.L2TxHash = &l2TxHash.String
		}
		if l2BlockNumber.Valid {
			msg.L2BlockNumber = &l2BlockNumber.Int64
		}
		if confirmedAt.Valid {
			msg.ConfirmedAt = &confirmedAt.String
		}
		if errorMsg.Valid {
			msg.ErrorMsg = &errorMsg.String
		}

		c.JSON(http.StatusOK, msg)
	}
}

// GetUserBridgeHistory returns bridge history for a user
func GetUserBridgeHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		limit := c.DefaultQuery("limit", "50")
		offset := c.DefaultQuery("offset", "0")

		query := `
			SELECT message_hash, direction, user_address, amount, status,
			       l1_tx_hash, l2_tx_hash, l1_block_number, l2_block_number,
			       initiated_at, confirmed_at, retry_count, error_msg
			FROM bridge_messages
			WHERE user_address = $1
			ORDER BY initiated_at DESC
			LIMIT $2 OFFSET $3
		`

		rows, err := db.Query(query, address, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type BridgeMessage struct {
			MessageHash   string  `json:"message_hash"`
			Direction     string  `json:"direction"`
			UserAddress   string  `json:"user_address"`
			Amount        string  `json:"amount"`
			Status        string  `json:"status"`
			L1TxHash      string  `json:"l1_tx_hash"`
			L2TxHash      *string `json:"l2_tx_hash,omitempty"`
			L1BlockNumber int64   `json:"l1_block_number"`
			L2BlockNumber *int64  `json:"l2_block_number,omitempty"`
			InitiatedAt   string  `json:"initiated_at"`
			ConfirmedAt   *string `json:"confirmed_at,omitempty"`
			RetryCount    int     `json:"retry_count"`
			ErrorMsg      *string `json:"error_msg,omitempty"`
		}

		var messages []BridgeMessage
		for rows.Next() {
			var msg BridgeMessage
			var l2TxHash, confirmedAt, errorMsg sql.NullString
			var l2BlockNumber sql.NullInt64

			err := rows.Scan(
				&msg.MessageHash,
				&msg.Direction,
				&msg.UserAddress,
				&msg.Amount,
				&msg.Status,
				&msg.L1TxHash,
				&l2TxHash,
				&msg.L1BlockNumber,
				&l2BlockNumber,
				&msg.InitiatedAt,
				&confirmedAt,
				&msg.RetryCount,
				&errorMsg,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Handle nullable fields
			if l2TxHash.Valid {
				msg.L2TxHash = &l2TxHash.String
			}
			if l2BlockNumber.Valid {
				msg.L2BlockNumber = &l2BlockNumber.Int64
			}
			if confirmedAt.Valid {
				msg.ConfirmedAt = &confirmedAt.String
			}
			if errorMsg.Valid {
				msg.ErrorMsg = &errorMsg.String
			}

			messages = append(messages, msg)
		}

		// Get total count
		var total int
		db.QueryRow(`SELECT COUNT(*) FROM bridge_messages WHERE user_address = $1`, address).Scan(&total)

		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"messages": messages,
			"pagination": gin.H{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}

// InitiateBridgeL1ToL2 initiates a bridge from L1 to L2
func InitiateBridgeL1ToL2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"user_address" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
			Token       string `json:"token" binding:"required"`
			Signature   string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Verify signature
		// TODO: Submit transaction to L1 Gateway contract
		// For now, return a placeholder response

		c.JSON(http.StatusAccepted, gin.H{
			"status": "pending",
			"message": "Bridge transaction L1→L2 submitted",
			"user_address": req.UserAddress,
			"amount": req.Amount,
			"token": req.Token,
			"direction": "L1_TO_L2",
			"note": "This endpoint requires L1 Gateway contract integration",
		})
	}
}

// InitiateBridgeL2ToL1 initiates a bridge from L2 to L1
func InitiateBridgeL2ToL1(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"user_address" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
			Signature   string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Verify signature
		// TODO: Submit transaction to L2 IntegratedVault contract
		// For now, return a placeholder response

		c.JSON(http.StatusAccepted, gin.H{
			"status": "pending",
			"message": "Bridge transaction L2→L1 submitted",
			"user_address": req.UserAddress,
			"amount": req.Amount,
			"direction": "L2_TO_L1",
			"note": "This endpoint requires L2 contract integration",
		})
	}
}

// RetryBridgeMessage retries a failed bridge message
func RetryBridgeMessage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageHash := c.Param("messageHash")

		// Create monitor instance
		monitor := bridge.NewMonitor(db)

		// Retry the message
		err := monitor.RetryMessage(messageHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "retry_queued",
			"message": "Bridge message queued for retry",
			"message_hash": messageHash,
		})
	}
}

// GetBridgeStats returns bridge statistics
func GetBridgeStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		monitor := bridge.NewMonitor(db)
		stats, err := monitor.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pending_count": stats.PendingCount,
			"confirmed_count": stats.ConfirmedCount,
			"failed_count": stats.FailedCount,
			"avg_confirmation_time_seconds": stats.AvgConfirmationTime,
		})
	}
}
