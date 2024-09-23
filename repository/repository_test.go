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

	t.Run("repository", func(t *testing.T) {
		repo := repository.NewRepository(db, rdb)

		assert.NotNil(t, repo.User)
		assert.NotNil(t, repo.Auth)
		assert.NotNil(t, repo.Post)

		err := redisMock.ExpectationsWereMet()
		assert.NoError(t, err)
		err = mockDB.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
