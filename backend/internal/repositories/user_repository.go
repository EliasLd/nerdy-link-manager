package repositories

import (
	"context"
	"database/sql"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    string
}

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

type sqliteUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

// TODO: Implement each function
func (r *sqliteUserRepository) Create(ctx context.Context, email, passwordHash string) error
func (r *sqliteUserRepository) FindByEmail(ctx context.Context, email string) (*User, error)
func (r *sqliteUserRepository) FindByID(ctx context.Context, id int64) (*User, error)
