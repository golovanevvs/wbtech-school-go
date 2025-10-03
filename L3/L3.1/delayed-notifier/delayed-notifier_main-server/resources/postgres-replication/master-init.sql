-- This file runs on master initialization
-- Create replication user (if not exists)
CREATE USER IF NOT EXISTS repl_user WITH REPLICATION ENCRYPTED PASSWORD 'repl_password';

-- Create replication slots (if not exist)
SELECT pg_create_physical_replication_slot('replication_slot_slave1', true);
SELECT pg_create_physical_replication_slot('replication_slot_slave2', true);