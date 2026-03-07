package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

// POST /links
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

// GET /links
func (h *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	// Optional query param: ?stats=true to include links statistics
	includeStats := r.URL.Query().Get("stats") == "true"

	if includeStats {
		links, err := h.service.GetAllLinks(r.Context())
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch links"})
			return
		}
		respondJSON(w, http.StatusOK, links)
	} else {
	}
}

// GET /links/{id}
func (h *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/links/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid link ID"})
		return
	}

	includeStats := r.URL.Query().Get("stats") == "true"

	if includeStats {
		link, err := h.service.GetLinkWithStats(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch link"})
			return
		}
		respondJSON(w, http.StatusOK, link)
	} else {
		link, err := h.service.GetLink(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch link"})
			return
		}
		respondJSON(w, http.StatusOK, link)
	}
}

// PUT /links/{id}
func (h *LinkHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/links/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid link ID"})
		return
	}

	var req UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	link, err := h.service.UpdateLink(r.Context(), id, req.Title, req.URL, req.Description)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, link)
}

// Extracts numerical id from a path, like /links/123
func extractIDFromPath(path, prefix string) (int64, error) {
	// Remove prefix and clean path
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.TrimSuffix(idStr, "/click")
	idStr = strings.TrimSuffix(idStr, "/")

	// Parse id to int64
	return strconv.ParseInt(idStr, 10, 64)

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
