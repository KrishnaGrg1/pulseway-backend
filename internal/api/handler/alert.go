package handler

import (
	"net/http"
	"strconv"
	"strings"

	mw "github.com/KrishnaGrg1/pulseway/internal/api/middleware"
	db "github.com/KrishnaGrg1/pulseway/internal/db/sqlc"
	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/go-chi/chi/v5"
)

type AlertHandler struct {
	store *store.Store
}

func NewAlertHandler(s *store.Store) *AlertHandler {
	return &AlertHandler{
		store: s,
	}
}

// AlertResponse extends the alert with additional fields expected by frontend
type AlertResponse struct {
	ID        int64  `json:"id"`
	MonitorID int64  `json:"monitor_id"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// ListAllAlerts returns all alerts for the authenticated user's monitors
func (h *AlertHandler) ListAllAlerts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	alerts, err := h.store.Queries.ListAllAlertsByUser(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to fetch alerts")
		return
	}

	// Transform alerts to match frontend format
	alertResponses := make([]AlertResponse, len(alerts))
	for i, alert := range alerts {
		alertResponses[i] = AlertResponse{
			ID:        alert.ID,
			MonitorID: alert.MonitorID,
			Email:     alert.Destination,
			IsActive:  true, // All alerts are active by default
			CreatedAt: alert.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response.Success(w, http.StatusOK, "Successfully retrieved alerts", map[string]any{
		"alerts": alertResponses,
	})
}

// ListAlertsByMonitor returns all alerts for a specific monitor
func (h *AlertHandler) ListAlertsByMonitor(w http.ResponseWriter, r *http.Request) {
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

	alerts, err := h.store.Queries.ListAlertsByMonitor(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to fetch alerts")
		return
	}

	// Transform alerts to match frontend format
	alertResponses := make([]AlertResponse, len(alerts))
	for i, alert := range alerts {
		alertResponses[i] = AlertResponse{
			ID:        alert.ID,
			MonitorID: alert.MonitorID,
			Email:     alert.Destination,
			IsActive:  true,
			CreatedAt: alert.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response.Success(w, http.StatusOK, "Successfully retrieved monitor alerts", map[string]any{
		"alerts": alertResponses,
	})
}

// CreateAlert creates a new alert for a monitor
func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	var input struct {
		MonitorID int64  `json:"monitor_id"`
		Email     string `json:"email"`
	}

	if err := response.Read(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", "VALIDATION_001", "Request body must be valid JSON")
		return
	}

	// Validate email
	if input.Email == "" || !strings.Contains(input.Email, "@") {
		response.Error(w, http.StatusBadRequest, "Validation failed", "VALIDATION_002", "Valid email address is required")
		return
	}

	// Verify monitor belongs to user
	monitor, err := h.store.Queries.GetMonitorByID(r.Context(), input.MonitorID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Not found", "MONITOR_001", "Monitor not found")
		return
	}
	if monitor.UserID != userID {
		response.Error(w, http.StatusForbidden, "Forbidden", "AUTH_003", "Access denied to this monitor")
		return
	}

	// Create the alert
	alert, err := h.store.Queries.CreateAlert(r.Context(), db.CreateAlertParams{
		MonitorID:   input.MonitorID,
		Type:        "email",
		Destination: input.Email,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to create alert")
		return
	}

	alertResponse := AlertResponse{
		ID:        alert.ID,
		MonitorID: alert.MonitorID,
		Email:     alert.Destination,
		IsActive:  true,
		CreatedAt: alert.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(w, http.StatusCreated, "Alert created", alertResponse)
}

// DeleteAlert deletes an alert
func (h *AlertHandler) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")
	alertID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", "VALIDATION_003", "Invalid alert id")
		return
	}

	// Get alert to verify ownership
	alert, err := h.store.Queries.GetAlert(r.Context(), alertID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Not found", "ALERT_001", "Alert not found")
		return
	}

	// Verify monitor belongs to user
	monitor, err := h.store.Queries.GetMonitorByID(r.Context(), alert.MonitorID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Not found", "MONITOR_001", "Monitor not found")
		return
	}
	if monitor.UserID != userID {
		response.Error(w, http.StatusForbidden, "Forbidden", "AUTH_003", "Access denied to this alert")
		return
	}

	// Delete the alert
	err = h.store.Queries.DeleteAlert(r.Context(), alertID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to delete alert")
		return
	}

	response.Success(w, http.StatusOK, "Alert deleted", nil)
}
