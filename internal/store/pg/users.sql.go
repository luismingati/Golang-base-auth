// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package pg

import (
	"context"

	"github.com/google/uuid"
)

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, email, password, username FROM users WHERE email = $1
`

type FindUserByEmailRow struct {
	ID       uuid.UUID
	Email    string
	Password string
	Username string
}

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (FindUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, findUserByEmail, email)
	var i FindUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Username,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO users (id, email, password, username)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type InsertUserParams struct {
	ID       uuid.UUID
	Email    string
	Password string
	Username string
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.Username,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1
WHERE email = $2
`

type UpdateUserPasswordParams struct {
	Password string
	Email    string
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Password, arg.Email)
	return err
}
