-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fullname TEXT NOT NULL,
    firstname TEXT,
    lastname TEXT,
    email TEXT
);

-- +goose Down

DROP TABLE IF EXISTS users;