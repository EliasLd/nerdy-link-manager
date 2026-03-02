package services

import (
	"context"
	"errors"
	"net/url"
	"strings"

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

// Validates data and creates new link
func (s *LinkService) CreateLink(ctx context.Context, title, rawURL string, description *string) (*repositories.Link, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if err := validateURL(rawURL); err != nil {
		return nil, err
	}

	// Clean description if exists
	if description != nil {
		trimmed := strings.TrimSpace(*description)
		if trimmed == "" {
			description = nil
		} else {
			description = &trimmed
		}
	}

	return s.repo.Create(ctx, title, rawURL, description)
}

func (s *LinkService) GetLink(ctx context.Context, id int64) (*repositories.Link, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *LinkService) GetAllLinks(ctx context.Context) ([]repositories.Link, error) {
	return s.repo.FindAll(ctx)
}

// Checks that URL is valid and complete
func validateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ErrInvalidURL
	}

	// Parse URL to check its validity
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	// Check that url contains a valid scheme (http/https)
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}
