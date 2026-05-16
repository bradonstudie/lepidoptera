-- name: ListPublishedShows :many
SELECT
    s.id, s.title, s.slug, s.date, s.description, s.flyer_url, s.ticket_url,
    v.name AS venue_name, v.city AS venue_city, v.address AS venue_address
FROM shows s
LEFT JOIN venues v ON v.id = s.venue_id
WHERE s.is_published = true
  AND s.deleted_at IS NULL
  AND s.date >= NOW()
  AND (sqlc.narg('genre')::text IS NULL OR EXISTS (
      SELECT 1 FROM show_bands sb
      JOIN bands b ON b.id = sb.band_id
      WHERE sb.show_id = s.id AND b.genre = sqlc.narg('genre')
  ))
ORDER BY s.date ASC;

-- name: GetShowBySlug :one
SELECT
    s.id, s.title, s.slug, s.date, s.description, s.flyer_url, s.ticket_url,
    v.name AS venue_name, v.city AS venue_city, v.address AS venue_address
FROM shows s
LEFT JOIN venues v ON v.id = s.venue_id
WHERE s.slug = $1
  AND s.is_published = true
  AND s.deleted_at IS NULL;

-- name: ListAllShows :many
SELECT
    s.id, s.title, s.slug, s.date, s.is_published, s.published_at,
    v.name AS venue_name
FROM shows s
LEFT JOIN venues v ON v.id = s.venue_id
WHERE s.deleted_at IS NULL
ORDER BY s.date DESC;

-- name: CreateShow :one
INSERT INTO shows (title, slug, date, venue_id, description, flyer_url, ticket_url, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: PublishShow :one
UPDATE shows
SET is_published = true, published_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteShow :exec
UPDATE shows SET deleted_at = NOW() WHERE id = $1;

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
