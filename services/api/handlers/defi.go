package handlers

import (
	"database/sql"
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeFiPool represents a DeFi protocol pool
type DeFiPool struct {
	ID               string   `json:"id"`
	Protocol         string   `json:"protocol"`
	Name             string   `json:"name"`
	PoolType         string   `json:"type"`
	APR              string   `json:"apr"`
	TVL              string   `json:"tvl"`
	RiskLevel        string   `json:"risk"`
	Icon             string   `json:"icon"`
	Color            string   `json:"color"`
	PointsRequired   int      `json:"pointsRequired"`
	PointsMultiplier float64  `json:"pointsMultiplier"`
	Features         []string `json:"features"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// UserPosition represents a user's position in a pool
type UserPosition struct {
	PoolID       string `json:"poolId"`
	PoolName     string `json:"poolName"`
	Protocol     string `json:"protocol"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Deposited    string `json:"deposited"`
	Earned       string `json:"earned"`
	PointsEarned string `json:"pointsEarned"`
	LastUpdated  string `json:"lastUpdated"`
}

// GetDeFiPools returns all DeFi protocol pools
func GetDeFiPools(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, protocol, name, pool_type, apr, tvl, risk_level,
			       icon, color, points_required, points_multiplier, features, metadata
			FROM defi_pools
			ORDER BY
				CASE
					WHEN id = 'uniswap' THEN 1
					WHEN id = 'aave' THEN 2
					WHEN id = 'stablecoin' THEN 3
					WHEN id = 'staking' THEN 4
					ELSE 5
				END
		`)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var pools []DeFiPool
		for rows.Next() {
			var pool DeFiPool
			var featuresJSON, metadataJSON []byte
			err := rows.Scan(
				&pool.ID, &pool.Protocol, &pool.Name, &pool.PoolType,
				&pool.APR, &pool.TVL, &pool.RiskLevel,
				&pool.Icon, &pool.Color, &pool.PointsRequired, &pool.PointsMultiplier,
				&featuresJSON, &metadataJSON,
			)
			if err != nil {
				continue
			}
			if err := json.Unmarshal(featuresJSON, &pool.Features); err != nil {
				pool.Features = []string{}
			}
			if err := json.Unmarshal(metadataJSON, &pool.Metadata); err != nil {
				pool.Metadata = make(map[string]interface{})
			}
			pools = append(pools, pool)
		}

		c.JSON(200, gin.H{"pools": pools})
	}
}

// GetPoolDetail returns details of a specific pool
func GetPoolDetail(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		poolID := c.Param("id")

		var pool DeFiPool
		var featuresJSON, metadataJSON []byte
		err := db.QueryRow(`
			SELECT id, protocol, name, pool_type, apr, tvl, risk_level,
			       icon, color, points_required, points_multiplier, features, metadata
			FROM defi_pools
			WHERE id = $1
		`, poolID).Scan(
			&pool.ID, &pool.Protocol, &pool.Name, &pool.PoolType,
			&pool.APR, &pool.TVL, &pool.RiskLevel,
			&pool.Icon, &pool.Color, &pool.PointsRequired, &pool.PointsMultiplier,
			&featuresJSON, &metadataJSON,
		)

		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Pool not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err := json.Unmarshal(featuresJSON, &pool.Features); err != nil {
			pool.Features = []string{}
		}
		if err := json.Unmarshal(metadataJSON, &pool.Metadata); err != nil {
			pool.Metadata = make(map[string]interface{})
		}

		c.JSON(200, pool)
	}
}

// GetUserDeFiPositions returns user's positions in all pools
func GetUserDeFiPositions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		rows, err := db.Query(`
			SELECT p.pool_id, p.deposited, p.earned, p.points_earned, p.last_updated,
			       d.name, d.protocol, d.icon, d.color
			FROM user_defi_positions p
			JOIN defi_pools d ON p.pool_id = d.id
			WHERE p.user_address = $1
		`, address)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var positions []UserPosition
		for rows.Next() {
			var pos UserPosition
			rows.Scan(
				&pos.PoolID, &pos.Deposited, &pos.Earned, &pos.PointsEarned,
				&pos.LastUpdated, &pos.PoolName, &pos.Protocol, &pos.Icon, &pos.Color,
			)
			positions = append(positions, pos)
		}

		// If no positions, return empty array
		if positions == nil {
			positions = []UserPosition{}
		}

		c.JSON(200, gin.H{
			"address":   address,
			"positions": positions,
		})
	}
}

// DepositToPool handles deposit to a DeFi pool
func DepositToPool(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			PoolID      string `json:"poolId" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Verify user has enough points
		var userPoints string
		err := db.QueryRow(`SELECT points FROM points WHERE user_address = $1`, req.UserAddress).Scan(&userPoints)
		if err == sql.ErrNoRows {
			userPoints = "0"
		}

		// Get pool requirements
		var pointsRequired int
		err = db.QueryRow(`SELECT points_required FROM defi_pools WHERE id = $1`, req.PoolID).Scan(&pointsRequired)
		if err != nil {
			c.JSON(404, gin.H{"error": "Pool not found"})
			return
		}

		// Check if user has enough points
		points, _ := strconv.ParseFloat(userPoints, 64)
		if int(points) < pointsRequired {
			c.JSON(400, gin.H{
				"error":    "Insufficient points",
				"required": pointsRequired,
				"current":  int(points),
			})
			return
		}

		// Update or insert position
		_, err = db.Exec(`
			INSERT INTO user_defi_positions (user_address, pool_id, deposited, last_updated)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (user_address, pool_id)
			DO UPDATE SET
				deposited = (CAST(user_defi_positions.deposited AS NUMERIC) + CAST($3 AS NUMERIC))::TEXT,
				last_updated = NOW()
		`, req.UserAddress, req.PoolID, req.Amount)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Record transaction
		db.Exec(`
			INSERT INTO defi_transactions (user_address, pool_id, tx_type, amount, status, timestamp)
			VALUES ($1, $2, 'deposit', $3, 'confirmed', NOW())
		`, req.UserAddress, req.PoolID, req.Amount)

		c.JSON(200, gin.H{
			"success": true,
			"message": "Deposit successful",
			"amount":  req.Amount,
			"poolId":  req.PoolID,
		})
	}
}

// WithdrawFromPool handles withdrawal from a DeFi pool
func WithdrawFromPool(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			PoolID      string `json:"poolId" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Get current deposited amount
		var deposited string
		err := db.QueryRow(`
			SELECT deposited FROM user_defi_positions
			WHERE user_address = $1 AND pool_id = $2
		`, req.UserAddress, req.PoolID).Scan(&deposited)

		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "No position found"})
			return
		}

		// Check if user has enough deposited
		depositedVal, _ := new(big.Float).SetString(deposited)
		withdrawVal, _ := new(big.Float).SetString(req.Amount)

		if depositedVal.Cmp(withdrawVal) < 0 {
			c.JSON(400, gin.H{"error": "Insufficient deposited amount"})
			return
		}

		// Update position
		newDeposited := new(big.Float).Sub(depositedVal, withdrawVal)
		_, err = db.Exec(`
			UPDATE user_defi_positions
			SET deposited = $1, last_updated = NOW()
			WHERE user_address = $2 AND pool_id = $3
		`, newDeposited.Text('f', 18), req.UserAddress, req.PoolID)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Record transaction
		db.Exec(`
			INSERT INTO defi_transactions (user_address, pool_id, tx_type, amount, status, timestamp)
			VALUES ($1, $2, 'withdraw', $3, 'confirmed', NOW())
		`, req.UserAddress, req.PoolID, req.Amount)

		c.JSON(200, gin.H{
			"success": true,
			"message": "Withdrawal successful",
			"amount":  req.Amount,
			"poolId":  req.PoolID,
		})
	}
}

// ClaimRewards handles claiming rewards from a pool
func ClaimRewards(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			PoolID      string `json:"poolId" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Get current earned amount
		var earned string
		err := db.QueryRow(`
			SELECT earned FROM user_defi_positions
			WHERE user_address = $1 AND pool_id = $2
		`, req.UserAddress, req.PoolID).Scan(&earned)

		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "No position found"})
			return
		}

		// Reset earned to 0
		_, err = db.Exec(`
			UPDATE user_defi_positions
			SET earned = '0', last_updated = NOW()
			WHERE user_address = $1 AND pool_id = $2
		`, req.UserAddress, req.PoolID)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Record transaction
		db.Exec(`
			INSERT INTO defi_transactions (user_address, pool_id, tx_type, amount, status, timestamp)
			VALUES ($1, $2, 'claim', $3, 'confirmed', NOW())
		`, req.UserAddress, req.PoolID, earned)

		c.JSON(200, gin.H{
			"success": true,
			"message": "Rewards claimed successfully",
			"amount":  earned,
			"poolId":  req.PoolID,
		})
	}
}

// GetDeFiHistory returns transaction history for a user
func GetDeFiHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		limit := c.DefaultQuery("limit", "20")

		rows, err := db.Query(`
			SELECT t.id, t.pool_id, t.tx_type, t.amount, t.tx_hash, t.status, t.timestamp, d.name, d.icon
			FROM defi_transactions t
			JOIN defi_pools d ON t.pool_id = d.id
			WHERE t.user_address = $1
			ORDER BY t.timestamp DESC
			LIMIT $2
		`, address, limit)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Transaction struct {
			ID        int    `json:"id"`
			PoolID    string `json:"poolId"`
			PoolName  string `json:"poolName"`
			Icon      string `json:"icon"`
			TxType    string `json:"type"`
			Amount    string `json:"amount"`
			TxHash    string `json:"txHash"`
			Status    string `json:"status"`
			Timestamp string `json:"timestamp"`
		}

		var history []Transaction
		for rows.Next() {
			var tx Transaction
			var txHash sql.NullString
			rows.Scan(&tx.ID, &tx.PoolID, &tx.TxType, &tx.Amount, &txHash, &tx.Status, &tx.Timestamp, &tx.PoolName, &tx.Icon)
			if txHash.Valid {
				tx.TxHash = txHash.String
			}
			history = append(history, tx)
		}

		if history == nil {
			history = []Transaction{}
		}

		c.JSON(200, gin.H{
			"address": address,
			"history": history,
		})
	}
}

// GetDeFiStats returns global DeFi statistics
func GetDeFiStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get total TVL
		var totalTVL float64
		db.QueryRow(`
			SELECT SUM(CAST(REPLACE(REPLACE(tvl, '$', ''), ',', '') AS NUMERIC))
			FROM defi_pools
		`).Scan(&totalTVL)

		// Get total users
		var totalUsers int
		db.QueryRow(`SELECT COUNT(DISTINCT user_address) FROM user_defi_positions`).Scan(&totalUsers)

		// Get total deposited
		var totalDeposited float64
		db.QueryRow(`
			SELECT SUM(CAST(deposited AS NUMERIC))
			FROM user_defi_positions
		`).Scan(&totalDeposited)

		c.JSON(200, gin.H{
			"totalTVL":       totalTVL,
			"totalUsers":     totalUsers,
			"totalDeposited": totalDeposited,
			"poolCount":      4,
		})
	}
}
