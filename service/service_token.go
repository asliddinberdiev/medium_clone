package service

import (
	"errors"
	"log"
	"strconv"
	"strings"
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
		log.Println("service_token:  generate - signedString error: ", err)
		return "", err
	}

	return signedToken, nil
}

func (s *TokenService) AccessTokenGenerate(userID, userRole string) (string, error) {
	accessTime, err := strconv.Atoi(s.cfg.AccessTime)
	if err != nil {
		log.Println("service_token: accessGenerate - time error: ", err)
		return "", err
	}
	return generate(userID, userRole, "access", s.cfg.TokenKey, time.Minute*time.Duration(accessTime))
}

func (s *TokenService) RefreshTokenGenerate(userID, userRole string) (string, error) {
	refreshTime, err := strconv.Atoi(s.cfg.RefreshTime)
	if err != nil {
		log.Println("service_token: refreshGenerate - time error: ", err)
		return "", err
	}
	return generate(userID, userRole, "refresh", s.cfg.TokenKey, time.Hour*time.Duration(refreshTime))
}

func (s *TokenService) Parse(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("service_token: parse - method error")
			return nil, errors.New("invalid token")
		}
		return []byte(s.cfg.TokenKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			log.Println("service_token: parse - token is expired: ", err)
			return nil, errors.New("token is expired")
		}
		log.Println("service_token: parse - error: ", err)
		return nil, errors.New("invalid token")
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

	log.Println("service_token: invalid")
	return nil, errors.New("invalid token")
}
