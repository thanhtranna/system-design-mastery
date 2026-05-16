.PHONY: build test run clean

build:
	go build -o bin/capest ./cmd/capest

test:
	go test -v -race -cover ./...

bench:
	go test -bench=. -benchmem ./estimator

run: build
	./bin/capest estimate --preset twitter

run-markdown: build
	./bin/capest estimate --preset whatsapp --format markdown

clean:
	rm -rf bin/
