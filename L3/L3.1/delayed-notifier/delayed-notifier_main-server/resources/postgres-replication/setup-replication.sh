#!/bin/bash

# Get env file path from argument or use default
ENV_FILE="../../.env"

echo "Looking for env file at: $ENV_FILE"
echo "Current directory: $(pwd)"

# Load environment variables
if [ -f "$ENV_FILE" ]; then
    echo "Loading environment from: $ENV_FILE"
    set -a
    source "$ENV_FILE"
    set +a
else
    echo "Error: Environment file $ENV_FILE not found!"
    echo "Available files in current directory:"
    ls -la
    echo "Available files in parent directory:"
    ls -la ../
    exit 1
fi

echo "Starting replication setup..."

# Wait for master to be ready
echo "Waiting for master to be ready..."
until docker exec postgres-master pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; do
    echo "Waiting for master..."
    sleep 2
done

# Wait a bit more for full initialization
sleep 5

# Setup slaves
for SLAVE in postgres-slave-1 postgres-slave-2; do
    echo "Setting up ${SLAVE}..."
    
    # Stop the slave container
    docker-compose --env-file ${ENV_FILE} -f docker-compose-postgres.yml stop ${SLAVE}
    
    # Remove existing data
    docker exec ${SLAVE} rm -rf /var/lib/postgresql/data/*
    
    # Create base backup
    echo "Creating base backup for ${SLAVE}..."
    docker exec postgres-master pg_basebackup -D /tmp/${SLAVE} -U ${POSTGRES_USER} -P -v -R -X stream -c fast
    
    # Copy backup to slave
    echo "Copying backup to ${SLAVE}..."
    docker cp postgres-master:/tmp/${SLAVE}/. ${SLAVE}:/var/lib/postgresql/data/
    
    # Set proper replication slot name
    if [ "${SLAVE}" = "postgres-slave-1" ]; then
        SLOT_NAME="${POSTGRES_REPLICATION_SLOT_SLAVE1}"
    else
        SLOT_NAME="${POSTGRES_REPLICATION_SLOT_SLAVE2}"
    fi
    
    # Update postgresql.conf with correct slot name
    echo "Setting replication slot to: ${SLOT_NAME}"
    docker exec ${SLAVE} bash -c "echo \"primary_slot_name = '${SLOT_NAME}'\" >> /var/lib/postgresql/data/postgresql.conf"
    
    # Cleanup
    docker exec postgres-master rm -rf /tmp/${SLAVE}
    
    # Start the slave
    docker-compose --env-file ${ENV_FILE} -f docker-compose-postgres.yml start ${SLAVE}
    
    echo "${SLAVE} setup completed with slot: ${SLOT_NAME}"
done

echo "Replication setup completed!"
echo "Run 'make status' to check replication status"