package repositories

import (
	"context"
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
	Create(ctx context.Context, title, url, string, description *string) (*Link, error)
	FindByID(ctx context.Context, id int64) (*Link, error)
	FindAll(ctx context.Context) ([]Link, error)
	Update(ctx context.Context, id int64, title, url, string, description *string) (*Link, error)
	Delete(ctx context.Context, id int64) error

	// Stats operations
	RecordClick(ctx context.Context, linkID int64) error
	GetLinkStats(ctx context.Context, linkID int64) (*LinkWithStats, error)
	GetAllLinkWithStats(ctx context.Context) ([]LinkWithStats, error)
}
