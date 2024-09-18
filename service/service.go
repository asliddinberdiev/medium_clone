package service

import (
	"github.com/asliddinberdiev/medium_clone/config"
	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
)

type User interface {
	Create(user models.UserCreate) (*models.User, error)
	GetByID(id string)(*models.User, error)
	Update(id string, user models.UpdateUser) (*models.User, error)
	Delete(id string) error
}

type Post interface {
}

type Token interface {
	AccessTokenGenerate(userID, userRole string) (string, error)
	RefreshTokenGenerate(userID, userRole string) (string, error)
	Parse(tokenString string) (map[string]interface{}, error)
}

type Service struct {
	User
	Post
	Token
}

func NewService(repo *repository.Repository, cfg config.App) *Service {
	return &Service{
		User:  NewUserService(repo.User),
		Post:  NewPostService(repo.Post),
		Token: NewTokenService(cfg),
	}
}
