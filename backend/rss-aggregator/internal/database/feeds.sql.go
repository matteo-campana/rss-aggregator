// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, url, name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, url, name
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Url       string
	Name      string
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Url,
		arg.Name,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.Name,
	)
	return i, err
}

const deleteFeed = `-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1
`

func (q *Queries) DeleteFeed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeed, id)
	return err
}

const getFeed = `-- name: GetFeed :one
SELECT id, created_at, updated_at, url, name FROM feeds WHERE id = $1
`

func (q *Queries) GetFeed(ctx context.Context, id uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.Name,
	)
	return i, err
}

const getFeedByName = `-- name: GetFeedByName :one
SELECT id, created_at, updated_at, url, name FROM feeds WHERE name = $1
`

func (q *Queries) GetFeedByName(ctx context.Context, name string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByName, name)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.Name,
	)
	return i, err
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, url, name FROM feeds WHERE url = $1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.Name,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, created_at, updated_at, url, name FROM feeds
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Url,
			&i.Name,
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

const updateFeed = `-- name: UpdateFeed :one
UPDATE feeds SET updated_at = $2, url = $3, name = $4
WHERE id = $1
RETURNING id, created_at, updated_at, url, name
`

type UpdateFeedParams struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Url       string
	Name      string
}

func (q *Queries) UpdateFeed(ctx context.Context, arg UpdateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, updateFeed,
		arg.ID,
		arg.UpdatedAt,
		arg.Url,
		arg.Name,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.Name,
	)
	return i, err
}
