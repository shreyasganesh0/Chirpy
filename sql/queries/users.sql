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

-- name: UpdateEmailPasswordUser :one
UPDATE users
SET email = $1, hashed_password = $2
WHERE id = $3
RETURNING *;

-- name: DeleteChirpByIdForUser :one
WITH delete_attempt AS (
    DELETE FROM chirps
    WHERE chirps.id = $1 AND chirps.user_id = $2
    RETURNING chirps.id, chirps.user_id
)
SELECT 
    COALESCE((SELECT delete_attempt.user_id FROM delete_attempt), (SELECT chirps.user_id FROM chirps WHERE chirps.id = $1)) AS user_id;

-- name: UpgradeUserToRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
