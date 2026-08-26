-- Every query is scoped by user_id: a feed follow is only ever readable or
-- deletable by the user it belongs to.

-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowByIdAndUserId :one
SELECT * FROM feed_follows WHERE id = $1 AND user_id = $2;

-- name: GetFeedsFollowsByUserId :many
SELECT * FROM feed_follows WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetFeedFollowsByUserIdAndFeedId :one
SELECT * FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

-- name: DeleteFeedFollowByIdAndUserId :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;

-- name: DeleteFeedFollowsByUserIdAndFeedId :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
