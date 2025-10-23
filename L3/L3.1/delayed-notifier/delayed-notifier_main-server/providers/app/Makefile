MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

# Load environment variables
ENV_FILE ?= .env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export APP_NAME APP_TRANSPORT_HTTP_PORT
endif

ifeq ($(OS),Windows_NT)
    RUN_CMD = start "" "/C/Program Files/Git/git-bash.exe" -c
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        # macOS
        RUN_CMD = open -a Terminal.app
    else
        # Linux (GNOME, KDE, etc.)
        RUN_CMD = x-terminal-emulator -e
    endif
endif

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

PROJECT_NAME = ${APP_NAME}
SERVICE_NAME = ${APP_NAME}

# FOR LOCAL UP
local-build:
	@echo "Building app $(APP_NAME)..."
	@cd $(PROJECT_DIR) && go build ./cmd/$(APP_NAME)
	@echo "Building app $(APP_NAME) completed"

local-up:
	@echo "Running app $(APP_NAME)..."
	@cd $(PROJECT_DIR) && $(RUN_CMD) "./$(APP_NAME) $(ARGS)"
	@sleep 1
	@echo "Running app $(APP_NAME) completed"

local-build-up: local-build local-up

# FOR DOCKER UP
up:
	@echo "Building and starting container $(APP_NAME)..."
	docker compose -p $(PROJECT_NAME) up --build -d
	@echo "Containers are up"

logs:
	@echo "Showing logs for ${APP_NAME}..."
	docker compose -p $(PROJECT_NAME) logs -f

stop:
	@echo "Stopping ${APP_NAME}..."
	docker compose -p $(PROJECT_NAME) stop
	@echo "${APP_NAME} stopped."

restart:
	@echo "Restarting ${APP_NAME}..."
	$(MAKE_CMD) stop
	$(MAKE_CMD) up

rm:
	@echo "➡ Removing ${APP_NAME}..."
	docker compose -p $(PROJECT_NAME) rm -f
	@echo "${NAME} removed."

down:
	@echo "Full cleanup ${APP_NAME}: containers, volumes, networks, data..."
	docker compose -p $(PROJECT_NAME) down -v --remove-orphans
	@echo "${NAME} cleaned."

ps:
    @echo "Status for ${APP_NAME}:"
	docker compose -p $(PROJECT_NAME) ps

exec:
	@echo "Connecting to container ${APP_NAME}..."
	docker compose -p $(PROJECT_NAME) exec $(SERVICE_NAME) sh


# ========================
# APP
# ========================

# APP = providers/app

# .PHONY: app-local-build app-local-up app-local-build-up app-up app-logs app-stop app-restart app-rm app-down app-ps app-exec

# app-local-build:
# 	$(MAKE_CMD) -C $(APP) local-build
# app-local-up:
# 	$(MAKE_CMD) -C $(APP) local-up
# app-local-build-up:
# 	$(MAKE_CMD) -C $(APP) local-build-up
# app-up:
# 	$(MAKE_CMD) -C $(APP) up
# app-logs:
# 	$(MAKE_CMD) -C $(APP) logs
# app-stop:
# 	$(MAKE_CMD) -C $(APP) stop
# app-restart:
# 	$(MAKE_CMD) -C $(APP) restart
# app-rm:
# 	$(MAKE_CMD) -C $(APP) rm
# app-down:
# 	$(MAKE_CMD) -C $(APP) down
# app-ps:
# 	$(MAKE_CMD) -C $(APP) ps
# app-exec:
# 	$(MAKE_CMD) -C $(APP) exec