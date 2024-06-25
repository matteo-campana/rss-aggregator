-- +goose Up

CREATE TABLE channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    link VARCHAR(255),
    atom_link VARCHAR(255),
    feed_id UUID REFERENCES feeds(id)
);

-- +goose Down

DROP TABLE IF EXISTS channels;