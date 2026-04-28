package repositories

import (
	"context"
	"database/sql"
	"time"
)

type Link struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description *string   `json:"description,omitempty"`
	FolderID    *int64    `json:"folder_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LinkWithStats struct {
	Link
	ClickCount int64 `json:"click_count"`
}

type LinkRepository interface {
	// CRUD operations
	Create(ctx context.Context, userID int64, title, url string, description *string, folderID *int64) (*Link, error)
	FindByID(ctx context.Context, userID int64, id int64) (*Link, error)
	FindAll(ctx context.Context, userID int64, folderID *int64) ([]Link, error)
	Update(ctx context.Context, userID int64, id int64, title, url string, description *string, folderID *int64) (*Link, error)
	Delete(ctx context.Context, userID int64, id int64) error

	// Stats operations
	RecordClick(ctx context.Context, userID int64, linkID int64) error
	GetLinkStats(ctx context.Context, userID int64, linkID int64) (*LinkWithStats, error)
	GetAllLinksWithStats(ctx context.Context, userID int64, folderID *int64) ([]LinkWithStats, error)
}

type sqliteLinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return &sqliteLinkRepository{db: db}
}

func (r *sqliteLinkRepository) Create(ctx context.Context, userID int64, title, url string, description *string, folderID *int64) (*Link, error) {
	query := `
		INSERT INTO links(title, url, description, user_id, folder_id)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, title, url, description, userID, folderID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, userID, id)
}

func (r *sqliteLinkRepository) FindByID(ctx context.Context, userID int64, id int64) (*Link, error) {
	query := `
		SELECT id, title, url, description, folder_id, created_at, updated_at
		FROM links
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id, userID)

	var link Link
	var description sql.NullString
	var folderID sql.NullInt64

	err := row.Scan(&link.ID, &link.Title, &link.URL, &description, &folderID, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		link.Description = &description.String
	}
	if folderID.Valid {
		link.FolderID = &folderID.Int64
	}

	return &link, nil
}

// Retrieves all links, sorted by creation date DESC
func (r *sqliteLinkRepository) FindAll(ctx context.Context, userID int64, folderID *int64) ([]Link, error) {
	query := `
		SELECT id, title, url, description, folder_id, created_at, updated_at
		FROM links
		WHERE user_id = ?
	`
	args := []any{userID}

	if folderID != nil {
		query += ` AND folder_id = ?`
		args = append(args, *folderID)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []Link{}

	for rows.Next() {
		var link Link
		var description sql.NullString
		var folderIDNull sql.NullInt64

		err := rows.Scan(&link.ID, &link.Title, &link.URL, &description, &folderIDNull, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			link.Description = &description.String
		}
		if folderIDNull.Valid {
			link.FolderID = &folderIDNull.Int64
		}

		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *sqliteLinkRepository) Update(ctx context.Context, userID int64, id int64, title, url string, description *string, folderID *int64) (*Link, error) {
	query := `
		UPDATE links
		SET title = ?, url = ?, description = ?, folder_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`

	res, err := r.db.ExecContext(ctx, query, title, url, description, folderID, id, userID)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindByID(ctx, userID, id)
}

// Deletes a link (Cascade deletes its clicks too)
func (r *sqliteLinkRepository) Delete(ctx context.Context, userID int64, id int64) error {
	query := `DELETE FROM links WHERE id = ? AND user_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *sqliteLinkRepository) RecordClick(ctx context.Context, userID int64, linkID int64) error {
	query := `
		INSERT INTO link_clicks(link_id)
		SELECT id
		FROM links
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, linkID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Retrieves a link with its clicks data
func (r *sqliteLinkRepository) GetLinkStats(ctx context.Context, userID int64, linkID int64) (*LinkWithStats, error) {
	query := `
		SELECT
			l.id, l.title, l.url, l.description, l.folder_id, l.created_at, l.updated_at,
			COUNT(lc.id) as click_count
		FROM links l
		LEFT JOIN link_clicks lc ON l.id = lc.link_id
		WHERE l.id = ? AND l.user_id = ?
		GROUP BY l.id
	`

	row := r.db.QueryRowContext(ctx, query, linkID, userID)

	var linkStats LinkWithStats
	var description sql.NullString
	var folderID sql.NullInt64

	err := row.Scan(
		&linkStats.ID,
		&linkStats.Title,
		&linkStats.URL,
		&description,
		&folderID,
		&linkStats.CreatedAt,
		&linkStats.UpdatedAt,
		&linkStats.ClickCount,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		linkStats.Description = &description.String
	}
	if folderID.Valid {
		linkStats.FolderID = &folderID.Int64
	}

	return &linkStats, nil
}

// Retrieves all links with their data
func (r *sqliteLinkRepository) GetAllLinksWithStats(ctx context.Context, userID int64, folderID *int64) ([]LinkWithStats, error) {
	query := `
		SELECT
			l.id, l.title, l.url, l.description, l.folder_id, l.created_at, l.updated_at,
			COUNT(lc.id) as click_count
		FROM links l
		LEFT JOIN link_clicks lc ON l.id = lc.link_id
		WHERE l.user_id = ?
	`
	args := []any{userID}

	if folderID != nil {
		query += ` AND l.folder_id = ?`
		args = append(args, *folderID)
	}

	query += `
		GROUP BY l.id
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []LinkWithStats{}

	for rows.Next() {
		var linkStats LinkWithStats
		var description sql.NullString
		var folderIDNull sql.NullInt64

		err := rows.Scan(
			&linkStats.ID,
			&linkStats.Title,
			&linkStats.URL,
			&description,
			&folderIDNull,
			&linkStats.CreatedAt,
			&linkStats.UpdatedAt,
			&linkStats.ClickCount,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			linkStats.Description = &description.String
		}
		if folderIDNull.Valid {
			linkStats.FolderID = &folderIDNull.Int64
		}

		links = append(links, linkStats)
	}

	return links, rows.Err()
}

