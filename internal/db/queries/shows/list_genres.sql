-- name: ListGenres :many
SELECT DISTINCT b.genre
FROM bands b
JOIN show_bands sb ON sb.band_id = b.id
JOIN shows s ON s.id = sb.show_id
WHERE b.genre IS NOT NULL
  AND b.deleted_at IS NULL
  AND s.is_published = true
  AND s.deleted_at IS NULL
ORDER BY b.genre;
