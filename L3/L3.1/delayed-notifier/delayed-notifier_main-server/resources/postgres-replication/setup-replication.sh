#!/bin/bash
set -e

# Загружаем переменные
export $(cat .env | grep -v '^#' | xargs)

echo "Starting replication setup..."

# Ждем готовности мастера
echo "Waiting for master to be ready..."
until docker exec postgres-master pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; do
    sleep 5
done

# Создаем пользователя и слоты репликации
echo "Creating replication user and slots..."
docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${POSTGRES_REPLICATION_USER}') THEN
        CREATE USER ${POSTGRES_REPLICATION_USER} WITH REPLICATION ENCRYPTED PASSWORD '${POSTGRES_REPLICATION_PASSWORD}';
    END IF;
END
\$\$;"

docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
SELECT pg_create_physical_replication_slot('${POSTGRES_REPLICATION_SLOT_SLAVE1}', true) 
WHERE NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = '${POSTGRES_REPLICATION_SLOT_SLAVE1}');"

docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
SELECT pg_create_physical_replication_slot('${POSTGRES_REPLICATION_SLOT_SLAVE2}', true) 
WHERE NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = '${POSTGRES_REPLICATION_SLOT_SLAVE2}');"

# Настраиваем pg_hba.conf
echo "Configuring replication access..."
docker exec postgres-master bash -c "
echo 'host replication ${POSTGRES_REPLICATION_USER} 0.0.0.0/0 scram-sha-256' >> /var/lib/postgresql/data/pg_hba.conf"
docker exec postgres-master psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "SELECT pg_reload_conf();"
sleep 3

# Настраиваем слейвы
for SLAVE_NUM in 1 2; do
    SLAVE_NAME="postgres-slave-${SLAVE_NUM}"
    VOLUME_NAME="postgres-replication_postgres-slave-${SLAVE_NUM}-data"
    SLOT_NAME="replication_slot_slave${SLAVE_NUM}"
    
    echo "Setting up ${SLAVE_NAME}..."
    
    # Останавливаем слейв
    docker-compose -f docker-compose-postgres.yml stop ${SLAVE_NAME}
    
    # Удаляем старые данные
    docker volume rm -f ${VOLUME_NAME} || true
    docker volume create ${VOLUME_NAME}
    
    # Создаем бэкап
    echo "Creating base backup for ${SLAVE_NAME}..."
    docker run --rm --network postgres-replication_postgres-replication \
      -v ${VOLUME_NAME}:/backup_data \
      -e PGPASSWORD=${POSTGRES_REPLICATION_PASSWORD} \
      postgres:latest \
      pg_basebackup -h postgres-master -U ${POSTGRES_REPLICATION_USER} -D /backup_data -P -v -R -X stream -c fast
    
    # Запускаем слейв
    docker-compose -f docker-compose-postgres.yml start ${SLAVE_NAME}
    
    echo "${SLAVE_NAME} setup completed"
done

echo "Replication setup completed successfully!"