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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LinkWithStats struct {
	Link
	ClickCount int64 `json:"click_count"`
}

type LinkRepository interface {
	// CRUD operations
	Create(ctx context.Context, title, url string, description *string) (*Link, error)
	FindByID(ctx context.Context, id int64) (*Link, error)
	FindAll(ctx context.Context) ([]Link, error)
	Update(ctx context.Context, id int64, title, url string, description *string) (*Link, error)
	Delete(ctx context.Context, id int64) error

	// Stats operations
	RecordClick(ctx context.Context, linkID int64) error
	GetLinkStats(ctx context.Context, linkID int64) (*LinkWithStats, error)
	GetAllLinksWithStats(ctx context.Context) ([]LinkWithStats, error)
}

type sqliteLinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return &sqliteLinkRepository{db: db}
}

func (r *sqliteLinkRepository) Create(ctx context.Context, title, url string, description *string) (*Link, error) {
	query := `
		INSERT INTO links(title, url, description)
		VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, title, url, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

func (r *sqliteLinkRepository) FindByID(ctx context.Context, id int64) (*Link, error) {
	query := `
		SELECT id, title, url, description, created_at, updated_at
		FROM links
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var link Link
	// Description field can be omitted
	var description sql.NullString

	err := row.Scan(&link.ID, &link.Title, &link.URL, &description, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		link.Description = &description.String
	}

	return &link, nil
}

// Retrieves all links, sorted by creation date DESC
func (r *sqliteLinkRepository) FindAll(ctx context.Context) ([]Link, error) {
	query := `
		SELECT id, title, url, description, created_at, updated_at
		FROM links
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []Link{}

	for rows.Next() {
		var link Link
		// Description field can be omitted
		var description sql.NullString

		err := rows.Scan(&link.ID, &link.Title, &link.URL, &description, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			link.Description = &description.String
		}

		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *sqliteLinkRepository) Update(ctx context.Context, id int64, title, url string, description *string) (*Link, error) {
	query := `
		UPDATE links
		SET title = ?, url = ?, description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, title, url, description, id)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

// Deletes a link (Cascade deletes its clicks too)
func (r *sqliteLinkRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM links WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *sqliteLinkRepository) RecordClick(ctx context.Context, linkID int64) error {
	query := `
		INSERT INTO link_clicks(link_id)
		VALUES (?)
	`

	_, err := r.db.ExecContext(ctx, query, linkID)
	return err
}

// Retrieves a link with its clicks data
func (r *sqliteLinkRepository) GetLinkStats(ctx context.Context, linkID int64) (*LinkWithStats, error) {
	query := `
		SELECT
			l.id, l.title, l.url, l.description, l.created_at, l.updated_at,
			COUNT(lc.id) as click_count
		FROM links l
		LEFT JOIN link_clicks lc ON l.id = lc.link_id
		WHERE l.id = ?
		GROUP BY l.id
	`

	row := r.db.QueryRowContext(ctx, query, linkID)

	var linkStats LinkWithStats
	var description sql.NullString

	err := row.Scan(
		&linkStats.ID,
		&linkStats.Title,
		&linkStats.URL,
		&description,
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

	return &linkStats, nil

}

// Retrieves all links with their data
func (r *sqliteLinkRepository) GetAllLinksWithStats(ctx context.Context) ([]LinkWithStats, error) {
	query := `
		SELECT
			l.id, l.title, l.url, l.description, l.created_at, l.updated_at,
			COUNT(lc.id) as click_count
		FROM links l
		LEFT JOIN link_clicks lc ON l.id = lc.link_id
		GROUP BY l.id
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []LinkWithStats{}

	for rows.Next() {
		var linkStats LinkWithStats
		var description sql.NullString

		err := rows.Scan(
			&linkStats.ID,
			&linkStats.Title,
			&linkStats.URL,
			&description,
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

		links = append(links, linkStats)
	}

	return links, rows.Err()
}
