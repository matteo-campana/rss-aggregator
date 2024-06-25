-- +goose Up

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255),
    link VARCHAR(255),
    guid VARCHAR(255) UNIQUE NOT NULL,
    pubdate VARCHAR(255),
    seeders INT,
    leechers INT,
    downloads INT,
    infohash VARCHAR(255),
    category_id VARCHAR(255),
    category VARCHAR(255),
    size VARCHAR(255),
    comments INT,
    trusted VARCHAR(255),
    remake VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE IF EXISTS items;