-- name: CreateAlert :one
INSERT INTO alerts (monitor_id, type, destination)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListAlertsByMonitor :many
SELECT * FROM alerts
WHERE monitor_id = $1;

-- name: DeleteAlert :exec
DELETE FROM alerts
WHERE id = $1 AND monitor_id = $2;