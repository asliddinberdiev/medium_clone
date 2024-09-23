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

type Comment interface {
	Create(comment models.Comment) (*models.Comment, error)
	GetAll(post_id string) ([]*models.Comment, error)
	GetByID(id string) (*models.Comment, error)
	Update(id, body string) (*models.Comment, error)
	Delete(id string) error
}

type SavedPost interface {
	Add(savedPost models.SavedPost) error
	Remove(post_id string) error
	GetByID(user_id, post_id string) (*models.SavedPost, error)
	GetAll(user_id string) ([]*models.Post, error)
}

type Repository struct {
	User
	Auth
	Post
	Comment
	SavedPost
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Auth:    NewAuthRepository(db, rdb),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
		SavedPost: NewSavedPostRepository(db),
	}
}
