MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

ENV_FILE ?= $(PROJECT_DIR)/.env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

.PHONY: build up build-up

MAIN_DIR := ./cmd/$(APP_NAME)

build:
	@echo "Building app #$(APP_NAME)..."
	@cd $(PROJECT_DIR) && go build $(MAIN_DIR)
	@echo "Building app #$(APP_NAME) completed"

up:
	@echo "Running app #$(APP_NAME)..."
	@cd $(PROJECT_DIR) && start "" "/C/Program Files/Git/git-bash.exe" -c "./$(APP_NAME)"
	@sleep 1
	@echo "Running app #$(APP_NAME) completed"

build-up: build up

#============================================================

# ========================
# APP
# ========================

# app-build:
# 	mingw32-make -C $(APP) build
# app-up:
# 	mingw32-make -C $(APP) up
# app-build-up:
# 	mingw32-make -C $(APP) build-up