-- name: GetFeed :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, url, name)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateFeed :one
UPDATE feeds SET updated_at = $2, url = $3, name = $4
WHERE id = $1
RETURNING *;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1;