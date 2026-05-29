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


-- name: GetMonitorStats :one
SELECT
  COUNT(*) AS total_monitors,
  COUNT(*) FILTER (WHERE is_active = true) AS healthy_monitors
FROM monitors
WHERE user_id = $1;

-- name: ListMonitorsWithStats :many
SELECT
  m.id,
  m.user_id,
  m.name,
  m.url,
  m.interval_secs,
  m.is_active,
  m.created_at,
  COALESCE(latest.status, 'unknown') AS current_status,
  COALESCE(stats.uptime_percentage, 0) AS uptime_percentage,
  COALESCE(stats.avg_latency_ms, 0)::INT AS avg_latency_ms,
  latest.checked_at AS last_checked_at,
  latest.status AS last_check_status
FROM monitors m
LEFT JOIN LATERAL (
  SELECT status, checked_at
  FROM check_results
  WHERE monitor_id = m.id
  ORDER BY checked_at DESC
  LIMIT 1
) latest ON true
LEFT JOIN LATERAL (
  SELECT
    COUNT(*) FILTER (WHERE status = 'up') * 100 / NULLIF(COUNT(*), 0) AS uptime_percentage,
    AVG(latency_ms) AS avg_latency_ms
  FROM check_results
  WHERE monitor_id = m.id
  AND checked_at > now() - INTERVAL '24 hours'
) stats ON true
WHERE m.user_id = $1 AND m.is_active = true
ORDER BY m.created_at DESC;