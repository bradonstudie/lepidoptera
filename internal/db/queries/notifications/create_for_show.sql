-- name: CreateNotificationsForShow :exec
INSERT INTO notifications (show_id, subscriber_id, type)
SELECT $1, id, $2
FROM subscribers
WHERE confirmed_at IS NOT NULL
  AND unsubscribed_at IS NULL
  AND deleted_at IS NULL;
