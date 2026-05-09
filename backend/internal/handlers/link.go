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

func NewLinkHandler(service *services.LinkService) *LinkHandler {
	return &LinkHandler{service: service}
}

type CreateLinkRequest struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description *string `json:"description,omitempty"`
	FolderID    *int64  `json:"folder_id,omitempty"`
	CustomIcon  *string `json:"custom_icon,omitempty"`
	FaviconURL  *string `json:"favicon_url,omitempty"`
}

type UpdateLinkRequest struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description *string `json:"description,omitempty"`
	FolderID    *int64  `json:"folder_id,omitempty"`
	CustomIcon  *string `json:"custom_icon,omitempty"`
	FaviconURL  *string `json:"favicon_url,omitempty"`
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

	link, err := h.service.CreateLink(
		r.Context(),
		userID,
		req.Title,
		req.URL,
		req.Description,
		req.FolderID,
		req.CustomIcon,
		req.FaviconURL,
	)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, link)
}

// GET /links
func (h *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	includeStats := r.URL.Query().Get("stats") == "true"

	var folderID *int64
	rawFolderID := strings.TrimSpace(r.URL.Query().Get("folder_id"))
	if rawFolderID != "" {
		parsed, err := strconv.ParseInt(rawFolderID, 10, 64)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid folder_id"})
			return
		}
		folderID = &parsed
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	if includeStats {
		links, err := h.service.GetAllLinksWithStats(r.Context(), userID, folderID)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch links"})
			return
		}
		respondJSON(w, http.StatusOK, links)
	} else {
		links, err := h.service.GetAllLinks(r.Context(), userID, folderID)
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

	link, err := h.service.UpdateLink(
		r.Context(),
		userID,
		id,
		req.Title,
		req.URL,
		req.Description,
		req.FolderID,
		req.CustomIcon,
		req.FaviconURL,
	)
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

func extractIDFromPath(path, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	if idStr == path {
		return 0, strconv.ErrSyntax
	}

	idStr = strings.TrimSuffix(idStr, "/click")
	idStr = strings.TrimSuffix(idStr, "/")
	return strconv.ParseInt(idStr, 10, 64)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidURL):
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid URL format"})
	case errors.Is(err, services.ErrEmptyTitle):
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "title cannot be empty"})
	case errors.Is(err, services.ErrLinkNotFound):
		respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
	case errors.Is(err, services.ErrFolderNotFound):
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "folder not found"})
	case errors.Is(err, services.ErrIconTooLarge):
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "custom icon too large"})
	default:
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

