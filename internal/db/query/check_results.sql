-- name: CreateCheckResult :one
INSERT INTO check_results (monitor_id, status, latency_ms, status_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListCheckResultsByMonitor :many
SELECT * FROM check_results
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT 100;

-- name: GetUptimePercentage :one
SELECT
  COUNT(*) FILTER (WHERE status = 'up') * 100 / COUNT(*) AS uptime_percentage
FROM check_results
WHERE monitor_id = $1
AND checked_at > now() - INTERVAL '24 hours';

-- name: GetStatsForUser :one
SELECT
  COUNT(*) FILTER (WHERE cr.status = 'up') * 100 / NULLIF(COUNT(*), 0) AS uptime_percentage,
  AVG(cr.latency_ms) AS avg_latency_ms
FROM check_results cr
JOIN monitors m ON cr.monitor_id = m.id
WHERE m.user_id = $1
AND cr.checked_at > now() - INTERVAL '24 hours';