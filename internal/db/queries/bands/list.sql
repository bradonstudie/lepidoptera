-- name: ListBands :many
SELECT * FROM bands WHERE deleted_at IS NULL ORDER BY name;
