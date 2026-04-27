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
	repo       repositories.LinkRepository
	folderRepo repositories.FolderRepository
}

func NewLinkService(repo repositories.LinkRepository, folderRepo repositories.FolderRepository) *LinkService {
	return &LinkService{
		repo:       repo,
		folderRepo: folderRepo,
	}
}

// Validates data and creates new link
// folderID is optional (nil => no folder)
func (s *LinkService) CreateLink(
	ctx context.Context,
	userID int64,
	title, rawURL string,
	description *string,
	folderID *int64,
) (*repositories.Link, error) {
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

	// Validate folder ownership if provided
	if folderID != nil {
		if _, err := s.folderRepo.FindByID(ctx, userID, *folderID); err != nil {
			return nil, ErrFolderNotFound
		}
	}

	return s.repo.Create(ctx, userID, title, rawURL, description, folderID)
}

func (s *LinkService) GetLink(ctx context.Context, userID int64, id int64) (*repositories.Link, error) {
	return s.repo.FindByID(ctx, userID, id)
}

func (s *LinkService) GetAllLinks(ctx context.Context, userID int64, folderID *int64) ([]repositories.Link, error) {
	return s.repo.FindAll(ctx, userID, folderID)
}

func (s *LinkService) UpdateLink(
	ctx context.Context,
	userID int64,
	id int64,
	title, rawURL string,
	description *string,
	folderID *int64,
) (*repositories.Link, error) {
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

	// folderID nil is valid => unlink from folder
	if folderID != nil {
		if _, err := s.folderRepo.FindByID(ctx, userID, *folderID); err != nil {
			return nil, ErrFolderNotFound
		}
	}

	return s.repo.Update(ctx, userID, id, title, rawURL, description, folderID)
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

func (s *LinkService) GetAllLinksWithStats(ctx context.Context, userID int64, folderID *int64) ([]repositories.LinkWithStats, error) {
	return s.repo.GetAllLinksWithStats(ctx, userID, folderID)
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
