-- name: UnsubscribeSubscriber :exec
UPDATE subscribers
SET unsubscribed_at = NOW()
WHERE id = $1;
