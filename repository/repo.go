package repository

import (
	"context"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user models.User) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error)
	Update(id string, req models.UpdateUser) (*models.User, error)
	Delete(id string) error
}

type Post interface {
	Create(ctx context.Context, req *models.Post) (*models.Post, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*models.Post, error)
	GetAllPersonal(ctx context.Context, userID string, limit int, offset int) ([]*models.PersonalPost, error)
	GetByID(ctx context.Context, id string) (*models.Post, error)
	Update(ctx context.Context, req *models.UpdatePost) error
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	User
	Post
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		Post: NewPostRepository(db),
	}
}
