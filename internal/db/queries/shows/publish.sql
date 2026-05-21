-- name: PublishShow :one
UPDATE shows
SET is_published = true, published_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
