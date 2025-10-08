MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

ENV_FILE ?= $(PROJECT_DIR)/.env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

.PHONY: build up build-up

MAIN_DIR := ./cmd/$(APP_NAME)

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

build:
	@echo "Building app #$(APP_NAME)..."
	@cd $(PROJECT_DIR) && go build $(MAIN_DIR)
	@echo "Building app #$(APP_NAME) completed"

up:
	@echo "Running app #$(APP_NAME)..."
	@cd $(PROJECT_DIR) && $(RUN_CMD) "./$(APP_NAME)"
	@sleep 1
	@echo "Running app #$(APP_NAME) completed"

build-up: build up

#============================================================


# ========================
# APP
# ========================

# APP = infra/app

# app-build:
# 	$(MAKE_CMD) -C $(APP) build
# app-up:
# 	$(MAKE_CMD) -C $(APP) up
# app-build-up:
# 	$(MAKE_CMD) -C $(APP) build-up