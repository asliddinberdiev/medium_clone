package service

import (
	"context"

	"log"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user models.UserCreate) (*models.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("service: uuid create: %v\n", err)
		return nil, err
	}

	hashPassword, err := utils.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Printf("service: password hashed: %v", err)
		return nil, err
	}

	if user.Role == "" {
		user.Role = "user"
	}

	newUser, err := s.repo.Create(ctx, models.User{ID: id.String(), FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: hashPassword, Role: user.Role})
	if err != nil {
		log.Printf("service: user create: %v\n", err)
		return nil, err
	}

	return newUser, nil
}
