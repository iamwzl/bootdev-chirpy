-- name: CreateUser :one
INSERT INTO users (email)
VALUES ($1)
RETURNING *;

-- name: UsersClear :exec
TRUNCATE users, messages;