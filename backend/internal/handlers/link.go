package handlers

import (
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
