DEV_COMPOSE_FILE := docker-compose.dev.yaml

.PHONY: compose-up-build
compose-up-build:
	@docker compose -f $(DEV_COMPOSE_FILE) up --build -d

.PHONY: compose-up
compose-up:
	@docker compose -f $(DEV_COMPOSE_FILE) up -d

.PHONY: compose-restart
compose-restart:
	@docker compose -f $(DEV_COMPOSE_FILE) restart

.PHONY: compose-down-prune
compose-down-prune:
	@docker compose -f $(DEV_COMPOSE_FILE) down -v

.PHONY: compose-down
compose-down:
	@docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: compose-logs
compose-logs:
	@docker compose -f $(DEV_COMPOSE_FILE) logs -f

.PHONY: psql
psql:
	@docker compose -f $(DEV_COMPOSE_FILE) exec -it postgres psql -d minisocial -U postgres

.PHONY: help
help:
	@echo "Available command:"
	@echo " * compose-up-build"
	@echo " * compose-up"
	@echo " * compose-restart"
	@echo " * compose-down-prune"
	@echo " * compose-down"
	@echo " * compose-logs"
	@echo " * psql"
	@echo " * help"
