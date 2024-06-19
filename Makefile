include .env
export

# Build the application
all: build

build-dev:
	@echo "Set GIN_MODE=debug"
	@export GIN_MODE=debug
	@echo "Building..."
	@start_time=$$(date +%s); \
	cd backend/rss-aggregator && go mod tidy && go build -o main cmd/api/main.go; \
	end_time=$$(date +%s); \
	echo "Building took $$((end_time - start_time)) seconds."

build-prod:
	@echo "Set GIN_MODE=release"
	@export GIN_MODE=release
	@echo "Building..."
	@start_time=$$(date +%s); \
	cd backend/rss-aggregator && go mod tidy && go build -o main cmd/api/main.go; \
	end_time=$$(date +%s); \
	echo "Building took $$((end_time - start_time)) seconds."

# Run the application
run:
	@cd backend/rss-aggregator && go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Migrate DB

migrate-up:
	@echo "Migrating up..."
	@goose -dir backend/rss-aggregator/sql/schema postgres \
	 postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable up
	@echo "...Migration up done."

migrate-down:
	@echo "Migrating down..."
	@goose -dir backend/rss-aggregator/sql/schema postgres \
	 postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable down
	@echo "...Migration down done."

# Test the application
test:
	@echo "Testing..."
	@go test backend/rss-aggregator/tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f backend/rss-aggregator/main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
