package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/internal/services"
)

type LinkHandler struct {
	service *services.LinkService
}

func NewLinkHandler(serivce *services.LinkService) *LinkHandler {
	return &LinkHandler{service: serivce}
}

type CreateLinkRequest struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description *string `json:"description,omitempty"`
}

type UpdateLinkRequest struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description *string `json:"description,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	link, err := h.service.CreateLink(r.Context(), req.Title, req.URL, req.Description)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, link)
}

// Encodes a response in JSON with the appropriate status
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Mapps service's errors to HTTP status
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidURL):
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid URL format"})
	case errors.Is(err, services.ErrEmptyTitle):
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "title cannot be empty"})
	case errors.Is(err, services.ErrLinkNotFound):
		respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
	default:
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
