CREATE TABLE comment (
    id          SERIAL PRIMARY KEY,
    parent_id   INTEGER REFERENCES comment(id) ON DELETE CASCADE,
    text        TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);