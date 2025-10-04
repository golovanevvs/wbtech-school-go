#!/bin/bash
set -e

# Определяем корневую директорию проекта
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env"
DOCKER_COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"

echo "Project root: $PROJECT_ROOT"
echo "Env file: $ENV_FILE"
echo "Docker compose: $DOCKER_COMPOSE_FILE"

# Загружаем переменные из .env
if [ -f "$ENV_FILE" ]; then
    echo "Loading environment from $ENV_FILE"
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    echo "ERROR: .env file not found at $ENV_FILE"
    exit 1
fi

echo "Starting replication setup..."
echo "POSTGRES_USER: $POSTGRES_USER"
echo "POSTGRES_REPLICATION_USER: $POSTGRES_REPLICATION_USER"

# Функция для проверки и создания БД
ensure_database_exists() {
    echo "Checking if database exists..."
    until docker exec postgres-master psql -U ${POSTGRES_USER} -d postgres -c "\l" | grep -q "${POSTGRES_DB}"; do
        echo "Database ${POSTGRES_DB} not found, creating..."
        docker exec postgres-master psql -U ${POSTGRES_USER} -d postgres -c "CREATE DATABASE ${POSTGRES_DB};"
        sleep 2
    done
    echo "Database ${POSTGRES_DB} confirmed"
}

# Функция для временного отключения синхронной репликации
disable_sync_replication() {
    echo "Temporarily disabling synchronous replication for setup..."
    docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "ALTER SYSTEM SET synchronous_standby_names TO '';" 2>/dev/null || true
    docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT pg_reload_conf();" 2>/dev/null || true
    sleep 3
}

# Функция для восстановления режима репликации
restore_replication_mode() {
    echo "Restoring original replication mode..."
    # Режим восстановится автоматически при перезапуске контейнера
    docker-compose -f "$DOCKER_COMPOSE_FILE" restart postgres-master
    sleep 5
}

# Ждем готовности мастера
echo "Waiting for master to be ready..."
until docker exec postgres-master pg_isready -U ${POSTGRES_USER} > /dev/null 2>&1; do
    echo "Waiting for master..."
    sleep 5
done

# Создаем БД если её нет
ensure_database_exists

# Временно отключаем синхронную репликацию чтобы избежать deadlock
disable_sync_replication

# Создаем пользователя и слоты репликации
echo "Creating replication user and slots..."
docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${POSTGRES_REPLICATION_USER}') THEN
        CREATE USER ${POSTGRES_REPLICATION_USER} WITH REPLICATION ENCRYPTED PASSWORD '${POSTGRES_REPLICATION_PASSWORD}';
    END IF;
END
\$\$;" 2>/dev/null || echo "Replication user already exists or error"

docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
SELECT pg_create_physical_replication_slot('${POSTGRES_REPLICATION_SLOT_SLAVE1}', true) 
WHERE NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = '${POSTGRES_REPLICATION_SLOT_SLAVE1}');" 2>/dev/null || echo "Slot creation failed"

docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
SELECT pg_create_physical_replication_slot('${POSTGRES_REPLICATION_SLOT_SLAVE2}', true) 
WHERE NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = '${POSTGRES_REPLICATION_SLOT_SLAVE2}');" 2>/dev/null || echo "Slot creation failed"

# Настраиваем pg_hba.conf
echo "Configuring replication access..."
docker exec postgres-master bash -c "
echo 'host replication ${POSTGRES_REPLICATION_USER} 0.0.0.0/0 scram-sha-256' >> /var/lib/postgresql/data/pg_hba.conf" 2>/dev/null || echo "pg_hba.conf configuration failed"
docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT pg_reload_conf();" 2>/dev/null || echo "Config reload failed"
sleep 3

# Настраиваем слейвы
for SLAVE_NUM in 1 2; do
    SLAVE_NAME="postgres-slave-${SLAVE_NUM}"       # имя контейнера
    SYNC_NAME="slave${SLAVE_NUM}"                  # имя для synchronous_standby_names
    VOLUME_NAME="postgres-replication_postgres-slave-${SLAVE_NUM}-data"

    echo "Setting up ${SLAVE_NAME} with application_name=${SYNC_NAME}..."

    # Останавливаем слейв и чистим старые данные
    docker-compose -f "$DOCKER_COMPOSE_FILE" stop ${SLAVE_NAME} 2>/dev/null || true
    docker volume rm -f ${VOLUME_NAME} 2>/dev/null || true
    docker volume create ${VOLUME_NAME} 2>/dev/null || true

    # Создаем базовый бэкап с правильным application_name
    docker run --rm --network postgres-replication_postgres-replication \
      -v ${VOLUME_NAME}:/var/lib/postgresql/data \
      -e PGPASSWORD=${POSTGRES_REPLICATION_PASSWORD} \
      postgres:latest \
      bash -c "
        rm -rf /var/lib/postgresql/data/* &&
        pg_basebackup -h postgres-master -U ${POSTGRES_REPLICATION_USER} -D /var/lib/postgresql/data -P -v -X stream -R &&
        echo \"primary_conninfo = 'host=postgres-master port=5432 user=${POSTGRES_REPLICATION_USER} password=${POSTGRES_REPLICATION_PASSWORD} application_name=${SYNC_NAME}'\" >> /var/lib/postgresql/data/postgresql.auto.conf &&
        touch /var/lib/postgresql/data/standby.signal
      "

    # Запускаем слейв
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d ${SLAVE_NAME}
    echo "${SLAVE_NAME} setup completed"
done


# Ждем когда слейвы подключатся
echo "Waiting for slaves to connect..."
for i in {1..10}; do
    CONNECTED_SLAVES=$(docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -t -c "SELECT count(*) FROM pg_stat_replication;" 2>/dev/null || echo "0")
    if [ "$CONNECTED_SLAVES" -eq 2 ]; then
        echo "All slaves connected successfully"
        break
    fi
    echo "Connected slaves: $CONNECTED_SLAVES/2, waiting..."
    sleep 2
done

# Восстанавливаем исходный режим репликации
restore_replication_mode

echo "Replication setup completed successfully!"
echo "Master: localhost:${POSTGRES_MASTER_PORT}"
echo "Slave 1: localhost:${POSTGRES_SLAVE1_PORT}"
echo "Slave 2: localhost:${POSTGRES_SLAVE2_PORT}"