-- name: ListAllShows :many
SELECT
    s.id, s.title, s.slug, s.date, s.is_published, s.published_at,
    v.name AS venue_name
FROM shows s
LEFT JOIN venues v ON v.id = s.venue_id
WHERE s.deleted_at IS NULL
ORDER BY s.date DESC;
