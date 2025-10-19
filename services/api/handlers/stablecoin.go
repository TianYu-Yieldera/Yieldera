package handlers

import (
	"database/sql"
	"math/big"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StablecoinPosition represents a user's stablecoin position
type StablecoinPosition struct {
	UserAddress      string  `json:"userAddress"`
	CollateralAmount string  `json:"collateralAmount"`
	DebtAmount       string  `json:"debtAmount"`
	CollateralRatio  float64 `json:"collateralRatio"`
	HealthStatus     string  `json:"healthStatus"`
	MaxMintable      string  `json:"maxMintable"`
	MaxRedeemable    string  `json:"maxRedeemable"`
	LastUpdated      string  `json:"lastUpdated"`
}

// MintSimulation represents a mint operation preview
type MintSimulation struct {
	Amount              string  `json:"amount"`
	Fee                 string  `json:"fee"`
	TotalCost           string  `json:"totalCost"`
	CollateralRequired  string  `json:"collateralRequired"`
	NewCollateralRatio  float64 `json:"newCollateralRatio"`
	NewHealthStatus     string  `json:"newHealthStatus"`
}

// RedeemSimulation represents a redeem operation preview
type RedeemSimulation struct {
	Amount              string  `json:"amount"`
	Fee                 string  `json:"fee"`
	CollateralReturned  string  `json:"collateralReturned"`
	NewCollateralRatio  float64 `json:"newCollateralRatio"`
	NewHealthStatus     string  `json:"newHealthStatus"`
}

const (
	COLLATERAL_RATIO       = 150.0
	LIQUIDATION_THRESHOLD  = 120.0
	MINT_FEE_PERCENT       = 0.2
	REDEEM_FEE_PERCENT     = 0.2
)

// GetStablecoinPosition returns user's stablecoin position
func GetStablecoinPosition(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		var pos StablecoinPosition
		err := db.QueryRow(`
			SELECT user_address, collateral_amount, debt_amount,
			       collateral_ratio, health_status, last_updated
			FROM stablecoin_positions
			WHERE user_address = $1
		`, address).Scan(
			&pos.UserAddress, &pos.CollateralAmount, &pos.DebtAmount,
			&pos.CollateralRatio, &pos.HealthStatus, &pos.LastUpdated,
		)

		if err == sql.ErrNoRows {
			// Return empty position
			c.JSON(200, StablecoinPosition{
				UserAddress:      address,
				CollateralAmount: "0",
				DebtAmount:       "0",
				CollateralRatio:  0,
				HealthStatus:     "safe",
				MaxMintable:      "0",
				MaxRedeemable:    "0",
				LastUpdated:      "",
			})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Calculate max mintable
		collateral, _ := new(big.Float).SetString(pos.CollateralAmount)
		debt, _ := new(big.Float).SetString(pos.DebtAmount)

		maxTotal := new(big.Float).Quo(
			new(big.Float).Mul(collateral, big.NewFloat(100)),
			big.NewFloat(COLLATERAL_RATIO),
		)
		maxMintable := new(big.Float).Sub(maxTotal, debt)
		if maxMintable.Cmp(big.NewFloat(0)) < 0 {
			maxMintable = big.NewFloat(0)
		}
		pos.MaxMintable = maxMintable.Text('f', 6)

		// Max redeemable is current debt
		pos.MaxRedeemable = pos.DebtAmount

		c.JSON(200, pos)
	}
}

// SimulateMint simulates a mint operation
func SimulateMint(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Get current position
		var collateral, debt string
		err := db.QueryRow(`
			SELECT collateral_amount, debt_amount
			FROM stablecoin_positions
			WHERE user_address = $1
		`, req.UserAddress).Scan(&collateral, &debt)

		if err == sql.ErrNoRows {
			collateral = "0"
			debt = "0"
		}

		amount, _ := new(big.Float).SetString(req.Amount)
		fee := new(big.Float).Quo(
			new(big.Float).Mul(amount, big.NewFloat(MINT_FEE_PERCENT)),
			big.NewFloat(100),
		)
		totalCost := new(big.Float).Add(amount, fee)

		collateralRequired := new(big.Float).Quo(
			new(big.Float).Mul(totalCost, big.NewFloat(COLLATERAL_RATIO)),
			big.NewFloat(100),
		)

		// Calculate new ratios
		currentCollateral, _ := new(big.Float).SetString(collateral)
		currentDebt, _ := new(big.Float).SetString(debt)
		newDebt := new(big.Float).Add(currentDebt, totalCost)

		var newRatio float64
		if newDebt.Cmp(big.NewFloat(0)) > 0 {
			ratio := new(big.Float).Quo(
				new(big.Float).Mul(currentCollateral, big.NewFloat(100)),
				newDebt,
			)
			newRatio, _ = ratio.Float64()
		} else {
			newRatio = 0
		}

		healthStatus := getHealthStatus(newRatio)

		sim := MintSimulation{
			Amount:             req.Amount,
			Fee:                fee.Text('f', 6),
			TotalCost:          totalCost.Text('f', 6),
			CollateralRequired: collateralRequired.Text('f', 6),
			NewCollateralRatio: newRatio,
			NewHealthStatus:    healthStatus,
		}

		c.JSON(200, sim)
	}
}

// SimulateRedeem simulates a redeem operation
func SimulateRedeem(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Get current position
		var collateral, debt string
		err := db.QueryRow(`
			SELECT collateral_amount, debt_amount
			FROM stablecoin_positions
			WHERE user_address = $1
		`, req.UserAddress).Scan(&collateral, &debt)

		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "No position found"})
			return
		}

		amount, _ := new(big.Float).SetString(req.Amount)
		fee := new(big.Float).Quo(
			new(big.Float).Mul(amount, big.NewFloat(REDEEM_FEE_PERCENT)),
			big.NewFloat(100),
		)

		collateralReturned := new(big.Float).Quo(
			new(big.Float).Mul(amount, big.NewFloat(COLLATERAL_RATIO)),
			big.NewFloat(100),
		)

		// Calculate new ratios
		currentCollateral, _ := new(big.Float).SetString(collateral)
		currentDebt, _ := new(big.Float).SetString(debt)
		newCollateral := new(big.Float).Sub(currentCollateral, collateralReturned)
		newDebt := new(big.Float).Sub(currentDebt, amount)

		var newRatio float64
		if newDebt.Cmp(big.NewFloat(0)) > 0 {
			ratio := new(big.Float).Quo(
				new(big.Float).Mul(newCollateral, big.NewFloat(100)),
				newDebt,
			)
			newRatio, _ = ratio.Float64()
		} else {
			newRatio = 0
		}

		healthStatus := getHealthStatus(newRatio)

		sim := RedeemSimulation{
			Amount:             req.Amount,
			Fee:                fee.Text('f', 6),
			CollateralReturned: collateralReturned.Text('f', 6),
			NewCollateralRatio: newRatio,
			NewHealthStatus:    healthStatus,
		}

		c.JSON(200, sim)
	}
}

// MintLUSD executes a mint operation
func MintLUSD(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Get or create position
		var collateral, debt string
		err := db.QueryRow(`
			SELECT collateral_amount, debt_amount
			FROM stablecoin_positions
			WHERE user_address = $1
		`, req.UserAddress).Scan(&collateral, &debt)

		if err == sql.ErrNoRows {
			collateral = "0"
			debt = "0"
		}

		// Calculate amounts
		amount, _ := new(big.Float).SetString(req.Amount)
		fee := new(big.Float).Quo(
			new(big.Float).Mul(amount, big.NewFloat(MINT_FEE_PERCENT)),
			big.NewFloat(100),
		)
		totalCost := new(big.Float).Add(amount, fee)

		collateralRequired := new(big.Float).Quo(
			new(big.Float).Mul(totalCost, big.NewFloat(COLLATERAL_RATIO)),
			big.NewFloat(100),
		)

		// Check if user has enough points for collateral
		var userPoints string
		err = db.QueryRow(`SELECT points FROM points WHERE user_address = $1`, req.UserAddress).Scan(&userPoints)
		if err == sql.ErrNoRows {
			userPoints = "0"
		}

		points, _ := strconv.ParseFloat(userPoints, 64)
		requiredPoints, _ := collateralRequired.Float64()

		if points < requiredPoints {
			c.JSON(400, gin.H{
				"error":    "Insufficient points for collateral",
				"required": requiredPoints,
				"current":  points,
			})
			return
		}

		// Update position
		currentCollateral, _ := new(big.Float).SetString(collateral)
		currentDebt, _ := new(big.Float).SetString(debt)

		newCollateral := new(big.Float).Add(currentCollateral, collateralRequired)
		newDebt := new(big.Float).Add(currentDebt, totalCost)

		var newRatio float64
		if newDebt.Cmp(big.NewFloat(0)) > 0 {
			ratio := new(big.Float).Quo(
				new(big.Float).Mul(newCollateral, big.NewFloat(100)),
				newDebt,
			)
			newRatio, _ = ratio.Float64()
		}

		healthStatus := getHealthStatus(newRatio)

		_, err = db.Exec(`
			INSERT INTO stablecoin_positions (user_address, collateral_amount, debt_amount, collateral_ratio, health_status, last_updated)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (user_address)
			DO UPDATE SET
				collateral_amount = $2,
				debt_amount = $3,
				collateral_ratio = $4,
				health_status = $5,
				last_updated = NOW()
		`, req.UserAddress, newCollateral.Text('f', 18), newDebt.Text('f', 18), newRatio, healthStatus)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Record transaction
		db.Exec(`
			INSERT INTO stablecoin_transactions (user_address, tx_type, amount, fee, collateral_ratio_after, status, timestamp)
			VALUES ($1, 'mint', $2, $3, $4, 'confirmed', NOW())
		`, req.UserAddress, req.Amount, fee.Text('f', 6), newRatio)

		c.JSON(200, gin.H{
			"success":          true,
			"message":          "LUSD minted successfully",
			"amount":           req.Amount,
			"fee":              fee.Text('f', 6),
			"collateralUsed":   collateralRequired.Text('f', 6),
			"newCollateralRatio": newRatio,
			"healthStatus":     healthStatus,
		})
	}
}

// RedeemLUSD executes a redeem operation
func RedeemLUSD(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"userAddress" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Get current position
		var collateral, debt string
		err := db.QueryRow(`
			SELECT collateral_amount, debt_amount
			FROM stablecoin_positions
			WHERE user_address = $1
		`, req.UserAddress).Scan(&collateral, &debt)

		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "No position found"})
			return
		}

		// Check if user has enough debt
		amount, _ := new(big.Float).SetString(req.Amount)
		currentDebt, _ := new(big.Float).SetString(debt)

		if currentDebt.Cmp(amount) < 0 {
			c.JSON(400, gin.H{"error": "Insufficient debt to redeem"})
			return
		}

		// Calculate amounts
		fee := new(big.Float).Quo(
			new(big.Float).Mul(amount, big.NewFloat(REDEEM_FEE_PERCENT)),
			big.NewFloat(100),
		)

		collateralReturned := new(big.Float).Quo(
			new(big.Float).Mul(amount, big.NewFloat(COLLATERAL_RATIO)),
			big.NewFloat(100),
		)

		// Update position
		currentCollateral, _ := new(big.Float).SetString(collateral)
		newCollateral := new(big.Float).Sub(currentCollateral, collateralReturned)
		newDebt := new(big.Float).Sub(currentDebt, amount)

		var newRatio float64
		if newDebt.Cmp(big.NewFloat(0)) > 0 {
			ratio := new(big.Float).Quo(
				new(big.Float).Mul(newCollateral, big.NewFloat(100)),
				newDebt,
			)
			newRatio, _ = ratio.Float64()
		}

		healthStatus := getHealthStatus(newRatio)

		_, err = db.Exec(`
			UPDATE stablecoin_positions
			SET collateral_amount = $1, debt_amount = $2, collateral_ratio = $3, health_status = $4, last_updated = NOW()
			WHERE user_address = $5
		`, newCollateral.Text('f', 18), newDebt.Text('f', 18), newRatio, healthStatus, req.UserAddress)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Record transaction
		db.Exec(`
			INSERT INTO stablecoin_transactions (user_address, tx_type, amount, fee, collateral_ratio_after, status, timestamp)
			VALUES ($1, 'redeem', $2, $3, $4, 'confirmed', NOW())
		`, req.UserAddress, req.Amount, fee.Text('f', 6), newRatio)

		c.JSON(200, gin.H{
			"success":           true,
			"message":           "LUSD redeemed successfully",
			"amount":            req.Amount,
			"fee":               fee.Text('f', 6),
			"collateralReturned": collateralReturned.Text('f', 6),
			"newCollateralRatio": newRatio,
			"healthStatus":      healthStatus,
		})
	}
}

// GetStablecoinHistory returns transaction history
func GetStablecoinHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		limit := c.DefaultQuery("limit", "20")

		rows, err := db.Query(`
			SELECT id, tx_type, amount, fee, collateral_ratio_before, collateral_ratio_after, tx_hash, status, timestamp
			FROM stablecoin_transactions
			WHERE user_address = $1
			ORDER BY timestamp DESC
			LIMIT $2
		`, address, limit)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Transaction struct {
			ID                      int     `json:"id"`
			TxType                  string  `json:"type"`
			Amount                  string  `json:"amount"`
			Fee                     string  `json:"fee"`
			CollateralRatioBefore   float64 `json:"collateralRatioBefore,omitempty"`
			CollateralRatioAfter    float64 `json:"collateralRatioAfter"`
			TxHash                  string  `json:"txHash"`
			Status                  string  `json:"status"`
			Timestamp               string  `json:"timestamp"`
		}

		var history []Transaction
		for rows.Next() {
			var tx Transaction
			var fee, txHash sql.NullString
			var ratioBefore sql.NullFloat64
			rows.Scan(&tx.ID, &tx.TxType, &tx.Amount, &fee, &ratioBefore, &tx.CollateralRatioAfter, &txHash, &tx.Status, &tx.Timestamp)
			if fee.Valid {
				tx.Fee = fee.String
			}
			if txHash.Valid {
				tx.TxHash = txHash.String
			}
			if ratioBefore.Valid {
				tx.CollateralRatioBefore = ratioBefore.Float64
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

// GetStablecoinStats returns global stablecoin statistics
func GetStablecoinStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get total collateral
		var totalCollateral float64
		db.QueryRow(`
			SELECT COALESCE(SUM(CAST(collateral_amount AS NUMERIC)), 0)
			FROM stablecoin_positions
		`).Scan(&totalCollateral)

		// Get total debt (LUSD supply)
		var totalDebt float64
		db.QueryRow(`
			SELECT COALESCE(SUM(CAST(debt_amount AS NUMERIC)), 0)
			FROM stablecoin_positions
		`).Scan(&totalDebt)

		// Get total users
		var totalUsers int
		db.QueryRow(`SELECT COUNT(*) FROM stablecoin_positions WHERE debt_amount != '0'`).Scan(&totalUsers)

		// Get average collateral ratio
		var avgRatio float64
		db.QueryRow(`
			SELECT COALESCE(AVG(collateral_ratio), 0)
			FROM stablecoin_positions
			WHERE debt_amount != '0'
		`).Scan(&avgRatio)

		// Get health status distribution
		var safeCount, warningCount, dangerCount int
		db.QueryRow(`SELECT COUNT(*) FROM stablecoin_positions WHERE health_status = 'safe' AND debt_amount != '0'`).Scan(&safeCount)
		db.QueryRow(`SELECT COUNT(*) FROM stablecoin_positions WHERE health_status = 'warning' AND debt_amount != '0'`).Scan(&warningCount)
		db.QueryRow(`SELECT COUNT(*) FROM stablecoin_positions WHERE health_status = 'danger' AND debt_amount != '0'`).Scan(&dangerCount)

		c.JSON(200, gin.H{
			"totalCollateral":    totalCollateral,
			"totalSupply":        totalDebt,
			"totalUsers":         totalUsers,
			"avgCollateralRatio": avgRatio,
			"healthDistribution": gin.H{
				"safe":    safeCount,
				"warning": warningCount,
				"danger":  dangerCount,
			},
			"lusdPrice": 1.00,
		})
	}
}

// Helper function to determine health status
func getHealthStatus(ratio float64) string {
	if ratio >= COLLATERAL_RATIO {
		return "safe"
	} else if ratio >= LIQUIDATION_THRESHOLD {
		return "warning"
	}
	return "danger"
}
