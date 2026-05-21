-- name: ListVenues :many
SELECT * FROM venues WHERE deleted_at IS NULL ORDER BY name;
