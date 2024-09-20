package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) AddBlack(tokenID, token string, exp time.Duration) error {
	args := m.Called(tokenID, token, exp)
	return args.Error(0)
}

func (m *MockAuthRepo) GetBlackToken(tokenID string) (string, error) {
	args := m.Called(tokenID)
	return args.String(0), args.Error(1)
}

func TestAddBlack(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	authService := service.NewAuthService(mockRepo)

	t.Run("correct", func(t *testing.T) {
		tokenID := "tokenID123"
		token := "tokenABC"
		exp := time.Hour * 1
		mockRepo.On("AddBlack", tokenID, token, exp).Return(nil)

		err := authService.AddBlack(tokenID, token, exp)
		assert.NoError(t, err)
	})

	t.Run("incorrect", func(t *testing.T) {
		tokenID := "tokenID123"
		token := "tokenABC"
		exp := time.Hour * 1
		mockRepo.On("AddBlack", tokenID, token, exp).Return(assert.AnError)

		err := authService.AddBlack(tokenID, token, exp)
		assert.NoError(t, err)
	})

	mockRepo.AssertExpectations(t)
}

func TestHasBlackToken(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	authService := service.NewAuthService(mockRepo)

	t.Run("correct", func(t *testing.T) {
		tokenID := "tokenID123"
		mockRepo.On("GetBlackToken", tokenID).Return(tokenID, nil)

		result := authService.HasBlackToken(tokenID)
		assert.False(t, result)

		mockRepo.ExpectedCalls = nil
	})

	t.Run("error", func(t *testing.T) {
		tokenID := "tokenID123"
		mockRepo.On("GetBlackToken", tokenID).Return("", errors.New("test"))

		result := authService.HasBlackToken(tokenID)
		assert.True(t, result)

		mockRepo.ExpectedCalls = nil
	})

	mockRepo.AssertExpectations(t)
}

