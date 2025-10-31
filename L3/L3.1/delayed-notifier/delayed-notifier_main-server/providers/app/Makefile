# SETUP
MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

ENV_FILES := .env ../email/.env ../rabbitmq/.env ../redis/.env ../telegram/.env 

$(foreach file,$(ENV_FILES), \
  $(if $(wildcard $(file)), $(eval include $(file))) \
)

export APP_NAME APP_TRANSPORT_HTTP_PORT APP_TRANSPORT_HTTP_PORT_EXT APP_TRANSPORT_HTTP_PUBLIC_HOST APP_TRANSPORT_HTTP_PUBLIC_HOST_DOCKER APP_TRANSPORT_HTTP_WEB_PUBLIC_HOST_DOCKER
export EMAIL_SMTP_HOST EMAIL_SMTP_PORT EMAIL_USERNAME EMAIL_PASSWORD EMAIL_FROM
export RABBITMQ_PORT RABBITMQ_VHOST RABBITMQ_USERNAME RABBITMQ_PASSWORD
export REDIS_PORT REDIS_PASSWORD
export TELEGRAM_TOKEN

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

NAME = ${APP_CONTAINER}
PROJECT_NAME = ${APP_CONTAINER}
SERVICE_NAME = ${APP_CONTAINER}

# FOR LOCAL UP
local-build:
	@echo "Building app $(NAME)..."
	@cd $(PROJECT_DIR) && go build ./cmd/$(APP_NAME)
	@echo "Building app $(NAME) completed"

ifeq ($(OS),Windows_NT)
	@echo "Creating run script $(APP_NAME).bat..."
	@echo @echo off > $(PROJECT_DIR)/$(APP_NAME).bat
	@echo set ENV=local >> $(PROJECT_DIR)/$(APP_NAME).bat
	@echo $(APP_NAME).exe >> $(PROJECT_DIR)/$(APP_NAME).bat
	@echo "File $(APP_NAME).bat created successfully"
endif

local-up:
	@echo "Running app $(NAME)..."
# 	@cd $(PROJECT_DIR) && $(RUN_CMD) "./$(APP_NAME) $(ARGS)"
# 	@echo "Running app $(NAME) completed"
ifeq ($(OS),Windows_NT)
	start powershell -NoExit -Command "cd '$(PROJECT_DIR)'; .\$(APP_NAME).exe $(ARGS)"
else
	@cd $(PROJECT_DIR) && $(RUN_CMD) "./$(APP_NAME) $(ARGS)"
endif
	@echo "Running app $(NAME) completed"

local-build-up: local-build local-up

# FOR DOCKER UP
up:
	@echo "Building and starting container $(NAME)..."
	docker compose -p $(PROJECT_NAME) build
	docker compose -p $(PROJECT_NAME) up -d
	@echo "${NAME} is up and running."

up-no-cache:
	@echo "Building and starting container $(NAME)..."
	docker compose -p $(PROJECT_NAME) build --no-cache
	docker compose -p $(PROJECT_NAME) up -d
	@echo "${NAME} is up and running."

logs:
	@echo "Showing logs for ${NAME}..."
	docker compose -p $(PROJECT_NAME) logs -f

stop:
	@echo "Stopping ${NAME}..."
	docker compose -p $(PROJECT_NAME) stop
	@echo "${NAME} stopped."

restart:
	@echo "Restarting ${NAME}..."
	$(MAKE_CMD) stop
	$(MAKE_CMD) up

rm:
	@echo "Removing ${NAME}..."
	docker compose -p $(PROJECT_NAME) rm -f
	@echo "${NAME} removed."

down:
	@echo "Full cleanup ${NAME}: containers, volumes, networks, data..."
	docker compose -p $(PROJECT_NAME) down -v --remove-orphans
	@echo "${NAME} cleaned."

ps:
    @echo "Status for ${NAME}:"
	docker compose -p $(PROJECT_NAME) ps

exec:
	@echo "Connecting to container ${NAME}..."
	docker compose -p $(PROJECT_NAME) exec $(SERVICE_NAME) sh

# TESTS
COVER_DIR = $(PROJECT_DIR)/resources/coverage
COVER_TMP = $(COVER_DIR)/coverage.out.tmp
COVER_OUT = $(COVER_DIR)/coverage.out
COVER_HTML = $(COVER_DIR)/coverage.html

test-cover:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVER_DIR)
	@cd $(PROJECT_DIR) && go test -coverprofile=$(COVER_TMP) ./...

	@echo "Filtering mocks..."
	@cat $(COVER_TMP) | grep -v "_mock.go" > $(COVER_OUT)

	@echo "Coverage summary:"
	@go tool cover -func=$(COVER_OUT)

	@echo "Generating HTML report..."
	@go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

	@echo "Opening coverage report..."
	@if [ "$$(uname)" = "Linux" ]; then \
		xdg-open $(COVER_HTML); \
	elif [ "$$(uname)" = "Darwin" ]; then \
		open $(COVER_HTML); \
	else \
		start $(COVER_HTML); \
	fi

	@echo "Cleaning temporary files..."
	@rm -f $(COVER_TMP)

	@echo "Coverage report generated at: $(COVER_HTML)"

# MOCKS
MOCK = $(PROJECT_DIR)/internal/repository/rpRedis/rpRedisSaveNotice/mocks

.PHONY: mock-rd-savenotice-gen mock-rd-savenotice-clean

mock-rd-savenotice-gen:
	$(MAKE_CMD) -C $(MOCK) mock-rd-savenotice-gen
mock-rd-savenotice-clean:
	$(MAKE_CMD) -C $(MOCK) mock-rd-savenotice-clean

# ========================
# APP
# ========================

# APP = providers/app

# .PHONY: app-local-build app-local-up app-local-build-up app-up app-up-no-cache app-logs app-stop app-restart app-rm app-down app-ps app-exec app-test-cover

# app-local-build:
# 	$(MAKE_CMD) -C $(APP) local-build
# app-local-up:
# 	$(MAKE_CMD) -C $(APP) local-up
# app-local-build-up:
# 	$(MAKE_CMD) -C $(APP) local-build-up
# app-up:
# 	$(MAKE_CMD) -C $(APP) up
# app-up-no-cache:
# 	$(MAKE_CMD) -C $(APP) up-no-cache
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
# app-test-cover:
# 	$(MAKE_CMD) -C $(APP) test-cover

# MOCKS

# .PHONY: app-mock-rd-savenotice-gen app-mock-rd-savenotice-clean

# app-mock-rd-savenotice-gen:
# 	$(MAKE_CMD) -C $(APP) mock-rd-savenotice-gen
# app-mock-rd-savenotice-clean:
# 	$(MAKE_CMD) -C $(APP) mock-rd-savenotice-clean