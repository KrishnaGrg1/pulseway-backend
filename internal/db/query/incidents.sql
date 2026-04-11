-- name: CreateIncident :one
INSERT INTO incidents (monitor_id)
VALUES ($1)
RETURNING *;

-- name: ResolveIncident :one
UPDATE incidents
SET resolved_at = now()
WHERE monitor_id = $1
AND resolved_at IS NULL
RETURNING *;

-- name: GetActiveIncident :one
SELECT * FROM incidents
WHERE monitor_id = $1
AND resolved_at IS NULL
LIMIT 1;

-- name: ListIncidentsByMonitor :many
SELECT * FROM incidents
WHERE monitor_id = $1
ORDER BY started_at DESC;