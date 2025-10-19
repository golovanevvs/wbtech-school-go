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

.PHONY: up stop restart down purge logs ps version check-data

COMPOSE_FILE := docker-compose.yml

SERVICE := observability

DATA_DIRS := ./data/loki ./data/grafana ./data/promtail

up:
	@echo "Starting Loki, Promtail, Grafana..."
	@mkdir -p $(DATA_DIRS)
	docker compose -p $(SERVICE) -f $(COMPOSE_FILE) up -d
	@echo "Services are up and running."

stop:
	@echo "Stopping Loki, Promtail, Grafana..."
	docker compose -p $(SERVICE) -f $(COMPOSE_FILE) stop
	@echo "Loki, Promtail, Grafana stopped."

restart: stop up

down:
	@echo "Removing Loki, Promtail, Grafana (keeping data)..."
	docker compose -p $(SERVICE) -f $(COMPOSE_FILE) down
	@echo "Loki, Promtail, Grafana removed."

purge:
	@echo "Full cleanup Loki, Promtail, Grafana: containers, volumes, networks, data..."
	docker compose -p $(SERVICE) -f $(COMPOSE_FILE) down --volumes --remove-orphans
	docker network prune -f
	docker volume prune -f
	rm -rf $(DATA_DIRS)
	@echo "Loki, Promtail, Grafana cleaned."

logs:
	docker compose -p $(SERVICE) -f $(COMPOSE_FILE) logs -f

ps:
	docker compose -p $(SERVICE) -f $(COMPOSE_FILE) ps

version:
	docker compose version

check-data:
	@echo "Checking data directories:"
	@for dir in $(DATA_DIRS); do \
		echo " - $$dir"; \
		ls -la $$dir 2>/dev/null || echo "   (empty or not created yet)"; \
	done

# ========================
# GRAFANA
# ========================

# GRAFANA = providers/grafana

# grafana-up:
# 	$(MAKE_CMD) -C $(GRAFANA) up
# grafana-stop:
# 	$(MAKE_CMD) -C $(GRAFANA) stop
# grafana-restart:
# 	$(MAKE_CMD) -C $(GRAFANA) restart
# grafana-down:
# 	$(MAKE_CMD) -C $(GRAFANA) down
# grafana-purge:
# 	$(MAKE_CMD) -C $(GRAFANA) purge
# grafana-logs:
# 	$(MAKE_CMD) -C $(GRAFANA) logs
# grafana-ps:
# 	$(MAKE_CMD) -C $(GRAFANA) ps
# grafana-version:
# 	$(MAKE_CMD) -C $(GRAFANA) version
# grafana-check-data:
# 	$(MAKE_CMD) -C $(GRAFANA) check-data