-- name: GetPendingAnnouncementNotifications :many
SELECT
    n.id, n.show_id, n.subscriber_id,
    s.email AS subscriber_email,
    sh.title, sh.slug, sh.date, sh.description, sh.ticket_url,
    v.name AS venue_name
FROM notifications n
JOIN subscribers s ON s.id = n.subscriber_id
JOIN shows sh ON sh.id = n.show_id
LEFT JOIN venues v ON v.id = sh.venue_id
WHERE n.type = 'published'
  AND n.sent_at IS NULL
  AND s.unsubscribed_at IS NULL
  AND s.deleted_at IS NULL
LIMIT $1;
  