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

.PHONY: create-db list-tables show-table drop-db

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

# POSTGRES = providers/postgres

# .PHONY: postgres-create-db postgres-list-tables postgres-show-table postgres-drop-db

# postgres-create-db:
# 	$(MAKE_CMD) -C $(POSTGRES) create-db
# postgres-list-tables:
# 	$(MAKE_CMD) -C $(POSTGRES) list-tables
# postgres-show-table:
# 	$(MAKE_CMD) -C $(POSTGRES) show-table TABLE="$(TABLE)"
# postgres-drop-db:
# 	$(MAKE_CMD) -C $(POSTGRES) drop-db