-- +goose Up
CREATE TABLE admin_sessions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash  TEXT UNIQUE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS admin_sessions;
