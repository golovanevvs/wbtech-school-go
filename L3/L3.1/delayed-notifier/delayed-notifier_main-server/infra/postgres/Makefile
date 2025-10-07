MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

# Load environment variables
ENV_FILE ?= $(PROJECT_DIR)/.env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

# Project paths
POSTGRES_DIR = replication
PGAUTO_DIR = pg-auto-failover

# Replication mode (async, sync-all, sync-one)
MODE ?= async

.PHONY: help up up-async up-sync-all up-sync-one down restart status status-sync logs-cont setup check clean purge test backup monitor \
        list-db list-tables show-table \
        promote-slave1 promote-slave2 check-master check-lag \
        pgauto-up pgauto-down pgauto-status pgauto-status-sync pgauto-monitor pgauto-purge \
        clean-all mode-async mode-sync-all mode-sync-one mode-status

# Default target
help:
	@echo "PostgreSQL Replication Management"
	@echo "Using env file: $(ENV_FILE)"
	@echo "Current mode: $(MODE)"
	@echo ""
	@echo "Available commands:"
	@echo "  make up           - Start PostgreSQL cluster (master + 2 slaves)"
	@echo "  make up-async     - Start cluster in ASYNC mode"
	@echo "  make up-sync-all  - Start cluster in SYNC-ALL mode"
	@echo "  make up-sync-one  - Start cluster in SYNC-ONE mode"
	@echo "  make down         - Stop PostgreSQL cluster"
	@echo "  make restart      - Restart PostgreSQL cluster"
	@echo "  make status       - Show replication status"
	@echo "  make status-sync  - Show detailed sync status on master"
	@echo "  make logs-cont    - Show logs for a specific container (set NAME=<container>)"
	@echo "  make setup        - Setup replication (run after first start)"
	@echo "  make check        - Check replication health"
	@echo "  make clean        - Remove containers and networks"
	@echo "  make purge        - Remove containers, networks and volumes"
	@echo "  make test         - Test replication with sample data"
	@echo "  make backup       - Create backup from master"
	@echo "  make monitor      - Real-time replication monitoring"
	@echo "  make list-db      - List databases in container (set NAME=<container>)"
	@echo "  make list-tables  - List tables in database (set NAME=<container> DB=<database>)"
	@echo "  make show-table   - Show first 10 rows of a table (set NAME=<container> DB=<database> TABLE=<table>)"
	@echo ""
	@echo "Replication Modes:"
	@echo "  make mode-async           - Switch to ASYNC replication"
	@echo "  make mode-sync-all        - Switch to SYNC replication (both slaves)"
	@echo "  make mode-sync-one        - Switch to SYNC replication (slave1 only)"
	@echo "  make mode-status          - Show current replication mode"
	@echo "  make setup-with-mode MODE=<mode> - Setup with specific mode"
	@echo ""
	@echo "Emergency Operations:"
	@echo "  make promote-slave1       - Promote slave1 to master (emergency)"
	@echo "  make promote-slave2       - Promote slave2 to master (emergency)"
	@echo "  make check-master         - Check which node is master"
	@echo "  make check-lag            - Check replication lag"
	@echo ""
	@echo "Auto Failover:"
	@echo "  make pgauto-up            - Start pg_auto_failover cluster"
	@echo "  make pgauto-down          - Stop pg_auto_failover cluster"
	@echo "  make pgauto-status        - Show pg_auto_failover status"
	@echo "  make pgauto-status-sync   - Show replication status of pg_auto_failover"
	@echo "  make pgauto-monitor       - Real-time pg_auto_failover monitor"
	@echo "  make pgauto-purge         - Remove all pg_auto_failover data"
	@echo "  make clean-all            - Stop all clusters"


# ========================
# Basic Cluster Operations
# ========================

# Start the cluster
up:
	@echo "Starting PostgreSQL replication cluster..."
	@echo "Using env file: $(ENV_FILE)"
	@echo "Mode: $(MODE)"
	docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml up -d
	@echo "Cluster started. Run 'make setup' to configure replication."

up-async: up setup
	@echo "Starting PostgreSQL replication cluster in ASYNC mode..."

up-sync-all: up setup mode-sync-all
	@echo "Starting PostgreSQL replication cluster in SYNC-ALL mode..."

up-sync-one: up setup mode-sync-one
	@echo "Starting PostgreSQL replication cluster in SYNC-ONE mode..."

# Stop the cluster
down:
	@echo "Stopping PostgreSQL replication cluster..."
	docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml down

# Restart the cluster
restart: down up

# Show cluster status
status:
	@echo "=== Container Status ==="
	docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml ps
	@echo ""
	@echo "=== Master Replication Status ==="
	@docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT client_addr, application_name, state, sync_state FROM pg_stat_replication;" 2>/dev/null || echo "Master not ready or no replicas connected"
	@echo ""
	@echo "=== Replication Mode ==="
	@docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SHOW synchronous_standby_names;" 2>/dev/null || echo "Master not ready"
	@echo ""
	@echo "=== Slave Status ==="
	@for slave in $(POSTGRES_SLAVE1_CONTAINER) $(POSTGRES_SLAVE2_CONTAINER); do \
		echo "$$slave:"; \
		docker exec $$slave psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT pg_is_in_recovery(), pg_last_wal_receive_lsn(), pg_last_wal_replay_lsn();" 2>/dev/null || echo "  Not ready"; \
		echo ""; \
	done

status-sync:
	@echo "=== Replication Sync Status ==="
	@docker exec -it $(POSTGRES_MASTER_CONTAINER) \
		psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-c "SELECT application_name, state, sync_state FROM pg_stat_replication;"

# Show logs
logs:
	docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml logs -f

# Setup replication (run after first start)
setup:
	@echo "Setting up replication..."
	@cd $(POSTGRES_DIR) && bash setup-replication.sh
	@echo "Replication setup completed!"

# Check replication health
check:
	@echo "=== Replication Health Check ==="
	@echo "Master:"
	@docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT count(*) as connected_slaves FROM pg_stat_replication;" 2>/dev/null || echo "Master not ready"
	@echo ""
	@echo "Slaves:"
	@for slave in $(POSTGRES_SLAVE1_CONTAINER) $(POSTGRES_SLAVE2_CONTAINER); do \
		echo "$$slave:"; \
		docker exec $$slave psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT now() - pg_last_xact_replay_timestamp() as replication_lag;" 2>/dev/null || echo "  Not ready"; \
	done

# Clean up containers and networks
clean:
	docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml down

# Remove everything including volumes
purge:
	@echo "WARNING: This will remove ALL data including databases!"
	@read -p "Are you sure? (y/N) " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "Stopping and removing containers..."; \
		docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml down; \
		echo "Removing volumes..."; \
		docker volume rm -f $(POSTGRES_REPLICATION_PROJECT_NAME)_master-data $(POSTGRES_REPLICATION_PROJECT_NAME)_slave-1-data $(POSTGRES_REPLICATION_PROJECT_NAME)_slave-2-data 2>/dev/null || true; \
		echo "Removing network..."; \
		docker network rm $(POSTGRES_REPLICATION_PROJECT_NAME)_network 2>/dev/null || true; \
		echo "All data removed."; \
	else \
		echo "Purge cancelled."; \
	fi

# ========================
# Replication Mode Management
# ========================

# Switch to ASYNC replication mode
mode-async:
	@echo "Switching to ASYNC replication mode..."
	@docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml -f $(POSTGRES_DIR)/docker-compose-async.yml up -d postgres-master
	@sleep 3
	@echo "Async mode activated"
	mingw32-make mode-status

# Switch to SYNC-ALL replication mode (both slaves)
mode-sync-all:
	@echo "Switching to SYNC-ALL replication mode (both slaves)..."
	@docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml -f $(POSTGRES_DIR)/docker-compose-sync-all.yml up -d postgres-master
	@sleep 3
	@echo "Sync-all mode activated"
	mingw32-make mode-status

# Switch to SYNC-ONE replication mode (slave1 only)
mode-sync-one:
	@echo "Switching to SYNC-ONE replication mode (slave1 only)..."
	@docker-compose --env-file $(ENV_FILE) -f $(POSTGRES_DIR)/docker-compose.yml -f $(POSTGRES_DIR)/docker-compose-sync-one.yml up -d postgres-master
	@sleep 3
	@echo "Sync-one mode activated"
	mingw32-make mode-status

# Show current replication mode
mode-status:
	@echo "=== Current Replication Mode ==="
	@docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SHOW synchronous_standby_names;" 2>/dev/null || echo "Master not ready"

# ========================
# Testing and Monitoring
# ========================

# Test replication with sample data
test:
	@echo "Testing replication..."
	@echo "Creating test data on master..."
	@docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c " \
		CREATE TABLE IF NOT EXISTS replication_test ( \
			id serial PRIMARY KEY, \
			data text, \
			created_at timestamptz DEFAULT now() \
		); \
		INSERT INTO replication_test (data) VALUES ('test_data_' || generate_series(1,5)); \
		SELECT count(*) as records_created FROM replication_test;" 2>/dev/null || echo "Master not ready"
	@echo ""
	@echo "Checking data on slaves:"
	@for slave in $(POSTGRES_SLAVE1_CONTAINER) $(POSTGRES_SLAVE2_CONTAINER); do \
		echo "$$slave:"; \
		docker exec $$slave psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT count(*) as records_replicated FROM replication_test;" 2>/dev/null || echo "  Not ready"; \
	done

# Create backup from master
backup:
	@echo "Creating backup from master..."
	@mkdir -p backups
	@docker exec $(POSTGRES_MASTER_CONTAINER) pg_dump -U ${POSTGRES_USER} -d ${POSTGRES_DB} -F c -f /tmp/backup.dump
	@docker cp $(POSTGRES_MASTER_CONTAINER):/tmp/backup.dump backups/backup_$$(date +%Y%m%d_%H%M%S).dump
	@docker exec $(POSTGRES_MASTER_CONTAINER) rm -f /tmp/backup.dump
	@echo "Backup created in backups/ directory"

# Monitor replication in real-time
monitor:
	@echo "Starting real-time replication monitor (Ctrl+C to stop)..."
	@while true; do \
		clear; \
		echo "=== PostgreSQL Replication Monitor ==="; \
		echo "Time: $$(date)"; \
		echo ""; \
		echo "Replication Mode:"; \
		docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SHOW synchronous_standby_names;" 2>/dev/null | head -3 | tail -1 || echo "  Master not ready"; \
		echo ""; \
		echo "Active Slaves:"; \
		docker exec $(POSTGRES_MASTER_CONTAINER) psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT count(*) as active_slaves FROM pg_stat_replication;" 2>/dev/null | head -3 | tail -1 || echo "  Master not ready"; \
		echo ""; \
		echo "Slave Status:"; \
		for slave in $(POSTGRES_SLAVE1_CONTAINER) $(POSTGRES_SLAVE2_CONTAINER); do \
			echo "$$slave:"; \
			docker exec $$slave psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT pg_is_in_recovery(), now() - pg_last_xact_replay_timestamp() as lag;" 2>/dev/null | grep -v row | grep -v "-" | head -1 || echo "  Not ready"; \
		done; \
		sleep 5; \
	done

logs-cont:
	@echo "Showing logs for container \"$(NAME)\"..."
	@docker logs -f $(NAME)

list-db:
	@echo "Listing databases of container \"$(NAME)\"..."
	@docker exec -it $(NAME) psql -U $(POSTGRES_USER) -l

list-tables:
	@echo "Listing tables in database \"$(DB)\" of container \"$(NAME)\"..."
	@docker exec -it $(NAME) psql -U $(POSTGRES_USER) -d $(DB) -c "\dt"

show-table:
	@echo "Showing 10 rows of table \"$(TABLE)\" of database \"$(DB)\" of container \"$(NAME)\"..."
	@docker exec -it $(NAME) psql -U $(POSTGRES_USER) -d $(DB) \
	-c "SELECT * FROM $(TABLE) LIMIT 10"

# ========================
# Emergency Operations
# ========================

promote-slave1:
	@echo "EMERGENCY: Promoting slave1 to master..."
	@docker exec $(POSTGRES_SLAVE1_CONTAINER) pg_ctl promote -D /var/lib/postgresql/data
	@echo "Slave1 promoted to master."
	@echo "Update your applications to connect to slave1 (port ${POSTGRES_SLAVE1_PORT}) for writes."

promote-slave2:
	@echo "EMERGENCY: Promoting slave2 to master..."
	@docker exec $(POSTGRES_SLAVE2_CONTAINER) pg_ctl promote -D /var/lib/postgresql/data
	@echo "Slave2 promoted to master."
	@echo "Update your applications to connect to slave2 (port ${POSTGRES_SLAVE2_PORT}) for writes."

check-master:
	@echo "Checking which nodes are in recovery mode (slaves)..."
	@for node in $(POSTGRES_MASTER_CONTAINER) $(POSTGRES_SLAVE1_CONTAINER) $(POSTGRES_SLAVE2_CONTAINER); do \
		echo -n "$$node: "; \
		docker exec $$node psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -t -c "SELECT pg_is_in_recovery();" 2>/dev/null | tr -d '[:space:]' && echo " (in recovery = slave)" || echo " (not in recovery = master)"; \
	done

check-lag:
	@echo "Replication lag check:"
	@for slave in $(POSTGRES_SLAVE1_CONTAINER) $(POSTGRES_SLAVE2_CONTAINER); do \
		echo -n "$$slave: "; \
		docker exec $$slave psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -t -c "SELECT now() - pg_last_xact_replay_timestamp() AS replication_lag;" 2>/dev/null | tr -d '[:space:]' || echo "unavailable"; \
	done

# ========================
# Auto Failover
# ========================

pgauto-up:
	@echo "Starting PG_AUTO_FAILOVER cluster..."
	docker-compose --env-file $(ENV_FILE) -f $(PGAUTO_DIR)/docker-compose-pgauto.yml up -d
	@echo "pg_auto_failover cluster starting..."
	@echo "Monitor UI: http://localhost:3000"
	@echo "Wait 30 seconds for full initialization"

pgauto-down:
	@echo "Stopping pg_auto_failover cluster..."
	docker-compose -f $(PGAUTO_DIR)/docker-compose-pgauto.yml down

pgauto-status:
	@echo "=== pg_auto_failover Status ==="
	@docker exec pgauto-monitor pg_autoctl show state 2>/dev/null || echo "pg_auto_failover not running"

pgauto-status-sync:
	@echo "=== PG_AUTO_FAILOVER Replication Status ==="
	@docker exec -it pgauto-primary psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-c "SELECT application_name, state, sync_state FROM pg_stat_replication;" 2>/dev/null || \
		echo "Primary not ready or no replicas connected"

pgauto-monitor:
	@echo "=== pg_auto_failover Real-time Monitor ==="
	@echo "Web UI: http://localhost:3000"
	@while true; do \
		clear; \
		echo "Time: $$(date)"; \
		docker exec pgauto-monitor pg_autoctl show state 2>/dev/null || echo "pg_auto_failover not available"; \
		sleep 5; \
	done

pgauto-purge:
	@echo "WARNING: This will remove ALL pg_auto_failover data!"
	@read -p "Are you sure? (y/N) " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "Stopping and removing pg_auto_failover containers..."; \
		docker-compose --env-file $(ENV_FILE) -f $(PGAUTO_DIR)/docker-compose-pgauto.yml down; \
		echo "Removing pg_auto_failover volumes..."; \
		docker volume rm -f pg-auto-failover_pgauto-monitor-data pg-auto-failover_pgauto-primary-data pg-auto-failover_pgauto-secondary-data 2>/dev/null || true; \
		echo "Removing pg_auto_failover network..."; \
		docker network rm pg-auto-failover_pgauto-network 2>/dev/null || true; \
		echo "All pg_auto_failover data removed."; \
	else \
		echo "Purge cancelled."; \
	fi

# Clean everything
clean-all:
	@echo "Cleaning ALL PostgreSQL clusters..."
	mingw32-make down
	mingw32-make pgauto-down
	@echo "All clusters stopped"

#============================================================

# POSTGRES = infra/postgres

# ========================
# POSTGRES
# ========================

# # Default target
# postgres-help:
# 	mingw32-make -C $(POSTGRES) help
# # Basic Cluster Operations
# postgres-up:
# 	mingw32-make -C $(POSTGRES) up
# postgres-up-async:
# 	mingw32-make -C $(POSTGRES) up-async
# postgres-up-sync-all:
# 	mingw32-make -C $(POSTGRES) up-sync-all
# postgres-up-sync-one:
# 	mingw32-make -C $(POSTGRES) up-sync-one
# postgres-down:
# 	mingw32-make -C $(POSTGRES) down
# postgres-restart:
# 	mingw32-make -C $(POSTGRES) restart
# postgres-status:
# 	mingw32-make -C $(POSTGRES) status
# postgres-status-sync:
# 	mingw32-make -C $(POSTGRES) status-sync
# postgres-logs:
# 	mingw32-make -C $(POSTGRES) logs
# postgres-setup:
# 	mingw32-make -C $(POSTGRES) setup
# postgres-check:
# 	mingw32-make -C $(POSTGRES) check
# postgres-clean:
# 	mingw32-make -C $(POSTGRES) clean
# postgres-purge:
# 	mingw32-make -C $(POSTGRES) purge
# # Replication Mode Management
# postgres-mode-async:
# 	mingw32-make -C $(POSTGRES) mode-async
# postgres-mode-sync-all:
# 	mingw32-make -C $(POSTGRES) mode-sync-all
# postgres-mode-sync-one:
# 	mingw32-make -C $(POSTGRES) mode-sync-one
# postgres-mode-status:
# 	mingw32-make -C $(POSTGRES) mode-status
# # Testing and Monitoring
# postgres-test:
# 	mingw32-make -C $(POSTGRES) test
# postgres-backup:
# 	mingw32-make -C $(POSTGRES) backup
# postgres-monitor:
# 	mingw32-make -C $(POSTGRES) monitor
# postgres-log-cont:
# 	mingw32-make -C $(POSTGRES) log-cont
# postgres-list-db:
# 	mingw32-make -C $(POSTGRES) list-db
# postgres-list-tables:
# 	mingw32-make -C $(POSTGRES) list-tables
# postgres-show-table:
# 	mingw32-make -C $(POSTGRES) show-table
# # Emergency Operations
# postgres-promote-slave1:
# 	mingw32-make -C $(POSTGRES) promote-slave1
# postgres-promote-slave2:
# 	mingw32-make -C $(POSTGRES) promote-slave2
# postgres-check-master:
# 	mingw32-make -C $(POSTGRES) check-master
# postgres-check-lag:
# 	mingw32-make -C $(POSTGRES) check-lag
# # Auto Failover
# postgres-pgauto-up:
# 	mingw32-make -C $(POSTGRES) pgauto-up
# postgres-pgauto-down:
# 	mingw32-make -C $(POSTGRES) pgauto-down
# postgres-pgauto-status:
# 	mingw32-make -C $(POSTGRES) pgauto-status
# postgres-pgauto-status-sync:
# 	mingw32-make -C $(POSTGRES) pgauto-status-sync
# postgres-pgauto-monitor:
# 	mingw32-make -C $(POSTGRES) pgauto-monitor
# postgres-pgauto-purge:
# 	mingw32-make -C $(POSTGRES) pgauto-purge
# postgres-clean-all:
# 	mingw32-make -C $(POSTGRES) clean-all