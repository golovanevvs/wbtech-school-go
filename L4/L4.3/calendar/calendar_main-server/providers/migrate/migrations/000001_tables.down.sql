-- Drop trigger and function
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_events_start;
DROP INDEX IF EXISTS idx_events_start_end;

-- Drop events table
DROP TABLE IF EXISTS events;
