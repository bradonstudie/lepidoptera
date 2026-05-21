-- name: ConfirmSubscriber :exec
UPDATE subscribers
SET confirmed_at = NOW()
WHERE id = $1 AND confirmed_at IS NULL;