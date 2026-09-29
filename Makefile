.PHONY: help dev-db dev-backend dev-frontend migrate-up migrate-down migrate-create sqlc test seed

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── Development ──

dev-db: ## Start PostgreSQL in Docker
	docker compose up db -d

dev-backend: ## Run Go backend with hot reload (requires air)
	cd backend && go run ./cmd/server

dev-frontend: ## Run Vue dev server
	cd frontend && npm run dev

# ── Database ──

migrate-up: ## Run all pending migrations
	cd backend && go run -tags migrate ./cmd/migrate up

migrate-down: ## Rollback last migration
	cd backend && go run -tags migrate ./cmd/migrate down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=create_products)
	migrate create -ext sql -dir backend/internal/database/migrations -seq $(name)

sqlc: ## Generate sqlc code
	cd backend && sqlc generate

seed: ## Run seed data
	docker compose exec db psql -U supermarket -d supermarket_db -f /seed/seed.sql

# ── Testing ──

test: ## Run all Go tests
	cd backend && go test ./... -v -count=1

test-unit: ## Run unit tests only (no DB)
	cd backend && go test ./... -short -v -count=1

test-integration: ## Run integration tests (requires test DB)
	cd backend && go test ./... -v -count=1 -tags=integration

# ── Docker ──

up: ## Start all services
	docker compose up --build -d

down: ## Stop all services
	docker compose down

logs: ## Tail all logs
	docker compose logs -f

