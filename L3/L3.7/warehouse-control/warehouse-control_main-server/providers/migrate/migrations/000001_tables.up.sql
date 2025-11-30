CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(255) NOT NULL CHECK (user_role IN ('Кладовщик', 'Менеджер', 'Аудитор')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, token)
);

-- Таблица товаров склада
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица истории изменений товаров
CREATE TABLE IF NOT EXISTS item_actions (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    action_type VARCHAR(10) NOT NULL CHECK (action_type IN ('create', 'update', 'delete')),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    changes JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Функция для создания записи в истории изменений
CREATE OR REPLACE FUNCTION log_item_changes()
RETURNS TRIGGER AS $$
DECLARE
    current_user_id INTEGER;
    current_user_name TEXT;
BEGIN
    -- Получаем информацию о пользователе из пользовательских параметров
    -- Используем NULLIF и COALESCE для безопасности
    current_user_id := NULLIF(current_setting('app.current_user_id', true), '')::INTEGER;
    current_user_name := COALESCE(current_setting('app.current_user_name', true), 'Unknown User');
    
    -- Если user_id все еще NULL, используем значение по умолчанию
    IF current_user_id IS NULL THEN
        current_user_id := 0; -- Или другое значение по умолчанию
    END IF;
    
    -- Для операции INSERT
    IF TG_OP = 'INSERT' THEN
        INSERT INTO item_actions (item_id, action_type, user_id, user_name, changes)
        VALUES (
            NEW.id,
            'create',
            current_user_id,
            current_user_name,
            jsonb_build_object(
                'name', NEW.name,
                'price', NEW.price,
                'quantity', NEW.quantity
            )
        );
        RETURN NEW;
    END IF;
    
    -- Для операции UPDATE
    IF TG_OP = 'UPDATE' THEN
        INSERT INTO item_actions (item_id, action_type, user_id, user_name, changes)
        VALUES (
            NEW.id,
            'update',
            current_user_id,
            current_user_name,
            jsonb_build_object(
                'name', jsonb_build_object('old', OLD.name, 'new', NEW.name),
                'price', jsonb_build_object('old', OLD.price, 'new', NEW.price),
                'quantity', jsonb_build_object('old', OLD.quantity, 'new', NEW.quantity)
            )
        );
        RETURN NEW;
    END IF;
    
    -- Для операции DELETE
    IF TG_OP = 'DELETE' THEN
        INSERT INTO item_actions (item_id, action_type, user_id, user_name, changes)
        VALUES (
            OLD.id,
            'delete',
            current_user_id,
            current_user_name,
            jsonb_build_object(
                'name', OLD.name,
                'price', OLD.price,
                'quantity', OLD.quantity
            )
        );
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Создание триггера для автоматического логирования изменений
DROP TRIGGER IF EXISTS item_changes_trigger ON items;
CREATE TRIGGER item_changes_trigger
    AFTER INSERT OR UPDATE OR DELETE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_changes();

-- Индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_item_actions_item_id ON item_actions(item_id);
CREATE INDEX IF NOT EXISTS idx_item_actions_created_at ON item_actions(created_at);
CREATE INDEX IF NOT EXISTS idx_item_actions_user_id ON item_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
CREATE INDEX IF NOT EXISTS idx_items_updated_at ON items(updated_at);