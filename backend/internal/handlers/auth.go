package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/internal/auth"
	"github.com/EliasLd/nerdy-link-manager/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(us *services.UserService, jm *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: us,
		jwtManager:  jm,
	}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "could not register", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtManager.Generate(user.ID, user.Email)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
