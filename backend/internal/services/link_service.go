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
	ErrIconTooLarge = errors.New("custom icon too large")
)

const maxIconLen = 300_000

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

func (s *LinkService) CreateLink(
	ctx context.Context,
	userID int64,
	title, rawURL string,
	description *string,
	folderID *int64,
	customIcon *string,
	faviconURL *string,
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

	if customIcon != nil {
		trimmed := strings.TrimSpace(*customIcon)
		if trimmed == "" {
			customIcon = nil
		} else if len(trimmed) > maxIconLen {
			return nil, ErrIconTooLarge
		} else {
			customIcon = &trimmed
		}
	}

	if faviconURL != nil {
		trimmed := strings.TrimSpace(*faviconURL)
		if trimmed == "" {
			faviconURL = nil
		} else {
			faviconURL = &trimmed
		}
	}

	if folderID != nil {
		if _, err := s.folderRepo.FindByID(ctx, userID, *folderID); err != nil {
			return nil, ErrFolderNotFound
		}
	}

	return s.repo.Create(ctx, userID, title, rawURL, description, folderID, customIcon, faviconURL)
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
	customIcon *string,
	faviconURL *string,
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

	if customIcon != nil {
		trimmed := strings.TrimSpace(*customIcon)
		if trimmed == "" {
			customIcon = nil
		} else if len(trimmed) > maxIconLen {
			return nil, ErrIconTooLarge
		} else {
			customIcon = &trimmed
		}
	}

	if faviconURL != nil {
		trimmed := strings.TrimSpace(*faviconURL)
		if trimmed == "" {
			faviconURL = nil
		} else {
			faviconURL = &trimmed
		}
	}

	if folderID != nil {
		if _, err := s.folderRepo.FindByID(ctx, userID, *folderID); err != nil {
			return nil, ErrFolderNotFound
		}
	}

	return s.repo.Update(ctx, userID, id, title, rawURL, description, folderID, customIcon, faviconURL)
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

func validateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ErrInvalidURL
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

