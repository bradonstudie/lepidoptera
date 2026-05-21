-- name: SoftDeleteShow :exec
UPDATE shows SET deleted_at = NOW() WHERE id = $1;
