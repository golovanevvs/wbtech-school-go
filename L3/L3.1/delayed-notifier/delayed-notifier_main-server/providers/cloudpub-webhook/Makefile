MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

ENV_FILE ?= .env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

.PHONY: up logs stop restart rm down ps exec

PROJECT_NAME = ${CLOUDPUB_WEBHOOK_CONTAINER}
NAME = ${CLOUDPUB_WEBHOOK_CONTAINER}
SERVICE_NAME = ${CLOUDPUB_WEBHOOK_CONTAINER}

up:
	@echo "Starting ${NAME}..."
	docker compose -p $(PROJECT_NAME) up -d
	@echo "${NAME} is up and running."

logs:
	docker compose -p $(PROJECT_NAME) logs -f

stop:
	@echo "Stopping ${NAME}..."
	docker compose -p $(PROJECT_NAME) stop
	@echo "${NAME} stopped."

restart: stop up

rm:
	@echo "Removing ${NAME} (keeping data)..."
	docker compose -p $(PROJECT_NAME) rm -f
	@echo "${NAME} removed."

down:
	@echo "Full cleanup ${NAME}: containers, volumes, networks, data..."
	docker compose -p $(PROJECT_NAME) down -v --remove-orphans
	@echo "${NAME} cleaned."

ps:
	docker compose -p $(PROJECT_NAME) ps

exec:
	docker compose -p $(PROJECT_NAME) exec $(SERVICE_NAME) sh


# ========================
# CLOUDPUB-WEBHOOK
# ========================

# CLOUDPUB-WEBHOOK = providers/cloudpub-webhook

# .PHONY: cloudpub-webhook-up cloudpub-webhook-logs cloudpub-webhook-stop cloudpub-webhook-restart cloudpub-webhook-rm cloudpub-webhook-down cloudpub-webhook-ps cloudpub-webhook-exec

# cloudpub-webhook-up:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) up
# cloudpub-webhook-logs:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) logs
# cloudpub-webhook-stop:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) stop
# cloudpub-webhook-restart:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) restart
# cloudpub-webhook-rm:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) rm
# cloudpub-webhook-down:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) down
# cloudpub-webhook-ps:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) ps
# cloudpub-webhook-exec:
# 	$(MAKE_CMD) -C $(CLOUDPUB-WEBHOOK) exec