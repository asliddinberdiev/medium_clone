package service_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	rdb, _ := redismock.NewClientMock()
	defer rdb.Close()
	sqlDB, _, _ := sqlmock.New()
	defer sqlDB.Close()
	db := sqlx.NewDb(sqlDB, "postgres")
	mockCfg := config.App{Port: "8000", Version: "v1", TokenKey: "secret", AccessTime: "15", RefreshTime: "30"}
	repo := repository.NewRepository(db, rdb)
	service := service.NewService(repo, mockCfg)

	t.Run("service", func(t *testing.T) {

		assert.NotNil(t, service.Auth)
		assert.NotNil(t, service.Token)
		assert.NotNil(t, service.User)
	})

}
