-- name: AddBandToShow :exec
INSERT INTO show_bands (show_id, band_id, billing_order, is_headliner)
VALUES ($1, $2, $3, $4)
ON CONFLICT (show_id, band_id) DO NOTHING;
