-- name: CreateAdminSession :one
INSERT INTO admin_sessions (token_hash, expires_at)
VALUES ($1, $2)
RETURNING *;
