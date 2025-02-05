-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id)
VALUES ($1, $2);

-- name: GetRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE (token = $1 AND expires_at>CURRENT_TIMESTAMP AND revoked_at IS NULL)
LIMIT 1;

-- name: RevokeRefreshToken :execrows
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP
WHERE token = $1;