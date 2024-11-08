INFRA_DIR_PATH = ./infra

## dc_up: Run docker compose
dc_up:
	@echo "Running docker compose up in detach mode"
	docker compose -f ${INFRA_DIR_PATH}/docker-compose.yaml up -d