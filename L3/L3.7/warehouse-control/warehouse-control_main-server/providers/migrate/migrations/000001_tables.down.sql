DROP TRIGGER IF EXISTS item_changes_trigger ON items;
DROP FUNCTION IF EXISTS log_item_changes();

DROP INDEX IF EXISTS idx_item_actions_item_id;
DROP INDEX IF EXISTS idx_item_actions_created_at;
DROP INDEX IF EXISTS idx_item_actions_user_id;
DROP INDEX IF EXISTS idx_items_name;
DROP INDEX IF EXISTS idx_items_updated_at;

DROP TABLE IF EXISTS item_actions;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS items;