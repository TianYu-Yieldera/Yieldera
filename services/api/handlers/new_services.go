package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// Yield Calculation Service Routes
// ============================================================================

// GetTreasuryRates returns current treasury yield rates
func GetTreasuryRates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT bond_type, annual_yield, effective_date, source
			FROM treasury_yield_rates
			WHERE effective_date = CURRENT_DATE
			ORDER BY bond_type
		`)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var rates []map[string]interface{}
		for rows.Next() {
			var bondType, source string
			var annualYield float64
			var effectiveDate string
			if err := rows.Scan(&bondType, &annualYield, &effectiveDate, &source); err == nil {
				rates = append(rates, map[string]interface{}{
					"bond_type":      bondType,
					"annual_yield":   annualYield,
					"effective_date": effectiveDate,
					"source":         source,
				})
			}
		}

		c.JSON(200, gin.H{"rates": rates})
	}
}

// GetUserYieldHistory returns user's yield history
func GetUserYieldHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		bondType := c.Query("bond_type")
		limit := c.DefaultQuery("limit", "30")

		query := `
			SELECT time, bond_type, yield_amount_usd, cumulative_yield, compounded
			FROM treasury_yield_accruals
			WHERE user_id = $1
		`
		args := []interface{}{userId}

		if bondType != "" {
			query += ` AND bond_type = $2`
			args = append(args, bondType)
			query += ` ORDER BY time DESC LIMIT $3`
			args = append(args, limit)
		} else {
			query += ` ORDER BY time DESC LIMIT $2`
			args = append(args, limit)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var history []map[string]interface{}
		for rows.Next() {
			var time, bondType string
			var yieldAmount, cumulativeYield float64
			var compounded bool
			if err := rows.Scan(&time, &bondType, &yieldAmount, &cumulativeYield, &compounded); err == nil {
				history = append(history, map[string]interface{}{
					"time":             time,
					"bond_type":        bondType,
					"yield_amount":     yieldAmount,
					"cumulative_yield": cumulativeYield,
					"compounded":       compounded,
				})
			}
		}

		c.JSON(200, gin.H{"history": history})
	}
}

// GetUserTotalYield returns user's total accumulated yield
func GetUserTotalYield(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		rows, err := db.Query(`
			SELECT bond_type, total_yield_earned
			FROM treasury_holdings
			WHERE user_id = $1
		`, userId)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		totalYield := 0.0
		yieldByType := make(map[string]float64)

		for rows.Next() {
			var bondType string
			var yield float64
			if err := rows.Scan(&bondType, &yield); err == nil {
				yieldByType[bondType] = yield
				totalYield += yield
			}
		}

		c.JSON(200, gin.H{
			"total_yield":     totalYield,
			"yield_by_type":   yieldByType,
		})
	}
}

// ProjectYield calculates projected yield for investment
func ProjectYield(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			BondType     string  `json:"bond_type" binding:"required"`
			PrincipalUsd float64 `json:"principal_usd" binding:"required"`
			DurationDays int     `json:"duration_days" binding:"required"`
			Compounding  bool    `json:"compounding"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Get current yield rate
		var annualYield float64
		err := db.QueryRow(`
			SELECT annual_yield
			FROM treasury_yield_rates
			WHERE bond_type = $1 AND effective_date = CURRENT_DATE
		`, req.BondType).Scan(&annualYield)

		if err != nil {
			c.JSON(404, gin.H{"error": "Bond type not found"})
			return
		}

		// Calculate projected yield (simplified version)
		years := float64(req.DurationDays) / 365.0
		var totalYield, finalValue float64

		if req.Compounding {
			// Compound interest: A = P(1 + r/n)^(nt), daily compounding
			finalValue = req.PrincipalUsd * pow(1 + annualYield/365, 365*years)
			totalYield = finalValue - req.PrincipalUsd
		} else {
			// Simple interest: I = P * r * t
			totalYield = req.PrincipalUsd * annualYield * years
			finalValue = req.PrincipalUsd + totalYield
		}

		dailyYield := (annualYield / 365) * req.PrincipalUsd
		effectiveAPY := pow(finalValue/req.PrincipalUsd, 365.0/float64(req.DurationDays)) - 1

		c.JSON(200, gin.H{
			"total_yield":     totalYield,
			"effective_apy":   effectiveAPY,
			"daily_yield":     dailyYield,
			"projected_value": finalValue,
		})
	}
}

// ============================================================================
// Notification Service Routes
// ============================================================================

// GetUserNotifications returns user's notifications
func GetUserNotifications(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		limit := c.DefaultQuery("limit", "20")
		unreadOnly := c.Query("unread") == "true"

		query := `
			SELECT time, type, priority, title, message, data, channels, read_at
			FROM notifications
			WHERE user_id = $1
		`
		args := []interface{}{userId}

		if unreadOnly {
			query += ` AND read_at IS NULL`
		}

		query += ` ORDER BY time DESC LIMIT $2`
		args = append(args, limit)

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var notifications []map[string]interface{}
		for rows.Next() {
			var time, notifType, priority, title, message string
			var data, channels sql.NullString
			var readAt sql.NullTime

			if err := rows.Scan(&time, &notifType, &priority, &title, &message, &data, &channels, &readAt); err == nil {
				notif := map[string]interface{}{
					"time":     time,
					"type":     notifType,
					"priority": priority,
					"title":    title,
					"message":  message,
				}
				if data.Valid {
					notif["data"] = data.String
				}
				if channels.Valid {
					notif["channels"] = channels.String
				}
				if readAt.Valid {
					notif["read_at"] = readAt.Time
				}
				notifications = append(notifications, notif)
			}
		}

		c.JSON(200, gin.H{"notifications": notifications})
	}
}

// MarkNotificationRead marks a notification as read
func MarkNotificationRead(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		timestamp := c.Param("timestamp")

		_, err := db.Exec(`
			UPDATE notifications
			SET read_at = NOW()
			WHERE user_id = $1 AND time = $2
		`, userId, timestamp)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	}
}

// GetNotificationPreferences returns user's notification preferences
func GetNotificationPreferences(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var channels, enabledTypes sql.NullString
		var minPriority, frequency string
		var quietStart, quietEnd sql.NullInt64

		err := db.QueryRow(`
			SELECT channels, min_priority, quiet_hours_start, quiet_hours_end, enabled_types, frequency
			FROM notification_preferences
			WHERE user_id = $1
		`, userId).Scan(&channels, &minPriority, &quietStart, &quietEnd, &enabledTypes, &frequency)

		if err == sql.ErrNoRows {
			// Return defaults
			c.JSON(200, gin.H{
				"channels":      []string{"email"},
				"min_priority":  "medium",
				"enabled_types": []string{"liquidation_warning", "high_risk_position", "daily_yield_report"},
				"frequency":     "realtime",
			})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		prefs := map[string]interface{}{
			"min_priority": minPriority,
			"frequency":    frequency,
		}
		if channels.Valid {
			prefs["channels"] = channels.String
		}
		if enabledTypes.Valid {
			prefs["enabled_types"] = enabledTypes.String
		}
		if quietStart.Valid {
			prefs["quiet_hours_start"] = quietStart.Int64
		}
		if quietEnd.Valid {
			prefs["quiet_hours_end"] = quietEnd.Int64
		}

		c.JSON(200, gin.H{"preferences": prefs})
	}
}

// UpdateNotificationPreferences updates user's notification preferences
func UpdateNotificationPreferences(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var req struct {
			Channels      []string `json:"channels"`
			MinPriority   string   `json:"min_priority"`
			QuietStart    *int     `json:"quiet_hours_start"`
			QuietEnd      *int     `json:"quiet_hours_end"`
			EnabledTypes  []string `json:"enabled_types"`
			Frequency     string   `json:"frequency"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Convert arrays to PostgreSQL array format
		channelsStr := "{" + joinStrings(req.Channels, ",") + "}"
		typesStr := "{" + joinStrings(req.EnabledTypes, ",") + "}"

		_, err := db.Exec(`
			INSERT INTO notification_preferences (user_id, channels, min_priority, quiet_hours_start, quiet_hours_end, enabled_types, frequency)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (user_id) DO UPDATE SET
				channels = EXCLUDED.channels,
				min_priority = EXCLUDED.min_priority,
				quiet_hours_start = EXCLUDED.quiet_hours_start,
				quiet_hours_end = EXCLUDED.quiet_hours_end,
				enabled_types = EXCLUDED.enabled_types,
				frequency = EXCLUDED.frequency,
				updated_at = NOW()
		`, userId, channelsStr, req.MinPriority, req.QuietStart, req.QuietEnd, typesStr, req.Frequency)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	}
}

// ============================================================================
// Auto Hedge Executor Routes
// ============================================================================

// GetHedgeHistory returns user's hedge execution history
func GetHedgeHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		limit := c.DefaultQuery("limit", "20")

		rows, err := db.Query(`
			SELECT time, strategy, status, tx_hash, gas_used, resulting_health_factor, error
			FROM hedge_executions
			WHERE user_id = $1
			ORDER BY time DESC
			LIMIT $2
		`, userId, limit)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var executions []map[string]interface{}
		for rows.Next() {
			var time, strategy, status string
			var txHash, gasUsed, errorMsg sql.NullString
			var resultingHF sql.NullFloat64

			if err := rows.Scan(&time, &strategy, &status, &txHash, &gasUsed, &resultingHF, &errorMsg); err == nil {
				exec := map[string]interface{}{
					"time":     time,
					"strategy": strategy,
					"status":   status,
				}
				if txHash.Valid {
					exec["tx_hash"] = txHash.String
				}
				if gasUsed.Valid {
					exec["gas_used"] = gasUsed.String
				}
				if resultingHF.Valid {
					exec["resulting_health_factor"] = resultingHF.Float64
				}
				if errorMsg.Valid {
					exec["error"] = errorMsg.String
				}
				executions = append(executions, exec)
			}
		}

		c.JSON(200, gin.H{"executions": executions})
	}
}

// GetHedgeStats returns hedge execution statistics
func GetHedgeStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		days := c.DefaultQuery("days", "7")
		daysInt, _ := strconv.Atoi(days)

		rows, err := db.Query(`
			SELECT day, total_executions, successful, failed, avg_resulting_hf
			FROM daily_hedge_stats
			WHERE day >= NOW() - INTERVAL '1 day' * $1
			ORDER BY day DESC
		`, daysInt)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var stats []map[string]interface{}
		for rows.Next() {
			var day string
			var total, successful, failed int
			var avgHF sql.NullFloat64

			if err := rows.Scan(&day, &total, &successful, &failed, &avgHF); err == nil {
				stat := map[string]interface{}{
					"day":         day,
					"total":       total,
					"successful":  successful,
					"failed":      failed,
					"success_rate": float64(successful) / float64(total),
				}
				if avgHF.Valid {
					stat["avg_health_factor"] = avgHF.Float64
				}
				stats = append(stats, stat)
			}
		}

		c.JSON(200, gin.H{"stats": stats})
	}
}

// GetUserSettings returns user's auto-hedge settings
func GetUserSettings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var autoHedgeEnabled bool
		var maxAmount, minHF, targetHF float64

		err := db.QueryRow(`
			SELECT auto_hedge_enabled, max_hedge_amount, min_health_factor, target_health_factor
			FROM user_settings
			WHERE user_id = $1
		`, userId).Scan(&autoHedgeEnabled, &maxAmount, &minHF, &targetHF)

		if err == sql.ErrNoRows {
			// Return defaults
			c.JSON(200, gin.H{
				"auto_hedge_enabled":   false,
				"max_hedge_amount":     10000,
				"min_health_factor":    1.5,
				"target_health_factor": 2.0,
			})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"auto_hedge_enabled":   autoHedgeEnabled,
			"max_hedge_amount":     maxAmount,
			"min_health_factor":    minHF,
			"target_health_factor": targetHF,
		})
	}
}

// UpdateUserSettings updates user's auto-hedge settings
func UpdateUserSettings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var req struct {
			AutoHedgeEnabled  bool    `json:"auto_hedge_enabled"`
			MaxHedgeAmount    float64 `json:"max_hedge_amount"`
			MinHealthFactor   float64 `json:"min_health_factor"`
			TargetHealthFactor float64 `json:"target_health_factor"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(`
			INSERT INTO user_settings (user_id, auto_hedge_enabled, max_hedge_amount, min_health_factor, target_health_factor)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (user_id) DO UPDATE SET
				auto_hedge_enabled = EXCLUDED.auto_hedge_enabled,
				max_hedge_amount = EXCLUDED.max_hedge_amount,
				min_health_factor = EXCLUDED.min_health_factor,
				target_health_factor = EXCLUDED.target_health_factor,
				updated_at = NOW()
		`, userId, req.AutoHedgeEnabled, req.MaxHedgeAmount, req.MinHealthFactor, req.TargetHealthFactor)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	}
}

// ============================================================================
// Yield Distribution Routes
// ============================================================================

// GetDistributionStats returns yield distribution statistics
func GetDistributionStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		days := c.DefaultQuery("days", "30")

		rows, err := db.Query(`
			SELECT date, total_distributed, recipient_count, success_rate, average_yield, total_gas_cost
			FROM distribution_reports
			WHERE date >= CURRENT_DATE - INTERVAL '1 day' * $1
			ORDER BY date DESC
		`, days)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var stats []map[string]interface{}
		for rows.Next() {
			var date string
			var totalDistributed, avgYield, gasCost float64
			var recipientCount int
			var successRate float64

			if err := rows.Scan(&date, &totalDistributed, &recipientCount, &successRate, &avgYield, &gasCost); err == nil {
				stats = append(stats, map[string]interface{}{
					"date":              date,
					"total_distributed": totalDistributed,
					"recipient_count":   recipientCount,
					"success_rate":      successRate,
					"average_yield":     avgYield,
					"total_gas_cost":    gasCost,
				})
			}
		}

		c.JSON(200, gin.H{"stats": stats})
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
