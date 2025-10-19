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

.PHONY: up stop restart down clean logs

SERVICE := redis

.PHONY: up down logs restart status shell purge

up:
	@echo "Starting Redis..."
	@docker compose up -d $(SERVICE)
	@echo "Redis are up and running."

stop:
	@echo "Stopping Redis..."
	@docker compose stop -d $(SERVICE)
	@echo "Redis stopped."

down:
	@echo "Down Redis..."
	@docker compose down

purge:
	@echo "Cleaning Redis completely..."
	@docker compose down -v --remove-orphans
	@rm -rf ./data/*
	@echo "Done. Containers, networks, and volumes removed."

logs:
	@docker compose logs -f $(SERVICE)

# ========================
# REDIS
# ========================

# REDIS = providers/redis

# redis-up:
# 	$(MAKE_CMD) -C $(REDIS) up
# redis-stop:
# 	$(MAKE_CMD) -C $(REDIS) stop
# redis-down:
# 	$(MAKE_CMD) -C $(REDIS) down
# redis-purge:
# 	$(MAKE_CMD) -C $(REDIS) purge
# redis-logs:
# 	$(MAKE_CMD) -C $(REDIS) logs