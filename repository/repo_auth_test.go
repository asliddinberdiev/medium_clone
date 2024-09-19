package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asliddinberdiev/medium_clone/repository"
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
		// create mock values
		tokenID := "token_id"
		token := "token_value"
		exp := 10 * time.Minute

		// create value in db
		redisMock.ExpectSet(tokenID, token, exp).SetVal("OK")

		// set values in db
		err := repo.AddBlack(tokenID, token, exp)
		// checking error
		assert.NoError(t, err)
		// waiting for all actions to be completed
		redisMock.ExpectationsWereMet()
	})

	t.Run("GetBlackToken", func(t *testing.T) {
		// create mock values
		tokenID := "token_id"
		token := "token_value"

		// create value in db
		redisMock.ExpectGet(tokenID).SetVal(token)

		// get value in db
		result, err := repo.GetBlackToken(tokenID)
		// checking error
		assert.NoError(t, err)
		// checking set value and get value
		assert.Equal(t, token, result)
		// waiting for all actions to be completed
		redisMock.ExpectationsWereMet()
	})

	t.Run("GetBlackTokenNotFound", func(t *testing.T) {
		// create invalid mock value
		tokenID := "not_found_token_id"

		// create nil value in db
		redisMock.ExpectGet(tokenID).RedisNil()

		// get value in db
		result, err := repo.GetBlackToken(tokenID)
		// checking error
		assert.Error(t, err)
		// checking for a nil value
		assert.Empty(t, result)
		// waiting for all actions to be completed
		redisMock.ExpectationsWereMet()
	})
}
