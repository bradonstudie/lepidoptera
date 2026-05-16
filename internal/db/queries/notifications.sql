-- name: CreateNotificationsForShow :exec
INSERT INTO notifications (show_id, subscriber_id, type)
SELECT $1, id, $2
FROM subscribers
WHERE confirmed_at IS NOT NULL
  AND unsubscribed_at IS NULL
  AND deleted_at IS NULL;

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
  AND s.deleted_at IS NULL;

-- name: GetPendingReminderNotifications :many
SELECT
    n.id, n.show_id, n.subscriber_id,
    s.email AS subscriber_email,
    sh.title, sh.slug, sh.date, sh.description, sh.ticket_url,
    v.name AS venue_name
FROM notifications n
JOIN subscribers s ON s.id = n.subscriber_id
JOIN shows sh ON sh.id = n.show_id
LEFT JOIN venues v ON v.id = sh.venue_id
WHERE n.type = 'reminder'
  AND n.sent_at IS NULL
  AND sh.date::date = (NOW() + interval '1 day')::date
  AND s.unsubscribed_at IS NULL
  AND s.deleted_at IS NULL;

-- name: MarkNotificationSent :exec
UPDATE notifications SET sent_at = NOW() WHERE id = $1;

-- Bands
-- name: CreateBand :one
INSERT INTO bands (name, genre, description, logo_url, website_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListBands :many
SELECT * FROM bands WHERE deleted_at IS NULL ORDER BY name;

-- name: GetBandsByShow :many
SELECT b.*, sb.billing_order, sb.is_headliner
FROM bands b
JOIN show_bands sb ON sb.band_id = b.id
WHERE sb.show_id = $1 AND b.deleted_at IS NULL
ORDER BY sb.billing_order;

-- name: AddBandToShow :exec
INSERT INTO show_bands (show_id, band_id, billing_order, is_headliner)
VALUES ($1, $2, $3, $4)
ON CONFLICT (show_id, band_id) DO NOTHING;

-- Venues
-- name: CreateVenue :one
INSERT INTO venues (name, address, city, state, capacity)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListVenues :many
SELECT * FROM venues WHERE deleted_at IS NULL ORDER BY name;
