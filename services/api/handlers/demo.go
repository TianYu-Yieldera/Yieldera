package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const demoGrantPoints int64 = 10000

// DemoSummary consolidates demo-mode data for a user
type DemoSummary struct {
	Address       string             `json:"address"`
	IsDemo        bool               `json:"is_demo"`
	DemoActive    bool               `json:"demo_active"`
	DemoExpiresAt *time.Time         `json:"demo_expires_at,omitempty"`
	Points        string             `json:"points"`
	TokenBalance  string             `json:"token_balance"`
	Stablecoin    DemoStablecoin     `json:"stablecoin"`
	DefiPositions []DemoDefiPosition `json:"defi_positions"`
}

// DemoStablecoin summarises stablecoin reserves in demo mode
type DemoStablecoin struct {
	Collateral string `json:"collateral"`
	Debt       string `json:"debt"`
}

// DemoDefiPosition represents a simulated DeFi position tied to demo data
type DemoDefiPosition struct {
	PoolID       string     `json:"pool_id"`
	Deposited    string     `json:"deposited"`
	Earned       string     `json:"earned"`
	PointsEarned string     `json:"points_earned"`
	LastUpdated  *time.Time `json:"last_updated,omitempty"`
}

// CreateDemoUser creates or updates a demo user and allocates the default points grant
func CreateDemoUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WalletAddress string `json:"wallet_address"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if err := validateWalletAddress(req.WalletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		now := time.Now().UTC()
		expires := now.Add(24 * time.Hour)

		tx, err := db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "begin transaction failed"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		var userID int64
		err = tx.QueryRowContext(ctx, `SELECT id FROM users WHERE address = $1`, req.WalletAddress).Scan(&userID)
		switch {
		case err == sql.ErrNoRows:
			_, err = tx.ExecContext(ctx, `INSERT INTO users (address, is_demo, demo_expires_at, created_at)
				VALUES ($1, TRUE, $2, $3)`, req.WalletAddress, expires, now)
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users"})
			return
		default:
			_, err = tx.ExecContext(ctx, `UPDATE users
				SET is_demo = TRUE,
				    demo_expires_at = $1
				WHERE address = $2`, expires, req.WalletAddress)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist demo user"})
			return
		}

		if err := upsertDemoPoints(ctx, tx, req.WalletAddress, demoGrantPoints, "demo_grant"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := ensureDemoBalance(ctx, tx, req.WalletAddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "commit transaction failed"})
			return
		}

		summary, err := loadDemoSummary(ctx, db, req.WalletAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load demo summary"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "demo user created successfully",
			"demo_user": gin.H{
				"address":         req.WalletAddress,
				"points":          demoGrantPoints,
				"is_demo":         true,
				"demo_expires_at": expires,
			},
			"summary": summary,
		})
	}
}

// GetDemoStatus returns whether an address is currently in demo mode together with summary info
func GetDemoStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletAddress := c.Query("address")
		if walletAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wallet address is required"})
			return
		}
		if err := validateWalletAddress(walletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		summary, err := loadDemoSummary(c.Request.Context(), db, walletAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load demo status"})
			return
		}

		resp := gin.H{
			"is_demo":     summary.IsDemo,
			"demo_active": summary.DemoActive,
			"points":      summary.Points,
			"summary":     summary,
		}
		if summary.DemoExpiresAt != nil {
			resp["demo_expires_at"] = summary.DemoExpiresAt
		}

		c.JSON(http.StatusOK, resp)
	}
}

// ResetDemoUser resets demo allocations and extends the window by 24h
func ResetDemoUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WalletAddress string `json:"wallet_address"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if err := validateWalletAddress(req.WalletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		now := time.Now().UTC()
		expires := now.Add(24 * time.Hour)

		tx, err := db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "begin transaction failed"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		var isDemo bool
		err = tx.QueryRowContext(ctx, `SELECT is_demo FROM users WHERE address = $1`, req.WalletAddress).Scan(&isDemo)
		switch {
		case err == sql.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
			return
		case !isDemo:
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is not in demo mode"})
			return
		}

		if _, err = tx.ExecContext(ctx, `UPDATE users
			SET demo_expires_at = $1
			WHERE address = $2`, expires, req.WalletAddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset user"})
			return
		}

		if err := upsertDemoPointsExact(ctx, tx, req.WalletAddress, demoGrantPoints, "demo_reset"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := ensureDemoBalance(ctx, tx, req.WalletAddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "commit transaction failed"})
			return
		}

		summary, err := loadDemoSummary(ctx, db, req.WalletAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load demo summary"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"message":         "demo user reset successfully",
			"points":          summary.Points,
			"demo_expires_at": summary.DemoExpiresAt,
			"summary":         summary,
		})
	}
}

// ExitDemoMode removes demo mode flags from all related tables
func ExitDemoMode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WalletAddress string `json:"wallet_address"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if err := validateWalletAddress(req.WalletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()

		tx, err := db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "begin transaction failed"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		if _, err = tx.ExecContext(ctx, `UPDATE users
			SET is_demo = FALSE,
			    demo_expires_at = NULL
			WHERE address = $1`, req.WalletAddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		// Core tables (must exist)
		coreStatements := []string{
			`UPDATE points SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE balances SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE points_events SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE balance_events SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE badges SET is_demo = FALSE WHERE user_address = $1`,
		}
		for _, stmt := range coreStatements {
			if _, err := tx.ExecContext(ctx, stmt, req.WalletAddress); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear demo flags"})
				return
			}
		}

		// Optional tables (may not exist, so ignore errors)
		optionalStatements := []string{
			`UPDATE user_defi_positions SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE defi_transactions SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE stablecoin_positions SET is_demo = FALSE WHERE user_address = $1`,
			`UPDATE stablecoin_transactions SET is_demo = FALSE WHERE user_address = $1`,
		}
		for _, stmt := range optionalStatements {
			_, _ = tx.ExecContext(ctx, stmt, req.WalletAddress) // Ignore errors for optional tables
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "commit transaction failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "successfully exited demo mode",
		})
	}
}

// GetDemoSummary provides a single endpoint to retrieve consolidated demo stats
func GetDemoSummary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletAddress := c.Query("address")
		if walletAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wallet address is required"})
			return
		}
		if err := validateWalletAddress(walletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		summary, err := loadDemoSummary(c.Request.Context(), db, walletAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load demo summary"})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

func validateWalletAddress(addr string) error {
	if len(addr) != 42 || !strings.HasPrefix(strings.ToLower(addr), "0x") {
		return fmt.Errorf("invalid wallet address format")
	}
	for _, ch := range addr[2:] {
		if !strings.ContainsRune("0123456789abcdefABCDEF", ch) {
			return fmt.Errorf("invalid wallet address format")
		}
	}
	return nil
}

func upsertDemoPoints(ctx context.Context, tx *sql.Tx, address string, delta int64, reason string) error {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO points (user_address, points, is_demo, updated_at)
		VALUES ($1, $2, TRUE, NOW())
		ON CONFLICT (user_address)
		DO UPDATE SET
			points = points.points + EXCLUDED.points,
			is_demo = TRUE,
			updated_at = NOW()
	`, address, delta); err != nil {
		return err
	}
	return insertPointsEvent(ctx, tx, address, delta, reason)
}

func upsertDemoPointsExact(ctx context.Context, tx *sql.Tx, address string, value int64, reason string) error {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO points (user_address, points, is_demo, updated_at)
		VALUES ($1, $2, TRUE, NOW())
		ON CONFLICT (user_address)
		DO UPDATE SET
			points = EXCLUDED.points,
			is_demo = TRUE,
			updated_at = NOW()
	`, address, value); err != nil {
		return err
	}
	return insertPointsEvent(ctx, tx, address, value, reason)
}

func ensureDemoBalance(ctx context.Context, tx *sql.Tx, address string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO balances (user_address, balance, is_demo, updated_at)
		VALUES ($1, 0, TRUE, NOW())
		ON CONFLICT (user_address)
		DO UPDATE SET is_demo = TRUE
	`, address)
	return err
}

func insertPointsEvent(ctx context.Context, tx *sql.Tx, address string, delta int64, reason string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO points_events (user_address, points_delta, reason, is_demo)
		VALUES ($1, $2, $3, TRUE)
	`, address, delta, reason)
	return err
}

func loadDemoSummary(ctx context.Context, db *sql.DB, address string) (*DemoSummary, error) {
	summary := &DemoSummary{
		Address:      address,
		Points:       "0",
		TokenBalance: "0",
		Stablecoin:   DemoStablecoin{Collateral: "0", Debt: "0"},
	}

	var isDemo sql.NullBool
	var expires sql.NullTime
	if err := db.QueryRowContext(ctx, `SELECT is_demo, demo_expires_at FROM users WHERE address = $1`, address).Scan(&isDemo, &expires); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == nil {
		summary.IsDemo = isDemo.Bool
		if expires.Valid {
			t := expires.Time.UTC()
			summary.DemoExpiresAt = &t
			summary.DemoActive = summary.IsDemo && time.Now().Before(t)
		}
	}

	if err := db.QueryRowContext(ctx, `SELECT points::text FROM points WHERE user_address = $1`, address).Scan(&summary.Points); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows || summary.Points == "" {
		summary.Points = "0"
	}

	if err := db.QueryRowContext(ctx, `SELECT balance::text FROM balances WHERE user_address = $1`, address).Scan(&summary.TokenBalance); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows || summary.TokenBalance == "" {
		summary.TokenBalance = "0"
	}

	var collateral, debt sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT collateral_amount, debt_amount FROM stablecoin_positions WHERE user_address = $1`, address).Scan(&collateral, &debt); err == nil {
		if collateral.String != "" {
			summary.Stablecoin.Collateral = collateral.String
		}
		if debt.String != "" {
			summary.Stablecoin.Debt = debt.String
		}
	}
	// Ignore errors for optional stablecoin_positions table

	rows, err := db.QueryContext(ctx, `SELECT pool_id, deposited, earned, points_earned, last_updated FROM user_defi_positions WHERE user_address = $1 AND is_demo = TRUE`, address)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var pos DemoDefiPosition
			var lastUpdated sql.NullTime
			if err := rows.Scan(&pos.PoolID, &pos.Deposited, &pos.Earned, &pos.PointsEarned, &lastUpdated); err == nil {
				if lastUpdated.Valid {
					lu := lastUpdated.Time.UTC()
					pos.LastUpdated = &lu
				}
				summary.DefiPositions = append(summary.DefiPositions, pos)
			}
		}
		// Ignore errors for rows
	}
	// Ignore errors for optional user_defi_positions table

	return summary, nil
}
