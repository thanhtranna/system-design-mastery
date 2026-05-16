.PHONY: test bench chaos build clean

test:
	go test -v -race -cover ./...

bench:
	go test -bench=. -benchmem ./circuitbreaker

chaos:
	go run ./cmd/chaos

build:
	go build -o bin/chaos ./cmd/chaos

clean:
	rm -rf bin/
