-- name: ListActiveSubscribers :many
SELECT * FROM subscribers
WHERE confirmed_at IS NOT NULL
  AND unsubscribed_at IS NULL
  AND deleted_at IS NULL;
  