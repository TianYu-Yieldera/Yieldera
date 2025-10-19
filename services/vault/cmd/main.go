package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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

	// Configure connection pool for optimal performance
	// MaxOpenConns: Maximum number of open connections to the database
	// MaxIdleConns: Maximum number of idle connections in the pool
	// ConnMaxLifetime: Maximum amount of time a connection may be reused
	maxOpenConns := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	log.Printf("Database pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v",
		maxOpenConns, maxIdleConns, connMaxLifetime)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Initialize Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Initialize services and handlers
	vaultService := NewVaultService(db)
	vaultHandler := NewVaultHandler(vaultService)

	// API routes
	api := r.Group("/api/vault")
	{
		api.POST("/deposit", vaultHandler.Deposit)
		api.POST("/withdraw", vaultHandler.Withdraw)
		api.GET("/balance/:address", vaultHandler.GetBalance)
		api.GET("/strategies", vaultHandler.GetStrategies)
		api.POST("/stake", vaultHandler.StakeToProtocol)
		api.POST("/unstake", vaultHandler.UnstakeFromProtocol)
		api.GET("/earnings/:address", vaultHandler.GetEarnings)
		api.GET("/positions/:address", vaultHandler.GetPositions)
		api.GET("/protocols", vaultHandler.GetProtocols)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "vault",
			"version": "1.0.0",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("VAULT_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Vault Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// VaultService handles business logic
type VaultService struct {
	db *sql.DB
}

func NewVaultService(db *sql.DB) *VaultService {
	return &VaultService{db: db}
}

// VaultHandler handles HTTP requests
type VaultHandler struct {
	service *VaultService
}

func NewVaultHandler(service *VaultService) *VaultHandler {
	return &VaultHandler{service: service}
}

// Deposit handles deposit requests
func (h *VaultHandler) Deposit(c *gin.Context) {
	var req struct {
		UserAddress    string  `json:"user_address"`
		Amount         float64 `json:"amount"`
		Mode           string  `json:"mode"`     // "smart" or "manual"
		Strategy       string  `json:"strategy"` // "conservative", "balanced", "aggressive"
		IdempotencyKey string  `json:"idempotency_key"` // Optional: for preventing duplicate requests
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate inputs
	if req.Amount <= 0 {
		c.JSON(400, gin.H{"error": "Amount must be positive"})
		return
	}

	if req.Mode != "smart" && req.Mode != "manual" {
		req.Mode = "smart" // Default to smart mode
	}

	if req.Strategy == "" {
		req.Strategy = "balanced" // Default strategy
	}

	// Generate idempotency key if not provided
	if req.IdempotencyKey == "" {
		// Fallback: generate from user address, amount, and timestamp
		req.IdempotencyKey = fmt.Sprintf("%s-%f-%d", req.UserAddress, req.Amount, time.Now().Unix())
	}

	// Process deposit with idempotency protection
	result, err := h.service.ProcessDeposit(req.UserAddress, req.Amount, req.Mode, req.Strategy, req.IdempotencyKey)
	if err != nil {
		// Check if this is a duplicate request error
		if err.Error() == "duplicate_request" {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Deposit already processed",
				"duplicate": true,
			})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// ProcessDeposit processes a deposit with idempotency protection
func (s *VaultService) ProcessDeposit(userAddress string, amount float64, mode, strategy, idempotencyKey string) (map[string]interface{}, error) {
	// First check if this idempotency key was already used
	var existingDepositId int
	var existingStatus string
	err := s.db.QueryRow(`
		SELECT id, status FROM vault_deposits WHERE idempotency_key = $1
	`, idempotencyKey).Scan(&existingDepositId, &existingStatus)

	if err == nil {
		// Idempotency key already exists
		if existingStatus == "confirmed" {
			// Return success but mark as duplicate
			return nil, errors.New("duplicate_request")
		}
		// If status is pending or failed, we can retry
	} else if err != sql.ErrNoRows {
		// Database error
		return nil, err
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check user points balance (with row lock to prevent concurrent updates)
	var userPoints float64
	err = tx.QueryRow(`
		SELECT COALESCE(points, 0) FROM points
		WHERE user_address = $1
		FOR UPDATE
	`, userAddress).Scan(&userPoints)

	if err == sql.ErrNoRows || userPoints < amount {
		return nil, fmt.Errorf("insufficient points: have %.2f, need %.2f", userPoints, amount)
	}

	// Deduct points
	_, err = tx.Exec(`
		UPDATE points SET points = points - $1, updated_at = NOW()
		WHERE user_address = $2
	`, amount, userAddress)
	if err != nil {
		return nil, err
	}

	// Record deposit with idempotency key
	var depositId int
	err = tx.QueryRow(`
		INSERT INTO vault_deposits (user_address, amount, mode, strategy, status, idempotency_key)
		VALUES ($1, $2, $3, $4, 'confirmed', $5)
		ON CONFLICT (idempotency_key) DO NOTHING
		RETURNING id
	`, userAddress, amount, mode, strategy, idempotencyKey).Scan(&depositId)

	if err == sql.ErrNoRows {
		// Conflict occurred - another request with same key was processed
		return nil, errors.New("duplicate_request")
	}
	if err != nil {
		return nil, err
	}

	// Allocate to protocols based on mode
	if mode == "smart" {
		allocations := getStrategyAllocations(strategy)

		// Validate that allocations sum to 100%
		totalPercentage := 0.0
		for _, percentage := range allocations {
			totalPercentage += percentage
		}
		if totalPercentage != 100 {
			return nil, fmt.Errorf("allocations must sum to 100%%, got %.2f%%", totalPercentage)
		}

		for protocol, percentage := range allocations {
			protocolAmount := amount * percentage / 100

			_, err = tx.Exec(`
				INSERT INTO vault_positions (user_address, protocol, amount, apy)
				VALUES ($1, $2, $3, (SELECT current_apy FROM defi_protocols WHERE name = $2))
				ON CONFLICT (user_address, protocol)
				DO UPDATE SET
					amount = vault_positions.amount + $3,
					updated_at = NOW()
			`, userAddress, protocol, protocolAmount)

			if err != nil {
				return nil, err
			}
		}
	}

	// Record in transaction history
	_, err = tx.Exec(`
		INSERT INTO transaction_history (user_address, transaction_type, amount, status)
		VALUES ($1, 'vault_deposit', $2, 'confirmed')
	`, userAddress, amount)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":    true,
		"deposit_id": depositId,
		"amount":     amount,
		"mode":       mode,
		"strategy":   strategy,
		"message":    "Deposit successful",
	}, nil
}

// Withdraw handles withdrawal requests
func (h *VaultHandler) Withdraw(c *gin.Context) {
	var req struct {
		UserAddress string  `json:"user_address"`
		Amount      float64 `json:"amount"`
		Emergency   bool    `json:"emergency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	result, err := h.service.ProcessWithdraw(req.UserAddress, req.Amount, req.Emergency)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// ProcessWithdraw processes a withdrawal
func (s *VaultService) ProcessWithdraw(userAddress string, amount float64, emergency bool) (map[string]interface{}, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Calculate total deposited
	var totalDeposited float64
	err = tx.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM vault_positions WHERE user_address = $1
	`, userAddress).Scan(&totalDeposited)

	if err != nil || totalDeposited < amount {
		return nil, err
	}

	// Calculate fee
	fee := amount * 0.005 // 0.5% fee
	if emergency {
		fee = amount * 0.02 // 2% emergency fee
	}
	netAmount := amount - fee

	// Reduce positions proportionally
	_, err = tx.Exec(`
		UPDATE vault_positions
		SET amount = amount * (1 - $1::numeric / $2::numeric)
		WHERE user_address = $3 AND amount > 0
	`, amount, totalDeposited, userAddress)
	if err != nil {
		return nil, err
	}

	// Return points to user
	_, err = tx.Exec(`
		UPDATE points
		SET points = points + $1, updated_at = NOW()
		WHERE user_address = $2
	`, netAmount, userAddress)
	if err != nil {
		return nil, err
	}

	// Record in transaction history
	_, err = tx.Exec(`
		INSERT INTO transaction_history (user_address, transaction_type, amount, fee, status)
		VALUES ($1, 'vault_withdraw', $2, $3, 'confirmed')
	`, userAddress, amount, fee)

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":    true,
		"amount":     amount,
		"fee":        fee,
		"net_amount": netAmount,
		"message":    "Withdrawal successful",
	}, nil
}

// GetBalance returns user's vault balance
func (h *VaultHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")

	// Get available points
	var available float64
	err := h.service.db.QueryRow(`
		SELECT COALESCE(points, 0) FROM points WHERE user_address = $1
	`, address).Scan(&available)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Get locked amount
	var locked float64
	err = h.service.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM vault_positions WHERE user_address = $1
	`, address).Scan(&locked)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Get total earned
	var earned float64
	err = h.service.db.QueryRow(`
		SELECT COALESCE(SUM(earned), 0) FROM vault_positions WHERE user_address = $1
	`, address).Scan(&earned)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Get positions
	positions := []map[string]interface{}{}
	rows, err := h.service.db.Query(`
		SELECT protocol, amount, earned, apy
		FROM vault_positions
		WHERE user_address = $1 AND amount > 0
		ORDER BY amount DESC
	`, address)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var protocol string
			var amount, earned, apy float64
			rows.Scan(&protocol, &amount, &earned, &apy)
			positions = append(positions, map[string]interface{}{
				"protocol": protocol,
				"amount":   amount,
				"earned":   earned,
				"apy":      apy,
			})
		}
	}

	c.JSON(200, gin.H{
		"available": available,
		"locked":    locked,
		"earned":    earned,
		"positions": positions,
	})
}

// GetStrategies returns available strategies
func (h *VaultHandler) GetStrategies(c *gin.Context) {
	rows, err := h.service.db.Query(`
		SELECT name, mode, protocol_allocations, min_amount, max_amount
		FROM vault_strategies
		WHERE is_active = true
		ORDER BY name
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	strategies := []map[string]interface{}{}
	for rows.Next() {
		var name, mode string
		var allocations interface{}
		var minAmount, maxAmount sql.NullFloat64
		rows.Scan(&name, &mode, &allocations, &minAmount, &maxAmount)

		strategy := map[string]interface{}{
			"name":        name,
			"mode":        mode,
			"allocations": allocations,
		}
		if minAmount.Valid {
			strategy["min_amount"] = minAmount.Float64
		}
		if maxAmount.Valid {
			strategy["max_amount"] = maxAmount.Float64
		}
		strategies = append(strategies, strategy)
	}

	c.JSON(200, gin.H{"strategies": strategies})
}

// StakeToProtocol handles manual staking
func (h *VaultHandler) StakeToProtocol(c *gin.Context) {
	var req struct {
		UserAddress string  `json:"user_address"`
		Protocol    string  `json:"protocol"`
		Amount      float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Implementation similar to deposit but for specific protocol
	c.JSON(200, gin.H{
		"success": true,
		"message": "Staked to " + req.Protocol,
	})
}

// UnstakeFromProtocol handles manual unstaking
func (h *VaultHandler) UnstakeFromProtocol(c *gin.Context) {
	var req struct {
		UserAddress string  `json:"user_address"`
		Protocol    string  `json:"protocol"`
		Amount      float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Implementation similar to withdraw but for specific protocol
	c.JSON(200, gin.H{
		"success": true,
		"message": "Unstaked from " + req.Protocol,
	})
}

// GetEarnings returns user's earnings history
func (h *VaultHandler) GetEarnings(c *gin.Context) {
	address := c.Param("address")

	rows, err := h.service.db.Query(`
		SELECT protocol, amount, apy, created_at
		FROM vault_earnings
		WHERE user_address = $1
		ORDER BY created_at DESC
		LIMIT 50
	`, address)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	earnings := []map[string]interface{}{}
	for rows.Next() {
		var protocol string
		var amount, apy float64
		var createdAt string
		rows.Scan(&protocol, &amount, &apy, &createdAt)
		earnings = append(earnings, map[string]interface{}{
			"protocol":   protocol,
			"amount":     amount,
			"apy":        apy,
			"created_at": createdAt,
		})
	}

	c.JSON(200, gin.H{"earnings": earnings})
}

// GetPositions returns user's current positions
func (h *VaultHandler) GetPositions(c *gin.Context) {
	address := c.Param("address")

	rows, err := h.service.db.Query(`
		SELECT p.protocol, p.amount, p.earned, p.apy, d.protocol_type, d.risk_level
		FROM vault_positions p
		JOIN defi_protocols d ON p.protocol = d.name
		WHERE p.user_address = $1 AND p.amount > 0
		ORDER BY p.amount DESC
	`, address)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	positions := []map[string]interface{}{}
	for rows.Next() {
		var protocol, protocolType, riskLevel string
		var amount, earned, apy float64
		rows.Scan(&protocol, &amount, &earned, &apy, &protocolType, &riskLevel)
		positions = append(positions, map[string]interface{}{
			"protocol":      protocol,
			"amount":        amount,
			"earned":        earned,
			"apy":           apy,
			"protocol_type": protocolType,
			"risk_level":    riskLevel,
		})
	}

	c.JSON(200, gin.H{"positions": positions})
}

// GetProtocols returns available DeFi protocols
func (h *VaultHandler) GetProtocols(c *gin.Context) {
	rows, err := h.service.db.Query(`
		SELECT name, protocol_type, risk_level, current_apy, tvl
		FROM defi_protocols
		WHERE is_active = true
		ORDER BY tvl DESC
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
		rows.Scan(&name, &protocolType, &riskLevel, &apy, &tvl)
		protocols = append(protocols, map[string]interface{}{
			"name":          name,
			"protocol_type": protocolType,
			"risk_level":    riskLevel,
			"current_apy":   apy,
			"tvl":           tvl,
		})
	}

	c.JSON(200, gin.H{"protocols": protocols})
}

// Helper function to get strategy allocations
func getStrategyAllocations(strategy string) map[string]float64 {
	switch strategy {
	case "conservative":
		return map[string]float64{
			"Aave V3":      40,
			"Compound V3":  30,
			"Curve":        20,
			"Lido":         10,
		}
	case "balanced":
		return map[string]float64{
			"Aave V3":       25,
			"Compound V3":   20,
			"Uniswap V3":    20,
			"Yearn Finance": 20,
			"Curve":         15,
		}
	case "aggressive":
		return map[string]float64{
			"Uniswap V3":    25,
			"GMX":           20,
			"Yearn Finance": 25,
			"Aave V3":       15,
			"Rocket Pool":   15,
		}
	default:
		return getStrategyAllocations("balanced")
	}
}

// Helper function to get int from environment variable with default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Helper function to get duration from environment variable with default
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}