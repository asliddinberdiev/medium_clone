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

type Comment interface {
	Create(user_id string, comment models.CreateComment) (*models.Comment, error)
	GetAll(post_id string) ([]*models.Comment, error)
	GetByID(id string) (*models.Comment, error)
	Update(id, body string) (*models.Comment, error)
	Delete(comment_id string) error
}

type SavedPost interface {
	Add(savedPost models.SavedPostAction) error
	Remove(post_id string) error
	GetByID(user_id, post_id string) (*models.SavedPost, error)
	GetAll(user_id string) ([]*models.Post, error)
}

type Service struct {
	User
	Auth
	Token
	Post
	Comment
	SavedPost
}

func NewService(repo *repository.Repository, cfg config.App) *Service {
	return &Service{
		User:      NewUserService(repo.User),
		Auth:      NewAuthService(repo.Auth),
		Token:     NewTokenService(cfg),
		Post:      NewPostService(repo.Post),
		Comment:   NewCommentService(repo.Comment),
		SavedPost: NewSavedPostService(repo.SavedPost),
	}
}
