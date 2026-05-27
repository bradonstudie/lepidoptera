-- +goose Up
CREATE TABLE used_login_tokens (
    token_hash TEXT PRIMARY KEY,
    used_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS used_login_tokens;
