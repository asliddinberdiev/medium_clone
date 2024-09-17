package repository

import (
	"context"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type User interface {
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, req *models.User) error
	Delete(ctx context.Context, id string) error
}


type Post interface {
	Create(ctx context.Context, req *models.Post) (*models.Post, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*models.Post, error)
	GetAllPersonal(ctx context.Context, userID string, limit int, offset int) ([]*models.PersonalPost, error)
	Get(ctx context.Context, id string) (*models.Post, error)
	Update(ctx context.Context, req *models.UpdatePost) error
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	User
	Post
}

func NewRepository(db *sqlx.DB, log *zap.Logger) *Repository {
	return &Repository{
		User: NewUserRepository(db, log),
		Post: NewPostRepository(db, log),
	}
}
