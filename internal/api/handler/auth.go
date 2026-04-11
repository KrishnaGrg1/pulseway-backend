package handler

import (
	"errors"
	"net/http"
	"time"

	db "github.com/KrishnaGrg1/pulseway/internal/db/sqlc"
	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store     *store.Store
	jwtSecret string
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  db.User `json:"user"`
}

func NewAuthHandler(s *store.Store, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		store:     s,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := response.Read(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", "VALIDATION_001", "Request body must be valid JSON")
		return
	}
	if input.Email == "" || input.Password == "" {
		response.Error(w, http.StatusBadRequest, "Validation failed", "VALIDATION_002", "email and password are required")
		return
	}
	_, err := h.store.Queries.GetUserByEmail(r.Context(), input.Email)
	if err == nil {
		response.Error(w, http.StatusConflict, "Email already in use", "AUTH_001", "email already used")
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		response.Error(w, http.StatusInternalServerError, "Failed to check existing user", "INTERNAL_002", "Database query failed")
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to hash password", "INTERNAL_001", "Hashing failed")
		return
	}
	user, err := h.store.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Email:    input.Email,
		Password: string(hashed),
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create user", "INTERNAL_001", "Database insert failed")
		return
	}
	response.Success(w, http.StatusAccepted, "Successfully signup", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := response.Read(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", "VALIDATION_001", "Request body must be valid JSON")
		return
	}
	if input.Email == "" || input.Password == "" {
		response.Error(w, http.StatusBadRequest, "Validation failed", "VALIDATION_002", "email and password are required")
		return
	}
	user, error := h.store.Queries.GetUserByEmail(r.Context(), input.Email)
	if error != nil {
		response.Error(w, http.StatusNotFound, "Invalid credentials", "AUTH_001", "Invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid credentials", "AUTH_002", "Invalid credentials")
		return
	}
	token, err := generateToken(int(user.ID), h.jwtSecret)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate token", "INTERNAL_001", "Token generation failed")
		return
	}
	secure := r.TLS != nil
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})
	response.Success(w, http.StatusOK, "Login successful", LoginResponse{
		Token: token,
		User:  user,
	})
}

func generateToken(userId int, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
