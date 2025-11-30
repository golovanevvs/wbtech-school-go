-- Удаляем триггеры и функции
DROP TRIGGER IF EXISTS item_changes_trigger ON items;
DROP FUNCTION IF EXISTS log_item_changes();

-- Удаляем таблицы в правильном порядке (с учетом внешних ключей)
DROP TABLE IF EXISTS item_actions;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS items;