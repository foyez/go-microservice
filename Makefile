FRONTEND_BINARY=frontApp
BROKER_BINARY=brokerApp
INFRA_DIR_PATH=./infra

## up: starts all container in background without forcing build
up:
	@echo "starting docker containers..."
	docker compose -f ${INFRA_DIR_PATH}/docker-compose.yaml up -d
	@echo "Docker containers started!"

## up_build: stops containers (if running), builds all images and starts all containers in the background
up_build:
	@echo "Stopping docker containers (if running...)"
	docker compose -f ${INFRA_DIR_PATH}/docker-compose.yaml down
	@echo "Building (when required) and starting docker containers..."
	docker compose -f ${INFRA_DIR_PATH}/docker-compose.yaml up --build -d
	@echo "Docker images built and started all containers!"

## down: stops all containers
down:
	@echo "Stopping all dontainers..."
	docker compose -f ${INFRA_DIR_PATH}/docker-compose.yaml down
	@echo "Docker containers stopped"

.PHONY: up up_build down