-- name: CreateChannel :one

INSERT INTO channels (id, created_at, updated_at, title, description, link, atom_link, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetChannelById :one

SELECT * FROM channels WHERE id = $1;

-- name: GetChannelByTitle :one
SELECT * FROM channels WHERE title = $1;

-- name: GetChannelByFeedId :many

SELECT * FROM channels WHERE feed_id = $1;

-- name: GetChannels :many

SELECT * FROM channels;

