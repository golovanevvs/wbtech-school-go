MAKEFLAGS += --no-print-directory

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

.PHONY: create-db list-db list-tables show-table

NAME = ${POSTGRES_CONTAINER}
PROJECT_NAME = ${POSTGRES_CONTAINER}
SERVICE_NAME = ${POSTGRES_CONTAINER}

VORTEX = ../../../../../../../vortex

create-db:
	$(MAKE_CMD) -C $(VORTEX) postgres-create-db DB="$(POSTGRES_DB)"

list-tables:
	$(MAKE_CMD) -C $(VORTEX) postgres-list-tables DB="$(POSTGRES_DB)"

show-table:
	$(MAKE_CMD) -C $(VORTEX) postgres-show-table DB="$(POSTGRES_DB)" TABLE="$(TABLE)"

drop-db:
	$(MAKE_CMD) -C $(VORTEX) postgres-drop-db DB="$(POSTGRES_DB)"

# ========================
# POSTGRES
# ========================

# POSTGRES = postgres

# .PHONY: postgres-create-db postgres-list-tables postgres-show-table postgres-drop-db

# postgres-up:
# 	$(MAKE_CMD) -C $(POSTGRES) up
# postgres-logs:
# 	$(MAKE_CMD) -C $(POSTGRES) logs
# postgres-stop:
# 	$(MAKE_CMD) -C $(POSTGRES) stop
# postgres-restart:
# 	$(MAKE_CMD) -C $(POSTGRES) restart
# postgres-rm:
# 	$(MAKE_CMD) -C $(POSTGRES) rm
# postgres-down:
# 	$(MAKE_CMD) -C $(POSTGRES) down
# postgres-ps:
# 	$(MAKE_CMD) -C $(POSTGRES) ps
# postgres-exec:
# 	$(MAKE_CMD) -C $(POSTGRES) exec
# postgres-create-db:
# 	$(MAKE_CMD) -C $(POSTGRES) create-db DB="$(DB)"
# postgres-list-db:
# 	$(MAKE_CMD) -C $(POSTGRES) list-db"
# postgres-list-tables:
# 	$(MAKE_CMD) -C $(POSTGRES) list-tables DB="$(DB)"
# postgres-show-table:
# 	$(MAKE_CMD) -C $(POSTGRES) show-table DB="$(DB)" TABLE="$(TABLE)"
# postgres-drop-db:
# 	$(MAKE_CMD) -C $(POSTGRES) drop-db DB="$(DB)"
