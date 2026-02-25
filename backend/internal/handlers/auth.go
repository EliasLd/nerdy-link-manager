package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/internal/auth"
	"github.com/EliasLd/nerdy-link-manager/internal/services"
)

type AuthHandler struct {
	userSerive *services.UserService
	jwtManager *auth.JWTManager
}

func NewAuthHandler(us *services.UserService, jm *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userSerive: us,
		jwtManager: jm,
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

	err := h.userSerive.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "could not register", http.StatusInternalServerError)
		return
	}
}
