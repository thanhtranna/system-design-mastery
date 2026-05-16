.PHONY: up down logs test demo

up:
	docker compose up -d --build
	@echo "Wait ~30s for Kafka to be healthy, then: make demo"

down:
	docker compose down -v

logs:
	docker compose logs -f outbox-relay order-events-consumer

demo:
	@echo "Creating 3 orders..."
	@for i in 1 2 3; do \
		curl -sX POST http://localhost:8080/orders \
			-H 'Content-Type: application/json' \
			-d "{\"customer_id\": \"customer-$$i\", \"amount_cents\": $$((1000 * i))}" ; \
		echo ""; \
	done
	@echo ""
	@echo "Watch the logs with: make logs"

test:
	go test -v ./...

inspect:
	docker compose exec postgres psql -U outbox -c \
		"SELECT id, event_type, published_at FROM outbox ORDER BY created_at DESC LIMIT 10;"
