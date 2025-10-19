# ========================
# Detect OS
# ========================

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

# ========================
# ALL
# ========================

all-up:
	$(MAKE_CMD) all-deps-up
	@sleep 3
	$(MAKE_CMD) app-build-up

all-deps-up: rabbitmq-up redis-up

all-down: rabbitmq-down redis-down

all-purge: rabbitmq-purge redis-purge

# ========================
# APP
# ========================

APP = providers/app

app-build:
	$(MAKE_CMD) -C $(APP) build
app-up:
	$(MAKE_CMD) -C $(APP) up
app-build-up:
	$(MAKE_CMD) -C $(APP) build-up

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

# Default target
postgres-help:
	$(MAKE_CMD) -C $(POSTGRES) help
# Basic Cluster Operations
postgres-up:
	$(MAKE_CMD) -C $(POSTGRES) up
postgres-up-async:
	$(MAKE_CMD) -C $(POSTGRES) up-async
postgres-up-sync-all:
	$(MAKE_CMD) -C $(POSTGRES) up-sync-all
postgres-up-sync-one:
	$(MAKE_CMD) -C $(POSTGRES) up-sync-one
postgres-down:
	$(MAKE_CMD) -C $(POSTGRES) down
postgres-restart:
	$(MAKE_CMD) -C $(POSTGRES) restart
postgres-status:
	$(MAKE_CMD) -C $(POSTGRES) status
postgres-status-sync:
	$(MAKE_CMD) -C $(POSTGRES) status-sync
postgres-logs:
	$(MAKE_CMD) -C $(POSTGRES) logs
postgres-setup:
	$(MAKE_CMD) -C $(POSTGRES) setup
postgres-check:
	$(MAKE_CMD) -C $(POSTGRES) check
postgres-clean:
	$(MAKE_CMD) -C $(POSTGRES) clean
postgres-purge:
	$(MAKE_CMD) -C $(POSTGRES) purge
# Replication Mode Management
postgres-mode-async:
	$(MAKE_CMD) -C $(POSTGRES) mode-async
postgres-mode-sync-all:
	$(MAKE_CMD) -C $(POSTGRES) mode-sync-all
postgres-mode-sync-one:
	$(MAKE_CMD) -C $(POSTGRES) mode-sync-one
postgres-mode-status:
	$(MAKE_CMD) -C $(POSTGRES) mode-status
# Testing and Monitoring
postgres-test:
	$(MAKE_CMD) -C $(POSTGRES) test
postgres-backup:
	$(MAKE_CMD) -C $(POSTGRES) backup
postgres-monitor:
	$(MAKE_CMD) -C $(POSTGRES) monitor
postgres-log-cont:
	$(MAKE_CMD) -C $(POSTGRES) log-cont
postgres-list-db:
	$(MAKE_CMD) -C $(POSTGRES) list-db
postgres-list-tables:
	$(MAKE_CMD) -C $(POSTGRES) list-tables
postgres-show-table:
	$(MAKE_CMD) -C $(POSTGRES) show-table
# Emergency Operations
postgres-promote-slave1:
	$(MAKE_CMD) -C $(POSTGRES) promote-slave1
postgres-promote-slave2:
	$(MAKE_CMD) -C $(POSTGRES) promote-slave2
postgres-check-master:
	$(MAKE_CMD) -C $(POSTGRES) check-master
postgres-check-lag:
	$(MAKE_CMD) -C $(POSTGRES) check-lag
# Auto Failover
postgres-pgauto-up:
	$(MAKE_CMD) -C $(POSTGRES) pgauto-up
postgres-pgauto-down:
	$(MAKE_CMD) -C $(POSTGRES) pgauto-down
postgres-pgauto-status:
	$(MAKE_CMD) -C $(POSTGRES) pgauto-status
postgres-pgauto-status-sync:
	$(MAKE_CMD) -C $(POSTGRES) pgauto-status-sync
postgres-pgauto-monitor:
	$(MAKE_CMD) -C $(POSTGRES) pgauto-monitor
postgres-pgauto-purge:
	$(MAKE_CMD) -C $(POSTGRES) pgauto-purge
postgres-clean-all:
	$(MAKE_CMD) -C $(POSTGRES) clean-all

# ========================
# RABBITMQ
# ========================

RABBITMQ = providers/rabbitmq

rabbitmq-up:
	$(MAKE_CMD) -C $(RABBITMQ) up
rabbitmq-down:
	$(MAKE_CMD) -C $(RABBITMQ) down
rabbitmq-restart:
	$(MAKE_CMD) -C $(RABBITMQ) restart
rabbitmq-logs:
	$(MAKE_CMD) -C $(RABBITMQ) logs
rabbitmq-status:
	$(MAKE_CMD) -C $(RABBITMQ) status
rabbitmq-shell:
	$(MAKE_CMD) -C $(RABBITMQ) shell
rabbitmq-purge:
	$(MAKE_CMD) -C $(RABBITMQ) purge

# ========================
# Ngrok
# ========================

NGROK = providers/ngrok

ngrok-build:
	$(MAKE_CMD) -C $(NGROK) build
ngrok-up:
	$(MAKE_CMD) -C $(NGROK) up
ngrok-url:
	$(MAKE_CMD) -C $(NGROK) url
ngrok-down:
	$(MAKE_CMD) -C $(NGROK) down
ngrok-webhook:
	$(MAKE_CMD) -C $(NGROK) webhook

# ========================
# REDIS
# ========================

REDIS = providers/redis

redis-up:
	$(MAKE_CMD) -C $(REDIS) up
redis-stop:
	$(MAKE_CMD) -C $(REDIS) stop
redis-down:
	$(MAKE_CMD) -C $(REDIS) down
redis-purge:
	$(MAKE_CMD) -C $(REDIS) purge
redis-logs:
	$(MAKE_CMD) -C $(REDIS) logs

# ========================
# GRAFANA
# ========================

GRAFANA = providers/grafana

grafana-up:
	$(MAKE_CMD) -C $(GRAFANA) up
grafana-stop:
	$(MAKE_CMD) -C $(GRAFANA) stop
grafana-restart:
	$(MAKE_CMD) -C $(GRAFANA) restart
grafana-down:
	$(MAKE_CMD) -C $(GRAFANA) down
grafana-purge:
	$(MAKE_CMD) -C $(GRAFANA) purge
grafana-logs:
	$(MAKE_CMD) -C $(GRAFANA) logs
grafana-ps:
	$(MAKE_CMD) -C $(GRAFANA) ps
grafana-version:
	$(MAKE_CMD) -C $(GRAFANA) version
grafana-check-data:
	$(MAKE_CMD) -C $(GRAFANA) check-data