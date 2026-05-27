-- name: MarkLoginTokenUsed :one
INSERT INTO used_login_tokens (token_hash)
VALUES ($1)
ON CONFLICT (token_hash) DO NOTHING
RETURNING token_hash;
