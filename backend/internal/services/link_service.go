package services

import (
	"errors"

	"github.com/EliasLd/nerdy-link-manager/internal/repositories"
)

var (
	ErrInvalidURL   = errors.New("invalid URL format")
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrLinkNotFound = errors.New("link not found")
)

type LinkService struct {
	repo repositories.LinkRepository
}

func NewLinkService(repo repositories.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}
