package kafka

import (
  "context"
  "time"
  "github.com/segmentio/kafka-go"
  "net"
  "strconv"
)

type Message = kafka.Message

func NewWriter(brokers []string, topic string) *kafka.Writer {
  return &kafka.Writer{
    Addr:         kafka.TCP(brokers...),
    Topic:        topic,
    Balancer:     &kafka.LeastBytes{},
    RequiredAcks: kafka.RequireAll,
    BatchTimeout: 50 * time.Millisecond,
  }
}

func EnsureTopic(ctx context.Context, broker string, topic string, partitions int, replication int) error {
  d := &kafka.Dialer{Timeout: 5 * time.Second, DualStack: true}
  conn, err := d.DialContext(ctx, "tcp", broker)
  if err != nil { return err }
  defer conn.Close()
  controller, err := conn.Controller()
  if err != nil { return err }
  caddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
  cconn, err := d.DialContext(ctx, "tcp", caddr)
  if err != nil { return err }
  defer cconn.Close()
  _ = cconn.CreateTopics(kafka.TopicConfig{
    Topic: topic, NumPartitions: partitions, ReplicationFactor: replication,
  })
  return nil
}
