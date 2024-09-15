package repo

import (
	"context"
	"time"
)

type PostStorageI interface {
	Create(ctx context.Context, req *Post) (*Post, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*Post, error)
	GetAllPersonal(ctx context.Context, userID string, limit int, offset int) ([]*PostPersonal, error)
	Get(ctx context.Context, id string) (*Post, error)
	Update(ctx context.Context, req *UpdatePost) error
	Delete(ctx context.Context, id string) error
}

type Post struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	Published bool      `db:"published"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostPersonal struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	Published bool      `db:"published"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdatePost struct {
	ID        string `db:"id"`
	UserID    string `db:"user_id"`
	Title     string `db:"title"`
	Body      string `db:"body"`
	Published bool   `db:"published"`
}
