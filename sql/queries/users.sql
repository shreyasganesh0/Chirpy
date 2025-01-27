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

