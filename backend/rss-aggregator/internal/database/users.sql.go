// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, fullname, firstname, lastname, email)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, fullname, firstname, lastname, email
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Fullname  string
	Firstname sql.NullString
	Lastname  sql.NullString
	Email     sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Fullname,
		arg.Firstname,
		arg.Lastname,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Fullname,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, updated_at, fullname, firstname, lastname, email FROM users WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Fullname,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, created_at, updated_at, fullname, firstname, lastname, email FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Fullname,
			&i.Firstname,
			&i.Lastname,
			&i.Email,
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

const updateUser = `-- name: UpdateUser :one
UPDATE users SET updated_at = $2, fullname = $3, firstname = $4, lastname = $5, email = $6
WHERE id = $1
RETURNING id, created_at, updated_at, fullname, firstname, lastname, email
`

type UpdateUserParams struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Fullname  string
	Firstname sql.NullString
	Lastname  sql.NullString
	Email     sql.NullString
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.UpdatedAt,
		arg.Fullname,
		arg.Firstname,
		arg.Lastname,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Fullname,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
	)
	return i, err
}
