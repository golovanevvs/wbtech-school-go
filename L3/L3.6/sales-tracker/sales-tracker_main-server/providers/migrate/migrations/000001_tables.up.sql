CREATE TABLE IF NOT EXISTS sales_records (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL,
    category VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    amount INTEGER NOT NULL CHECK (amount > 0)
);