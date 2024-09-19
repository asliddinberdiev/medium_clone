package repository_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	gofakeit.Seed(0)
	user := models.User{
		ID:        gofakeit.UUID(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, true, false, 16),
		Role:      "user",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO users`).WithArgs(user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
			AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Role, "2024-01-01", "2024-01-01"))
	mock.ExpectCommit()

	createdUser, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if createdUser.ID != user.ID {
		t.Fatalf("expected user ID %s, but got %s", user.ID, createdUser.ID)
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, role, created_at, updated_at FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
			AddRow("1", "John", "Doe", "john@example.com", "user", "2024-01-01", "2024-01-01").
			AddRow("2", "Jane", "Doe", "jane@example.com", "admin", "2024-01-01", "2024-01-01"))

	users, err := userRepo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, but got %d", len(users))
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, role, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
			AddRow("1", "John", "Doe", "john@example.com", "user", "2024-01-01", "2024-01-01").
			AddRow("2", "John", "Doe", "john@example.com", "user", "2024-01-01", "2024-01-01"))

	user, err := userRepo.GetByID("1")
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if user.ID != "1" {
		t.Fatalf("expected user ID %s, but got %s", "1", user.ID)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, password, role, created_at, updated_at FROM users WHERE email = \$1`).
		WithArgs("john@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "role", "created_at", "updated_at"}).
			AddRow("1", "John", "Doe", "john@example.com", "password", "user", "2024-01-01", "2024-01-01"))

	user, err := userRepo.GetByEmail("john@example.com")
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if user.Email != "john@example.com" {
		t.Fatalf("expected user email %s, but got %s", "john@example.com", user.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	gofakeit.Seed(0)
	updateUser := models.UpdateUser{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Role:      "user",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE users SET first_name = \$1, last_name = \$2, role = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, first_name, last_name, email, role, created_at, updated_at`).
		WithArgs(updateUser.FirstName, updateUser.LastName, updateUser.Role, sqlmock.AnyArg(), "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role", "created_at", "updated_at"}).
			AddRow("1", updateUser.FirstName, updateUser.LastName, "john@example.com", updateUser.Role, "2024-01-01", time.Now()))
	mock.ExpectCommit()

	updatedUser, err := userRepo.Update("1", updateUser)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if updatedUser.FirstName != updateUser.FirstName {
		t.Fatalf("expected user first name %s, but got %s", updateUser.FirstName, updatedUser.FirstName)
	}
	if updatedUser.LastName != updateUser.LastName {
		t.Fatalf("expected user last name %s, but got %s", updateUser.LastName, updatedUser.LastName)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Delete("1")
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	userRepo := repository.NewUserRepository(sqlxDB)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs("999").                          
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := userRepo.Delete("999")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, but got %v", err)
	}
}
