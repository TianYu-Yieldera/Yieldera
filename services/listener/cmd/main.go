package main

import (
  "context"
  "encoding/json"
  "log"
  "net/http"
  "os"

  "loyalty-points-system/internal/chain"
  konf "loyalty-points-system/internal/config"
  "loyalty-points-system/internal/kafka"
)

func main() {
  if konf.Env("LISTENER_MODE", "real") != "real" {
    log.Fatal("LISTENER_MODE must be 'real' for on-chain listener")
  }
  brokers := []string{konf.Env("KAFKA_BROKERS", "kafka:9092")}
  topic   := konf.Env("KAFKA_TOPIC_RAW", "events.raw")
  w := kafka.NewWriter(brokers, topic)
  defer w.Close()
  _ = kafka.EnsureTopic(context.Background(), brokers[0], topic, 1, 1)

  var cfgs []chain.ChainCfg
  if err := json.Unmarshal([]byte(os.Getenv("CHAINS_JSON")), &cfgs); err != nil || len(cfgs) == 0 {
    log.Fatalf("invalid CHAINS_JSON: %v", err)
  }

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
  for _, cfg := range cfgs {
    if err := chain.NewWorker(cfg, w).Start(ctx); err != nil {
      log.Fatalf("start worker %s err: %v", cfg.Name, err)
    }
  }
  http.HandleFunc("/healthz", func(wr http.ResponseWriter, r *http.Request) { wr.WriteHeader(200); wr.Write([]byte("ok")) })
  log.Println("Listener (real) on :8090")
  http.ListenAndServe(":8090", nil)
}
