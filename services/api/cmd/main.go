package main

import (
  "context"
  "database/sql"
  "errors"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
  "os"
  "os/signal"
  "strconv"
  "strings"
  "syscall"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/graphql-go/graphql"
  "loyalty-points-system/internal/config"
  "loyalty-points-system/internal/db"
  "loyalty-points-system/internal/airdrop"
  "loyalty-points-system/services/api/handlers"
  "loyalty-points-system/services/api/middleware"
)

func main() {
  log.Println("🚀 Starting API Service...")

  // Load configuration
  cfg, err := config.LoadConfig()
  if err != nil {
    log.Fatalf("Failed to load config: %v", err)
  }

  database, err := db.Open(cfg.DatabaseURL)
  if err != nil { log.Fatal(err) }
  defer database.Close()

  r := gin.Default()
  r.Use(cors(cfg.APIAllowOrigin))
  r.Use(timeoutMiddleware(30 * time.Second))
  r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

  // Authentication routes (public, no auth required)
  auth := r.Group("/auth")
  {
    auth.GET("/message", middleware.GetAuthMessageHandler)
    auth.POST("/authenticate", middleware.AuthenticateHandler)
  }

  // Proxy routes to microservices
  // Vault Service proxy
  vaultTarget := cfg.VaultServiceURL
  if vaultURL, err := url.Parse(vaultTarget); err == nil {
    vaultProxy := httputil.NewSingleHostReverseProxy(vaultURL)
    originalDirector := vaultProxy.Director
    vaultProxy.Director = func(req *http.Request) {
      originalDirector(req)
      req.Host = vaultURL.Host
      req.URL.Host = vaultURL.Host
      req.URL.Scheme = vaultURL.Scheme
    }
    vaultAPI := r.Group("/api/vault")
    vaultAPI.Any("/*path", func(c *gin.Context) {
      vaultProxy.ServeHTTP(c.Writer, c.Request)
    })
  }

  // RWA Service proxy
  rwaTarget := cfg.RWAServiceURL
  if rwaURL, err := url.Parse(rwaTarget); err == nil {
    rwaProxy := httputil.NewSingleHostReverseProxy(rwaURL)
    originalDirector := rwaProxy.Director
    rwaProxy.Director = func(req *http.Request) {
      originalDirector(req)
      req.Host = rwaURL.Host
      req.URL.Host = rwaURL.Host
      req.URL.Scheme = rwaURL.Scheme
    }
    rwaAPI := r.Group("/api/rwa")
    rwaAPI.Any("/*path", func(c *gin.Context) {
      rwaProxy.ServeHTTP(c.Writer, c.Request)
    })
  }

  // Oracle Service proxy
  oracleTarget := cfg.OracleServiceURL
  if oracleURL, err := url.Parse(oracleTarget); err == nil {
    oracleProxy := httputil.NewSingleHostReverseProxy(oracleURL)
    originalDirector := oracleProxy.Director
    oracleProxy.Director = func(req *http.Request) {
      originalDirector(req)
      req.Host = oracleURL.Host
      req.URL.Host = oracleURL.Host
      req.URL.Scheme = oracleURL.Scheme
    }
    oracleAPI := r.Group("/api/oracle")
    oracleAPI.Any("/*path", func(c *gin.Context) {
      oracleProxy.ServeHTTP(c.Writer, c.Request)
    })
  }

  r.GET("/users/:addr/balance", func(c *gin.Context) {
    addr := c.Param("addr")
    if err := validateEthereumAddress(addr); err != nil {
      c.JSON(400, gin.H{"error": err.Error()})
      return
    }
    var bal string
    err := database.QueryRow(`SELECT balance FROM balances WHERE user_address=$1`, addr).Scan(&bal)
    if err == sql.ErrNoRows { c.JSON(404, gin.H{"error":"not found"}); return }
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    c.JSON(200, gin.H{"address": addr, "balance": bal})
  })

  r.GET("/users/:addr/points", func(c *gin.Context) {
    addr := c.Param("addr")
    if err := validateEthereumAddress(addr); err != nil {
      c.JSON(400, gin.H{"error": err.Error()})
      return
    }
    var pts string
    err := database.QueryRow(`SELECT points FROM points WHERE user_address=$1`, addr).Scan(&pts)
    if err == sql.ErrNoRows { pts = "0" } else if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    c.JSON(200, gin.H{"address": addr, "points": pts})
  })

  r.GET("/users/:addr/badges", func(c *gin.Context) {
    addr := c.Param("addr")
    if err := validateEthereumAddress(addr); err != nil {
      c.JSON(400, gin.H{"error": err.Error()})
      return
    }
    rows, err := database.Query(`SELECT badge_code FROM badges WHERE user_address=$1 ORDER BY created_at`, addr)
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    defer rows.Close()
    var codes []string
    for rows.Next(){ var code string; _=rows.Scan(&code); codes = append(codes, code) }
    c.JSON(200, gin.H{"address": addr, "badges": codes})
  })

  r.GET("/leaderboard", func(c *gin.Context) {
    offset := c.DefaultQuery("offset", "0")
    limit := c.DefaultQuery("limit", "20")

    // Validate and limit max page size
    limitInt, _ := strconv.Atoi(limit)
    if limitInt > 100 {
      limitInt = 100
    }

    rows, err := database.Query(`SELECT user_address, points FROM points ORDER BY CAST(points AS NUMERIC) DESC LIMIT $1 OFFSET $2`, limitInt, offset)
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    defer rows.Close()

    type item struct{ Address, Points string }
    var out []item
    for rows.Next(){ var it item; if err := rows.Scan(&it.Address, &it.Points); err==nil { out = append(out, it) } }

    // Get total count
    var total int
    database.QueryRow(`SELECT COUNT(*) FROM points`).Scan(&total)

    c.JSON(200, gin.H{
      "items": out,
      "pagination": gin.H{
        "total": total,
        "offset": offset,
        "limit": limitInt,
      },
    })
  })

  // DeFi pool routes
  defi := r.Group("/api/defi")
  {
    defi.GET("/pools", handlers.GetDeFiPools(database))
    defi.GET("/pools/:id", handlers.GetPoolDetail(database))
    defi.GET("/positions/:address", handlers.GetUserDeFiPositions(database))
    defi.POST("/deposit", handlers.DepositToPool(database))
    defi.POST("/withdraw", handlers.WithdrawFromPool(database))
    defi.POST("/claim", handlers.ClaimRewards(database))
    defi.GET("/history/:address", handlers.GetDeFiHistory(database))
    defi.GET("/stats", handlers.GetDeFiStats(database))
  }

  // Stablecoin routes
  stable := r.Group("/api/stablecoin")
  {
    stable.GET("/position/:address", handlers.GetStablecoinPosition(database))
    stable.POST("/simulate-mint", handlers.SimulateMint(database))
    stable.POST("/simulate-redeem", handlers.SimulateRedeem(database))
    stable.POST("/mint", handlers.MintLUSD(database))
    stable.POST("/redeem", handlers.RedeemLUSD(database))
    stable.GET("/history/:address", handlers.GetStablecoinHistory(database))
    stable.GET("/stats", handlers.GetStablecoinStats(database))
  }

  // Demo mode routes (for hackathon and new users)
  demo := r.Group("/api/demo")
  {
    demo.POST("/create", handlers.CreateDemoUser(database))
    demo.GET("/status", handlers.GetDemoStatus(database))
    demo.GET("/summary", handlers.GetDemoSummary(database))
    demo.POST("/reset", handlers.ResetDemoUser(database))
    demo.POST("/exit", handlers.ExitDemoMode(database))
  }

  // Airdrop routes - Admin (requires admin authentication)
  adminAirdrop := r.Group("/api/admin/airdrop")
  adminAirdrop.Use(airdrop.AdminAuthMiddleware(database))
  {
    adminAirdrop.POST("/campaigns", airdrop.CreateCampaignHandler(database))
    adminAirdrop.PUT("/campaigns/:id", airdrop.UpdateCampaignHandler(database))
    adminAirdrop.POST("/campaigns/:id/allocations/import", airdrop.ImportAllocationsHandler(database))
    adminAirdrop.POST("/campaigns/:id/activate", airdrop.ActivateCampaignHandler(database))
    adminAirdrop.POST("/campaigns/:id/close", airdrop.CloseCampaignHandler(database))
    adminAirdrop.GET("/campaigns/:id/stats", airdrop.GetCampaignStatsHandler(database))
  }

  // Airdrop routes - Public (no auth required for listing and checking eligibility)
  publicAirdrop := r.Group("/api/airdrop")
  {
    publicAirdrop.GET("/campaigns", airdrop.GetCampaignsHandler(database))
    publicAirdrop.GET("/campaigns/:id", airdrop.GetCampaignHandler(database))
    publicAirdrop.GET("/campaigns/:id/eligibility", airdrop.CheckEligibilityHandler(database))
    publicAirdrop.POST("/campaigns/:id/claim", airdrop.ClaimAirdropHandler(database))
  }

  // L1 routes (Layer 1 collateral management)
  l1 := r.Group("/api/v1/l1")
  {
    l1.GET("/user/:address/balance", handlers.GetL1Balance(database))
    l1.GET("/user/:address/deposits", handlers.GetL1Deposits(database))
    l1.POST("/deposit", handlers.InitiateL1Deposit(database))
    l1.POST("/withdraw", handlers.InitiateL1Withdrawal(database))
    l1.GET("/state/snapshots", handlers.GetL1StateSnapshots(database))
  }

  // L2 routes (Layer 2 vault and RWA management)
  l2 := r.Group("/api/v1/l2")
  {
    l2.GET("/user/:address/position", handlers.GetL2VaultPosition(database))
    l2.GET("/vault/stats", handlers.GetL2VaultStats(database))
    l2.GET("/strategies", handlers.GetL2Strategies(database))
    l2.POST("/deposit", handlers.DepositToL2Vault(database))
    l2.POST("/withdraw", handlers.WithdrawFromL2Vault(database))
    l2.GET("/rwa/assets", handlers.GetL2RWAAssets(database))
    l2.GET("/rwa/user/:address/holdings", handlers.GetL2RWAHoldings(database))
    l2.GET("/rwa/marketplace/listings", handlers.GetL2RWAListings(database))
    l2.GET("/rwa/governance/proposals", handlers.GetL2RWAProposals(database))
  }

  // Bridge routes (L1 <-> L2 cross-chain operations)
  bridgeAPI := r.Group("/api/v1/bridge")
  {
    bridgeAPI.GET("/status/:messageHash", handlers.GetBridgeStatus(database))
    bridgeAPI.GET("/user/:address/messages", handlers.GetUserBridgeHistory(database))
    bridgeAPI.POST("/l1-to-l2", handlers.InitiateBridgeL1ToL2(database))
    bridgeAPI.POST("/l2-to-l1", handlers.InitiateBridgeL2ToL1(database))
    bridgeAPI.POST("/retry/:messageHash", handlers.RetryBridgeMessage(database))
    bridgeAPI.GET("/stats", handlers.GetBridgeStats(database))
  }

  schema, err := buildSchema(database)
  if err != nil {
    log.Fatal("GraphQL schema error: ", err)
  }
  r.POST("/graphql", func(c *gin.Context) {
    var body struct{ Query string `json:"query"` }
    if err := c.BindJSON(&body); err != nil { c.JSON(400, gin.H{"error":"bad request"}); return }
    result := graphql.Do(graphql.Params{Schema: schema, RequestString: body.Query})
    if len(result.Errors) > 0 { c.JSON(400, result); return }
    c.JSON(200, result)
  })

  // Create HTTP server
  srv := &http.Server{
    Addr:    ":" + cfg.APIPort,
    Handler: r,
  }

  // Graceful shutdown
  go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
    <-sigChan
    log.Println("🛑 Shutting down gracefully...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
      log.Printf("Server shutdown error: %v", err)
    }

    database.Close()
    log.Println("✅ Shutdown complete")
  }()

  log.Printf("🚀 API on :%s", cfg.APIPort)
  if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
  }
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
    defer cancel()
    c.Request = c.Request.WithContext(ctx)
    c.Next()
  }
}

func validateEthereumAddress(addr string) error {
  if len(addr) != 42 {
    return errors.New("invalid ethereum address length")
  }
  if !strings.HasPrefix(strings.ToLower(addr), "0x") {
    return errors.New("ethereum address must start with 0x")
  }
  for _, c := range addr[2:] {
    if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
      return errors.New("invalid hex characters in address")
    }
  }
  return nil
}

func cors(allow string) gin.HandlerFunc {
  return func(c *gin.Context) {
    origin := c.Request.Header.Get("Origin")
    // If allow is "*" or origin matches, set it
    if allow == "*" || origin == allow {
      c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
    } else if strings.Contains(allow, ",") {
      // Support multiple origins separated by comma
      for _, allowed := range strings.Split(allow, ",") {
        if strings.TrimSpace(allowed) == origin {
          c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
          break
        }
      }
    }
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    if c.Request.Method == http.MethodOptions { c.AbortWithStatus(204); return }
    c.Next()
  }
}

func buildSchema(db *sql.DB) (graphql.Schema, error) {
  balanceType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Balance", Fields: graphql.Fields{"address": &graphql.Field{Type: graphql.String}, "balance": &graphql.Field{Type: graphql.String}},
  })
  pointsType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Points", Fields: graphql.Fields{"address": &graphql.Field{Type: graphql.String}, "points": &graphql.Field{Type: graphql.String}},
  })
  badgeType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Badge", Fields: graphql.Fields{"address": &graphql.Field{Type: graphql.String}, "code": &graphql.Field{Type: graphql.String}},
  })
  lbItem := graphql.NewObject(graphql.ObjectConfig{
    Name: "LeaderboardItem", Fields: graphql.Fields{"address": &graphql.Field{Type: graphql.String}, "points": &graphql.Field{Type: graphql.String}},
  })
  query := graphql.NewObject(graphql.ObjectConfig{
    Name: "Query",
    Fields: graphql.Fields{
      "balance": &graphql.Field{
        Type: balanceType,
        Args: graphql.FieldConfigArgument{"address": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)}},
        Resolve: func(p graphql.ResolveParams) (any, error) {
          addr := p.Args["address"].(string)
          var bal string
          err := db.QueryRow(`SELECT balance FROM balances WHERE user_address=$1`, addr).Scan(&bal)
          if err == sql.ErrNoRows { bal = "0" } else if err != nil { return nil, err }
          return map[string]any{"address": addr, "balance": bal}, nil
        },
      },
      "points": &graphql.Field{
        Type: pointsType,
        Args: graphql.FieldConfigArgument{"address": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)}},
        Resolve: func(p graphql.ResolveParams) (any, error) {
          addr := p.Args["address"].(string)
          var pts string
          err := db.QueryRow(`SELECT points FROM points WHERE user_address=$1`, addr).Scan(&pts)
          if err == sql.ErrNoRows { pts = "0" } else if err != nil { return nil, err }
          return map[string]any{"address": addr, "points": pts}, nil
        },
      },
      "badges": &graphql.Field{
        Type: graphql.NewList(badgeType),
        Args: graphql.FieldConfigArgument{"address": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)}},
        Resolve: func(p graphql.ResolveParams) (any, error) {
          addr := p.Args["address"].(string)
          rows, err := db.Query(`SELECT badge_code FROM badges WHERE user_address=$1 ORDER BY created_at`, addr)
          if err != nil { return nil, err }
          defer rows.Close()
          var res []map[string]any
          for rows.Next(){ var code string; _=rows.Scan(&code); res = append(res, map[string]any{"address": addr, "code": code}) }
          return res, nil
        },
      },
      "leaderboard": &graphql.Field{
        Type: graphql.NewList(lbItem),
        Resolve: func(p graphql.ResolveParams) (any, error) {
          rows, err := db.Query(`SELECT user_address, points FROM points ORDER BY CAST(points AS NUMERIC) DESC LIMIT 20`)
          if err != nil { return nil, err }
          defer rows.Close()
          var out []map[string]any
          for rows.Next(){ var a, pts string; _=rows.Scan(&a, &pts); out = append(out, map[string]any{"address": a, "points": pts}) }
          return out, nil
        },
      },
    },
  })
  return graphql.NewSchema(graphql.SchemaConfig{Query: query})
}
