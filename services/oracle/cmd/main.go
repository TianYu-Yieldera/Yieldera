package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize random seed for price simulation
	rand.Seed(time.Now().UnixNano())

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

	// Initialize services
	priceService := NewPriceService(db)

	// Start background price updater
	go priceService.StartPriceUpdater()

	// Initialize Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api/oracle")
	{
		api.GET("/price/:ticker", priceService.GetPrice)
		api.POST("/prices", priceService.GetMultiplePrices)
		api.GET("/apy/:protocol", priceService.GetAPY)
		api.GET("/apys", priceService.GetAllAPYs)
		api.GET("/stats", priceService.GetMarketStats)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "oracle",
			"version": "1.0.0",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("ORACLE_PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Oracle Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// PriceService handles price feeds and APY calculations
type PriceService struct {
	db *sql.DB
}

func NewPriceService(db *sql.DB) *PriceService {
	return &PriceService{db: db}
}

// StartPriceUpdater runs in background to update prices
func (s *PriceService) StartPriceUpdater() {
	log.Println("Starting price updater...")

	// Update immediately on start
	s.updatePrices()
	s.updateAPYs()

	// Create tickers for periodic updates
	priceTicker := time.NewTicker(60 * time.Second)      // Update prices every minute
	apyTicker := time.NewTicker(5 * time.Minute)         // Update APYs every 5 minutes

	for {
		select {
		case <-priceTicker.C:
			s.updatePrices()
		case <-apyTicker.C:
			s.updateAPYs()
		}
	}
}

// updatePrices updates RWA asset prices
func (s *PriceService) updatePrices() {
	// Get all active assets
	rows, err := s.db.Query(`
		SELECT ticker, current_price, asset_type
		FROM rwa_assets
		WHERE is_active = true
	`)
	if err != nil {
		log.Printf("Error fetching assets: %v", err)
		return
	}
	defer rows.Close()

	// Collect all price updates
	type priceUpdate struct {
		ticker         string
		newPrice       float64
		priceChange24h float64
	}
	updates := []priceUpdate{}

	for rows.Next() {
		var ticker, assetType string
		var currentPrice float64

		err := rows.Scan(&ticker, &currentPrice, &assetType)
		if err != nil {
			continue
		}

		// Simulate price movement (in production, fetch from real sources)
		newPrice := s.simulatePriceMovement(currentPrice, assetType)
		priceChange24h := ((newPrice - currentPrice) / currentPrice) * 100

		updates = append(updates, priceUpdate{
			ticker:         ticker,
			newPrice:       newPrice,
			priceChange24h: priceChange24h,
		})
	}

	// Use transaction for atomic updates
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}
	defer tx.Rollback()

	// Apply all updates within transaction
	for _, update := range updates {
		// Update asset price
		_, err = tx.Exec(`
			UPDATE rwa_assets
			SET current_price = $1,
				price_change_24h = $2,
				updated_at = NOW()
			WHERE ticker = $3
		`, update.newPrice, update.priceChange24h, update.ticker)

		if err != nil {
			log.Printf("Error updating price for %s: %v", update.ticker, err)
			return
		}

		// Record in price history
		_, err = tx.Exec(`
			INSERT INTO price_history (
				asset_ticker, price, high_24h, low_24h,
				volume_24h, source
			) VALUES ($1, $2, $3, $4, $5, 'internal')
		`, update.ticker, update.newPrice, update.newPrice*1.02, update.newPrice*0.98, rand.Float64()*1000000)

		if err != nil {
			log.Printf("Error recording price history for %s: %v", update.ticker, err)
			return
		}
	}

	// Commit all updates atomically
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing price updates: %v", err)
		return
	}

	log.Printf("Updated prices for all assets")
}

// updateAPYs updates DeFi protocol APYs
func (s *PriceService) updateAPYs() {
	// Get all active protocols
	rows, err := s.db.Query(`
		SELECT name, current_apy, protocol_type
		FROM defi_protocols
		WHERE is_active = true
	`)
	if err != nil {
		log.Printf("Error fetching protocols: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name, protocolType string
		var currentAPY float64

		err := rows.Scan(&name, &currentAPY, &protocolType)
		if err != nil {
			continue
		}

		// Simulate APY changes (in production, fetch from real sources)
		newAPY := s.simulateAPYChange(currentAPY, protocolType)

		// Update protocol APY
		_, err = s.db.Exec(`
			UPDATE defi_protocols
			SET current_apy = $1,
				updated_at = NOW()
			WHERE name = $2
		`, newAPY, name)

		if err != nil {
			log.Printf("Error updating APY for %s: %v", name, err)
			continue
		}

		// Record in APY history
		_, err = s.db.Exec(`
			INSERT INTO apy_history (protocol, apy, tvl)
			VALUES ($1, $2, $3)
		`, name, newAPY, rand.Float64()*1000000000)

		if err != nil {
			log.Printf("Error recording APY history for %s: %v", name, err)
		}
	}

	log.Printf("Updated APYs for all protocols")
}

// simulatePriceMovement simulates realistic price movements
func (s *PriceService) simulatePriceMovement(currentPrice float64, assetType string) float64 {
	// Different volatility for different asset types
	var volatility float64
	switch assetType {
	case "stock":
		volatility = 0.02 // 2% volatility
	case "bond":
		volatility = 0.001 // 0.1% volatility
	case "commodity":
		volatility = 0.015 // 1.5% volatility
	default:
		volatility = 0.01
	}

	// Generate random price movement
	change := (rand.Float64() - 0.5) * 2 * volatility
	newPrice := currentPrice * (1 + change)

	// Ensure price doesn't go negative
	if newPrice < 0 {
		newPrice = currentPrice * 0.99
	}

	return math.Round(newPrice*100) / 100
}

// simulateAPYChange simulates APY fluctuations
func (s *PriceService) simulateAPYChange(currentAPY float64, protocolType string) float64 {
	// Different volatility for different protocol types
	var volatility float64
	switch protocolType {
	case "lending":
		volatility = 0.05 // 5% volatility
	case "dex":
		volatility = 0.1 // 10% volatility
	case "derivatives":
		volatility = 0.15 // 15% volatility
	case "yield":
		volatility = 0.08 // 8% volatility
	default:
		volatility = 0.05
	}

	// Generate random APY movement
	change := (rand.Float64() - 0.5) * 2 * volatility
	newAPY := currentAPY * (1 + change)

	// Keep APY within reasonable bounds
	if newAPY < 0.5 {
		newAPY = 0.5
	} else if newAPY > 50 {
		newAPY = 50
	}

	return math.Round(newAPY*100) / 100
}

// GetPrice returns current price for a single asset
func (s *PriceService) GetPrice(c *gin.Context) {
	ticker := c.Param("ticker")

	var result struct {
		Ticker         string
		Name           string
		CurrentPrice   float64
		PriceChange24h float64
		High24h        float64
		Low24h         float64
		UpdatedAt      string
	}

	err := s.db.QueryRow(`
		SELECT
			a.ticker, a.name, a.current_price, a.price_change_24h,
			COALESCE(
				(SELECT MAX(price) FROM price_history
				 WHERE asset_ticker = a.ticker
				 AND timestamp > NOW() - INTERVAL '24 hours'),
				a.current_price * 1.02
			) as high_24h,
			COALESCE(
				(SELECT MIN(price) FROM price_history
				 WHERE asset_ticker = a.ticker
				 AND timestamp > NOW() - INTERVAL '24 hours'),
				a.current_price * 0.98
			) as low_24h,
			a.updated_at
		FROM rwa_assets a
		WHERE a.ticker = $1 AND a.is_active = true
	`, ticker).Scan(
		&result.Ticker, &result.Name, &result.CurrentPrice,
		&result.PriceChange24h, &result.High24h, &result.Low24h,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Asset not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, gin.H{
		"ticker":           result.Ticker,
		"name":             result.Name,
		"current_price":    result.CurrentPrice,
		"price_change_24h": result.PriceChange24h,
		"high_24h":         result.High24h,
		"low_24h":          result.Low24h,
		"updated_at":       result.UpdatedAt,
		"source":           "internal",
	})
}

// GetMultiplePrices returns prices for multiple assets
func (s *PriceService) GetMultiplePrices(c *gin.Context) {
	var req struct {
		Tickers []string `json:"tickers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	prices := []map[string]interface{}{}

	for _, ticker := range req.Tickers {
		var price float64
		var priceChange24h float64

		err := s.db.QueryRow(`
			SELECT current_price, price_change_24h
			FROM rwa_assets
			WHERE ticker = $1 AND is_active = true
		`, ticker).Scan(&price, &priceChange24h)

		if err != nil {
			continue
		}

		prices = append(prices, map[string]interface{}{
			"ticker":           ticker,
			"current_price":    price,
			"price_change_24h": priceChange24h,
		})
	}

	c.JSON(200, gin.H{"prices": prices})
}

// GetAPY returns current APY for a protocol
func (s *PriceService) GetAPY(c *gin.Context) {
	protocol := c.Param("protocol")

	var result struct {
		Protocol     string
		ProtocolType string
		RiskLevel    string
		CurrentAPY   float64
		TVL          float64
		UpdatedAt    string
	}

	err := s.db.QueryRow(`
		SELECT
			name, protocol_type, risk_level,
			current_apy, tvl, updated_at
		FROM defi_protocols
		WHERE name = $1 AND is_active = true
	`, protocol).Scan(
		&result.Protocol, &result.ProtocolType, &result.RiskLevel,
		&result.CurrentAPY, &result.TVL, &result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Protocol not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	// Get historical APY for comparison
	var apy7dAgo float64
	s.db.QueryRow(`
		SELECT apy FROM apy_history
		WHERE protocol = $1
		AND timestamp < NOW() - INTERVAL '7 days'
		ORDER BY timestamp DESC
		LIMIT 1
	`, protocol).Scan(&apy7dAgo)

	apyChange7d := 0.0
	if apy7dAgo > 0 {
		apyChange7d = result.CurrentAPY - apy7dAgo
	}

	c.JSON(200, gin.H{
		"protocol":      result.Protocol,
		"protocol_type": result.ProtocolType,
		"risk_level":    result.RiskLevel,
		"current_apy":   result.CurrentAPY,
		"apy_change_7d": apyChange7d,
		"tvl":           result.TVL,
		"updated_at":    result.UpdatedAt,
	})
}

// GetAllAPYs returns APYs for all protocols
func (s *PriceService) GetAllAPYs(c *gin.Context) {
	rows, err := s.db.Query(`
		SELECT
			name, protocol_type, risk_level,
			current_apy, tvl
		FROM defi_protocols
		WHERE is_active = true
		ORDER BY current_apy DESC
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	protocols := []map[string]interface{}{}
	for rows.Next() {
		var name, protocolType, riskLevel string
		var apy, tvl float64

		err := rows.Scan(&name, &protocolType, &riskLevel, &apy, &tvl)
		if err != nil {
			continue
		}

		protocols = append(protocols, map[string]interface{}{
			"protocol":      name,
			"protocol_type": protocolType,
			"risk_level":    riskLevel,
			"current_apy":   apy,
			"tvl":           tvl,
		})
	}

	// Calculate average APYs by risk level
	avgAPYs := map[string]float64{
		"low":    0,
		"medium": 0,
		"high":   0,
	}
	counts := map[string]int{
		"low":    0,
		"medium": 0,
		"high":   0,
	}

	for _, p := range protocols {
		risk := p["risk_level"].(string)
		apy := p["current_apy"].(float64)
		avgAPYs[risk] += apy
		counts[risk]++
	}

	for risk, count := range counts {
		if count > 0 {
			avgAPYs[risk] = math.Round((avgAPYs[risk]/float64(count))*100) / 100
		}
	}

	c.JSON(200, gin.H{
		"protocols": protocols,
		"averages":  avgAPYs,
		"count":     len(protocols),
	})
}

// GetMarketStats returns overall market statistics
func (s *PriceService) GetMarketStats(c *gin.Context) {
	var stats struct {
		TotalRWAMarketCap   float64
		TotalVaultTVL       float64
		AverageAPY          float64
		TopGainer           string
		TopGainerChange     float64
		TopLoser            string
		TopLoserChange      float64
		ActiveAssets        int
		ActiveProtocols     int
		TotalUsers          int
		TotalTransactions   int
	}

	// Get RWA market cap
	s.db.QueryRow(`
		SELECT COALESCE(SUM(market_cap), 0)
		FROM rwa_assets
		WHERE is_active = true
	`).Scan(&stats.TotalRWAMarketCap)

	// Get Vault TVL
	s.db.QueryRow(`
		SELECT COALESCE(SUM(tvl), 0)
		FROM defi_protocols
		WHERE is_active = true
	`).Scan(&stats.TotalVaultTVL)

	// Get average APY
	s.db.QueryRow(`
		SELECT COALESCE(AVG(current_apy), 0)
		FROM defi_protocols
		WHERE is_active = true
	`).Scan(&stats.AverageAPY)

	// Get top gainer
	s.db.QueryRow(`
		SELECT ticker, price_change_24h
		FROM rwa_assets
		WHERE is_active = true
		ORDER BY price_change_24h DESC
		LIMIT 1
	`).Scan(&stats.TopGainer, &stats.TopGainerChange)

	// Get top loser
	s.db.QueryRow(`
		SELECT ticker, price_change_24h
		FROM rwa_assets
		WHERE is_active = true
		ORDER BY price_change_24h ASC
		LIMIT 1
	`).Scan(&stats.TopLoser, &stats.TopLoserChange)

	// Count active assets and protocols
	s.db.QueryRow(`
		SELECT COUNT(*) FROM rwa_assets WHERE is_active = true
	`).Scan(&stats.ActiveAssets)

	s.db.QueryRow(`
		SELECT COUNT(*) FROM defi_protocols WHERE is_active = true
	`).Scan(&stats.ActiveProtocols)

	// Count users and transactions
	s.db.QueryRow(`
		SELECT COUNT(DISTINCT user_address) FROM vault_positions
	`).Scan(&stats.TotalUsers)

	s.db.QueryRow(`
		SELECT COUNT(*) FROM transaction_history
	`).Scan(&stats.TotalTransactions)

	c.JSON(200, gin.H{
		"market_stats": map[string]interface{}{
			"total_rwa_market_cap": fmt.Sprintf("$%.2fB", stats.TotalRWAMarketCap/1000000000),
			"total_vault_tvl":      fmt.Sprintf("$%.2fB", stats.TotalVaultTVL/1000000000),
			"average_apy":          fmt.Sprintf("%.2f%%", stats.AverageAPY),
			"top_gainer": map[string]interface{}{
				"ticker": stats.TopGainer,
				"change": fmt.Sprintf("+%.2f%%", stats.TopGainerChange),
			},
			"top_loser": map[string]interface{}{
				"ticker": stats.TopLoser,
				"change": fmt.Sprintf("%.2f%%", stats.TopLoserChange),
			},
			"active_assets":      stats.ActiveAssets,
			"active_protocols":   stats.ActiveProtocols,
			"total_users":        stats.TotalUsers,
			"total_transactions": stats.TotalTransactions,
		},
	})
}