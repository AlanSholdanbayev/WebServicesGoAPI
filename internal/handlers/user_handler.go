package handlers

import (
	"encoding/json"
	"finalproject/internal/config"
	"finalproject/internal/logger"
	"finalproject/internal/middleware"
	"finalproject/internal/models"
	"finalproject/internal/service"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// UserHandler хранит зависимости
type UserHandler struct {
	Service  *service.UserService
	Validate *validator.Validate
	Cfg      *config.Config
	Logger   *logger.LoggerWrapper
}

// NewUserHandler создаёт новый обработчик
func NewUserHandler(s *service.UserService, cfg *config.Config, lg *logger.LoggerWrapper) *UserHandler {
	return &UserHandler{
		Service:  s,
		Validate: validator.New(),
		Cfg:      cfg,
		Logger:   lg,
	}
}

// ---------------------------
// Register
// ---------------------------

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		h.Logger.Info().Msg("Failed to decode register request")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Info().Err(err).Msg("Validation failed for register request")
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := h.Service.Register(r.Context(), user); err != nil {
		http.Error(w, "could not create", http.StatusInternalServerError)
		h.Logger.Error().Err(err).Msg("Failed to register user")
		return
	}

	h.Logger.Info().Str("email", req.Email).Msg("User registered successfully")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ---------------------------
// Login
// ---------------------------

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		h.Logger.Info().Msg("Failed to decode login request")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Info().Err(err).Msg("Validation failed for login request")
		return
	}

	u, err := h.Service.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		h.Logger.Info().Str("email", req.Email).Msg("Failed login attempt")
		return
	}

	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	signed, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		h.Logger.Error().Err(err).Msg("Failed to sign JWT")
		return
	}

	h.Logger.Info().Int64("user_id", u.ID).Msg("User logged in successfully")
	json.NewEncoder(w).Encode(map[string]string{"token": signed})
}

// ---------------------------
// Me (protected endpoint)
// ---------------------------

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	idIf := r.Context().Value(middleware.UserIDKey)
	if idIf == nil {
		http.Error(w, "missing user", http.StatusUnauthorized)
		h.Logger.Info().Msg("Unauthorized access to /me")
		return
	}

	id := idIf.(int64)
	u, err := h.Service.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		h.Logger.Info().Int64("user_id", id).Msg("User not found in /me")
		return
	}

	u.Password = ""
	h.Logger.Info().Int64("user_id", id).Msg("Fetched user profile")
	json.NewEncoder(w).Encode(u)
}

// ---------------------------
// Update User
// ---------------------------

type updateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Name     string `json:"name"`
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idIf := r.Context().Value(middleware.UserIDKey)
	if idIf == nil {
		http.Error(w, "missing user", http.StatusUnauthorized)
		h.Logger.Info().Msg("Unauthorized update attempt")
		return
	}
	id := idIf.(int64)

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		h.Logger.Info().Msg("Failed to decode update request")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Info().Err(err).Msg("Validation failed for update request")
		return
	}

	user := &models.User{
		ID:       id,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := h.Service.Update(r.Context(), user); err != nil {
		http.Error(w, "could not update", http.StatusInternalServerError)
		h.Logger.Error().Err(err).Int64("user_id", id).Msg("Failed to update user")
		return
	}

	h.Logger.Info().Int64("user_id", id).Msg("User updated successfully")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ---------------------------
// Delete User
// ---------------------------

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idIf := r.Context().Value(middleware.UserIDKey)
	if idIf == nil {
		http.Error(w, "missing user", http.StatusUnauthorized)
		h.Logger.Info().Msg("Unauthorized delete attempt")
		return
	}
	id := idIf.(int64)

	if err := h.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, "could not delete", http.StatusInternalServerError)
		h.Logger.Error().Err(err).Int64("user_id", id).Msg("Failed to delete user")
		return
	}

	h.Logger.Info().Int64("user_id", id).Msg("User deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------
// Register all routes
// ---------------------------
func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/register", h.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", h.Login).Methods("POST")

	// защищённые маршруты
	r.Handle("/api/v1/me", middleware.JWTAuth(h.Cfg)(http.HandlerFunc(h.Me))).Methods("GET")
	r.Handle("/api/v1/users/me", middleware.JWTAuth(h.Cfg)(http.HandlerFunc(h.UpdateUser))).Methods("PUT")
	r.Handle("/api/v1/users/me", middleware.JWTAuth(h.Cfg)(http.HandlerFunc(h.DeleteUser))).Methods("DELETE")
}
