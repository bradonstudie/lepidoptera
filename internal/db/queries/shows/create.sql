-- name: CreateShow :one
INSERT INTO shows (title, slug, date, venue_id, description, flyer_url, ticket_url, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
