-- name: RevokeAllAdminSessions :exec
UPDATE admin_sessions
SET revoked_at = NOW()
WHERE revoked_at IS NULL;
