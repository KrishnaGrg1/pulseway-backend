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

-- name: GetCheckHistory :many
SELECT status, latency_ms, checked_at
FROM check_results
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT $2;

-- name: GetMonitorCurrentStatus :one
SELECT
  cr.status AS current_status,
  COUNT(*) FILTER (WHERE cr.status = 'up') * 100 / NULLIF(COUNT(*), 0) AS uptime_percentage,
  AVG(cr.latency_ms)::INT AS avg_latency_ms,
  MAX(cr.checked_at) AS last_checked_at
FROM check_results cr
WHERE cr.monitor_id = $1
AND cr.checked_at > now() - INTERVAL '24 hours'
GROUP BY cr.monitor_id
LIMIT 1;

-- name: GetLastCheckStatus :one
SELECT status, checked_at
FROM check_results
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT 1;

-- name: GetMetricsHistory :many
WITH daily_stats AS (
  SELECT
    DATE_TRUNC('day', cr.checked_at) AS day,
    COUNT(DISTINCT m.id) AS total_monitors,
    COUNT(DISTINCT CASE WHEN cr.status = 'up' THEN m.id END) AS healthy_monitors,
    COUNT(*) FILTER (WHERE cr.status = 'up') * 100 / NULLIF(COUNT(*), 0) AS uptime_percentage,
    AVG(cr.latency_ms)::INT AS avg_latency_ms
  FROM check_results cr
  JOIN monitors m ON cr.monitor_id = m.id
  WHERE m.user_id = sqlc.arg(user_id)
  AND cr.checked_at > now() - INTERVAL '1 day' * sqlc.arg(days)
  GROUP BY day
  ORDER BY day ASC
)
SELECT
  day AS timestamp,
  total_monitors AS value,
  CASE WHEN total_monitors > 0 THEN true ELSE false END AS healthy,
  'total_monitors' AS metric_type
FROM daily_stats
UNION ALL
SELECT
  day AS timestamp,
  healthy_monitors AS value,
  CASE WHEN healthy_monitors >= total_monitors * 0.8 THEN true ELSE false END AS healthy,
  'healthy_count' AS metric_type
FROM daily_stats
UNION ALL
SELECT
  day AS timestamp,
  uptime_percentage AS value,
  CASE WHEN uptime_percentage >= 95 THEN true ELSE false END AS healthy,
  'uptime_percentage' AS metric_type
FROM daily_stats
UNION ALL
SELECT
  day AS timestamp,
  avg_latency_ms AS value,
  CASE WHEN avg_latency_ms <= 200 THEN true ELSE false END AS healthy,
  'avg_latency_ms' AS metric_type
FROM daily_stats
ORDER BY metric_type, timestamp;