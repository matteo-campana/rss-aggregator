// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id,created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, feed_id, user_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.UserID,
	)
	return i, err
}

const deleteFeedFollows = `-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE id = $1
`

func (q *Queries) DeleteFeedFollows(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollows, id)
	return err
}

const deleteFeedFollowsByFeedId = `-- name: DeleteFeedFollowsByFeedId :exec
DELETE FROM feed_follows WHERE feed_id = $1
`

func (q *Queries) DeleteFeedFollowsByFeedId(ctx context.Context, feedID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowsByFeedId, feedID)
	return err
}

const deleteFeedFollowsByUserId = `-- name: DeleteFeedFollowsByUserId :exec
DELETE FROM feed_follows WHERE user_id = $1
`

func (q *Queries) DeleteFeedFollowsByUserId(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowsByUserId, userID)
	return err
}

const deleteFeedFollowsByUserIdAndFeedId = `-- name: DeleteFeedFollowsByUserIdAndFeedId :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2
`

type DeleteFeedFollowsByUserIdAndFeedIdParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) DeleteFeedFollowsByUserIdAndFeedId(ctx context.Context, arg DeleteFeedFollowsByUserIdAndFeedIdParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowsByUserIdAndFeedId, arg.UserID, arg.FeedID)
	return err
}

const getFeedFollowsById = `-- name: GetFeedFollowsById :one
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows WHERE id = $1
`

func (q *Queries) GetFeedFollowsById(ctx context.Context, id uuid.UUID) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, getFeedFollowsById, id)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.UserID,
	)
	return i, err
}

const getFeedFollowsByUserIdAndFeedId = `-- name: GetFeedFollowsByUserIdAndFeedId :one
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows WHERE user_id = $1 AND feed_id = $2
`

type GetFeedFollowsByUserIdAndFeedIdParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) GetFeedFollowsByUserIdAndFeedId(ctx context.Context, arg GetFeedFollowsByUserIdAndFeedIdParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, getFeedFollowsByUserIdAndFeedId, arg.UserID, arg.FeedID)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.UserID,
	)
	return i, err
}

const getFeedsFollows = `-- name: GetFeedsFollows :many
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows
`

func (q *Queries) GetFeedsFollows(ctx context.Context) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedsFollows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedsFollowsByFeedId = `-- name: GetFeedsFollowsByFeedId :many
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows WHERE feed_id = $1
`

func (q *Queries) GetFeedsFollowsByFeedId(ctx context.Context, feedID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedsFollowsByFeedId, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedsFollowsByUserId = `-- name: GetFeedsFollowsByUserId :many
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows WHERE user_id = $1
`

func (q *Queries) GetFeedsFollowsByUserId(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedsFollowsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}