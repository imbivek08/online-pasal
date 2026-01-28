.PHONY: help run build dev migrate-up migrate-down migrate-status migrate-create seed docker-up docker-down docker-restart clean install test

# Load environment variables
ifneq (,$(wildcard .env))
    include .env
    export
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Start the server
	@echo "Starting server..."
	go run cmd/nepify/main.go

dev: ## Start server with hot reload (requires air)
	@echo "Starting server with hot reload..."
	air

build: ## Build the application
	@echo "Building application..."
	go build -o bin/nepify cmd/nepify/main.go

install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-up: ## Start PostgreSQL container
	@echo "Starting PostgreSQL container..."
	cd ../.. && docker-compose up -d

docker-down: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	cd ../.. && docker-compose down

docker-restart: ## Restart PostgreSQL container
	@echo "Restarting PostgreSQL container..."
	cd ../.. && docker-compose restart

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	goose -dir ./internal/database/migrations/ postgres "$(DATABASE_URL)" up

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	goose -dir ./internal/database/migrations/ postgres "$(DATABASE_URL)" down

migrate-status: ## Check migration status
	@echo "Checking migration status..."
	goose -dir ./internal/database/migrations/ postgres "$(DATABASE_URL)" status

migrate-reset: ## Reset database (rollback all migrations)
	@echo "Resetting database..."
	goose -dir ./internal/database/migrations/ postgres "$(DATABASE_URL)" reset

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name using 'make migrate-create name=your_migration_name'"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	goose -dir ./internal/database/migrations/ create $(name) sql

seed: ## Load seed data into database
	@echo "Loading seed data..."
	@./seed.sh

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out
	go clean

all: docker-up migrate-up run ## Start everything (Docker, migrations, and server)

setup: docker-up install migrate-up ## Initial project setup
	@echo "Setup complete! Run 'make run' to start the server."
