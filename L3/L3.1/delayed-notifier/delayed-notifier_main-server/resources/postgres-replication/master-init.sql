-- Master initialization script
DO $$
BEGIN
    -- Create replication user if not exists
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'repl_user') THEN
        CREATE USER repl_user WITH REPLICATION ENCRYPTED PASSWORD 'repl_password';
    END IF;
    
    -- Create replication slots if not exist
    IF NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = 'replication_slot_slave1') THEN
        PERFORM pg_create_physical_replication_slot('replication_slot_slave1', true);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = 'replication_slot_slave2') THEN
        PERFORM pg_create_physical_replication_slot('replication_slot_slave2', true);
    END IF;
END
$$;