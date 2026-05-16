.PHONY: up down logs test demo concurrent

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f

test:
	go test -v -race ./...

demo:
	@echo "Write to node-1, read from node-2 (sequential)..."
	curl -sX PUT http://localhost:8001/kv/greeting -d '{"value": "hello"}'
	@echo ""
	@sleep 1
	curl -s http://localhost:8002/kv/greeting
	@echo ""

concurrent:
	@echo "Concurrent writes to node-1 and node-2..."
	curl -sX PUT http://localhost:8001/kv/score -d '{"value": "100"}' & \
	curl -sX PUT http://localhost:8002/kv/score -d '{"value": "200"}' & \
	wait
	@echo ""
	@sleep 1
	@echo "Read from node-3 — should see BOTH concurrent versions:"
	curl -s http://localhost:8003/kv/score
	@echo ""
