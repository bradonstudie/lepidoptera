-- name: GetSubscriberByEmail :one
SELECT * FROM subscribers
WHERE email = $1 AND deleted_at IS NULL;
