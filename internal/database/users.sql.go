// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps(body, user_id)
VALUES(
    $1,
    $2
) RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const createRefreshTokens = `-- name: CreateRefreshTokens :one
INSERT INTO refresh_tokens(token, user_id, expires_at, revoked_at)
VALUES(
    $1,
    $2,
    $3,
    $4
) RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokensParams struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

func (q *Queries) CreateRefreshTokens(ctx context.Context, arg CreateRefreshTokensParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshTokens,
		arg.Token,
		arg.UserID,
		arg.ExpiresAt,
		arg.RevokedAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users(email, hashed_password)
VALUES(
    $1,
    $2
) RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteChirpByIdForUser = `-- name: DeleteChirpByIdForUser :one
WITH delete_attempt AS (
    DELETE FROM chirps
    WHERE chirps.id = $1 AND chirps.user_id = $2
    RETURNING chirps.id, chirps.user_id
)
SELECT 
    COALESCE((SELECT delete_attempt.user_id FROM delete_attempt), (SELECT chirps.user_id FROM chirps WHERE chirps.id = $1)) AS user_id
`

type DeleteChirpByIdForUserParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteChirpByIdForUser(ctx context.Context, arg DeleteChirpByIdForUserParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, deleteChirpByIdForUser, arg.ID, arg.UserID)
	var user_id interface{}
	err := row.Scan(&user_id)
	return user_id, err
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at ASC
`

func (q *Queries) GetAllChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
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

const getAllChirpsByAuthor = `-- name: GetAllChirpsByAuthor :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetAllChirpsByAuthor(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsByAuthor, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
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

const getAllChirpsByAuthorDESC = `-- name: GetAllChirpsByAuthorDESC :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetAllChirpsByAuthorDESC(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsByAuthorDESC, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
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

const getAllChirpsDESC = `-- name: GetAllChirpsDESC :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at DESC
`

func (q *Queries) GetAllChirpsDESC(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsDESC)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
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

const getChirpByID = `-- name: GetChirpByID :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE id = $1
`

func (q *Queries) GetChirpByID(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirpByID, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
where email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT user_id, revoked_at from refresh_tokens
WHERE token = $1
`

type GetUserFromRefreshTokenRow struct {
	UserID    uuid.UUID
	RevokedAt sql.NullTime
}

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (GetUserFromRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i GetUserFromRefreshTokenRow
	err := row.Scan(&i.UserID, &i.RevokedAt)
	return i, err
}

const resetTables = `-- name: ResetTables :exec
DELETE FROM users
`

func (q *Queries) ResetTables(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetTables)
	return err
}

const revokeToken = `-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $1
WHERE token = $2
`

type RevokeTokenParams struct {
	RevokedAt sql.NullTime
	Token     string
}

func (q *Queries) RevokeToken(ctx context.Context, arg RevokeTokenParams) error {
	_, err := q.db.ExecContext(ctx, revokeToken, arg.RevokedAt, arg.Token)
	return err
}

const updateEmailPasswordUser = `-- name: UpdateEmailPasswordUser :one
UPDATE users
SET email = $1, hashed_password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateEmailPasswordUserParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateEmailPasswordUser(ctx context.Context, arg UpdateEmailPasswordUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateEmailPasswordUser, arg.Email, arg.HashedPassword, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeUserToRed = `-- name: UpgradeUserToRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
`

func (q *Queries) UpgradeUserToRed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUserToRed, id)
	return err
}
