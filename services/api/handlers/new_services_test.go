package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Test database setup
func setupTestDB(t *testing.T) *sql.DB {
	// Use environment variable or default test database
	dbURL := "postgres://loyalty_user:loyalty_pass@localhost:5432/loyalty_db?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to database: %v", err)
		return nil
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return nil
	}

	return db
}

// TestGetTreasuryRates tests the treasury rates endpoint
func TestGetTreasuryRates(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/yields/rates", GetTreasuryRates(db))

	req, _ := http.NewRequest("GET", "/api/v1/yields/rates", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "rates")
}

// TestProjectYield tests the yield projection endpoint
func TestProjectYield(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/yields/project", ProjectYield(db))

	payload := map[string]interface{}{
		"bond_type":      "TBILL_3M",
		"principal_usd":  1000,
		"duration_days":  90,
		"compounding":    true,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/yields/project", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "total_yield")
	assert.Contains(t, result, "effective_apy")
	assert.Contains(t, result, "projected_value")

	// Validate calculated values
	totalYield := result["total_yield"].(float64)
	assert.Greater(t, totalYield, 0.0, "Total yield should be positive")
}

// TestGetUserTotalYield tests fetching user's total yield
func TestGetUserTotalYield(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/yields/total/:userId", GetUserTotalYield(db))

	testAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"
	req, _ := http.NewRequest("GET", "/api/v1/yields/total/"+testAddress, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "total_yield")
	assert.Contains(t, result, "yield_by_type")
}

// TestGetNotificationPreferences tests getting notification preferences
func TestGetNotificationPreferences(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/notifications/:userId/preferences", GetNotificationPreferences(db))

	testAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"
	req, _ := http.NewRequest("GET", "/api/v1/notifications/"+testAddress+"/preferences", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
}

// TestUpdateNotificationPreferences tests updating notification preferences
func TestUpdateNotificationPreferences(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/v1/notifications/:userId/preferences", UpdateNotificationPreferences(db))

	testAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"
	payload := map[string]interface{}{
		"channels":      []string{"email", "push"},
		"min_priority":  "high",
		"enabled_types": []string{"liquidation_warning"},
		"frequency":     "realtime",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/v1/notifications/"+testAddress+"/preferences", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "success")
	assert.True(t, result["success"].(bool))
}

// TestGetHedgeSettings tests getting hedge settings
func TestGetHedgeSettings(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/hedge/settings/:userId", GetHedgeSettings(db))

	testAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"
	req, _ := http.NewRequest("GET", "/api/v1/hedge/settings/"+testAddress, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "auto_hedge_enabled")
	assert.Contains(t, result, "max_hedge_amount")
}

// TestUpdateHedgeSettings tests updating hedge settings
func TestUpdateHedgeSettings(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/v1/hedge/settings/:userId", UpdateHedgeSettings(db))

	testAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"
	payload := map[string]interface{}{
		"auto_hedge_enabled":   true,
		"max_hedge_amount":     5000,
		"min_health_factor":    1.5,
		"target_health_factor": 2.0,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/v1/hedge/settings/"+testAddress, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "success")
	assert.True(t, result["success"].(bool))
}

// TestGetUserNotifications tests getting user notifications
func TestGetUserNotifications(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/notifications/:userId", GetUserNotifications(db))

	testAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6"
	req, _ := http.NewRequest("GET", "/api/v1/notifications/"+testAddress, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "notifications")
}

// TestGetDistributionStats tests getting distribution statistics
func TestGetDistributionStats(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/distribution/stats", GetDistributionStats(db))

	req, _ := http.NewRequest("GET", "/api/v1/distribution/stats", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Contains(t, result, "stats")
}

// TestProjectYieldValidation tests input validation for yield projection
func TestProjectYieldValidation(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/yields/project", ProjectYield(db))

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Invalid bond type",
			payload: map[string]interface{}{
				"bond_type":     "INVALID",
				"principal_usd": 1000,
				"duration_days": 90,
			},
			expectedStatus: 400,
		},
		{
			name: "Negative principal",
			payload: map[string]interface{}{
				"bond_type":     "TBILL_3M",
				"principal_usd": -1000,
				"duration_days": 90,
			},
			expectedStatus: 400,
		},
		{
			name: "Zero duration",
			payload: map[string]interface{}{
				"bond_type":     "TBILL_3M",
				"principal_usd": 1000,
				"duration_days": 0,
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/yields/project", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

// Benchmark tests
func BenchmarkGetTreasuryRates(b *testing.B) {
	db := setupTestDB(&testing.T{})
	if db == nil {
		b.Skip("Database not available")
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/yields/rates", GetTreasuryRates(db))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/yields/rates", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

func BenchmarkProjectYield(b *testing.B) {
	db := setupTestDB(&testing.T{})
	if db == nil {
		b.Skip("Database not available")
		return
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/yields/project", ProjectYield(db))

	payload := map[string]interface{}{
		"bond_type":     "TBILL_3M",
		"principal_usd": 1000,
		"duration_days": 90,
		"compounding":   true,
	}
	jsonPayload, _ := json.Marshal(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/yields/project", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}
