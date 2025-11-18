CREATE TABLE image (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    original_path TEXT NOT NULL,
    processed_path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    format VARCHAR(10),
    operations JSONB DEFAULT '{"resize": true, "thumbnail": false, "watermark": false}'
);