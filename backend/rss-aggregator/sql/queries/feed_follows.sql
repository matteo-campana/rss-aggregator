-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id,created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowsById :one
SELECT * FROM feed_follows WHERE id = $1;

-- name: GetFeedsFollows :many
SELECT * FROM feed_follows;

-- name: GetFeedsFollowsByUserId :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: GetFeedsFollowsByFeedId :many
SELECT * FROM feed_follows WHERE feed_id = $1;

-- name: GetFeedFollowsByUserIdAndFeedId :one
SELECT * FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE id = $1;

-- name: DeleteFeedFollowsByUserId :exec
DELETE FROM feed_follows WHERE user_id = $1;

-- name: DeleteFeedFollowsByFeedId :exec
DELETE FROM feed_follows WHERE feed_id = $1;

-- name: DeleteFeedFollowsByUserIdAndFeedId :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;