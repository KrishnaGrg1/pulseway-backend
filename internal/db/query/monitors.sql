-- name: CreateMonitor :one
INSERT INTO monitors (user_id, name, url, interval_secs)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListMonitorsByUser :many
SELECT * FROM monitors
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: GetMonitorByID :one
SELECT * FROM monitors
WHERE id = $1;

-- name: UpdateMonitor :one
UPDATE monitors
SET name = $2, url = $3, interval_secs = $4
WHERE id = $1 AND user_id = $5
RETURNING *;

-- name: DeleteMonitor :exec
UPDATE monitors
SET is_active = false
WHERE id = $1 AND user_id = $2;

-- name: ListAllActiveMonitors :many
SELECT * FROM monitors
WHERE is_active = true;