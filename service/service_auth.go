package service

import (
	"log"
	"time"

	"github.com/asliddinberdiev/medium_clone/repository"
)

type AuthService struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) AddBlack(tokenID, token string, exp time.Duration) error {
	return s.repo.AddBlack(tokenID, token, exp)
}

func (s *AuthService) HasBlackToken(tokenID string) bool {
	_, err := s.repo.GetBlackToken(tokenID)
	if err == nil {
		return false
	}
	log.Println("service_auth: HasBlackToken - repo err: ", err)
	return true
}
