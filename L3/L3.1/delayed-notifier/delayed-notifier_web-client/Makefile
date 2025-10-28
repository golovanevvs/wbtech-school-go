MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath .)

ENV_FILE ?= .env.production
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

.PHONY: up logs stop restart rm down ps exec build rebuild

NAME = web-client
SERVICE_NAME = web-client

up:
	@echo "Starting ${NAME}..."
	docker compose -p $(PROJECT_NAME) up -d
	@echo "${NAME} is up and running."

logs:
	docker compose -p $(PROJECT_NAME) logs -f $(SERVICE_NAME)

stop:
	@echo "Stopping ${NAME}..."
	docker compose -p $(PROJECT_NAME) stop
	@echo "${NAME} stopped."

restart: stop up

rm:
	@echo "Removing ${NAME} (keeping volumes and images)..."
	docker compose -p $(PROJECT_NAME) rm -f
	@echo "${NAME} containers removed."

down:
	@echo "Full cleanup ${NAME}: containers, volumes, networks..."
	docker compose -p $(PROJECT_NAME) down -v --remove-orphans
	@echo "${NAME} fully cleaned."

ps:
	docker compose -p $(PROJECT_NAME) ps

exec:
	docker compose -p $(PROJECT_NAME) exec $(SERVICE_NAME) sh

build:
	@echo "Building ${NAME} image..."
	docker compose -p $(PROJECT_NAME) build
	@echo "Build completed."

rebuild:
	@echo "Rebuilding ${NAME} from scratch..."
	docker compose -p $(PROJECT_NAME) build --no-cache
	@echo "Rebuild done."
