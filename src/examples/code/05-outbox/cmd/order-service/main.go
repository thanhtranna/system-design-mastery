// Command order-service is an HTTP API that creates orders and writes
// outbox events atomically within the same transaction.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type createOrderRequest struct {
	CustomerID  string `json:"customer_id"`
	AmountCents int64  `json:"amount_cents"`
}

type orderCreatedEvent struct {
	OrderID     uuid.UUID `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	AmountCents int64     `json:"amount_cents"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", handleCreateOrder(pool))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("order-service listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func handleCreateOrder(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req createOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.CustomerID == "" || req.AmountCents <= 0 {
			http.Error(w, "customer_id and amount_cents>0 required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		orderID := uuid.New()
		eventID := uuid.New()
		now := time.Now().UTC()

		event := orderCreatedEvent{
			OrderID:     orderID,
			CustomerID:  req.CustomerID,
			AmountCents: req.AmountCents,
			CreatedAt:   now,
		}
		eventPayload, _ := json.Marshal(event)

		// THE CRITICAL TRANSACTION:
		// Order row + outbox event written atomically.
		// If either fails, both roll back.
		err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx,
				`INSERT INTO orders (id, customer_id, amount_cents, created_at)
                 VALUES ($1, $2, $3, $4)`,
				orderID, req.CustomerID, req.AmountCents, now,
			); err != nil {
				return err
			}
			_, err := tx.Exec(ctx,
				`INSERT INTO outbox (id, aggregate_id, event_type, payload, created_at)
                 VALUES ($1, $2, $3, $4, $5)`,
				eventID, orderID, "OrderCreated", eventPayload, now,
			)
			return err
		})
		if err != nil {
			log.Printf("create order failed: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"order_id": orderID,
			"status":   "created",
		})
	}
}
