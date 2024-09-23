package repository_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		gofakeit.Seed(0)
		user := models.User{
			ID:        gofakeit.UUID(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, true, false, 16),
			Role:      "user",
		}

		mock.ExpectQuery(`INSERT INTO users`).WithArgs(user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role).
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
				AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Role, "2024-01-01", "2024-01-01"))

		newUser, err := userRepo.Create(user)
		assert.NoError(t, err)
		assert.Equal(t, newUser.ID, user.ID)
		assert.Equal(t, newUser.FirstName, user.FirstName)
		assert.Equal(t, newUser.LastName, user.LastName)
		assert.Equal(t, newUser.Email, user.Email)
	})

	t.Run("incorrect", func(t *testing.T) {
		gofakeit.Seed(0)
		user := models.User{
			ID:        gofakeit.UUID(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, true, false, 16),
			Role:      "user",
		}

		mock.ExpectQuery(`INSERT INTO users`).WithArgs(user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role).
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
				AddRow("12345", "fakeError", "fakeError", "fakeError", user.Role, "2024-01-01", "2024-01-01"))

		newUser, err := userRepo.Create(user)
		assert.NoError(t, err)
		assert.NotEqual(t, newUser.ID, user.ID)
		assert.NotEqual(t, newUser.FirstName, user.FirstName)
		assert.NotEqual(t, newUser.LastName, user.LastName)
		assert.NotEqual(t, newUser.Email, user.Email)
	})

	t.Run("exec_error", func(t *testing.T) {
		gofakeit.Seed(0)
		user := models.User{
			ID:        gofakeit.UUID(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, true, false, 16),
			Role:      "user",
		}

		mock.ExpectQuery(`INSERT INTO users`).WithArgs(user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role).
			WillReturnError(sql.ErrTxDone)

		newUser, err := userRepo.Create(user)
		assert.Error(t, err)
		assert.Nil(t, newUser)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, first_name, last_name, email, role, created_at, updated_at FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
				AddRow("1", "John", "Doe", "john@example.com", "user", "2024-01-01", "2024-01-01").
				AddRow("2", "Jane", "Doe", "jane@example.com", "admin", "2024-01-01", "2024-01-01"))

		users, err := userRepo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("incorrect", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, first_name, last_name, email, role, created_at, updated_at FROM users`).
			WillReturnError(sql.ErrNoRows)

		users, err := userRepo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, users)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, first_name, last_name, email, role, created_at, updated_at FROM users WHERE id = \$1`).
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
				AddRow("1", "John", "Doe", "john@example.com", "user", "2024-01-01", "2024-01-01"))

		user, err := userRepo.GetByID("1")
		assert.NotEqual(t, err, sql.ErrNoRows)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, user.ID, "1")
	})

	t.Run("incorrect", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, first_name, last_name, email, role, created_at, updated_at FROM users WHERE id = \$1`).
			WithArgs("1").
			WillReturnError(sql.ErrNoRows)

		user, err := userRepo.GetByID("1")
		assert.Equal(t, err, sql.ErrNoRows)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		email := "john@example.com"
		mock.ExpectQuery(`SELECT id, first_name, last_name, email, password, role, created_at, updated_at FROM users WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "role", "created_at", "updated_at"}).
				AddRow("1", "John", "Doe", email, "password", "user", "2024-01-01", "2024-01-01"))

		user, err := userRepo.GetByEmail(email)
		assert.NoError(t, err)
		assert.NotEqual(t, err, sql.ErrNoRows)
		assert.Equal(t, email, user.Email)
	})

	t.Run("incorrect", func(t *testing.T) {
		email := "john@example.com"
		mock.ExpectQuery(`SELECT id, first_name, last_name, email, password, role, created_at, updated_at FROM users WHERE email = \$1`).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := userRepo.GetByEmail(email)
		assert.Error(t, err)
		assert.Equal(t, err, sql.ErrNoRows)
		assert.Nil(t, user)
	})

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		gofakeit.Seed(0)
		user := models.UpdateUser{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Role:      "user",
		}

		mock.ExpectQuery(`UPDATE users SET first_name = \$1, last_name = \$2, role = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, first_name, last_name, email, role, created_at, updated_at`).
			WithArgs(user.FirstName, user.LastName, user.Role, sqlmock.AnyArg(), "1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
				AddRow("1", user.FirstName, user.LastName, "john@example.com", user.Role, "2024-01-01", "2024-01-01"))

		newUser, err := userRepo.Update("1", user)
		assert.NoError(t, err)
		assert.NotNil(t, newUser)
		assert.Equal(t, user.FirstName, newUser.FirstName)
		assert.Equal(t, user.LastName, newUser.LastName)
	})

	t.Run("incorrect", func(t *testing.T) {
		gofakeit.Seed(0)
		user := models.UpdateUser{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Role:      "user",
		}

		mock.ExpectQuery(`UPDATE users SET first_name = \$1, last_name = \$2, role = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, first_name, last_name, email, role, created_at, updated_at`).
			WithArgs(user.FirstName, user.LastName, user.Role, sqlmock.AnyArg(), "1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
				AddRow("1", "test", "incorrect", "john@example.com", user.Role, "2024-01-01", "2024-01-01"))

		newUser, err := userRepo.Update("1", user)
		assert.NoError(t, err)
		assert.NotNil(t, newUser)
		assert.NotEqual(t, user.FirstName, newUser.FirstName)
		assert.NotEqual(t, user.LastName, newUser.LastName)
	})

	t.Run("not_found", func(t *testing.T) {
		gofakeit.Seed(0)
		user := models.UpdateUser{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Role:      "user",
		}

		mock.ExpectQuery(`UPDATE users SET first_name = \$1, last_name = \$2, role = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, first_name, last_name, email, role, created_at, updated_at`).
			WithArgs(user.FirstName, user.LastName, user.Role, sqlmock.AnyArg(), "1").
			WillReturnError(sql.ErrNoRows)

		newUser, err := userRepo.Update("1", user)
		assert.Error(t, err)
		assert.Nil(t, newUser)
		assert.Equal(t, err, sql.ErrNoRows)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
			WithArgs("1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := userRepo.Delete("1")
		assert.NoError(t, err)
	})

	t.Run("incorrect", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
			WithArgs("1").
			WillReturnError(sql.ErrNoRows)

		err := userRepo.Delete("1")
		assert.Error(t, err)
		assert.Equal(t, err, sql.ErrNoRows)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
