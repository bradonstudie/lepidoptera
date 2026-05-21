-- name: CreateBand :one
INSERT INTO bands (name, genre, description, logo_url, website_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
