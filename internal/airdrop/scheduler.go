package airdrop

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// UpdateCampaignStatuses updates campaign statuses based on time
func UpdateCampaignStatuses(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	// Update scheduled campaigns to active if start_time has passed
	result, err := db.ExecContext(ctx, `
		UPDATE airdrop_campaigns
		SET status = $1, updated_at = $2
		WHERE status = $3
		  AND start_time <= $4
	`, StatusActive, now, StatusScheduled, now)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("📅 Activated %d scheduled campaigns", rowsAffected)
	}

	// Update active campaigns to claimable if end_time has passed
	result, err = db.ExecContext(ctx, `
		UPDATE airdrop_campaigns
		SET status = $1, updated_at = $2
		WHERE status = $3
		  AND end_time <= $4
	`, StatusClaimable, now, StatusActive, now)

	if err != nil {
		return err
	}

	rowsAffected, _ = result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("✅ Changed %d active campaigns to claimable", rowsAffected)
	}

	// Auto-close claimable campaigns after 7 days (configurable)
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)
	result, err = db.ExecContext(ctx, `
		UPDATE airdrop_campaigns
		SET status = $1, updated_at = $2
		WHERE status = $3
		  AND end_time <= $4
	`, StatusClosed, now, StatusClaimable, sevenDaysAgo)

	if err != nil {
		return err
	}

	rowsAffected, _ = result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("🔒 Auto-closed %d claimable campaigns (7 days after end)", rowsAffected)
	}

	return nil
}

// CleanupExpiredAllocations can be used to clean up old data (optional)
func CleanupExpiredAllocations(ctx context.Context, db *sql.DB) error {
	// Clean up allocations for campaigns that are closed for more than 30 days
	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)

	result, err := db.ExecContext(ctx, `
		DELETE FROM airdrop_allocations
		WHERE campaign_id IN (
			SELECT id FROM airdrop_campaigns
			WHERE status = $1 AND updated_at < $2
		)
	`, StatusClosed, thirtyDaysAgo)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("🧹 Cleaned up %d expired allocations", rowsAffected)
	}

	return nil
}
