-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES ($1, $2)
RETURNING *;

-- name: UsersClear :exec
TRUNCATE users, messages, refresh_tokens;

-- name: LoginUser :one
SELECT id, created_at, updated_at, email, hashed_password
FROM users
WHERE (email = $1)
LIMIT 1;

-- name: UserUpdateSelf :one
UPDATE users
SET email = $2, hashed_password = $3
WHERE id = $1
Returning *;