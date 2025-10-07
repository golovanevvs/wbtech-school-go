CREATE TABLE IF NOT EXISTS notices(
    id SERIAL PRIMARY KEY,
    user_id INT,
    channels TEXT,
    created_at TIMESTAMP,
    sent_at TIMESTAMP,
    status TEXT
);
