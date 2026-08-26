-- name: GetFeed :one
SELECT * FROM feeds WHERE id = $1;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, url, name)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = $1;

-- name: UpdateFeed :one
UPDATE feeds SET updated_at = $2, url = $3, name = $4
WHERE id = $1
RETURNING *;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListFeeds :many
SELECT *, COUNT(*) OVER() AS total_count FROM feeds
ORDER BY name
LIMIT sqlc.arg('page_size')::int OFFSET sqlc.arg('page_offset')::int;
