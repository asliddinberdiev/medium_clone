package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	rdb, redisMock := redismock.NewClientMock()
	defer rdb.Close()
	sqlDB, mockDB, _ := sqlmock.New() 
	defer sqlDB.Close()

	db := sqlx.NewDb(sqlDB, "postgres")

	repo := repository.NewRepository(db, rdb)

	assert.NotNil(t, repo.User, "User repository should not be nil")
	assert.NotNil(t, repo.Auth, "Auth repository should not be nil")

	redisMock.ExpectationsWereMet()
	mockDB.ExpectationsWereMet()   
}
