package airdrop

import "time"

// Campaign status constants
const (
	StatusDraft     = "draft"
	StatusScheduled = "scheduled"
	StatusActive    = "active"
	StatusClaimable = "claimable"
	StatusClosed    = "closed"
	StatusArchived  = "archived"
)

// Asset type constants
const (
	AssetTypePoints = "points"
	AssetTypeTokens = "tokens"
	AssetTypeNative = "native"
)

// Campaign represents an airdrop campaign
type Campaign struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	AssetType        string    `json:"asset_type"`
	Status           string    `json:"status"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	TotalBudget      string    `json:"total_budget"`
	ClaimedAmount    string    `json:"claimed_amount"`
	ParticipantCount int       `json:"participant_count"`
	IsDemo           bool      `json:"is_demo"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Allocation represents a user's allocation in a campaign
type Allocation struct {
	ID         int       `json:"id"`
	CampaignID int       `json:"campaign_id"`
	UserAddress string   `json:"user_address"`
	Amount     string    `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// Claim represents a claim record
type Claim struct {
	ID          int       `json:"id"`
	CampaignID  int       `json:"campaign_id"`
	UserAddress string    `json:"user_address"`
	Amount      string    `json:"amount"`
	Nonce       string    `json:"nonce"`
	Signature   string    `json:"signature"`
	ClaimedAt   time.Time `json:"claimed_at"`
}

// EligibilityResponse represents the eligibility check response
type EligibilityResponse struct {
	Eligible bool   `json:"eligible"`
	Amount   string `json:"amount,omitempty"`
	Claimed  bool   `json:"claimed"`
	Reason   string `json:"reason,omitempty"`
}

// CreateCampaignRequest represents the request to create a campaign
type CreateCampaignRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	AssetType   string    `json:"asset_type"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	TotalBudget string    `json:"total_budget" binding:"required"`
	IsDemo      bool      `json:"is_demo"`
}

// UpdateCampaignRequest represents the request to update a campaign
type UpdateCampaignRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TotalBudget string    `json:"total_budget"`
}

// ClaimRequest represents the request to claim airdrop
type ClaimRequest struct {
	Address   string `json:"address" binding:"required"`
	Nonce     string `json:"nonce" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// AllocationImport represents a single allocation from CSV
type AllocationImport struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}
