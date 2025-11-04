package main

import (
  "context"
  "database/sql"
  "encoding/json"
  "log"
  "os"
  "os/signal"
  "strings"
  "sync"
  "syscall"
  "time"

  "loyalty-points-system/internal/config"
  "loyalty-points-system/internal/db"
  "loyalty-points-system/internal/models"

  k "github.com/segmentio/kafka-go"
)

func main() {
  log.Println("🚀 Starting Consumer Service...")

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

  // Create readers for all topics
  brokers := strings.Split(cfg.KafkaBrokers, ",")
  topicRaw := cfg.KafkaTopicRaw
  topicL1 := cfg.KafkaTopicL1
  topicL2 := cfg.KafkaTopicL2
  topicBridge := cfg.KafkaTopicBridge

  readerRaw := k.NewReader(k.ReaderConfig{
    Brokers:     brokers,
    GroupID:     "loyalty-consumer",
    Topic:       topicRaw,
    StartOffset: k.LastOffset,
  })
  defer readerRaw.Close()

  readerL1 := k.NewReader(k.ReaderConfig{
    Brokers:     brokers,
    GroupID:     "loyalty-consumer",
    Topic:       topicL1,
    StartOffset: k.LastOffset,
  })
  defer readerL1.Close()

  readerL2 := k.NewReader(k.ReaderConfig{
    Brokers:     brokers,
    GroupID:     "loyalty-consumer",
    Topic:       topicL2,
    StartOffset: k.LastOffset,
  })
  defer readerL2.Close()

  readerBridge := k.NewReader(k.ReaderConfig{
    Brokers:     brokers,
    GroupID:     "loyalty-consumer",
    Topic:       topicBridge,
    StartOffset: k.LastOffset,
  })
  defer readerBridge.Close()

  log.Printf("📥 Consumer started, reading from 4 topics...")
  log.Printf("   - %s (legacy balance events)", topicRaw)
  log.Printf("   - %s (L1 events)", topicL1)
  log.Printf("   - %s (L2 events)", topicL2)
  log.Printf("   - %s (bridge events)", topicBridge)

  // Context for graceful shutdown
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  // Listen for interrupt signals
  sigChan := make(chan os.Signal, 1)
  signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

  var wg sync.WaitGroup

  // Start goroutine for each topic
  wg.Add(4)

  go func() {
    defer wg.Done()
    consumeRaw(ctx, readerRaw, database)
  }()

  go func() {
    defer wg.Done()
    consumeL1(ctx, readerL1, database)
  }()

  go func() {
    defer wg.Done()
    consumeL2(ctx, readerL2, database)
  }()

  go func() {
    defer wg.Done()
    consumeBridge(ctx, readerBridge, database)
  }()

  // Wait for interrupt signal
  <-sigChan
  log.Println("🛑 Shutdown signal received, stopping consumer...")
  cancel()
  wg.Wait()
  log.Println("✅ Consumer stopped gracefully")
}

func applyEvent(dbx *sql.DB, evt *models.BalanceEvent) (err error) {
  tx, txErr := dbx.Begin()
  if txErr != nil {
    return txErr
  }

  // Proper transaction rollback/commit handling
  defer func() {
    if err != nil {
      if rbErr := tx.Rollback(); rbErr != nil {
        log.Printf("Transaction rollback error: %v (original error: %v)", rbErr, err)
      }
    } else {
      if cmErr := tx.Commit(); cmErr != nil {
        log.Printf("Transaction commit error: %v", cmErr)
        err = cmErr
      }
    }
  }()

  _, err = tx.Exec(`INSERT INTO users(address) VALUES($1) ON CONFLICT (address) DO NOTHING`, evt.UserAddress)
  if err != nil { return err }

  _, err = tx.Exec(`INSERT INTO balances(user_address, balance) VALUES($1, 0) ON CONFLICT (user_address) DO NOTHING`, evt.UserAddress)
  if err != nil { return err }

  _, err = tx.Exec(`INSERT INTO balance_events (user_address, amount, event_type, tx_hash, chain, block_number, confirmed)
                    VALUES ($1, $2, $3, $4, $5, $6, TRUE)`,
    evt.UserAddress, evt.Amount, evt.EventType, evt.TxHash, evt.Chain, evt.BlockNumber)
  if err != nil { return err }

  op := "+"
  if evt.EventType == "burn" || evt.EventType == "transfer_out" { op = "-" }

  _, err = tx.Exec(`
    UPDATE balances
    SET balance = CASE WHEN $3 = '-' THEN balance - CAST($2 AS NUMERIC)
                       ELSE balance + CAST($2 AS NUMERIC) END, updated_at = NOW()
    WHERE user_address = $1
  `, evt.UserAddress, evt.Amount, op)

  return err
}

// consumeRaw handles legacy balance events from events.raw topic
func consumeRaw(ctx context.Context, reader *k.Reader, database *sql.DB) {
  log.Println("📥 [RAW] Started consumer for legacy balance events")
  for {
    select {
    case <-ctx.Done():
      log.Println("🛑 [RAW] Consumer stopped")
      return
    default:
      m, err := reader.ReadMessage(ctx)
      if err != nil {
        if ctx.Err() != nil {
          return
        }
        log.Printf("❌ [RAW] Read error: %v", err)
        time.Sleep(time.Second)
        continue
      }

      var evt models.BalanceEvent
      if err := json.Unmarshal(m.Value, &evt); err != nil {
        log.Printf("❌ [RAW] JSON unmarshal error: %v", err)
        continue
      }

      if !evt.Confirmed {
        continue
      }

      if err := applyEvent(database, &evt); err != nil {
        log.Printf("❌ [RAW] Apply event error: %v", err)
      } else {
        log.Printf("✅ [RAW] Processed %s tx=%s user=%s", evt.EventType, evt.TxHash, evt.UserAddress)
      }
    }
  }
}

// consumeL1 handles L1 events from events.l1 topic
func consumeL1(ctx context.Context, reader *k.Reader, database *sql.DB) {
  log.Println("📥 [L1] Started consumer for L1 events")
  for {
    select {
    case <-ctx.Done():
      log.Println("🛑 [L1] Consumer stopped")
      return
    default:
      m, err := reader.ReadMessage(ctx)
      if err != nil {
        if ctx.Err() != nil {
          return
        }
        log.Printf("❌ [L1] Read error: %v", err)
        time.Sleep(time.Second)
        continue
      }

      var evt models.L1Event
      if err := json.Unmarshal(m.Value, &evt); err != nil {
        log.Printf("❌ [L1] JSON unmarshal error: %v", err)
        continue
      }

      if !evt.Confirmed {
        continue
      }

      if err := handleL1Event(database, &evt); err != nil {
        log.Printf("❌ [L1] Handle event error: %v", err)
      } else {
        log.Printf("✅ [L1] Processed %s tx=%s user=%s block=%d", evt.EventType, evt.TxHash, evt.UserAddress, evt.BlockNumber)
      }
    }
  }
}

// consumeL2 handles L2 events from events.l2 topic
func consumeL2(ctx context.Context, reader *k.Reader, database *sql.DB) {
  log.Println("📥 [L2] Started consumer for L2 events")
  for {
    select {
    case <-ctx.Done():
      log.Println("🛑 [L2] Consumer stopped")
      return
    default:
      m, err := reader.ReadMessage(ctx)
      if err != nil {
        if ctx.Err() != nil {
          return
        }
        log.Printf("❌ [L2] Read error: %v", err)
        time.Sleep(time.Second)
        continue
      }

      var evt models.L2Event
      if err := json.Unmarshal(m.Value, &evt); err != nil {
        log.Printf("❌ [L2] JSON unmarshal error: %v", err)
        continue
      }

      if !evt.Confirmed {
        continue
      }

      if err := handleL2Event(database, &evt); err != nil {
        log.Printf("❌ [L2] Handle event error: %v", err)
      } else {
        log.Printf("✅ [L2] Processed %s tx=%s user=%s block=%d", evt.EventType, evt.TxHash, evt.UserAddress, evt.BlockNumber)
      }
    }
  }
}

// consumeBridge handles bridge events from events.bridge topic
func consumeBridge(ctx context.Context, reader *k.Reader, database *sql.DB) {
  log.Println("📥 [BRIDGE] Started consumer for bridge events")
  for {
    select {
    case <-ctx.Done():
      log.Println("🛑 [BRIDGE] Consumer stopped")
      return
    default:
      m, err := reader.ReadMessage(ctx)
      if err != nil {
        if ctx.Err() != nil {
          return
        }
        log.Printf("❌ [BRIDGE] Read error: %v", err)
        time.Sleep(time.Second)
        continue
      }

      var evt models.BridgeEvent
      if err := json.Unmarshal(m.Value, &evt); err != nil {
        log.Printf("❌ [BRIDGE] JSON unmarshal error: %v", err)
        continue
      }

      if err := handleBridgeEvent(database, &evt); err != nil {
        log.Printf("❌ [BRIDGE] Handle event error: %v", err)
      } else {
        log.Printf("✅ [BRIDGE] Processed %s msg=%s user=%s status=%s", evt.Direction, evt.MessageHash, evt.UserAddress, evt.Status)
      }
    }
  }
}

// handleL1Event processes L1 events using the database operations
func handleL1Event(dbx *sql.DB, evt *models.L1Event) (err error) {
  tx, txErr := dbx.Begin()
  if txErr != nil {
    return txErr
  }

  defer func() {
    if err != nil {
      if rbErr := tx.Rollback(); rbErr != nil {
        log.Printf("Transaction rollback error: %v (original error: %v)", rbErr, err)
      }
    } else {
      if cmErr := tx.Commit(); cmErr != nil {
        log.Printf("Transaction commit error: %v", cmErr)
        err = cmErr
      }
    }
  }()

  return db.ProcessL1Event(tx, evt)
}

// handleL2Event processes L2 events using the database operations
func handleL2Event(dbx *sql.DB, evt *models.L2Event) (err error) {
  tx, txErr := dbx.Begin()
  if txErr != nil {
    return txErr
  }

  defer func() {
    if err != nil {
      if rbErr := tx.Rollback(); rbErr != nil {
        log.Printf("Transaction rollback error: %v (original error: %v)", rbErr, err)
      }
    } else {
      if cmErr := tx.Commit(); cmErr != nil {
        log.Printf("Transaction commit error: %v", cmErr)
        err = cmErr
      }
    }
  }()

  return db.ProcessL2Event(tx, evt)
}

// handleBridgeEvent processes bridge events using the database operations
func handleBridgeEvent(dbx *sql.DB, evt *models.BridgeEvent) (err error) {
  tx, txErr := dbx.Begin()
  if txErr != nil {
    return txErr
  }

  defer func() {
    if err != nil {
      if rbErr := tx.Rollback(); rbErr != nil {
        log.Printf("Transaction rollback error: %v (original error: %v)", rbErr, err)
      }
    } else {
      if cmErr := tx.Commit(); cmErr != nil {
        log.Printf("Transaction commit error: %v", cmErr)
        err = cmErr
      }
    }
  }()

  return db.ProcessBridgeEvent(tx, evt)
}
