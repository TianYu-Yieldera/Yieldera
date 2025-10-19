package main

import (
  "context"
  "database/sql"
  "encoding/json"
  "log"
  "os"
  "strings"
  "time"

  konf "loyalty-points-system/internal/config"
  "loyalty-points-system/internal/db"
  "loyalty-points-system/internal/models"

  k "github.com/segmentio/kafka-go"
)

func main() {
  brokers := strings.Split(konf.Env("KAFKA_BROKERS", "kafka:9092"), ",")
  topic := konf.Env("KAFKA_TOPIC_RAW", "events.raw")
  dsn := os.Getenv("DATABASE_URL")
  if dsn == "" { log.Fatal("DATABASE_URL required") }
  database, err := db.Open(dsn)
  if err != nil { log.Fatal(err) }
  defer database.Close()

  reader := k.NewReader(k.ReaderConfig{ Brokers: brokers, GroupID: "loyalty-consumer", Topic: topic, StartOffset: k.LastOffset })
  defer reader.Close()
  log.Printf("📥 Consumer reading %s ...", topic)

  for {
    m, err := reader.ReadMessage(context.Background())
    if err != nil { log.Printf("read err: %v", err); time.Sleep(time.Second); continue }
    var evt models.BalanceEvent
    if err := json.Unmarshal(m.Value, &evt); err != nil { log.Printf("json err: %v", err); continue }
    if !evt.Confirmed { continue }
    if err := applyEvent(database, &evt); err != nil { log.Printf("apply err: %v", err) }
  }
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
