package handler

import (
	"net/http"
	"strconv"

	mw "github.com/KrishnaGrg1/pulseway/internal/api/middleware"
	db "github.com/KrishnaGrg1/pulseway/internal/db/sqlc"
	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/store"
)

type StatsHandler struct {
	store *store.Store
}

func NewStatsHandler(s *store.Store) *StatsHandler {
	return &StatsHandler{store: s}
}

func (h *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	// Get monitor counts
	monitorStats, err := h.store.Queries.GetMonitorStats(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch monitor stats", "monitor_stats_error", err.Error())
		return
	}

	// Get uptime and latency
	checkStats, err := h.store.Queries.GetStatsForUser(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch check stats", "check_stats_error", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "stats fetched", map[string]any{
		"total_monitors":    monitorStats.TotalMonitors,
		"healthy_monitors":  monitorStats.HealthyMonitors,
		"uptime_percentage": checkStats.UptimePercentage,
		"avg_latency_ms":    checkStats.AvgLatencyMs,
	})
}

func (h *StatsHandler) GetMetricsHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	daysStr := r.URL.Query().Get("days")
	days := int32(7)
	if daysStr != "" {
		if daysVal, err := strconv.ParseInt(daysStr, 10, 32); err == nil && daysVal > 0 {
			days = int32(daysVal)
		}
	}

	history, err := h.store.Queries.GetMetricsHistory(r.Context(), db.GetMetricsHistoryParams{
		UserID: userID,
		Days:   days,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch metrics history", "metrics_history_error", err.Error())
		return
	}

	result := map[string][]map[string]any{
		"total_monitors":    {},
		"healthy_count":     {},
		"uptime_percentage": {},
		"avg_latency_ms":    {},
	}

	for _, record := range history {
		item := map[string]any{
			"timestamp": record.Timestamp,
			"value":     record.Value,
			"healthy":   record.Healthy,
		}
		result[record.MetricType] = append(result[record.MetricType], item)
	}

	response.Success(w, http.StatusOK, "metrics history fetched", result)
}
