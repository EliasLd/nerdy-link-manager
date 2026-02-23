package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/EliasLd/nerdy-link-manager/internal/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// TODO: Implement these functions
func (s *UserService) Register(ctx context.Context, email, password string) error
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*repositories.User, error)
