package handler

import (
	"net/http"
	"strconv"
	"time"

	mw "github.com/KrishnaGrg1/pulseway/internal/api/middleware"
	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/go-chi/chi/v5"
)

type IncidentHandler struct {
	store *store.Store
}

func NewIncidentHandler(s *store.Store) *IncidentHandler {
	return &IncidentHandler{
		store: s,
	}
}

// IncidentResponse extends the incident with a calculated duration
type IncidentResponse struct {
	ID              int64      `json:"id"`
	MonitorID       int64      `json:"monitor_id"`
	StartedAt       time.Time  `json:"started_at"`
	ResolvedAt      *time.Time `json:"resolved_at"`
	Notified        bool       `json:"notified"`
	DurationSeconds *int64     `json:"duration_seconds,omitempty"`
}

// ListAllIncidents returns all incidents for the authenticated user's monitors
func (h *IncidentHandler) ListAllIncidents(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	incidents, err := h.store.Queries.ListAllIncidentsByUser(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to fetch incidents")
		return
	}

	// Transform incidents to include duration
	incidentResponses := make([]IncidentResponse, len(incidents))
	for i, incident := range incidents {
		incidentResponses[i] = IncidentResponse{
			ID:         incident.ID,
			MonitorID:  incident.MonitorID,
			StartedAt:  incident.StartedAt.Time,
			Notified:   incident.Notified,
		}

		if incident.ResolvedAt.Valid {
			resolvedAt := incident.ResolvedAt.Time
			incidentResponses[i].ResolvedAt = &resolvedAt
			duration := int64(resolvedAt.Sub(incident.StartedAt.Time).Seconds())
			incidentResponses[i].DurationSeconds = &duration
		}
	}

	response.Success(w, http.StatusOK, "Successfully retrieved incidents", map[string]any{
		"incidents": incidentResponses,
	})
}

// ListIncidentsByMonitor returns all incidents for a specific monitor
func (h *IncidentHandler) ListIncidentsByMonitor(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")
	monitorID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", "VALIDATION_003", "Invalid monitor id")
		return
	}

	// Verify monitor belongs to user
	monitor, err := h.store.Queries.GetMonitorByID(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Not found", "MONITOR_001", "Monitor not found")
		return
	}
	if monitor.UserID != userID {
		response.Error(w, http.StatusForbidden, "Forbidden", "AUTH_003", "Access denied to this monitor")
		return
	}

	incidents, err := h.store.Queries.ListIncidentsByMonitor(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to fetch incidents")
		return
	}

	// Transform incidents to include duration
	incidentResponses := make([]IncidentResponse, len(incidents))
	for i, incident := range incidents {
		incidentResponses[i] = IncidentResponse{
			ID:         incident.ID,
			MonitorID:  incident.MonitorID,
			StartedAt:  incident.StartedAt.Time,
			Notified:   incident.Notified,
		}

		if incident.ResolvedAt.Valid {
			resolvedAt := incident.ResolvedAt.Time
			incidentResponses[i].ResolvedAt = &resolvedAt
			duration := int64(resolvedAt.Sub(incident.StartedAt.Time).Seconds())
			incidentResponses[i].DurationSeconds = &duration
		}
	}

	response.Success(w, http.StatusOK, "Successfully retrieved monitor incidents", map[string]any{
		"incidents": incidentResponses,
	})
}
