// Command outbox-relay polls the outbox table and publishes unpublished
// events to Kafka. Multiple instances can run concurrently — FOR UPDATE
// SKIP LOCKED handles contention.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

const defaultTopic = "order-events"

type outboxRow struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
}

func main() {
	pollInterval := envDurationMS("POLL_INTERVAL_MS", 500*time.Millisecond)
	batchSize := envInt("BATCH_SIZE", 100)
	brokers := strings.Split(envOrDefault("KAFKA_BROKERS", "kafka:9092"), ",")

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  defaultTopic,
		Balancer:               &kafka.Hash{}, // partition by key
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireAll,
		Async:                  false,
	}
	defer writer.Close()

	log.Printf("outbox-relay starting: poll=%v batch=%d brokers=%v topic=%s",
		pollInterval, batchSize, brokers, defaultTopic)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-ticker.C:
			if n, err := pollAndPublish(ctx, pool, writer, batchSize); err != nil {
				log.Printf("poll error: %v", err)
			} else if n > 0 {
				log.Printf("published %d events", n)
			}
		}
	}
}

func pollAndPublish(ctx context.Context, pool *pgxpool.Pool,
	writer *kafka.Writer, batchSize int) (int, error) {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// SKIP LOCKED lets multiple relays poll concurrently without contention.
	rows, err := tx.Query(ctx,
		`SELECT id, aggregate_id, event_type, payload
         FROM outbox
         WHERE published_at IS NULL
         ORDER BY created_at
         LIMIT $1
         FOR UPDATE SKIP LOCKED`,
		batchSize,
	)
	if err != nil {
		return 0, err
	}

	var batch []outboxRow
	for rows.Next() {
		var r outboxRow
		if err := rows.Scan(&r.ID, &r.AggregateID, &r.EventType, &r.Payload); err != nil {
			rows.Close()
			return 0, err
		}
		batch = append(batch, r)
	}
	rows.Close()

	if len(batch) == 0 {
		return 0, nil
	}

	// Publish to Kafka. Key = aggregate_id so a stream of events
	// for the same aggregate lands on the same partition (ordering).
	msgs := make([]kafka.Message, 0, len(batch))
	for _, r := range batch {
		msgs = append(msgs, kafka.Message{
			Key:   []byte(r.AggregateID.String()),
			Value: r.Payload,
			Headers: []kafka.Header{
				{Key: "event_id", Value: []byte(r.ID.String())},
				{Key: "event_type", Value: []byte(r.EventType)},
			},
		})
	}
	if err := writer.WriteMessages(ctx, msgs...); err != nil {
		return 0, err
	}

	// Mark as published — only after Kafka has acked.
	// If we crash here, the next iteration republishes (at-least-once).
	ids := make([]uuid.UUID, len(batch))
	for i, r := range batch {
		ids[i] = r.ID
	}
	if _, err := tx.Exec(ctx,
		`UPDATE outbox SET published_at = NOW() WHERE id = ANY($1)`, ids,
	); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(batch), nil
}

// --- small env helpers ---

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
func envDurationMS(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return time.Duration(n) * time.Millisecond
		}
	}
	return def
}

// suppress unused
var _ = pgx.ErrNoRows
