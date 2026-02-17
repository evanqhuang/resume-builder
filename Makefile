.PHONY: up down build logs restart clean dev-cli dev-server dev help

.DEFAULT_GOAL := help

up: ## Start all services
	docker compose up -d
	@echo ""
	@echo "  Frontend: http://localhost:5173"
	@echo "  Backend:  http://localhost:8080"
	@echo ""

down: ## Stop all services
	docker compose down

build: ## Build and start all services
	docker compose up --build -d
	@echo ""
	@echo "  Frontend: http://localhost:5173"
	@echo "  Backend:  http://localhost:8080"
	@echo ""

logs: ## Tail logs from all services
	docker compose logs -f

restart: ## Restart all services
	docker compose restart

dev-cli: ## Run CLI (pass args via ARGS=)
	cd cli && go run . $(ARGS)

dev-server: ## Start Go backend server for local development
	cd cli && go run . serve --port 8080

dev: ## Start backend and frontend for local development
	@echo "Starting backend server..."
	@cd cli && go run . serve --port 8080 &
	@sleep 2
	@echo "Starting frontend..."
	@cd web/frontend && npm run dev

clean: ## Stop services and remove volumes/images
	docker compose down --rmi local --volumes
	rm -f cli/resume-cli

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
