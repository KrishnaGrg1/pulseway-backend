package handler

import (
	"net/http"

	mw "github.com/KrishnaGrg1/pulseway/internal/api/middleware"
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
		"active_monitors":   monitorStats.ActiveMonitors,
		"uptime_percentage": checkStats.UptimePercentage,
		"avg_latency_ms":    checkStats.AvgLatencyMs,
	})
}
