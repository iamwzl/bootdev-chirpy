-- name: CreateMessage :one
INSERT INTO messages (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetMessages_CreatedAtASC :many
SELECT id, created_at, updated_at, body, user_id
FROM messages
ORDER BY created_at ASC;

-- name: GetMessage :one
SELECT id, created_at, updated_at, body, user_id
FROM messages
WHERE id = $1
LIMIT 1;