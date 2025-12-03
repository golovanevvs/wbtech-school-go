# ========================
# SETUP
# ========================

ENV_FILE ?= .env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export ENV PROJECT_NAME
endif

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

# ========================
# APP
# ========================

APP = providers/app

.PHONY: app-build app-up app-build-up app-test-cover app-create-network

app-build:
	$(MAKE_CMD) -C $(APP) build
app-up:
	$(MAKE_CMD) -C $(APP) up
app-build-up:
	$(MAKE_CMD) -C $(APP) build-up
app-test-cover:
	$(MAKE_CMD) -C $(APP) test-cover