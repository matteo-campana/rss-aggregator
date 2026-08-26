-- +goose Up

-- pubdate keeps the raw RSS string; published_at is the parsed value the API
-- can sort on.
ALTER TABLE items ADD COLUMN published_at TIMESTAMP;

CREATE INDEX items_published_at_idx ON items (published_at DESC NULLS LAST);

-- +goose Down

DROP INDEX IF EXISTS items_published_at_idx;

ALTER TABLE items DROP COLUMN published_at;
