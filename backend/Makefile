.PHONY: test test-verbose run dev db-up db-down db-reset build clean

# Test commands
test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

# Development commands
run:
	go run cmd/server/main.go

dev:
	DATABASE_URL="postgres://bremray:devsecret@localhost/bremray_dev?sslmode=disable" go run cmd/server/main.go

# Database commands
db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

db-reset: db-down db-up
	@echo "Waiting for database to be ready..."
	@sleep 3

# Build commands
build:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...