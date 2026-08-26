-- +goose Up

-- The column is filled in explicitly instead of relying on the default being
-- evaluated per row while the table is rewritten, so existing users cannot end
-- up sharing a key and colliding on the unique constraint.
ALTER TABLE users ADD COLUMN api_key VARCHAR(64);

UPDATE users SET api_key = encode(sha256(random()::text::bytea), 'hex') WHERE api_key IS NULL;

ALTER TABLE users ALTER COLUMN api_key SET DEFAULT encode(sha256(random()::text::bytea), 'hex');

ALTER TABLE users ALTER COLUMN api_key SET NOT NULL;

ALTER TABLE users ADD CONSTRAINT users_api_key_key UNIQUE (api_key);

-- +goose Down

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_api_key_key;

ALTER TABLE users DROP COLUMN api_key;
