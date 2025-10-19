package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	konf "loyalty-points-system/internal/config"
	appdb "loyalty-points-system/internal/db"
	"loyalty-points-system/internal/points"
	"loyalty-points-system/internal/airdrop"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL required")
	}
	rate := konf.Env("POINTS_RATE", "0.05")
	intervalStr := konf.Env("SCHEDULER_INTERVAL_SEC", "60")
	interval, _ := strconv.Atoi(intervalStr)

	db, err := appdb.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Printf("⏱️ Scheduler every %ds, rate=%s", interval, rate)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 先 ping，失败就重连
		if err := ping(db); err != nil {
			log.Printf("db ping failed: %v; reconnecting...", err)
			_ = db.Close()
			if db, err = appdb.Open(dsn); err != nil {
				log.Printf("reconnect failed: %v", err)
				continue
			}
		}
		// 跑一轮；如遇连接错误，重连并重试一次
		if err := runOnce(db, rate); err != nil {
			if isConnErr(err) {
				log.Printf("conn err: %v; reconnect & retry...", err)
				_ = db.Close()
				if db, err = appdb.Open(dsn); err != nil {
					log.Printf("reconnect failed: %v", err)
					continue
				}
				if err2 := runOnce(db, rate); err2 != nil {
					log.Printf("retry failed: %v", err2)
				}
			} else {
				log.Printf("points err: %v", err)
			}
		}
	}
}

func runOnce(db *sql.DB, rate string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Update points
	if err := points.AddPointsForAll(ctx, db, rate); err != nil {
		return err
	}

	// Update airdrop campaign statuses
	if err := airdrop.UpdateCampaignStatuses(ctx, db); err != nil {
		log.Printf("⚠️ Airdrop status update error: %v", err)
		// Don't fail the whole scheduler if airdrop update fails
	}

	return nil
}

func ping(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

func isConnErr(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "driver: bad connection") ||
		strings.Contains(s, "broken pipe") ||
		strings.Contains(s, "connection reset") ||
		errors.Is(err, context.DeadlineExceeded)
}
