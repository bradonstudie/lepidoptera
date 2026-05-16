-- name: CreateSubscriber :one
INSERT INTO subscribers (email)
VALUES ($1)
ON CONFLICT (email) DO NOTHING
RETURNING *;

-- name: GetSubscriberByEmail :one
SELECT * FROM subscribers
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetSubscriberByID :one
SELECT * FROM subscribers
WHERE id = $1 AND deleted_at IS NULL;

-- name: ConfirmSubscriber :exec
UPDATE subscribers
SET confirmed_at = NOW()
WHERE id = $1 AND confirmed_at IS NULL;

-- name: UnsubscribeSubscriber :exec
UPDATE subscribers
SET unsubscribed_at = NOW()
WHERE id = $1;

-- name: ListActiveSubscribers :many
SELECT * FROM subscribers
WHERE confirmed_at IS NOT NULL
  AND unsubscribed_at IS NULL
  AND deleted_at IS NULL;

-- name: SoftDeleteSubscriber :exec
UPDATE subscribers SET deleted_at = NOW() WHERE id = $1;
