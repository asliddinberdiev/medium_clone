package service

import (
	"context"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"go.uber.org/zap"
)

type User interface {
	Create(ctx context.Context, user models.UserCreate) (*models.User, error)
}

type Post interface {
	Create(ctx context.Context, user models.UserCreate) (*models.User, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository, log *zap.Logger) *Service {
	return &Service{
		User: NewUserService(repo.User, log),
	}
}
