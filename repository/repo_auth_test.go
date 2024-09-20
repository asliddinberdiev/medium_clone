package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository(t *testing.T) {
	rdb, redisMock := redismock.NewClientMock()
	defer rdb.Close()

	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "postgres")

	repo := repository.NewAuthRepository(db, rdb)

	t.Run("AddBlack", func(t *testing.T) {
		t.Run("correct", func(t *testing.T) {
			tokenID := "token_id"
			token := "token_value"
			exp := 10 * time.Minute
			redisMock.ExpectSet(tokenID, token, exp).SetVal("OK")

			err := repo.AddBlack(tokenID, token, exp)
			assert.NoError(t, err)
		})

		t.Run("incorrect", func(t *testing.T) {
			tokenID := "token_id"
			token := "token_value"
			exp := 10 * time.Minute
			redisMock.ExpectSet(tokenID, token, exp).SetErr(assert.AnError)

			err := repo.AddBlack(tokenID, token, exp)

			assert.Error(t, err)
			assert.Equal(t, err, assert.AnError)
		})
	})

	t.Run("GetBlackToken", func(t *testing.T) {
		t.Run("correct", func(t *testing.T) {
			tokenID := "token_id"
			token := "token_value"
			redisMock.ExpectGet(tokenID).SetVal(token)

			result, err := repo.GetBlackToken(tokenID)

			assert.NoError(t, err)
			assert.Equal(t, token, result)
		})

		t.Run("incorrect", func(t *testing.T) {
			tokenID := "token_id"
			token := "token_value"
			redisMock.ExpectGet(tokenID).SetErr(assert.AnError)

			result, err := repo.GetBlackToken(tokenID)

			assert.Error(t, err)
			assert.Equal(t, err, assert.AnError)
			assert.NotEqual(t, result, token)
			assert.Empty(t, result)
		})

		t.Run("not_found", func(t *testing.T) {
			tokenID := "token_id"
			redisMock.ExpectGet(tokenID).RedisNil()

			result, err := repo.GetBlackToken(tokenID)

			assert.Error(t, err)
			assert.Equal(t, err, redis.Nil)
			assert.Empty(t, result)

			err = redisMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})

	})

	err := redisMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
