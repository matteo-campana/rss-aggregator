-- +goose Up

CREATE TABLE channel (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    channel_id INT UNIQUE NOT NULL,
    title VARCHAR(255),
    description TEXT,
    link VARCHAR(255),
    atom_link VARCHAR(255),
    feed_id UUID REFERENCES feeds(id)
);

-- +goose Down

DROP TABLE IF EXISTS channel;