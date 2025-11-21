# SETUP
MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

ENV_FILES := ../../../../../.env .env ../redis/.env ../telegram/.env ../email/.env

$(foreach file,$(ENV_FILES), \
  $(if $(wildcard $(file)), $(eval include $(file))) \
)

export PUBLIC_HOST_DOCKER
export APP_NAME APP_TRANSPORT_HTTP_PORT APP_TRANSPORT_HTTP_PUBLIC_HOST
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
APP_NAME = ${APP_CONTAINER}
PROJECT_NAME = ${APP_CONTAINER}
SERVICE_NAME = ${APP_CONTAINER}

build:
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

up:
	@echo "Running app $(NAME)..."
# 	@cd $(PROJECT_DIR) && $(RUN_CMD) "./$(APP_NAME) $(ARGS)"
# 	@echo "Running app $(NAME) completed"
ifeq ($(OS),Windows_NT)
	start powershell -NoExit -Command "cd '$(PROJECT_DIR)'; .\$(APP_NAME).exe $(ARGS)"
else
	@cd $(PROJECT_DIR) && $(RUN_CMD) "./$(APP_NAME) $(ARGS)"
endif
	@echo "Running app $(NAME) completed"

build-up: build up

create-network:
	@echo "Creating Docker network: $(APP_CONTAINER)"
	@docker network create $(APP_CONTAINER) || echo "Network $(APP_CONTAINER) already exists"
	@echo "Docker network $(APP_CONTAINER) ready"

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

# .PHONY: app-build app-up app-build-up app-test-cover app-create-network

# app-build:
# 	$(MAKE_CMD) -C $(APP) build
# app-up:
# 	$(MAKE_CMD) -C $(APP) up
# app-build-up:
# 	$(MAKE_CMD) -C $(APP) build-up
# app-create-network:
# 	$(MAKE_CMD) -C $(APP) create-network
# app-test-cover:
# 	$(MAKE_CMD) -C $(APP) test-cover

# MOCKS

# .PHONY: app-mock-rd-savenotice-gen app-mock-rd-savenotice-clean

# app-mock-rd-savenotice-gen:
# 	$(MAKE_CMD) -C $(APP) mock-rd-savenotice-gen
# app-mock-rd-savenotice-clean:
# 	$(MAKE_CMD) -C $(APP) mock-rd-savenotice-clean