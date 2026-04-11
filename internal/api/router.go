package api

import (
	"net/http"

	"github.com/KrishnaGrg1/pulseway/internal/api/handler"
	mw "github.com/KrishnaGrg1/pulseway/internal/api/middleware"
	"github.com/KrishnaGrg1/pulseway/internal/config"
	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/sse"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(s *store.Store, cfg *config.Config, hub *sse.Hub) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusAccepted, "Health is good", nil)
	})
	r.Get("/sse", hub.ServeHTTP)
	AuthHandler := handler.NewAuthHandler(s, cfg.JWT_SECRET)
	MonitorHandler := handler.NewMonitorHandler(s)
	statsHandler := handler.NewStatsHandler(s)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", AuthHandler.Register)
			r.Post("/login", AuthHandler.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(mw.JwtAuth(s, cfg.JWT_SECRET))

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				userIDValue := r.Context().Value(mw.UserIDKey)
				userID, ok := userIDValue.(int64)
				if !ok {
					response.Error(w, http.StatusUnauthorized, "Unauthorized", "AUTH_002", "Invalid token")
					return
				}
				user, err := s.Queries.GetUserByID(r.Context(), userID)
				if err != nil {
					response.Error(w, http.StatusInternalServerError, "Failed to retrieve user", "INTERNAL_001", "Database query failed")
					return
				}

				response.Success(w, http.StatusAccepted, "Retrieved user data successfully", user)
			})
			r.Get("/dashboard/stats", statsHandler.Get)
			r.Route("/monitors", func(r chi.Router) {
				r.Post("/", MonitorHandler.CreateMonitor)
				r.Get("/", MonitorHandler.List)
				r.Get("/{id}", MonitorHandler.Get)
				r.Put("/{id}", MonitorHandler.Update)
				r.Delete("/{id}", MonitorHandler.Delete)
				r.Get("/{id}/results", MonitorHandler.GetResults)
			})
		})
	})
	return r
}
