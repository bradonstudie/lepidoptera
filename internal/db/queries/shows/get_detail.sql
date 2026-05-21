-- name: GetShowBySlug :one
SELECT
    s.id, s.title, s.slug, s.date, s.description, s.flyer_url, s.ticket_url,
    v.name AS venue_name, v.city AS venue_city, v.address AS venue_address
FROM shows s
LEFT JOIN venues v ON v.id = s.venue_id
WHERE s.slug = $1
  AND s.is_published = true
  AND s.deleted_at IS NULL;
  