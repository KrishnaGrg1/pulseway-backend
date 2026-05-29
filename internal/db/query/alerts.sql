-- name: CreateAlert :one
INSERT INTO alerts (monitor_id, type, destination)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListAlertsByMonitor :many
SELECT * FROM alerts
WHERE monitor_id = $1;

-- name: ListAllAlertsByUser :many
SELECT a.* FROM alerts a
JOIN monitors m ON a.monitor_id = m.id
WHERE m.user_id = $1;

-- name: GetAlert :one
SELECT * FROM alerts
WHERE id = $1;

-- name: DeleteAlert :exec
DELETE FROM alerts
WHERE id = $1;