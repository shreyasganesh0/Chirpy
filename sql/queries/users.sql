-- name: CreateUser :one
INSERT INTO users(email)
VALUES(
    $1
) RETURNING *;

-- name: ResetTables :exec
DELETE FROM users;
