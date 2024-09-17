package service

import (
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

func (s *UserService) Create(user models.UserCreate) (*models.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("service: user create uuid error: ", err)
		return nil, err
	}

	hashPassword, err := utils.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Println("service: user create password hashed error: ", err)
		return nil, err
	}

	if user.Role == "" {
		user.Role = "user"
	}

	newUser, err := s.repo.Create(models.User{ID: id.String(), FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: hashPassword, Role: user.Role})
	if err != nil {
		log.Println("service: user create repo: ", err)
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		log.Println("service: user getByID repo: ", err)
		return nil, err
	}

	return user, nil
}
