package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetL2VaultPosition returns user's vault position on L2
func GetL2VaultPosition(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		query := `
			SELECT user_address, deposited, shares, current_value, yield_earned, last_updated
			FROM l2_vault_positions
			WHERE user_address = $1
		`

		var position struct {
			UserAddress  string `json:"user_address"`
			Deposited    string `json:"deposited"`
			Shares       string `json:"shares"`
			CurrentValue string `json:"current_value"`
			YieldEarned  string `json:"yield_earned"`
			LastUpdated  string `json:"last_updated"`
		}

		err := db.QueryRow(query, address).Scan(
			&position.UserAddress,
			&position.Deposited,
			&position.Shares,
			&position.CurrentValue,
			&position.YieldEarned,
			&position.LastUpdated,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{
				"user_address": address,
				"deposited":    "0",
				"shares":       "0",
				"current_value": "0",
				"yield_earned": "0",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, position)
	}
}

// GetL2VaultStats returns overall vault statistics
func GetL2VaultStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats struct {
			TotalValueLocked string `json:"total_value_locked"`
			TotalUsers       int    `json:"total_users"`
			TotalYield       string `json:"total_yield"`
		}

		// Get TVL
		err := db.QueryRow(`
			SELECT COALESCE(SUM(current_value), 0)
			FROM l2_vault_positions
		`).Scan(&stats.TotalValueLocked)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get total users
		db.QueryRow(`SELECT COUNT(*) FROM l2_vault_positions WHERE deposited > 0`).Scan(&stats.TotalUsers)

		// Get total yield
		db.QueryRow(`
			SELECT COALESCE(SUM(yield_earned), 0)
			FROM l2_vault_positions
		`).Scan(&stats.TotalYield)

		c.JSON(http.StatusOK, stats)
	}
}

// GetL2Strategies returns available DeFi strategies
func GetL2Strategies(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT strategy, allocation_percentage, allocated_amount, current_value, apy, last_updated
			FROM l2_strategy_allocations
			ORDER BY allocated_amount DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Strategy struct {
			Name                 string `json:"name"`
			AllocationPercentage string `json:"allocation_percentage"`
			AllocatedAmount      string `json:"allocated_amount"`
			CurrentValue         string `json:"current_value"`
			APY                  string `json:"apy"`
			LastUpdated          string `json:"last_updated"`
		}

		var strategies []Strategy
		for rows.Next() {
			var s Strategy
			err := rows.Scan(&s.Name, &s.AllocationPercentage, &s.AllocatedAmount, &s.CurrentValue, &s.APY, &s.LastUpdated)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			strategies = append(strategies, s)
		}

		c.JSON(http.StatusOK, gin.H{
			"strategies": strategies,
		})
	}
}

// DepositToL2Vault deposits to the L2 vault
func DepositToL2Vault(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"user_address" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
			Signature   string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Verify signature
		// TODO: Submit transaction to L2 IntegratedVault contract
		// For now, return a placeholder response

		c.JSON(http.StatusAccepted, gin.H{
			"status": "pending",
			"message": "Deposit transaction submitted to L2 vault",
			"user_address": req.UserAddress,
			"amount": req.Amount,
			"note": "This endpoint requires L2 contract integration",
		})
	}
}

// WithdrawFromL2Vault withdraws from the L2 vault
func WithdrawFromL2Vault(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserAddress string `json:"user_address" binding:"required"`
			Amount      string `json:"amount" binding:"required"`
			Signature   string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user has sufficient balance
		var currentValue string
		err := db.QueryRow(`
			SELECT current_value FROM l2_vault_positions
			WHERE user_address = $1
		`, req.UserAddress).Scan(&currentValue)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No vault position found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: Verify signature
		// TODO: Submit transaction to L2 IntegratedVault contract
		// For now, return a placeholder response

		c.JSON(http.StatusAccepted, gin.H{
			"status": "pending",
			"message": "Withdrawal transaction submitted to L2 vault",
			"user_address": req.UserAddress,
			"amount": req.Amount,
			"current_value": currentValue,
			"note": "This endpoint requires L2 contract integration",
		})
	}
}

// GetL2RWAAssets returns list of RWA assets
func GetL2RWAAssets(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT asset_id, asset_name, asset_type, total_supply, price_per_token, valuation, status, created_at
			FROM l2_rwa_assets
			WHERE status = 'active'
			ORDER BY created_at DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type RWAAsset struct {
			AssetID       int64  `json:"asset_id"`
			AssetName     string `json:"asset_name"`
			AssetType     string `json:"asset_type"`
			TotalSupply   string `json:"total_supply"`
			PricePerToken string `json:"price_per_token"`
			Valuation     string `json:"valuation"`
			Status        string `json:"status"`
			CreatedAt     string `json:"created_at"`
		}

		var assets []RWAAsset
		for rows.Next() {
			var a RWAAsset
			err := rows.Scan(&a.AssetID, &a.AssetName, &a.AssetType, &a.TotalSupply, &a.PricePerToken, &a.Valuation, &a.Status, &a.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, a)
		}

		c.JSON(http.StatusOK, gin.H{
			"assets": assets,
		})
	}
}

// GetL2RWAHoldings returns user's RWA holdings
func GetL2RWAHoldings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		query := `
			SELECT h.user_address, h.asset_id, a.asset_name, h.amount, h.purchase_price, h.current_value, h.updated_at
			FROM l2_rwa_holdings h
			LEFT JOIN l2_rwa_assets a ON h.asset_id = a.asset_id
			WHERE h.user_address = $1 AND h.amount > 0
			ORDER BY h.updated_at DESC
		`

		rows, err := db.Query(query, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Holding struct {
			UserAddress   string `json:"user_address"`
			AssetID       int64  `json:"asset_id"`
			AssetName     string `json:"asset_name"`
			Amount        string `json:"amount"`
			PurchasePrice string `json:"purchase_price"`
			CurrentValue  string `json:"current_value"`
			UpdatedAt     string `json:"updated_at"`
		}

		var holdings []Holding
		for rows.Next() {
			var h Holding
			err := rows.Scan(&h.UserAddress, &h.AssetID, &h.AssetName, &h.Amount, &h.PurchasePrice, &h.CurrentValue, &h.UpdatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			holdings = append(holdings, h)
		}

		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"holdings": holdings,
		})
	}
}

// GetL2RWAListings returns active marketplace listings
func GetL2RWAListings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT l.listing_id, l.asset_id, a.asset_name, l.seller_address, l.amount, l.price_per_token, l.status, l.created_at
			FROM l2_rwa_listings l
			LEFT JOIN l2_rwa_assets a ON l.asset_id = a.asset_id
			WHERE l.status = 'active'
			ORDER BY l.created_at DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Listing struct {
			ListingID     int64  `json:"listing_id"`
			AssetID       int64  `json:"asset_id"`
			AssetName     string `json:"asset_name"`
			SellerAddress string `json:"seller_address"`
			Amount        string `json:"amount"`
			PricePerToken string `json:"price_per_token"`
			Status        string `json:"status"`
			CreatedAt     string `json:"created_at"`
		}

		var listings []Listing
		for rows.Next() {
			var l Listing
			err := rows.Scan(&l.ListingID, &l.AssetID, &l.AssetName, &l.SellerAddress, &l.Amount, &l.PricePerToken, &l.Status, &l.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			listings = append(listings, l)
		}

		c.JSON(http.StatusOK, gin.H{
			"listings": listings,
		})
	}
}

// GetL2RWAProposals returns governance proposals
func GetL2RWAProposals(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT p.proposal_id, p.asset_id, a.asset_name, p.proposer_address, p.description, p.votes_for, p.votes_against, p.status, p.created_at
			FROM l2_rwa_proposals p
			LEFT JOIN l2_rwa_assets a ON p.asset_id = a.asset_id
			WHERE p.status = 'active'
			ORDER BY p.created_at DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Proposal struct {
			ProposalID      int64  `json:"proposal_id"`
			AssetID         int64  `json:"asset_id"`
			AssetName       string `json:"asset_name"`
			ProposerAddress string `json:"proposer_address"`
			Description     string `json:"description"`
			VotesFor        string `json:"votes_for"`
			VotesAgainst    string `json:"votes_against"`
			Status          string `json:"status"`
			CreatedAt       string `json:"created_at"`
		}

		var proposals []Proposal
		for rows.Next() {
			var p Proposal
			err := rows.Scan(&p.ProposalID, &p.AssetID, &p.AssetName, &p.ProposerAddress, &p.Description, &p.VotesFor, &p.VotesAgainst, &p.Status, &p.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			proposals = append(proposals, p)
		}

		c.JSON(http.StatusOK, gin.H{
			"proposals": proposals,
		})
	}
}
