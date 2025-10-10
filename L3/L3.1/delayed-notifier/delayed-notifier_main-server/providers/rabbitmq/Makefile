MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

# Load environment variables
ENV_FILE ?= $(PROJECT_DIR)/.env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

RABBIT_SERVICE := rabbitmq

.PHONY: up down logs restart status shell purge

up:
	@echo "Starting RabbitMQ..."
	@docker compose up -d $(RABBIT_SERVICE)
	$(MAKE_CMD) status

down:
	@echo "Stopping RabbitMQ..."
	@docker compose down

restart:
	@echo "Restarting RabbitMQ..."
	@docker compose restart $(RABBIT_SERVICE)
	$(MAKE_CMD) status

logs:
	@echo "Showing RabbitMQ logs..."
	@docker compose logs -f $(RABBIT_SERVICE)

status:
	@echo "RabbitMQ container status:"
	@docker compose ps $(RABBIT_SERVICE)

shell:
	@echo "Opening shell inside RabbitMQ container..."
	@docker compose exec $(RABBIT_SERVICE) bash

purge:
	@echo "Cleaning RabbitMQ completely..."
	docker compose down -v --remove-orphans
	@echo "Done. Containers, networks, and volumes removed."

# ========================
# RABBITMQ
# ========================

# RABBITMQ = providers/rabbitmq

# rabbitmq-up:
# 	$(MAKE_CMD) -C $(RABBITMQ) up
# rabbitmq-down:
# 	$(MAKE_CMD) -C $(RABBITMQ) down
# rabbitmq-restart:
# 	$(MAKE_CMD) -C $(RABBITMQ) restart
# rabbitmq-logs:
# 	$(MAKE_CMD) -C $(RABBITMQ) logs
# rabbitmq-status:
# 	$(MAKE_CMD) -C $(RABBITMQ) status
# rabbitmq-shell:
# 	$(MAKE_CMD) -C $(RABBITMQ) shell
# rabbitmq-purge:
# 	$(MAKE_CMD) -C $(RABBITMQ) purge