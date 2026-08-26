-- name: CreateItem :one
INSERT INTO items (id, title, link, guid, pubdate, published_at, seeders, leechers, downloads, infohash,
 category_id, category, size, comments, trusted, remake, description, created_at, updated_at, channel_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
RETURNING *;

-- name: GetItemById :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemByChannelIdAndGuid :one
SELECT * FROM items WHERE channel_id = $1 AND guid = $2;

-- name: GetItemByChannelId :many
SELECT * FROM items WHERE channel_id = $1;

-- name: GetItems :many
SELECT * FROM items;

-- name: ListItems :many
SELECT *, COUNT(*) OVER() AS total_count FROM items
WHERE (sqlc.narg('search')::text IS NULL OR title ILIKE '%' || sqlc.narg('search')::text || '%')
  AND (sqlc.narg('category')::text IS NULL OR category = sqlc.narg('category')::text)
  AND (sqlc.narg('min_seeders')::int IS NULL OR seeders >= sqlc.narg('min_seeders')::int)
  AND (sqlc.narg('channel_id')::uuid IS NULL OR channel_id = sqlc.narg('channel_id')::uuid)
ORDER BY
  CASE WHEN sqlc.arg('sort')::text = 'seeders' THEN seeders END DESC NULLS LAST,
  CASE WHEN sqlc.arg('sort')::text = 'oldest' THEN published_at END ASC NULLS LAST,
  CASE WHEN sqlc.arg('sort')::text = 'oldest' THEN created_at END ASC,
  published_at DESC NULLS LAST,
  created_at DESC
LIMIT sqlc.arg('page_size')::int OFFSET sqlc.arg('page_offset')::int;

-- name: ListItemCategories :many
SELECT DISTINCT category FROM items
WHERE category IS NOT NULL AND category <> ''
ORDER BY category;

-- name: UpdateItem :one
UPDATE items SET title = $2, link = $3, guid = $4, pubdate = $5, published_at = $6, seeders = $7,
 leechers = $8, downloads = $9, infohash = $10, category_id = $11, category = $12, size = $13,
 comments = $14, trusted = $15, remake = $16, description = $17, updated_at = $18, channel_id = $19
WHERE id = $1
RETURNING *;
