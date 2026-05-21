-- name: CreateVenue :one
INSERT INTO venues (name, address, city, state, capacity)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
