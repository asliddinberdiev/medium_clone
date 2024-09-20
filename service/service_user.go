package service

import (
	"database/sql"
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
	_, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		log.Println("service_user: create - checking db email: ", err)
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("service_user: create - uuid error: ", err)
		return nil, err
	}

	hashPassword, err := utils.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Println("service_user: create - password hashed error: ", err)
		return nil, err
	}

	if user.Role == "" {
		user.Role = "user"
	}

	newUser, err := s.repo.Create(models.User{ID: id.String(), FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: hashPassword, Role: user.Role})
	if err != nil {
		log.Println("service_user: create - repo: ", err)
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetAll() ([]*models.User, error) {
	list, err := s.repo.GetAll()
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.User{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		log.Println("service_user: getByID - repo: ", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		log.Println("service_user: getByEmail - repo: ", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(id string, user models.UpdateUser) (*models.User, error) {
	dbUser, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("service_user: update - not found dbUser")
			return nil, err
		}
		return nil, err
	}

	if user.FirstName == "" {
		user.FirstName = dbUser.FirstName
	}
	if user.LastName == "" {
		user.LastName = dbUser.LastName
	}
	if user.Role == "" {
		user.Role = dbUser.Role
	}

	updateUser, err := s.repo.Update(id, user)
	if err != nil {
		log.Println("service_user: update - repo error: ", err)
		return nil, err
	}

	return updateUser, nil
}

func (s *UserService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		log.Println("service_user: delete - repo error: ", err)
		return err
	}
	return nil
}
