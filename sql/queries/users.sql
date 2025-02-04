-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES ($1, $2)
RETURNING *;

-- name: UsersClear :exec
TRUNCATE users, messages;

-- name: LoginUser :one
SELECT id, created_at, updated_at, email, hashed_password
FROM users
WHERE (email = $1)
LIMIT 1;