package service

import (
	"time"

	"github.com/asliddinberdiev/medium_clone/config"
	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
)

type User interface {
	Create(user models.UserCreate) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id string, user models.UpdateUser) (*models.User, error)
	Delete(id string) error
}

type Auth interface {
	AddBlack(tokenID, token string, exp time.Duration) error
	HasBlackToken(tokenID string) bool
}

type Token interface {
	AccessTokenGenerate(userID string) (string, error)
	RefreshTokenGenerate(userID string) (string, error)
	Parse(tokenString string) (map[string]interface{}, error)
}

type Post interface {
	Create(user_id string, post models.CreatePost) (*models.Post, error)
	GetByID(id string) (*models.Post, error)
	GetPersonal(user_id string) ([]*models.Post, error)
	GetAll() ([]*models.Post, error)
	Update(id string, post models.UpdatePost) (*models.Post, error)
	Delete(id string) error
}

type Service struct {
	User
	Auth
	Token
	Post
}

func NewService(repo *repository.Repository, cfg config.App) *Service {
	return &Service{
		User:  NewUserService(repo.User),
		Auth:  NewAuthService(repo.Auth),
		Token: NewTokenService(cfg),
		Post:  NewPostService(repo.Post),
	}
}
