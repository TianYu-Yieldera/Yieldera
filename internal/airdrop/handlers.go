package airdrop

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"loyalty-points-system/services/api/middleware"
)

// AdminAuthMiddleware validates that the requester is an admin
func AdminAuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First use the standard wallet auth middleware
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// For admin operations, we expect the token to be the wallet address
		// (simplified for now, can be enhanced with JWT later)
		address := strings.ToLower(parts[1])

		// Check if address is in admin whitelist
		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM admin_whitelist WHERE LOWER(address) = $1)`, address).Scan(&exists)
		if err != nil {
			log.Printf("Admin check error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		// Set admin address in context
		c.Set("adminAddress", address)
		c.Next()
	}
}

// CreateCampaignHandler creates a new airdrop campaign
func CreateCampaignHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateCampaignRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate dates
		if req.EndTime.Before(req.StartTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time must be after start_time"})
			return
		}

		// Get admin address from context
		adminAddr, _ := c.Get("adminAddress")
		createdBy := adminAddr.(string)

		// Set default asset type
		if req.AssetType == "" {
			req.AssetType = AssetTypePoints
		}

		// Insert campaign
		var campaignID int
		err := db.QueryRow(`
			INSERT INTO airdrop_campaigns
			(name, description, asset_type, status, start_time, end_time, total_budget, created_by, is_demo)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, req.Name, req.Description, req.AssetType, StatusDraft, req.StartTime, req.EndTime, req.TotalBudget, createdBy, req.IsDemo).Scan(&campaignID)

		if err != nil {
			log.Printf("Create campaign error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":      campaignID,
			"message": "Campaign created successfully",
		})
	}
}

// UpdateCampaignHandler updates an existing campaign
func UpdateCampaignHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")
		var req UpdateCampaignRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if campaign exists and is in draft status
		var status string
		err := db.QueryRow(`SELECT status FROM airdrop_campaigns WHERE id = $1`, campaignID).Scan(&status)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if status != StatusDraft {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can only update campaigns in draft status"})
			return
		}

		// Build update query dynamically
		updates := []string{}
		args := []interface{}{}
		argCount := 1

		if req.Name != "" {
			updates = append(updates, fmt.Sprintf("name = $%d", argCount))
			args = append(args, req.Name)
			argCount++
		}
		if req.Description != "" {
			updates = append(updates, fmt.Sprintf("description = $%d", argCount))
			args = append(args, req.Description)
			argCount++
		}
		if !req.StartTime.IsZero() {
			updates = append(updates, fmt.Sprintf("start_time = $%d", argCount))
			args = append(args, req.StartTime)
			argCount++
		}
		if !req.EndTime.IsZero() {
			updates = append(updates, fmt.Sprintf("end_time = $%d", argCount))
			args = append(args, req.EndTime)
			argCount++
		}
		if req.TotalBudget != "" {
			updates = append(updates, fmt.Sprintf("total_budget = $%d", argCount))
			args = append(args, req.TotalBudget)
			argCount++
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
		args = append(args, time.Now())
		argCount++

		args = append(args, campaignID)

		query := fmt.Sprintf("UPDATE airdrop_campaigns SET %s WHERE id = $%d", strings.Join(updates, ", "), argCount)
		_, err = db.Exec(query, args...)
		if err != nil {
			log.Printf("Update campaign error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Campaign updated successfully"})
	}
}

// ImportAllocationsHandler imports allocations from CSV
func ImportAllocationsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")

		// Check if campaign exists and is in draft status
		var status string
		err := db.QueryRow(`SELECT status FROM airdrop_campaigns WHERE id = $1`, campaignID).Scan(&status)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if status != StatusDraft {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can only import allocations for campaigns in draft status"})
			return
		}

		// Parse multipart form
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		defer file.Close()

		// Parse CSV
		reader := csv.NewReader(file)

		// Skip header
		header, err := reader.Read()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV file"})
			return
		}

		// Validate header
		if len(header) < 2 || strings.ToLower(header[0]) != "address" || strings.ToLower(header[1]) != "amount" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV must have 'address' and 'amount' columns"})
			return
		}

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback()

		// Delete existing allocations
		_, err = tx.Exec(`DELETE FROM airdrop_allocations WHERE campaign_id = $1`, campaignID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing allocations"})
			return
		}

		// Insert new allocations
		importedCount := 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("CSV parse error: %v", err)})
				return
			}

			if len(record) < 2 {
				continue
			}

			address := strings.TrimSpace(strings.ToLower(record[0]))
			amount := strings.TrimSpace(record[1])

			// Validate address
			if len(address) != 42 || !strings.HasPrefix(address, "0x") {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid address: %s", address)})
				return
			}

			// Insert allocation
			_, err = tx.Exec(`
				INSERT INTO airdrop_allocations (campaign_id, user_address, amount)
				VALUES ($1, $2, $3)
				ON CONFLICT (campaign_id, user_address) DO UPDATE SET amount = $3
			`, campaignID, address, amount)

			if err != nil {
				log.Printf("Insert allocation error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import allocations"})
				return
			}

			importedCount++
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Allocations imported successfully",
			"count":   importedCount,
		})
	}
}

// ActivateCampaignHandler activates a campaign
func ActivateCampaignHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")

		// Check current status
		var status string
		var startTime time.Time
		err := db.QueryRow(`SELECT status, start_time FROM airdrop_campaigns WHERE id = $1`, campaignID).Scan(&status, &startTime)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if status != StatusDraft && status != StatusScheduled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can only activate campaigns in draft or scheduled status"})
			return
		}

		// Determine new status based on start time
		newStatus := StatusActive
		if time.Now().Before(startTime) {
			newStatus = StatusScheduled
		}

		// Update status
		_, err = db.Exec(`UPDATE airdrop_campaigns SET status = $1, updated_at = $2 WHERE id = $3`, newStatus, time.Now(), campaignID)
		if err != nil {
			log.Printf("Activate campaign error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate campaign"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Campaign activated successfully",
			"status":  newStatus,
		})
	}
}

// CloseCampaignHandler closes a campaign
func CloseCampaignHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")

		// Update status to closed
		_, err := db.Exec(`UPDATE airdrop_campaigns SET status = $1, updated_at = $2 WHERE id = $3`, StatusClosed, time.Now(), campaignID)
		if err != nil {
			log.Printf("Close campaign error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close campaign"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Campaign closed successfully"})
	}
}

// GetCampaignsHandler returns list of campaigns
func GetCampaignsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		limit := c.DefaultQuery("limit", "20")
		offset := c.DefaultQuery("offset", "0")

		query := `SELECT id, name, description, asset_type, status, start_time, end_time,
		          total_budget, claimed_amount, participant_count, is_demo, created_by, created_at, updated_at
		          FROM airdrop_campaigns WHERE 1=1`
		args := []interface{}{}
		argCount := 1

		if status != "" {
			query += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, status)
			argCount++
		}

		query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
		args = append(args, limit, offset)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Printf("Get campaigns error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		campaigns := []Campaign{}
		for rows.Next() {
			var campaign Campaign
			err := rows.Scan(
				&campaign.ID, &campaign.Name, &campaign.Description, &campaign.AssetType,
				&campaign.Status, &campaign.StartTime, &campaign.EndTime, &campaign.TotalBudget,
				&campaign.ClaimedAmount, &campaign.ParticipantCount, &campaign.IsDemo,
				&campaign.CreatedBy, &campaign.CreatedAt, &campaign.UpdatedAt,
			)
			if err != nil {
				log.Printf("Scan campaign error: %v", err)
				continue
			}
			campaigns = append(campaigns, campaign)
		}

		c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
	}
}

// GetCampaignHandler returns a single campaign
func GetCampaignHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")

		var campaign Campaign
		err := db.QueryRow(`
			SELECT id, name, description, asset_type, status, start_time, end_time,
			       total_budget, claimed_amount, participant_count, is_demo, created_by, created_at, updated_at
			FROM airdrop_campaigns WHERE id = $1
		`, campaignID).Scan(
			&campaign.ID, &campaign.Name, &campaign.Description, &campaign.AssetType,
			&campaign.Status, &campaign.StartTime, &campaign.EndTime, &campaign.TotalBudget,
			&campaign.ClaimedAmount, &campaign.ParticipantCount, &campaign.IsDemo,
			&campaign.CreatedBy, &campaign.CreatedAt, &campaign.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			log.Printf("Get campaign error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, campaign)
	}
}

// CheckEligibilityHandler checks if a user is eligible for an airdrop
func CheckEligibilityHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")
		address := strings.ToLower(c.Query("address"))

		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
			return
		}

		// Check if campaign is active or claimable
		var status string
		err := db.QueryRow(`SELECT status FROM airdrop_campaigns WHERE id = $1`, campaignID).Scan(&status)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if status != StatusActive && status != StatusClaimable {
			c.JSON(http.StatusOK, EligibilityResponse{
				Eligible: false,
				Claimed:  false,
				Reason:   "Campaign is not active",
			})
			return
		}

		// Check if already claimed
		var claimed bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM airdrop_claims WHERE campaign_id = $1 AND user_address = $2)`, campaignID, address).Scan(&claimed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if claimed {
			c.JSON(http.StatusOK, EligibilityResponse{
				Eligible: false,
				Claimed:  true,
				Reason:   "Already claimed",
			})
			return
		}

		// Check allocation
		var amount string
		err = db.QueryRow(`SELECT amount FROM airdrop_allocations WHERE campaign_id = $1 AND user_address = $2`, campaignID, address).Scan(&amount)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, EligibilityResponse{
				Eligible: false,
				Claimed:  false,
				Reason:   "Not in whitelist",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, EligibilityResponse{
			Eligible: true,
			Amount:   amount,
			Claimed:  false,
		})
	}
}

// ClaimAirdropHandler processes an airdrop claim
func ClaimAirdropHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")
		var req ClaimRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		address := strings.ToLower(req.Address)

		// Verify signature
		message := fmt.Sprintf("Claim airdrop from campaign %s with nonce %s", campaignID, req.Nonce)
		if err := middleware.VerifySignature(address, message, req.Signature); err != nil {
			log.Printf("Signature verification failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback()

		// Check campaign status
		var status, assetType string
		err = tx.QueryRow(`SELECT status, asset_type FROM airdrop_campaigns WHERE id = $1`, campaignID).Scan(&status, &assetType)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if status != StatusActive && status != StatusClaimable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign is not active"})
			return
		}

		// Check if already claimed
		var claimed bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM airdrop_claims WHERE campaign_id = $1 AND user_address = $2)`, campaignID, address).Scan(&claimed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if claimed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already claimed"})
			return
		}

		// Get allocation amount
		var amount string
		err = tx.QueryRow(`SELECT amount FROM airdrop_allocations WHERE campaign_id = $1 AND user_address = $2`, campaignID, address).Scan(&amount)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not eligible for this airdrop"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Insert claim record
		_, err = tx.Exec(`
			INSERT INTO airdrop_claims (campaign_id, user_address, amount, nonce, signature)
			VALUES ($1, $2, $3, $4, $5)
		`, campaignID, address, amount, req.Nonce, req.Signature)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Already claimed"})
				return
			}
			log.Printf("Insert claim error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record claim"})
			return
		}

		// For points, update user points
		if assetType == AssetTypePoints {
			// Ensure user exists in points table
			_, err = tx.Exec(`
				INSERT INTO points (user_address, points)
				VALUES ($1, $2)
				ON CONFLICT (user_address) DO UPDATE SET points = points.points + $2
			`, address, amount)
			if err != nil {
				log.Printf("Update points error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to award points"})
				return
			}

			// Insert points event
			_, err = tx.Exec(`
				INSERT INTO points_events (user_address, points_delta, reason)
				VALUES ($1, $2, $3)
			`, address, amount, fmt.Sprintf("Airdrop claim from campaign #%s", campaignID))
			if err != nil {
				log.Printf("Insert points event error: %v", err)
				// Non-critical, continue
			}
		}

		// Update campaign stats
		_, err = tx.Exec(`
			UPDATE airdrop_campaigns
			SET claimed_amount = claimed_amount + $1,
			    participant_count = participant_count + 1,
			    updated_at = $2
			WHERE id = $3
		`, amount, time.Now(), campaignID)
		if err != nil {
			log.Printf("Update campaign stats error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
			return
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Airdrop claimed successfully",
			"amount":  amount,
		})
	}
}

// GetCampaignStatsHandler returns campaign statistics
func GetCampaignStatsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignID := c.Param("id")

		var stats struct {
			TotalBudget      string `json:"total_budget"`
			ClaimedAmount    string `json:"claimed_amount"`
			RemainingAmount  string `json:"remaining_amount"`
			ParticipantCount int    `json:"participant_count"`
			TotalAllocations int    `json:"total_allocations"`
			ClaimRate        string `json:"claim_rate"`
		}

		// Get campaign stats
		err := db.QueryRow(`
			SELECT total_budget, claimed_amount, participant_count
			FROM airdrop_campaigns WHERE id = $1
		`, campaignID).Scan(&stats.TotalBudget, &stats.ClaimedAmount, &stats.ParticipantCount)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Get total allocations
		err = db.QueryRow(`SELECT COUNT(*) FROM airdrop_allocations WHERE campaign_id = $1`, campaignID).Scan(&stats.TotalAllocations)
		if err != nil {
			stats.TotalAllocations = 0
		}

		// Calculate remaining and claim rate
		budgetFloat, _ := strconv.ParseFloat(stats.TotalBudget, 64)
		claimedFloat, _ := strconv.ParseFloat(stats.ClaimedAmount, 64)
		stats.RemainingAmount = fmt.Sprintf("%.2f", budgetFloat-claimedFloat)

		if stats.TotalAllocations > 0 {
			rate := float64(stats.ParticipantCount) / float64(stats.TotalAllocations) * 100
			stats.ClaimRate = fmt.Sprintf("%.2f%%", rate)
		} else {
			stats.ClaimRate = "0%"
		}

		c.JSON(http.StatusOK, stats)
	}
}
