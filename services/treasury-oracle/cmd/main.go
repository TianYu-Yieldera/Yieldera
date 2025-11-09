package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"loyalty-points-system/internal/config"
	"loyalty-points-system/internal/db"
)

// TreasuryRate represents Treasury rate data from external API
type TreasuryRate struct {
	CUSIP        string  `json:"cusip"`
	SecurityType string  `json:"security_type"` // BILL, NOTE, BOND
	MaturityTerm string  `json:"maturity_term"` // 4W, 13W, 26W, 2Y, 10Y, 30Y
	Rate         float64 `json:"rate"`          // Yield rate
	Price        float64 `json:"price"`         // Price
	IssueDate    string  `json:"issue_date"`
	MaturityDate string  `json:"maturity_date"`
}

// TreasuryAPIResponse represents the API response structure
type TreasuryAPIResponse struct {
	Data      []TreasuryRate `json:"data"`
	Timestamp string         `json:"timestamp"`
}

func main() {
	log.Println("🚀 Starting Treasury Oracle Service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Get update interval from environment (default: 1 hour)
	updateInterval := 1 * time.Hour
	if intervalStr := os.Getenv("ORACLE_UPDATE_INTERVAL"); intervalStr != "" {
		if duration, err := time.ParseDuration(intervalStr); err == nil {
			updateInterval = duration
		}
	}

	log.Printf("📊 Update interval: %v", updateInterval)

	// Start HTTP server for health checks and metrics
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	go func() {
		log.Println("🏥 Metrics server on :8085")
		if err := http.ListenAndServe(":8085", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create ticker for periodic updates
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	// Run initial update
	if err := updateTreasuryPrices(ctx, database); err != nil {
		log.Printf("❌ Initial update failed: %v", err)
	}

	// Main loop
	for {
		select {
		case <-ticker.C:
			if err := updateTreasuryPrices(ctx, database); err != nil {
				log.Printf("❌ Update failed: %v", err)
			}
		case <-sigChan:
			log.Println("🛑 Shutdown signal received, stopping oracle...")
			cancel()
			return
		case <-ctx.Done():
			log.Println("✅ Oracle stopped gracefully")
			return
		}
	}
}

// updateTreasuryPrices fetches latest Treasury rates and updates database
func updateTreasuryPrices(ctx context.Context, database *sql.DB) error {
	log.Println("📥 Fetching Treasury rates...")

	// Get Treasury data from API
	rates, err := fetchTreasuryRates(ctx)
	if err != nil {
		return fmt.Errorf("fetch rates failed: %w", err)
	}

	log.Printf("📊 Fetched %d Treasury rates", len(rates))

	// Update database
	for _, rate := range rates {
		if err := updateAssetPrice(database, &rate); err != nil {
			log.Printf("⚠️  Failed to update %s: %v", rate.CUSIP, err)
			continue
		}
		log.Printf("✅ Updated %s: $%.2f (%.2f%%)", rate.CUSIP, rate.Price, rate.Rate)
	}

	return nil
}

// fetchTreasuryRates fetches Treasury rates from external API
func fetchTreasuryRates(ctx context.Context) ([]TreasuryRate, error) {
	// Option 1: Use Treasury Direct API (Free but limited)
	// Option 2: Use Financial data provider (requires API key)
	// For this implementation, we'll use a mock/fallback with real-world structure

	// Check for custom API URL
	apiURL := os.Getenv("TREASURY_API_URL")
	if apiURL == "" {
		// Use mock data for development/demo
		return getMockTreasuryRates(), nil
	}

	// Fetch from real API
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Add API key if provided
	if apiKey := os.Getenv("TREASURY_API_KEY"); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var apiResp TreasuryAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

// getMockTreasuryRates returns mock data for development
func getMockTreasuryRates() []TreasuryRate {
	now := time.Now()

	return []TreasuryRate{
		{
			CUSIP:        "912796YZ1",
			SecurityType: "BILL",
			MaturityTerm: "13W",
			Rate:         5.25,
			Price:        980.50,
			IssueDate:    now.AddDate(0, -1, 0).Format("2006-01-02"),
			MaturityDate: now.AddDate(0, 2, 0).Format("2006-01-02"),
		},
		{
			CUSIP:        "91282CHX6",
			SecurityType: "NOTE",
			MaturityTerm: "2Y",
			Rate:         4.50,
			Price:        985.75,
			IssueDate:    now.AddDate(0, -6, 0).Format("2006-01-02"),
			MaturityDate: now.AddDate(2, -6, 0).Format("2006-01-02"),
		},
		{
			CUSIP:        "912828YK4",
			SecurityType: "NOTE",
			MaturityTerm: "10Y",
			Rate:         4.25,
			Price:        952.30,
			IssueDate:    now.AddDate(-1, 0, 0).Format("2006-01-02"),
			MaturityDate: now.AddDate(9, 0, 0).Format("2006-01-02"),
		},
		{
			CUSIP:        "912810TT4",
			SecurityType: "BOND",
			MaturityTerm: "30Y",
			Rate:         4.00,
			Price:        920.15,
			IssueDate:    now.AddDate(-2, 0, 0).Format("2006-01-02"),
			MaturityDate: now.AddDate(28, 0, 0).Format("2006-01-02"),
		},
	}
}

// updateAssetPrice updates Treasury asset price and yield in database
func updateAssetPrice(database *sql.DB, rate *TreasuryRate) error {
	// Convert rate percentage to basis points (e.g., 5.25% -> "0.0525")
	yieldStr := fmt.Sprintf("%.4f", rate.Rate/100.0)
	priceStr := fmt.Sprintf("%.2f", rate.Price)

	// Update using the treasury_operations functions
	// First, try to find the asset by CUSIP
	var assetID int64
	err := database.QueryRow(`
		SELECT asset_id FROM treasury_assets
		WHERE cusip = $1 AND status = 'active'
	`, rate.CUSIP).Scan(&assetID)

	if err == sql.ErrNoRows {
		// Asset doesn't exist, create it
		log.Printf("📝 Creating new Treasury asset: %s", rate.CUSIP)
		return createNewAsset(database, rate)
	} else if err != nil {
		return err
	}

	// Update existing asset
	return db.UpdateTreasuryPrice(database, assetID, priceStr, yieldStr, "api")
}

// createNewAsset creates a new Treasury asset in the database
func createNewAsset(database *sql.DB, rate *TreasuryRate) error {
	treasuryType := rate.SecurityType
	if treasuryType == "BILL" {
		treasuryType = "T-BILL"
	} else if treasuryType == "NOTE" {
		treasuryType = "T-NOTE"
	} else if treasuryType == "BOND" {
		treasuryType = "T-BOND"
	}

	// Parse dates
	issueDate, err := time.Parse("2006-01-02", rate.IssueDate)
	if err != nil {
		return fmt.Errorf("invalid issue date: %w", err)
	}

	maturityDate, err := time.Parse("2006-01-02", rate.MaturityDate)
	if err != nil {
		return fmt.Errorf("invalid maturity date: %w", err)
	}

	// Convert rate to basis points for coupon rate
	couponRate := fmt.Sprintf("%.4f", rate.Rate/100.0)
	faceValue := "1000" // Standard face value
	priceStr := fmt.Sprintf("%.2f", rate.Price)
	yieldStr := fmt.Sprintf("%.4f", rate.Rate/100.0)

	query := `
		INSERT INTO treasury_assets (
			treasury_type, maturity_term, cusip, issue_date, maturity_date,
			face_value, coupon_rate, current_price, current_yield,
			tokens_issued, tokens_outstanding, status, last_price_update
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 0, 0, 'active', NOW())
		ON CONFLICT (cusip) DO UPDATE SET
			current_price = $8,
			current_yield = $9,
			last_price_update = NOW(),
			updated_at = NOW()
	`

	_, err = database.Exec(query,
		treasuryType,
		rate.MaturityTerm,
		rate.CUSIP,
		issueDate,
		maturityDate,
		faceValue,
		couponRate,
		priceStr,
		yieldStr,
	)

	return err
}
