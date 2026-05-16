// Command order-events-consumer reads from Kafka and processes events
// with dedup-based idempotency.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	brokers := strings.Split(envOrDefault("KAFKA_BROKERS", "kafka:9092"), ",")
	dbURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "order-events",
		GroupID:     "order-events-consumer",
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer reader.Close()

	log.Printf("consumer starting brokers=%v", brokers)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("read error: %v", err)
			continue
		}

		eventID, err := extractEventID(msg.Headers)
		if err != nil {
			log.Printf("missing event_id header, skipping: %v", err)
			continue
		}

		processed, err := tryProcess(ctx, pool, eventID, msg.Value)
		if err != nil {
			// Genuine error — don't commit offset; reader will retry.
			log.Printf("process error: %v", err)
			continue
		}
		if processed {
			log.Printf("processed event %s key=%s", eventID, msg.Key)
		} else {
			log.Printf("duplicate event %s — skipped", eventID)
		}
	}
}

// tryProcess returns (true, nil) if this was a new event we processed,
// (false, nil) if it was a duplicate, or (false, err) on failure.
//
// The idempotency check + business action live in the SAME transaction.
// If the business action fails, the dedup row also rolls back, so the
// retry will re-attempt cleanly.
func tryProcess(ctx context.Context, pool *pgxpool.Pool,
	eventID uuid.UUID, payload []byte) (bool, error) {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var insertedID uuid.UUID
	err = tx.QueryRow(ctx,
		`INSERT INTO processed_events (event_id) VALUES ($1)
         ON CONFLICT (event_id) DO NOTHING
         RETURNING event_id`,
		eventID,
	).Scan(&insertedID)

	if err != nil {
		// pgx returns ErrNoRows when ON CONFLICT prevents insert.
		// That's our "duplicate" signal.
		if err.Error() == "no rows in result set" {
			_ = tx.Commit(ctx) // close the tx cleanly
			return false, nil
		}
		return false, err
	}

	// Business action goes here. In a real consumer:
	//   - update materialized view
	//   - send notification
	//   - call downstream API (with its own idempotency)
	_ = payload

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func extractEventID(headers []kafka.Header) (uuid.UUID, error) {
	for _, h := range headers {
		if h.Key == "event_id" {
			return uuid.Parse(string(h.Value))
		}
	}
	return uuid.Nil, errMissingHeader
}

var errMissingHeader = &headerErr{}

type headerErr struct{}

func (e *headerErr) Error() string { return "missing event_id header" }

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
