-- name: CreateUser :one
INSERT INTO users(email, hashed_password)
VALUES(
    $1,
    $2
) RETURNING *;

-- name: ResetTables :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps(body, user_id)
VALUES(
    $1,
    $2
) RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
where email = $1;

-- name: CreateRefreshTokens :one
INSERT INTO refresh_tokens(token, user_id, expires_at, revoked_at)
VALUES(
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id, revoked_at from refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $1
WHERE token = $2; 
