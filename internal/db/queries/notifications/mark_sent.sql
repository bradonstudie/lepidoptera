-- name: MarkNotificationSent :exec
UPDATE notifications SET sent_at = NOW() WHERE id = $1;
