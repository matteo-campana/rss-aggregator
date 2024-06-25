-- name: CreateItem :one
INSERT INTO items (id, title, link, guid, pubdate, seeders, leechers, downloads, infohash, category_id, category, size, comments, trusted, remake, description, created_at, updated_at, channel_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
RETURNING *;

-- name: GetItemById :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemByGuid :one
SELECT * FROM items WHERE guid = $1;

-- name: GetItemByChannelId :many
SELECT * FROM items WHERE channel_id = $1;

-- name: GetItems :many
SELECT * FROM items;


