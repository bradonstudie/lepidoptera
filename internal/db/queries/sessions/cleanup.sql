-- name: DeleteExpiredSessions :exec
DELETE FROM admin_sessions
WHERE expires_at < NOW() - INTERVAL '1 month' OR revoked_at < NOW() - INTERVAL '1 month';

-- name: DeleteUsedLoginTokens :exec
DELETE FROM used_login_tokens
WHERE used_at < NOW() - INTERVAL '1 month';
