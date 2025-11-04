package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"loyalty-points-system/internal/db"

	"github.com/gin-gonic/gin"
)

// GetTreasuryAssets returns all treasury assets
// GET /api/v1/treasury/assets?type=T-BILL
func GetTreasuryAssets(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		treasuryType := c.Query("type")

		var typePtr *string
		if treasuryType != "" {
			typePtr = &treasuryType
		}

		assets, err := db.GetAllTreasuryAssets(database, typePtr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch treasury assets",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count":  len(assets),
			"assets": assets,
		})
	}
}

// GetTreasuryAssetDetail returns detailed info for a single asset
// GET /api/v1/treasury/assets/:assetId
func GetTreasuryAssetDetail(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetIdStr := c.Param("assetId")
		assetId, err := strconv.ParseInt(assetIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid asset ID",
			})
			return
		}

		asset, err := db.GetTreasuryAsset(database, assetId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Asset not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch asset",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"asset": asset,
		})
	}
}

// GetTreasuryPriceHistory returns price history for an asset
// GET /api/v1/treasury/assets/:assetId/price-history?limit=100
func GetTreasuryPriceHistory(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetIdStr := c.Param("assetId")
		assetId, err := strconv.ParseInt(assetIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid asset ID",
			})
			return
		}

		limitStr := c.DefaultQuery("limit", "100")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			limit = 100
		}

		history, err := db.GetTreasuryPriceHistory(database, assetId, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch price history",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"asset_id": assetId,
			"count":    len(history),
			"history":  history,
		})
	}
}

// GetUserTreasuryHoldings returns user's treasury holdings
// GET /api/v1/treasury/user/:address/holdings
func GetUserTreasuryHoldings(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address required",
			})
			return
		}

		holdings, err := db.GetUserTreasuryHoldings(database, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch holdings",
				"details": err.Error(),
			})
			return
		}

		// Calculate total portfolio value
		totalValue := "0"
		totalInvested := "0"
		totalGainLoss := "0"

		// In production, you'd sum these properly
		// For now, just return the holdings

		c.JSON(http.StatusOK, gin.H{
			"user_address":    address,
			"holdings_count":  len(holdings),
			"holdings":        holdings,
			"total_value":     totalValue,
			"total_invested":  totalInvested,
			"total_gain_loss": totalGainLoss,
		})
	}
}

// GetUserTreasuryYield returns user's yield information
// GET /api/v1/treasury/user/:address/yield
func GetUserTreasuryYield(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address required",
			})
			return
		}

		// Query yield distributions for user's holdings
		query := `
			SELECT
				d.id,
				d.asset_id,
				a.cusip,
				a.treasury_type,
				d.distribution_date,
				d.distribution_type,
				d.total_yield,
				d.yield_per_token,
				h.tokens_held,
				(h.tokens_held::numeric * d.yield_per_token::numeric / 1e18) AS user_yield,
				d.status
			FROM treasury_yield_distributions d
			JOIN treasury_assets a ON d.asset_id = a.asset_id
			JOIN treasury_holdings h ON d.asset_id = h.asset_id
			WHERE h.user_address = $1
			ORDER BY d.distribution_date DESC
			LIMIT 50
		`

		rows, err := database.Query(query, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch yield data",
				"details": err.Error(),
			})
			return
		}
		defer rows.Close()

		var distributions []map[string]interface{}
		totalYieldEarned := "0"

		for rows.Next() {
			var (
				id               int64
				assetID          int64
				cusip            string
				treasuryType     string
				distDate         string
				distType         string
				totalYield       string
				yieldPerToken    string
				tokensHeld       string
				userYield        string
				status           string
			)

			err := rows.Scan(
				&id,
				&assetID,
				&cusip,
				&treasuryType,
				&distDate,
				&distType,
				&totalYield,
				&yieldPerToken,
				&tokensHeld,
				&userYield,
				&status,
			)
			if err != nil {
				continue
			}

			distributions = append(distributions, map[string]interface{}{
				"id":                id,
				"asset_id":          assetID,
				"cusip":             cusip,
				"treasury_type":     treasuryType,
				"distribution_date": distDate,
				"distribution_type": distType,
				"total_yield":       totalYield,
				"yield_per_token":   yieldPerToken,
				"tokens_held":       tokensHeld,
				"user_yield":        userYield,
				"status":            status,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"user_address":        address,
			"distributions_count": len(distributions),
			"distributions":       distributions,
			"total_yield_earned":  totalYieldEarned,
		})
	}
}

// GetTreasuryMarketOrders returns market orders for an asset
// GET /api/v1/treasury/market/:assetId/orders?type=BUY
func GetTreasuryMarketOrders(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetIdStr := c.Param("assetId")
		assetId, err := strconv.ParseInt(assetIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid asset ID",
			})
			return
		}

		orderType := c.Query("type")
		var typePtr *string
		if orderType != "" {
			typePtr = &orderType
		}

		orders, err := db.GetTreasuryMarketOrders(database, assetId, typePtr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch market orders",
				"details": err.Error(),
			})
			return
		}

		// Separate buy and sell orders
		buyOrders := []interface{}{}
		sellOrders := []interface{}{}

		for _, order := range orders {
			if order.OrderType == "BUY" {
				buyOrders = append(buyOrders, order)
			} else {
				sellOrders = append(sellOrders, order)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"asset_id":   assetId,
			"buy_orders": buyOrders,
			"sell_orders": sellOrders,
			"total_buy_orders": len(buyOrders),
			"total_sell_orders": len(sellOrders),
		})
	}
}

// GetTreasuryTradeHistory returns recent trades for an asset
// GET /api/v1/treasury/market/:assetId/trades?limit=50
func GetTreasuryTradeHistory(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetIdStr := c.Param("assetId")
		assetId, err := strconv.ParseInt(assetIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid asset ID",
			})
			return
		}

		limitStr := c.DefaultQuery("limit", "50")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 500 {
			limit = 50
		}

		trades, err := db.GetTreasuryTrades(database, assetId, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch trades",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"asset_id": assetId,
			"count":    len(trades),
			"trades":   trades,
		})
	}
}

// CreateTreasuryOrder creates a buy or sell order (requires signature verification)
// POST /api/v1/treasury/market/order
func CreateTreasuryOrder(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AssetID       int64  `json:"asset_id" binding:"required"`
			OrderType     string `json:"order_type" binding:"required"`
			TokenAmount   string `json:"token_amount" binding:"required"`
			PricePerToken string `json:"price_per_token" binding:"required"`
			Signature     string `json:"signature" binding:"required"`
			ExpiresIn     int64  `json:"expires_in"` // seconds
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
				"details": err.Error(),
			})
			return
		}

		// Validate order type
		if req.OrderType != "BUY" && req.OrderType != "SELL" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid order type (must be BUY or SELL)",
			})
			return
		}

		// TODO: Verify signature and extract user address
		// For now, returning placeholder
		c.JSON(http.StatusOK, gin.H{
			"message": "Order creation requires smart contract interaction",
			"note":    "This endpoint should call TreasuryMarketplace contract",
			"next_steps": []string{
				"1. Verify user signature",
				"2. Call marketplace.createBuyOrder() or marketplace.createSellOrder()",
				"3. Wait for transaction confirmation",
				"4. Return order ID and tx hash",
			},
		})
	}
}

// CancelTreasuryOrder cancels an open order
// DELETE /api/v1/treasury/market/order/:orderId
func CancelTreasuryOrder(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIdStr := c.Param("orderId")
		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid order ID",
			})
			return
		}

		// TODO: Verify user owns this order
		// TODO: Call smart contract to cancel

		err = db.CancelTreasuryOrder(database, orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to cancel order",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Order cancelled",
			"order_id": orderId,
		})
	}
}

// ClaimTreasuryYield claims pending yield for user
// POST /api/v1/treasury/yield/claim
func ClaimTreasuryYield(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AssetID   int64  `json:"asset_id" binding:"required"`
			Signature string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
				"details": err.Error(),
			})
			return
		}

		// TODO: Verify signature and extract user address
		// TODO: Call TreasuryYieldDistributor.claimYield()

		c.JSON(http.StatusOK, gin.H{
			"message": "Yield claim requires smart contract interaction",
			"note":    "This endpoint should call TreasuryYieldDistributor contract",
			"next_steps": []string{
				"1. Verify user signature",
				"2. Calculate pending yield from contract",
				"3. Call yieldDistributor.claimYield(assetId)",
				"4. Return claimed amount and tx hash",
			},
		})
	}
}

// GetYieldDistributions returns yield distribution history
// GET /api/v1/treasury/yield/distributions?asset_id=1&limit=20
func GetYieldDistributions(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetIdStr := c.Query("asset_id")
		limitStr := c.DefaultQuery("limit", "20")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			limit = 20
		}

		query := `
			SELECT id, asset_id, distribution_date, distribution_type,
			       total_yield, yield_per_token, recipients_count,
			       total_distributed, status, created_at, distributed_at
			FROM treasury_yield_distributions
			WHERE 1=1
		`

		args := []interface{}{}
		argCount := 1

		if assetIdStr != "" {
			assetId, err := strconv.ParseInt(assetIdStr, 10, 64)
			if err == nil {
				query += " AND asset_id = $" + strconv.Itoa(argCount)
				args = append(args, assetId)
				argCount++
			}
		}

		query += " ORDER BY distribution_date DESC LIMIT $" + strconv.Itoa(argCount)
		args = append(args, limit)

		rows, err := database.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch distributions",
				"details": err.Error(),
			})
			return
		}
		defer rows.Close()

		var distributions []map[string]interface{}
		for rows.Next() {
			var (
				id               int64
				assetID          int64
				distDate         string
				distType         string
				totalYield       string
				yieldPerToken    string
				recipientsCount  int
				totalDistributed string
				status           string
				createdAt        string
				distributedAt    *string
			)

			err := rows.Scan(
				&id,
				&assetID,
				&distDate,
				&distType,
				&totalYield,
				&yieldPerToken,
				&recipientsCount,
				&totalDistributed,
				&status,
				&createdAt,
				&distributedAt,
			)
			if err != nil {
				continue
			}

			distributions = append(distributions, map[string]interface{}{
				"id":                id,
				"asset_id":          assetID,
				"distribution_date": distDate,
				"distribution_type": distType,
				"total_yield":       totalYield,
				"yield_per_token":   yieldPerToken,
				"recipients_count":  recipientsCount,
				"total_distributed": totalDistributed,
				"status":            status,
				"created_at":        createdAt,
				"distributed_at":    distributedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"count":         len(distributions),
			"distributions": distributions,
		})
	}
}

// GetTreasuryStats returns overall treasury statistics
// GET /api/v1/treasury/stats
func GetTreasuryStats(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Total assets
		var totalAssets int
		database.QueryRow("SELECT COUNT(*) FROM treasury_assets WHERE status = 'active'").Scan(&totalAssets)

		// Total TVL (Total Value Locked)
		var totalTVL string
		database.QueryRow(`
			SELECT COALESCE(SUM(tokens_outstanding::numeric * COALESCE(current_price::numeric, face_value::numeric)), 0)
			FROM treasury_assets
			WHERE status = 'active'
		`).Scan(&totalTVL)

		// Total holders
		var totalHolders int
		database.QueryRow("SELECT COUNT(DISTINCT user_address) FROM treasury_holdings WHERE tokens_held > 0").Scan(&totalHolders)

		// Total trades (24h)
		var trades24h int
		database.QueryRow("SELECT COUNT(*) FROM treasury_trades WHERE executed_at > NOW() - INTERVAL '24 hours'").Scan(&trades24h)

		// Total volume (24h)
		var volume24h string
		database.QueryRow(`
			SELECT COALESCE(SUM(total_value::numeric), 0)
			FROM treasury_trades
			WHERE executed_at > NOW() - INTERVAL '24 hours'
		`).Scan(&volume24h)

		// Active orders
		var activeOrders int
		database.QueryRow("SELECT COUNT(*) FROM treasury_market_orders WHERE status IN ('open', 'partial') AND expires_at > NOW()").Scan(&activeOrders)

		c.JSON(http.StatusOK, gin.H{
			"total_assets":    totalAssets,
			"total_tvl":       totalTVL,
			"total_holders":   totalHolders,
			"trades_24h":      trades24h,
			"volume_24h":      volume24h,
			"active_orders":   activeOrders,
			"last_updated":    "real-time",
		})
	}
}
