package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/EliasLd/nerdy-link-manager/internal/middleware"
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

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	link, err := h.service.CreateLink(r.Context(), userID, req.Title, req.URL, req.Description)
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

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	if includeStats {
		links, err := h.service.GetAllLinksWithStats(r.Context(), userID)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch links"})
			return
		}
		respondJSON(w, http.StatusOK, links)
	} else {
		links, err := h.service.GetAllLinks(r.Context(), userID)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch links"})
			return
		}
		respondJSON(w, http.StatusOK, links)
	}
}

// GET /links/{id}
func (h *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/links/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid link ID"})
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	includeStats := r.URL.Query().Get("stats") == "true"

	if includeStats {
		link, err := h.service.GetLinkWithStats(r.Context(), userID, id)
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
		link, err := h.service.GetLink(r.Context(), userID, id)
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
	id, err := extractIDFromPath(r.URL.Path, "/api/links/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid link ID"})
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	link, err := h.service.UpdateLink(r.Context(), userID, id, req.Title, req.URL, req.Description)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, link)
}

// DELETE /links/{id}
func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/links/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid link ID"})
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	err = h.service.DeleteLink(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to delete link"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /links/{id}/click
func (h *LinkHandler) TrackClick(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/links/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid link ID"})
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	err = h.service.TrackClick(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, services.ErrLinkNotFound) {
			respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to track click"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "click recorded"})
}

// Extracts numerical id from a path, like /links/123
func extractIDFromPath(path, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	if idStr == path {
		return 0, strconv.ErrSyntax // prefix not found
	}

	idStr = strings.TrimSuffix(idStr, "/click")
	idStr = strings.TrimSuffix(idStr, "/")
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
