-- name: SoftDeleteSubscriber :exec
UPDATE subscribers SET deleted_at = NOW() WHERE id = $1;
