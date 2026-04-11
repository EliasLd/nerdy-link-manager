package services

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/EliasLd/nerdy-link-manager/internal/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, email, string(hash))
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*repositories.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) BootstrapInitialUser(ctx context.Context, email, password string) (bool, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	count, err := s.repo.Count(ctx)
	if err != nil {
		return false, err
	}

	if count > 0 {
		// Already initialized
		return false, nil
	}

	// No user exists: bootstrap is mandatory
	if email == "" || password == "" {
		return false, errors.New("No users found in database: INITIAL_ADMIN_EMAIL and INITIAL_ADMIN_PASASWORD are required environment variables")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	if err := s.repo.Create(ctx, email, string(hash)); err != nil {
		return false, err
	}

	return true, nil
}
