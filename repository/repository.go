package repository

import (
	"time"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type User interface {
	Create(user models.User) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id string, req models.UpdateUser) (*models.User, error)
	Delete(id string) error
}

type Auth interface {
	AddBlack(tokenID, token string, exp time.Duration) error
	GetBlackToken(tokenID string) (string, error)
}

type Post interface {
	Create(post models.Post) (*models.Post, error)
	GetByID(id string) (*models.Post, error)
	GetPersonal(user_id string) ([]*models.Post, error)
	GetAll() ([]*models.Post, error)
	Update(id string, post models.UpdatePost) (*models.Post, error)
	Delete(id string) error
}

type Repository struct {
	User
	Auth
	Post
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		Auth: NewAuthRepository(db, rdb),
		Post: NewPostRepository(db),
	}
}
