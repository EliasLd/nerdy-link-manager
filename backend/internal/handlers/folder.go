package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/EliasLd/nerdy-link-manager/internal/middleware"
	"github.com/EliasLd/nerdy-link-manager/internal/services"
)

type FolderHandler struct {
	service *services.FolderService
}

func NewFolderHandler(s *services.FolderService) *FolderHandler {
	return &FolderHandler{service: s}
}

type folderReq struct {
	Name string `json:"name"`
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var req folderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	f, err := h.service.CreateFolder(r.Context(), userID, req.Name)
	if err != nil {
		if errors.Is(err, services.ErrEmptyFolderName) {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "folder name cannot be empty"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create folder"})
		return
	}
	respondJSON(w, http.StatusCreated, f)
}

func (h *FolderHandler) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	folders, err := h.service.GetFolders(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch folders"})
		return
	}
	respondJSON(w, http.StatusOK, folders)
}

func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/folders/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid folder ID"})
		return
	}
	var req folderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	f, err := h.service.UpdateFolder(r.Context(), userID, id, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyFolderName):
			respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "folder name cannot be empty"})
		case errors.Is(err, services.ErrFolderNotFound):
			respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "folder not found"})
		default:
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to update folder"})
		}
		return
	}
	respondJSON(w, http.StatusOK, f)
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/folders/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid folder ID"})
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	err = h.service.DeleteFolder(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, services.ErrFolderNotFound) {
			respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "folder not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to delete folder"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseInt64(s string) (int64, error) { return strconv.ParseInt(strings.TrimSpace(s), 10, 64) }
