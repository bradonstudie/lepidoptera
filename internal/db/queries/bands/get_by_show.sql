-- name: GetBandsByShow :many
SELECT b.*, sb.billing_order, sb.is_headliner
FROM bands b
JOIN show_bands sb ON sb.band_id = b.id
WHERE sb.show_id = $1 AND b.deleted_at IS NULL
ORDER BY sb.billing_order;
