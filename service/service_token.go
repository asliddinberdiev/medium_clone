package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	cfg config.App
}

func NewTokenService(cfg config.App) *TokenService {
	return &TokenService{cfg: cfg}
}

func generate(userID, userRole, tokenType, sekretKey string, expireTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": userRole,
		"type": tokenType,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(expireTime).Unix(),
		"jti":  uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(sekretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *TokenService) AccessTokenGenerate(userID, userRole string) (string, error) {
	accessTime, err := strconv.Atoi(s.cfg.AccessTime)
	if err != nil {
		return "", err
	}
	return generate(userID, userRole, "access", s.cfg.TokenKey, time.Minute*time.Duration(accessTime))
}

func (s *TokenService) RefreshTokenGenerate(userID, userRole string) (string, error) {
	refreshTime, err := strconv.Atoi(s.cfg.RefreshTime)
	if err != nil {
		return "", err
	}
	return generate(userID, userRole, "refresh", s.cfg.TokenKey, time.Hour*time.Duration(refreshTime))
}

func (s *TokenService) Parse(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("an error occurred while validating the token")
		}
		return []byte(s.cfg.TokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := map[string]interface{}{
			"id":   claims["id"],
			"role": claims["role"],
			"type": claims["type"],
			"jti":  claims["jti"],
		}
		return result, nil
	}

	return nil, errors.New("invalid token")
}
