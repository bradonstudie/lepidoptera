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
