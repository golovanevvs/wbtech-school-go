# ========================
# APP
# ========================

APP = infra/app

app-build:
	mingw32-make -C $(APP) build
app-up:
	mingw32-make -C $(APP) up
app-build-up:
	mingw32-make -C $(APP) build-up

# ========================
# MIGRATE
# ========================

MIGRATE = infra/migrate

migrate-create:
	mingw32-make -C $(MIGRATE) create
migrate-up:
	mingw32-make -C $(MIGRATE) up
migrate-down:
	mingw32-make -C $(MIGRATE) down


# ========================
# POSTGRES
# ========================

POSTGRES = infra/postgres

# Default target
postgres-help:
	mingw32-make -C $(POSTGRES) help
# Basic Cluster Operations
postgres-up:
	mingw32-make -C $(POSTGRES) up
postgres-up-async:
	mingw32-make -C $(POSTGRES) up-async
postgres-up-sync-all:
	mingw32-make -C $(POSTGRES) up-sync-all
postgres-up-sync-one:
	mingw32-make -C $(POSTGRES) up-sync-one
postgres-down:
	mingw32-make -C $(POSTGRES) down
postgres-restart:
	mingw32-make -C $(POSTGRES) restart
postgres-status:
	mingw32-make -C $(POSTGRES) status
postgres-status-sync:
	mingw32-make -C $(POSTGRES) status-sync
postgres-logs:
	mingw32-make -C $(POSTGRES) logs
postgres-setup:
	mingw32-make -C $(POSTGRES) setup
postgres-check:
	mingw32-make -C $(POSTGRES) check
postgres-clean:
	mingw32-make -C $(POSTGRES) clean
postgres-purge:
	mingw32-make -C $(POSTGRES) purge
# Replication Mode Management
postgres-mode-async:
	mingw32-make -C $(POSTGRES) mode-async
postgres-mode-sync-all:
	mingw32-make -C $(POSTGRES) mode-sync-all
postgres-mode-sync-one:
	mingw32-make -C $(POSTGRES) mode-sync-one
postgres-mode-status:
	mingw32-make -C $(POSTGRES) mode-status
# Testing and Monitoring
postgres-test:
	mingw32-make -C $(POSTGRES) test
postgres-backup:
	mingw32-make -C $(POSTGRES) backup
postgres-monitor:
	mingw32-make -C $(POSTGRES) monitor
postgres-log-cont:
	mingw32-make -C $(POSTGRES) log-cont
postgres-list-db:
	mingw32-make -C $(POSTGRES) list-db
postgres-list-tables:
	mingw32-make -C $(POSTGRES) list-tables
postgres-show-table:
	mingw32-make -C $(POSTGRES) show-table
# Emergency Operations
postgres-promote-slave1:
	mingw32-make -C $(POSTGRES) promote-slave1
postgres-promote-slave2:
	mingw32-make -C $(POSTGRES) promote-slave2
postgres-check-master:
	mingw32-make -C $(POSTGRES) check-master
postgres-check-lag:
	mingw32-make -C $(POSTGRES) check-lag
# Auto Failover
postgres-pgauto-up:
	mingw32-make -C $(POSTGRES) pgauto-up
postgres-pgauto-down:
	mingw32-make -C $(POSTGRES) pgauto-down
postgres-pgauto-status:
	mingw32-make -C $(POSTGRES) pgauto-status
postgres-pgauto-status-sync:
	mingw32-make -C $(POSTGRES) pgauto-status-sync
postgres-pgauto-monitor:
	mingw32-make -C $(POSTGRES) pgauto-monitor
postgres-pgauto-purge:
	mingw32-make -C $(POSTGRES) pgauto-purge
postgres-clean-all:
	mingw32-make -C $(POSTGRES) clean-all