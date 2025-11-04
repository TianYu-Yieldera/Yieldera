package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetL1Balance returns L1 collateral balance for a user
func GetL1Balance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		// Query all L1 collateral balances
		query := `
			SELECT token, amount, usd_value, updated_at
			FROM l1_collateral_balances
			WHERE user_address = $1
			ORDER BY token
		`

		rows, err := db.Query(query, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type TokenBalance struct {
			Token     string `json:"token"`
			Amount    string `json:"amount"`
			USDValue  string `json:"usd_value"`
			UpdatedAt string `json:"updated_at"`
		}

		var balances []TokenBalance
		for rows.Next() {
			var b TokenBalance
			err := rows.Scan(&b.Token, &b.Amount, &b.USDValue, &b.UpdatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			balances = append(balances, b)
		}

		if len(balances) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"address":  address,
				"balances": []TokenBalance{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"address":  address,
			"balances": balances,
		})
	}
}

// GetL1Deposits returns L1 deposit history for a user
func GetL1Deposits(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		limit := c.DefaultQuery("limit", "50")
		offset := c.DefaultQuery("offset", "0")

		query := `
			SELECT user_address, token, amount, tx_hash, block_number, confirmed, created_at
			FROM l1_collateral_deposits
			WHERE user_address = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`

		rows, err := db.Query(query, address, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Deposit struct {
			UserAddress string `json:"user_address"`
			Token       string `json:"token"`
			Amount      string `json:"amount"`
			TxHash      string `json:"tx_hash"`
			BlockNumber int64  `json:"block_number"`
			Confirmed   bool   `json:"confirmed"`
			CreatedAt   string `json:"created_at"`
		}

		var deposits []Deposit
		for rows.Next() {
			var d Deposit
			err := rows.Scan(&d.UserAddress, &d.Token, &d.Amount, &d.TxHash, &d.BlockNumber, &d.Confirmed, &d.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			deposits = append(deposits, d)
		}

		// Get total count
		var total int
		db.QueryRow(`SELECT COUNT(*) FROM l1_collateral_deposits WHERE user_address = $1`, address).Scan(&total)

		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"deposits": deposits,
			"pagination": gin.H{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}

// InitiateL1Deposit initiates a deposit to L1 collateral vault
func InitiateL1Deposit(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"user_address" binding:"required"`
			Token       string `json:"token" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
			Signature   string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Verify signature
		// TODO: Submit transaction to L1 CollateralVault contract
		// For now, return a placeholder response

		c.JSON(http.StatusAccepted, gin.H{
			"status": "pending",
			"message": "Deposit transaction submitted to L1",
			"user_address": req.UserAddress,
			"token": req.Token,
			"amount": req.Amount,
			"note": "This endpoint requires L1 contract integration",
		})
	}
}

// InitiateL1Withdrawal initiates a withdrawal from L1 collateral vault
func InitiateL1Withdrawal(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"user_address" binding:"required"`
			Token       string `json:"token" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
			Signature   string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user has sufficient balance
		var currentBalance string
		err := db.QueryRow(`
			SELECT amount FROM l1_collateral_balances
			WHERE user_address = $1 AND token = $2
		`, req.UserAddress, req.Token).Scan(&currentBalance)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: Verify signature
		// TODO: Submit transaction to L1 CollateralVault contract
		// For now, return a placeholder response

		c.JSON(http.StatusAccepted, gin.H{
			"status": "pending",
			"message": "Withdrawal transaction submitted to L1",
			"user_address": req.UserAddress,
			"token": req.Token,
			"amount": req.Amount,
			"current_balance": currentBalance,
			"note": "This endpoint requires L1 contract integration",
		})
	}
}

// GetL1StateSnapshots returns L1 state snapshots
func GetL1StateSnapshots(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "20")
		offset := c.DefaultQuery("offset", "0")

		query := `
			SELECT l2_block_number, state_root, tx_hash, block_number, created_at
			FROM l1_state_snapshots
			ORDER BY l2_block_number DESC
			LIMIT $1 OFFSET $2
		`

		rows, err := db.Query(query, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type StateSnapshot struct {
			L2BlockNumber int64  `json:"l2_block_number"`
			StateRoot     string `json:"state_root"`
			TxHash        string `json:"tx_hash"`
			BlockNumber   int64  `json:"block_number"`
			CreatedAt     string `json:"created_at"`
		}

		var snapshots []StateSnapshot
		for rows.Next() {
			var s StateSnapshot
			err := rows.Scan(&s.L2BlockNumber, &s.StateRoot, &s.TxHash, &s.BlockNumber, &s.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			snapshots = append(snapshots, s)
		}

		// Get total count
		var total int
		db.QueryRow(`SELECT COUNT(*) FROM l1_state_snapshots`).Scan(&total)

		c.JSON(http.StatusOK, gin.H{
			"snapshots": snapshots,
			"pagination": gin.H{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}
