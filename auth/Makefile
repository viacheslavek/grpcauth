# Variables
COMPOSE_FILE=docker-compose.yml
DOCKER_COMPOSE=docker-compose -f $(COMPOSE_FILE)

PROJECT_PATH=./cmd/auth
MIGRATE_PATH=./cmd/migrator
TEST_PATH=tests/

# Docker commands
.PHONY: up down build logs

up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

build:
	@echo "Building Docker containers..."
	$(DOCKER_COMPOSE) build

logs:
	@echo "Fetching Docker container logs..."
	$(DOCKER_COMPOSE) logs -f

# Project setup
.PHONY: migrate

migrate:
	@echo "Running migrations..."
	go run $(MIGRATE_PATH)

# Run the application
.PHONY: run

run:
	@echo "Running the application..."
	go run $(PROJECT_PATH)
