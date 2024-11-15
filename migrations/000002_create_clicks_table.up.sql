CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    banner_id INTEGER NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    count INTEGER DEFAULT 1,
    CONSTRAINT fk_banner
        FOREIGN KEY (banner_id)
        REFERENCES banners(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_clicks_banner_timestamp ON clicks(banner_id, timestamp);
