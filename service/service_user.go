package service

import (
	"context"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	repo repository.User
	log  *zap.Logger
}

func NewUserService(repo repository.User, log *zap.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (s *UserService) Create(ctx context.Context, user models.UserCreate) (*models.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		s.log.Error("service: uuid create err", zap.Error(err))
		return nil, err
	}

	hashPassword, err := utils.GeneratePasswordHash(user.Password)
	if err != nil {
		s.log.Error("service: password hashed err", zap.Error(err))
		return nil, err
	}

	newuser, err := s.repo.Create(ctx, models.User{ID: id.String(), FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: hashPassword, Role: user.Role})
	if err != nil {
		s.log.Error("service: user create err", zap.Error(err))
		return nil, err
	}

	return newuser, nil
}
