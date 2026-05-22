-- name: RevokeAdminSession :exec
UPDATE admin_sessions
SET revoked_at = NOW()
WHERE token_hash = $1;
