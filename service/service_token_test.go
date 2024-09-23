package service_test

import (
	"io"
	"log"
	"testing"

	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenGenerate(t *testing.T) {
	log.SetOutput(io.Discard)

	t.Run("correct", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)

		userID := "user123"
		accessToken, err := tokenService.AccessTokenGenerate(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
	})

	t.Run("incorrect", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15.5", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)

		userID := "user123"
		accessToken, err := tokenService.AccessTokenGenerate(userID)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
	})
}

func TestRefreshTokenGenerate(t *testing.T) {
	log.SetOutput(io.Discard)

	t.Run("correct", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)

		userID := "user123"
		accessToken, err := tokenService.RefreshTokenGenerate(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
	})

	t.Run("incorrect", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30.3"}
		tokenService := service.NewTokenService(mockCfg)

		userID := "user123"
		refreshToken, err := tokenService.RefreshTokenGenerate(userID)

		assert.Error(t, err)
		assert.Empty(t, refreshToken)
	})
}

func TestParse(t *testing.T) {
	log.SetOutput(io.Discard)

	t.Run("access_parse", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)
		userID := "user123"

		accessToken, err := tokenService.AccessTokenGenerate(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)

		t.Run("correct", func(t *testing.T) {
			claims, err := tokenService.Parse(accessToken)
			assert.NoError(t, err)
			assert.NotEmpty(t, claims)
			assert.Equal(t, claims["type"], "access")
			assert.Equal(t, claims["id"], userID)
			assert.NotEmpty(t, claims["jti"])
		})

		t.Run("incorrect", func(t *testing.T) {
			claims, err := tokenService.Parse(accessToken)
			assert.NoError(t, err)
			assert.NotEmpty(t, claims)
			assert.NotEqual(t, claims["type"], "test")
			assert.NotEqual(t, claims["id"], "test_id")
		})
	})

	t.Run("refresh_parse", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)
		userID := "user123"

		refreshToken, err := tokenService.RefreshTokenGenerate(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshToken)

		t.Run("correct", func(t *testing.T) {
			claims, err := tokenService.Parse(refreshToken)
			assert.NoError(t, err)
			assert.NotEmpty(t, claims)
			assert.Equal(t, claims["type"], "refresh")
			assert.Equal(t, claims["id"], userID)
			assert.NotEmpty(t, claims["jti"])
		})

		t.Run("incorrect", func(t *testing.T) {
			claims, err := tokenService.Parse(refreshToken)
			assert.NoError(t, err)
			assert.NotEmpty(t, claims)
			assert.NotEqual(t, claims["type"], "test")
			assert.NotEqual(t, claims["id"], "test_id")
		})
	})

	t.Run("invalid", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)

		invalidToken := "invalid_token"
		claims, err := tokenService.Parse(invalidToken)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid token")
		assert.Nil(t, claims)
	})

	t.Run("expired", func(t *testing.T) {
		mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "0", RefreshTime: "30"}
		tokenService := service.NewTokenService(mockCfg)
		userID := "user123"

		accessToken, err := tokenService.AccessTokenGenerate(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)

		claims, err := tokenService.Parse(accessToken)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "token is expired")
		assert.Nil(t, claims)
	})
}
