DROP TRIGGER IF EXISTS item_changes_trigger ON items;
DROP FUNCTION IF EXISTS log_item_changes();

DROP TABLE IF EXISTS item_actions;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS items;