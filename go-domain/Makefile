.PHONY: help dev build run test clean migrate-up migrate-down migrate-create docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Run the application in development mode
	go run ./cmd/api

build: ## Build the application
	go build -o bin/api ./cmd/api

run: build ## Build and run the application
	./bin/api

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

migrate-up: ## Run database migrations up
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL is not set. Loading from .env file..."; \
		export $$(cat .env | xargs) && \
		migrate -path ./migrations -database "$$DATABASE_URL" up; \
	else \
		migrate -path ./migrations -database "$(DATABASE_URL)" up; \
	fi

migrate-down: ## Rollback last migration
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL is not set. Loading from .env file..."; \
		export $$(cat .env | xargs) && \
		migrate -path ./migrations -database "$$DATABASE_URL" down 1; \
	else \
		migrate -path ./migrations -database "$(DATABASE_URL)" down 1; \
	fi

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_posts_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=create_posts_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./migrations -seq $(name)

docker-up: ## Start PostgreSQL using Docker Compose
	docker-compose up -d

docker-down: ## Stop PostgreSQL
	docker-compose down

docker-logs: ## View PostgreSQL logs
	docker-compose logs -f postgres

tidy: ## Tidy go modules
	go mod tidy

install-tools: ## Install required development tools
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed successfully!"

setup: docker-up ## Setup development environment
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "Creating .env file from .env.example if it doesn't exist..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Setup complete! Run 'make dev' to start the application."
