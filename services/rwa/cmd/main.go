package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/loyalty_points?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully with connection pooling")

	// Initialize Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Initialize services and handlers
	rwaService := NewRWAService(db)
	rwaHandler := NewRWAHandler(rwaService)

	// API routes
	api := r.Group("/api/rwa")
	{
		api.GET("/assets", rwaHandler.GetAssets)
		api.GET("/assets/:ticker", rwaHandler.GetAssetDetail)
		api.POST("/orders", rwaHandler.CreateOrder)
		api.GET("/orders/:address", rwaHandler.GetUserOrders)
		api.DELETE("/orders/:orderId", rwaHandler.CancelOrder)
		api.GET("/holdings/:address", rwaHandler.GetUserHoldings)
		api.GET("/prices/:ticker", rwaHandler.GetPriceHistory)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "rwa",
			"version": "1.0.0",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("RWA_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("RWA Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// RWAService handles business logic
type RWAService struct {
	db *sql.DB
}

func NewRWAService(db *sql.DB) *RWAService {
	return &RWAService{db: db}
}

// RWAHandler handles HTTP requests
type RWAHandler struct {
	service *RWAService
}

func NewRWAHandler(service *RWAService) *RWAHandler {
	return &RWAHandler{service: service}
}

// GetAssets returns list of available RWA assets
func (h *RWAHandler) GetAssets(c *gin.Context) {
	assetType := c.Query("type") // stock, bond, commodity
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	query := `
		SELECT
			ticker, name, asset_type, issuer,
			current_price, price_change_24h, price_change_7d,
			market_cap, total_supply, metadata
		FROM rwa_assets
		WHERE is_active = true
	`

	args := []interface{}{}
	argCount := 0

	if assetType != "" {
		argCount++
		query += fmt.Sprintf(" AND asset_type = $%d", argCount)
		args = append(args, assetType)
	}

	query += ` ORDER BY market_cap DESC`

	// Add pagination
	// Note: In production, use proper pagination with LIMIT and OFFSET
	// For now, we'll return all matching results

	rows, err := h.service.db.Query(query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	assets := []map[string]interface{}{}
	for rows.Next() {
		var ticker, name, assetType, issuer string
		var currentPrice, priceChange24h, priceChange7d, marketCap, totalSupply float64
		var metadataJSON []byte

		err := rows.Scan(
			&ticker, &name, &assetType, &issuer,
			&currentPrice, &priceChange24h, &priceChange7d,
			&marketCap, &totalSupply, &metadataJSON,
		)
		if err != nil {
			continue
		}

		var metadata map[string]interface{}
		json.Unmarshal(metadataJSON, &metadata)

		asset := map[string]interface{}{
			"ticker":           ticker,
			"name":             name,
			"asset_type":       assetType,
			"issuer":           issuer,
			"current_price":    currentPrice,
			"price_change_24h": priceChange24h,
			"price_change_7d":  priceChange7d,
			"market_cap":       marketCap,
			"total_supply":     totalSupply,
			"metadata":         metadata,
		}
		assets = append(assets, asset)
	}

	c.JSON(200, gin.H{
		"assets": assets,
		"page":   page,
		"limit":  limit,
		"total":  len(assets),
	})
}

// GetAssetDetail returns detailed information about a specific asset
func (h *RWAHandler) GetAssetDetail(c *gin.Context) {
	ticker := c.Param("ticker")

	var asset struct {
		Ticker            string
		Name              string
		AssetType         string
		Issuer            string
		CurrentPrice      float64
		PriceChange24h    float64
		PriceChange7d     float64
		MarketCap         float64
		TotalSupply       float64
		CirculatingSupply sql.NullFloat64
		ContractAddress   sql.NullString
		Metadata          []byte
	}

	err := h.service.db.QueryRow(`
		SELECT
			ticker, name, asset_type, issuer,
			current_price, price_change_24h, price_change_7d,
			market_cap, total_supply, circulating_supply,
			contract_address, metadata
		FROM rwa_assets
		WHERE ticker = $1 AND is_active = true
	`, ticker).Scan(
		&asset.Ticker, &asset.Name, &asset.AssetType, &asset.Issuer,
		&asset.CurrentPrice, &asset.PriceChange24h, &asset.PriceChange7d,
		&asset.MarketCap, &asset.TotalSupply, &asset.CirculatingSupply,
		&asset.ContractAddress, &asset.Metadata,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Asset not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	var metadata map[string]interface{}
	json.Unmarshal(asset.Metadata, &metadata)

	result := map[string]interface{}{
		"ticker":            asset.Ticker,
		"name":              asset.Name,
		"asset_type":        asset.AssetType,
		"issuer":            asset.Issuer,
		"current_price":     asset.CurrentPrice,
		"price_change_24h":  asset.PriceChange24h,
		"price_change_7d":   asset.PriceChange7d,
		"market_cap":        asset.MarketCap,
		"total_supply":      asset.TotalSupply,
		"metadata":          metadata,
	}

	if asset.CirculatingSupply.Valid {
		result["circulating_supply"] = asset.CirculatingSupply.Float64
	}

	if asset.ContractAddress.Valid {
		result["contract_address"] = asset.ContractAddress.String
	}

	c.JSON(200, result)
}

// CreateOrder creates a new buy/sell order
func (h *RWAHandler) CreateOrder(c *gin.Context) {
	var req struct {
		UserAddress string  `json:"user_address"`
		Ticker      string  `json:"ticker"`
		OrderType   string  `json:"type"`       // "buy" or "sell"
		OrderStyle  string  `json:"style"`      // "market" or "limit"
		Amount      float64 `json:"amount"`     // USDC for buy, tokens for sell
		LimitPrice  float64 `json:"limitPrice"` // For limit orders
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Validate inputs
	if req.Amount <= 0 {
		c.JSON(400, gin.H{"error": "Amount must be positive"})
		return
	}

	if req.OrderType != "buy" && req.OrderType != "sell" {
		c.JSON(400, gin.H{"error": "Invalid order type"})
		return
	}

	if req.OrderStyle != "market" && req.OrderStyle != "limit" {
		req.OrderStyle = "market" // Default to market order
	}

	// Process order
	result, err := h.service.ProcessOrder(req.UserAddress, req.Ticker, req.OrderType, req.OrderStyle, req.Amount, req.LimitPrice)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// ProcessOrder processes a buy/sell order
func (s *RWAService) ProcessOrder(userAddress, ticker, orderType, orderStyle string, amount, limitPrice float64) (map[string]interface{}, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get current asset price
	var currentPrice float64
	err = tx.QueryRow(`
		SELECT current_price FROM rwa_assets WHERE ticker = $1 AND is_active = true
	`, ticker).Scan(&currentPrice)

	if err != nil {
		return nil, err
	}

	// For market orders, use current price
	orderPrice := currentPrice
	if orderStyle == "limit" && limitPrice > 0 {
		orderPrice = limitPrice
	}

	// Check balance for buy orders
	if orderType == "buy" {
		// Check user has enough earnings (from vault earnings) with row lock
		var userEarnings float64
		err = tx.QueryRow(`
			SELECT COALESCE(SUM(earned), 0)
			FROM vault_positions
			WHERE user_address = $1
			FOR UPDATE
		`, userAddress).Scan(&userEarnings)

		if err != nil || userEarnings < amount {
			return nil, err
		}
	} else {
		// Check holdings for sell orders
		var currentHolding float64
		err = tx.QueryRow(`
			SELECT amount FROM rwa_holdings
			WHERE user_address = $1 AND asset_ticker = $2
		`, userAddress, ticker).Scan(&currentHolding)

		if err != nil || currentHolding < amount {
			return nil, err
		}
	}

	// Create order
	var orderId int
	err = tx.QueryRow(`
		INSERT INTO rwa_orders (
			user_address, asset_ticker, order_type, order_style,
			amount, price, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userAddress, ticker, orderType, orderStyle, amount, orderPrice, "pending").Scan(&orderId)

	if err != nil {
		return nil, err
	}

	// Execute market orders immediately
	if orderStyle == "market" {
		if orderType == "buy" {
			// Calculate tokens received
			tokensReceived := amount / currentPrice

			// Update or create holding
			// Note: Using 'amount' (actual USDC spent) for average cost calculation
			// rather than currentPrice * tokensReceived to avoid precision issues
			_, err = tx.Exec(`
				INSERT INTO rwa_holdings (
					user_address, asset_ticker, amount, average_cost, current_value
				) VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (user_address, asset_ticker)
				DO UPDATE SET
					amount = rwa_holdings.amount + $3,
					average_cost = (rwa_holdings.average_cost * rwa_holdings.amount + $5) /
								  (rwa_holdings.amount + $3),
					current_value = (rwa_holdings.amount + $3) * $4,
					last_updated = NOW()
			`, userAddress, ticker, tokensReceived, currentPrice, amount)

			if err != nil {
				return nil, err
			}

			// Deduct from earnings (simplified - in production, implement proper accounting)
			// PostgreSQL doesn't support LIMIT in UPDATE, using subquery instead
			_, err = tx.Exec(`
				UPDATE vault_positions
				SET earned = earned - $1
				WHERE id = (
					SELECT id FROM vault_positions
					WHERE user_address = $2 AND earned >= $1
					ORDER BY earned DESC
					LIMIT 1
				)
			`, amount, userAddress)

		} else { // sell
			// Update holdings
			_, err = tx.Exec(`
				UPDATE rwa_holdings
				SET amount = amount - $1,
					current_value = (amount - $1) * $2,
					last_updated = NOW()
				WHERE user_address = $3 AND asset_ticker = $4
			`, amount, currentPrice, userAddress, ticker)

			if err != nil {
				return nil, err
			}

			// Add proceeds to vault earnings
			proceeds := amount * currentPrice
			_, err = tx.Exec(`
				INSERT INTO vault_earnings (user_address, protocol, amount, apy)
				VALUES ($1, 'RWA_SALE', $2, 0)
			`, userAddress, proceeds)
		}

		// Update order status
		_, err = tx.Exec(`
			UPDATE rwa_orders
			SET status = 'executed',
				filled_amount = amount,
				average_price = $1,
				executed_at = NOW()
			WHERE id = $2
		`, currentPrice, orderId)

		if err != nil {
			return nil, err
		}
	}

	// Record in transaction history
	// Determine status based on order style (Go doesn't support ternary operator)
	historyStatus := "pending"
	if orderStyle == "market" {
		historyStatus = "executed"
	}
	_, err = tx.Exec(`
		INSERT INTO transaction_history (
			user_address, transaction_type, asset, amount, price, status
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, userAddress, "rwa_"+orderType, ticker, amount, orderPrice, historyStatus)

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	status := "pending"
	if orderStyle == "market" {
		status = "executed"
	}

	return map[string]interface{}{
		"success":  true,
		"order_id": orderId,
		"status":   status,
		"message":  "Order created successfully",
	}, nil
}

// GetUserOrders returns user's order history
func (h *RWAHandler) GetUserOrders(c *gin.Context) {
	address := c.Param("address")
	status := c.Query("status") // filter by status

	query := `
		SELECT
			id, asset_ticker, order_type, order_style,
			amount, price, status, filled_amount,
			average_price, created_at, executed_at
		FROM rwa_orders
		WHERE user_address = $1
	`

	args := []interface{}{address}

	if status != "" {
		query += ` AND status = $2`
		args = append(args, status)
	}

	query += ` ORDER BY created_at DESC LIMIT 100`

	rows, err := h.service.db.Query(query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	orders := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var ticker, orderType, orderStyle, status string
		var amount, price, filledAmount float64
		var averagePrice sql.NullFloat64
		var createdAt string
		var executedAt sql.NullString

		err := rows.Scan(
			&id, &ticker, &orderType, &orderStyle,
			&amount, &price, &status, &filledAmount,
			&averagePrice, &createdAt, &executedAt,
		)
		if err != nil {
			continue
		}

		order := map[string]interface{}{
			"id":            id,
			"ticker":        ticker,
			"type":          orderType,
			"style":         orderStyle,
			"amount":        amount,
			"price":         price,
			"status":        status,
			"filled_amount": filledAmount,
			"created_at":    createdAt,
		}

		if averagePrice.Valid {
			order["average_price"] = averagePrice.Float64
		}
		if executedAt.Valid {
			order["executed_at"] = executedAt.String
		}

		orders = append(orders, order)
	}

	c.JSON(200, gin.H{"orders": orders})
}

// CancelOrder cancels a pending order
func (h *RWAHandler) CancelOrder(c *gin.Context) {
	orderId := c.Param("orderId")

	result, err := h.service.db.Exec(`
		UPDATE rwa_orders
		SET status = 'cancelled', cancelled_at = NOW()
		WHERE id = $1 AND status = 'pending'
	`, orderId)

	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Order not found or already processed"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Order cancelled successfully",
	})
}

// GetUserHoldings returns user's RWA holdings
func (h *RWAHandler) GetUserHoldings(c *gin.Context) {
	address := c.Param("address")

	rows, err := h.service.db.Query(`
		SELECT
			h.asset_ticker, a.name, a.asset_type,
			h.amount, h.average_cost, h.current_value,
			h.pnl, h.pnl_percentage,
			a.current_price, a.price_change_24h
		FROM rwa_holdings h
		JOIN rwa_assets a ON h.asset_ticker = a.ticker
		WHERE h.user_address = $1 AND h.amount > 0
		ORDER BY h.current_value DESC
	`, address)

	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	holdings := []map[string]interface{}{}
	totalValue := 0.0
	totalPnl := 0.0

	for rows.Next() {
		var ticker, name, assetType string
		var amount, averageCost, currentValue, pnl, pnlPercentage, currentPrice, priceChange24h float64

		err := rows.Scan(
			&ticker, &name, &assetType,
			&amount, &averageCost, &currentValue,
			&pnl, &pnlPercentage,
			&currentPrice, &priceChange24h,
		)
		if err != nil {
			continue
		}

		holding := map[string]interface{}{
			"ticker":            ticker,
			"name":              name,
			"asset_type":        assetType,
			"amount":            amount,
			"average_cost":      averageCost,
			"current_value":     currentValue,
			"pnl":               pnl,
			"pnl_percentage":    pnlPercentage,
			"current_price":     currentPrice,
			"price_change_24h":  priceChange24h,
		}

		holdings = append(holdings, holding)
		totalValue += currentValue
		totalPnl += pnl
	}

	c.JSON(200, gin.H{
		"holdings":     holdings,
		"total_value":  totalValue,
		"total_pnl":    totalPnl,
		"total_assets": len(holdings),
	})
}

// GetPriceHistory returns price history for an asset
func (h *RWAHandler) GetPriceHistory(c *gin.Context) {
	ticker := c.Param("ticker")
	period := c.DefaultQuery("period", "7d") // 1d, 7d, 30d, 90d

	// Calculate time range based on period
	var interval string
	switch period {
	case "1d":
		interval = "1 day"
	case "7d":
		interval = "7 days"
	case "30d":
		interval = "30 days"
	case "90d":
		interval = "90 days"
	default:
		interval = "7 days"
	}

	rows, err := h.service.db.Query(`
		SELECT price, high_24h, low_24h, volume_24h, timestamp
		FROM price_history
		WHERE asset_ticker = $1
		AND timestamp > NOW() - INTERVAL '` + interval + `'
		ORDER BY timestamp DESC
		LIMIT 500
	`, ticker)

	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	prices := []map[string]interface{}{}
	for rows.Next() {
		var price, high24h, low24h, volume24h float64
		var timestamp string

		err := rows.Scan(&price, &high24h, &low24h, &volume24h, &timestamp)
		if err != nil {
			continue
		}

		prices = append(prices, map[string]interface{}{
			"price":      price,
			"high_24h":   high24h,
			"low_24h":    low24h,
			"volume_24h": volume24h,
			"timestamp":  timestamp,
		})
	}

	// Get current price
	var currentPrice float64
	h.service.db.QueryRow(`
		SELECT current_price FROM rwa_assets WHERE ticker = $1
	`, ticker).Scan(&currentPrice)

	c.JSON(200, gin.H{
		"ticker":        ticker,
		"period":        period,
		"current_price": currentPrice,
		"prices":        prices,
	})
}