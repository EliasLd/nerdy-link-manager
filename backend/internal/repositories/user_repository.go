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

func (r *sqliteUserRepository) Create(ctx context.Context, email, passwordHash string) error {
	query := `
		INSERT INTO users(email, password_hash)
		VALUES (?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, email, passwordHash)
	return err

}

func (r *sqliteUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = ?
	`

	row := r.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *sqliteUserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
