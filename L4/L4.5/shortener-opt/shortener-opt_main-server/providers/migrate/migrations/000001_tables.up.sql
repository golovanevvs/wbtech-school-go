CREATE TABLE IF NOT EXISTS short_url (
    id SERIAL PRIMARY KEY,
    original TEXT NOT NULL,
    short TEXT UNIQUE NOT NULL,
    custom BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS analytic (
    id SERIAL PRIMARY KEY,
    short TEXT NOT NULL,
    user_agent TEXT,
    ip INET,
    referer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE analytic
ADD CONSTRAINT fk_analytic_short
FOREIGN KEY (short) REFERENCES short_url(short)
ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_analytic_short ON analytic(short);
CREATE INDEX IF NOT EXISTS idx_analytic_created_at ON analytic(created_at);
CREATE INDEX IF NOT EXISTS idx_short_url_short ON short_url(short);
CREATE INDEX IF NOT EXISTS idx_analytic_created_at_day ON analytic(date_trunc('day',created_at));
CREATE INDEX IF NOT EXISTS idx_analytic_created_at_month ON analytic(date_trunc('month',created_at));