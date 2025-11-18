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

.PHONY: app-build app-up app-build-up app-test-cover

app-build:
	$(MAKE_CMD) -C $(APP) build
app-up:
	$(MAKE_CMD) -C $(APP) up
app-build-up:
	$(MAKE_CMD) -C $(APP) build-up
app-test-cover:
	$(MAKE_CMD) -C $(APP) test-cover

# MOCKS

.PHONY: app-mock-rd-savenotice-gen app-mock-rd-savenotice-clean

app-mock-rd-savenotice-gen:
	$(MAKE_CMD) -C $(APP) mock-rd-savenotice-gen
app-mock-rd-savenotice-clean:
	$(MAKE_CMD) -C $(APP) mock-rd-savenotice-clean

# ========================
# MIGRATE
# ========================

MIGRATE = providers/migrate

migrate-create:
	$(MAKE_CMD) -C $(MIGRATE) create
migrate-up:
	$(MAKE_CMD) -C $(MIGRATE) up
migrate-down:
	$(MAKE_CMD) -C $(MIGRATE) down


# ========================
# POSTGRES
# ========================

POSTGRES = providers/postgres

.PHONY: postgres-create-db postgres-list-tables postgres-show-table postgres-drop-db

postgres-create-db:
	$(MAKE_CMD) -C $(POSTGRES) create-db
postgres-list-tables:
	$(MAKE_CMD) -C $(POSTGRES) list-tables
postgres-show-table:
	$(MAKE_CMD) -C $(POSTGRES) show-table TABLE="$(TABLE)"
postgres-drop-db:
	$(MAKE_CMD) -C $(POSTGRES) drop-db

# ========================
# KAFKA
# ========================

KAFKA = providers/kafka

.PHONY: kafka-create-topic

kafka-create-topic:
	$(MAKE_CMD) -C $(KAFKA) create-topic