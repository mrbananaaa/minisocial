DEV_COMPOSE_FILE := docker-compose.dev.yaml

.PHONY: compose-dev-up-build
compose-dev-up-build:
	@docker compose -f $(DEV_COMPOSE_FILE) up --build -d

.PHONY: compose-dev-up
compose-dev-up:
	@docker compose -f $(DEV_COMPOSE_FILE) up -d

.PHONY: compose-dev-restart
compose-dev-restart:
	@docker compose -f $(DEV_COMPOSE_FILE) restart

.PHONY: compose-dev-down
compose-dev-down:
	@docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: compose-dev-logs
compose-dev-logs:
	@docker compose -f $(DEV_COMPOSE_FILE) logs -f

.PHONY: psql
psql:
	@docker compose -f $(DEV_COMPOSE_FILE) exec -it postgres psql -d minisocial -U postgres
