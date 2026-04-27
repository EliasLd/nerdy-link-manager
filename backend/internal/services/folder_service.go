package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/EliasLd/nerdy-link-manager/internal/repositories"
)

var (
	ErrEmptyFolderName = errors.New("folder name cannot be empty")
	ErrFolderNotFound  = errors.New("folder not found")
)

type FolderService struct {
	repo repositories.FolderRepository
}

func NewFolderService(repo repositories.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) CreateFolder(ctx context.Context, userID int64, name string) (*repositories.Folder, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyFolderName
	}
	return s.repo.Create(ctx, userID, name)
}

func (s *FolderService) GetFolders(ctx context.Context, userID int64) ([]repositories.Folder, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *FolderService) UpdateFolder(ctx context.Context, userID, folderID int64, name string) (*repositories.Folder, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyFolderName
	}
	f, err := s.repo.Update(ctx, userID, folderID, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}
	return f, nil
}

func (s *FolderService) DeleteFolder(ctx context.Context, userID, folderID int64) error {
	err := s.repo.Delete(ctx, userID, folderID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return ErrFolderNotFound
	}
	return err
}
