-- name: GetSubscriberByID :one
SELECT * FROM subscribers
WHERE id = $1 AND deleted_at IS NULL;
