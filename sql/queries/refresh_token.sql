-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id)
VALUES (
    $1,
    $2
)
RETURNING *;


-- name: GetRefreshTokenFromToken :one
SELECT * from refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshTokenInvokeFromToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;