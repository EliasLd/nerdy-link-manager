package repositories

import (
	"context"
	"database/sql"
	"time"
)

type Folder struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FolderRepository interface {
	Create(ctx context.Context, userID int64, name string) (*Folder, error)
	FindByID(ctx context.Context, userID, folderID int64) (*Folder, error)
	FindAll(ctx context.Context, userID int64) ([]Folder, error)
	Update(ctx context.Context, userID, folderID int64, name string) (*Folder, error)
	Delete(ctx context.Context, userID, folderID int64) error
}

type sqliteFolderRepository struct {
	db *sql.DB
}

func NewFolderRepository(db *sql.DB) FolderRepository {
	return &sqliteFolderRepository{db: db}
}

func (r *sqliteFolderRepository) Create(ctx context.Context, userID int64, name string) (*Folder, error) {
	q := `INSERT INTO folders(user_id, name) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, q, userID, name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, userID, id)
}

func (r *sqliteFolderRepository) FindByID(ctx context.Context, userID, folderID int64) (*Folder, error) {
	q := `
		SELECT id, user_id, name, created_at, updated_at
		FROM folders
		WHERE id = ? AND user_id = ?
	`
	row := r.db.QueryRowContext(ctx, q, folderID, userID)

	var f Folder
	if err := row.Scan(&f.ID, &f.UserID, &f.Name, &f.CreatedAt, &f.UpdatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *sqliteFolderRepository) FindAll(ctx context.Context, userID int64) ([]Folder, error) {
	q := `
		SELECT id, user_id, name, created_at, updated_at
		FROM folders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Folder{}
	for rows.Next() {
		var f Folder
		if err := rows.Scan(&f.ID, &f.UserID, &f.Name, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *sqliteFolderRepository) Update(ctx context.Context, userID, folderID int64, name string) (*Folder, error) {
	q := `
		UPDATE folders
		SET name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`
	res, err := r.db.ExecContext(ctx, q, name, folderID, userID)
	if err != nil {
		return nil, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if ra == 0 {
		return nil, sql.ErrNoRows
	}
	return r.FindByID(ctx, userID, folderID)
}

func (r *sqliteFolderRepository) Delete(ctx context.Context, userID, folderID int64) error {
	q := `DELETE FROM folders WHERE id = ? AND user_id = ?`
	res, err := r.db.ExecContext(ctx, q, folderID, userID)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return sql.ErrNoRows
	}
	return nil
}
