-- name: CreateSubscriber :one
INSERT INTO subscribers (email)
VALUES ($1)
ON CONFLICT (email) DO NOTHING
RETURNING *;
