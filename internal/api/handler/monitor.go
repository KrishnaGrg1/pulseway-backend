package handler

import (
	"errors"
	"net/http"
	"strconv"

	mw "github.com/KrishnaGrg1/pulseway/internal/api/middleware"
	db "github.com/KrishnaGrg1/pulseway/internal/db/sqlc"
	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type MonitorHandler struct {
	store *store.Store
}

func NewMonitorHandler(s *store.Store) *MonitorHandler {
	return &MonitorHandler{
		store: s,
	}
}

func (h *MonitorHandler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(mw.UserIDKey).(int64)
	if err := h.GetUserByIDCheck(w, r, userId); err != nil {
		return
	}
	var input struct {
		Name         string `json:"name"`
		URL          string `json:"url"`
		IntervalSecs int32  `json:"interval_secs"`
	}
	if err := response.Read(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", "VALIDATION_001", "Request body must be valid JSON")
		return
	}
	if input.Name == "" || input.URL == "" {
		response.Error(w, http.StatusBadRequest, "Validation failed", "VALIDATION_002", "Monitor's name and url are required")
		return
	}
	if input.IntervalSecs == 0 {
		input.IntervalSecs = 60
	}

	monitor, err := h.store.Queries.CreateMonitor(r.Context(), db.CreateMonitorParams{
		UserID:       userId,
		Name:         input.Name,
		Url:          input.URL,
		IntervalSecs: input.IntervalSecs,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to create monitor  ")
		return
	}
	response.Success(w, http.StatusAccepted, "Successfully created monitor", monitor)

}

func (h *MonitorHandler) List(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(mw.UserIDKey).(int64)
	if err := h.GetUserByIDCheck(w, r, userId); err != nil {
		return
	}
	monitor, err := h.store.Queries.ListMonitorsByUser(r.Context(), userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server issue", "INTERNAL_001", "Failed to fetch data")
		return
	}
	response.Success(w, http.StatusAccepted, "Successfully retreived monitors of user", monitor)
}

func (h *MonitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", "VALIDATION_003", "Invalid monitor id")
		return
	}

	monitor, err := h.store.Queries.GetMonitorByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Not found", "MONITOR_001", "Monitor not found")
		return
	}

	response.Success(w, http.StatusAccepted, "Successfully retrieved monitor details by id ", monitor)
}

func (h *MonitorHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)
	if err := h.GetUserByIDCheck(w, r, userID); err != nil {
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", "VALIDATION_003", "Invalid monitor id")
		return
	}

	var input struct {
		Name         string `json:"name"`
		URL          string `json:"url"`
		IntervalSecs int32  `json:"interval_secs"`
	}
	if err := response.Read(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", "VALIDATION_001", "Request body must be valid JSON")
		return
	}

	monitor, err := h.store.Queries.UpdateMonitor(r.Context(), db.UpdateMonitorParams{
		ID:           id,
		Name:         input.Name,
		Url:          input.URL,
		IntervalSecs: input.IntervalSecs,
		UserID:       userID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to update monitor")
		return
	}
	response.Success(w, http.StatusAccepted, "Successfully updated monitor", monitor)
}

func (h *MonitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)
	if err := h.GetUserByIDCheck(w, r, userID); err != nil {
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", "VALIDATION_003", "Invalid monitor id")
		return
	}
	existingMonitor, err := h.store.Queries.GetMonitorByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Not found", "MONITOR_001", "Monitor not found, cannot delete")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to fetch monitor")
		return
	}
	if existingMonitor.IsActive == false {
		response.Error(w, http.StatusConflict, "Inactive monitor", "MONITOR_002", "Inactive monitor cannot be deleted")
		return
	}
	err = h.store.Queries.DeleteMonitor(r.Context(), db.DeleteMonitorParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to delete monitor")
		return
	}
	response.Success(w, http.StatusOK, "Successfully deleted monitor", nil)
}

func (h *MonitorHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey).(int64)
	if err := h.GetUserByIDCheck(w, r, userID); err != nil {
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", "VALIDATION_003", "Invalid monitor id")
		return
	}

	results, err := h.store.Queries.ListCheckResultsByMonitor(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to fetch results")
		return
	}
	response.Success(w, http.StatusAccepted, "Successfully retrieved monitor results", results)
}

func (h *MonitorHandler) GetUserByIDCheck(w http.ResponseWriter, r *http.Request, userId int64) error {
	_, err := h.store.Queries.GetUserByID(r.Context(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "User not found", "AUTH_004", "User does not exist")
			return err
		}
		response.Error(w, http.StatusInternalServerError, "Server error", "INTERNAL_001", "Failed to check the existing user")
		return err
	}
	return nil
}
