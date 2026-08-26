-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, fullname, firstname, lastname, email)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET updated_at = $2, fullname = $3, firstname = $4, lastname = $5, email = $6
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: ListUsers :many
SELECT *, COUNT(*) OVER() AS total_count FROM users
ORDER BY created_at DESC
LIMIT sqlc.arg('page_size')::int OFFSET sqlc.arg('page_offset')::int;
