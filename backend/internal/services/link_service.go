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
func (s *LinkService) CreateLink(ctx context.Context, userID int64, title, rawURL string, description *string) (*repositories.Link, error) {
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

	return s.repo.Create(ctx, userID, title, rawURL, description)
}

func (s *LinkService) GetLink(ctx context.Context, userID int64, id int64) (*repositories.Link, error) {
	return s.repo.FindByID(ctx, userID, id)
}

func (s *LinkService) GetAllLinks(ctx context.Context, userID int64) ([]repositories.Link, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *LinkService) UpdateLink(ctx context.Context, userID int64, id int64, title, rawURL string, description *string) (*repositories.Link, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if err := validateURL(rawURL); err != nil {
		return nil, err
	}

	if description != nil {
		trimmed := strings.TrimSpace(*description)
		if trimmed == "" {
			description = nil
		} else {
			description = &trimmed
		}
	}

	return s.repo.Update(ctx, userID, id, title, rawURL, description)
}

func (s *LinkService) DeleteLink(ctx context.Context, userID int64, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *LinkService) TrackClick(ctx context.Context, userID int64, linkID int64) error {
	_, err := s.repo.FindByID(ctx, userID, linkID)
	if err != nil {
		return ErrLinkNotFound
	}

	return s.repo.RecordClick(ctx, userID, linkID)
}

func (s *LinkService) GetLinkWithStats(ctx context.Context, userID int64, linkID int64) (*repositories.LinkWithStats, error) {
	return s.repo.GetLinkStats(ctx, userID, linkID)
}

func (s *LinkService) GetAllLinksWithStats(ctx context.Context, userID int64) ([]repositories.LinkWithStats, error) {
	return s.repo.GetAllLinksWithStats(ctx, userID)
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
