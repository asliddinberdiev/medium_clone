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

func TestSavedPostRepository_Add(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewSavedPostRepository(sqlxDB)

	gofakeit.Seed(0)
	savedPost := models.SavedPost{
		ID:     gofakeit.UUID(),
		UserID: gofakeit.UUID(),
		PostID: gofakeit.UUID(),
	}

	mock.ExpectExec(`INSERT INTO saved_posts`).
		WithArgs(savedPost.ID, savedPost.UserID, savedPost.PostID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Add(savedPost)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSavedPostRepository_Remove(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewSavedPostRepository(sqlxDB)

	gofakeit.Seed(0)
	postID := gofakeit.UUID()

	mock.ExpectExec(`DELETE FROM saved_posts WHERE post_id = \$1`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Remove(postID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSavedPostRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewSavedPostRepository(sqlxDB)

	gofakeit.Seed(0)
	ID := gofakeit.UUID()
	userID := gofakeit.UUID()
	postID := gofakeit.UUID()

	mock.ExpectQuery(`SELECT id, user_id, post_id FROM saved_posts WHERE user_id = \$1 AND post_id = \$2`).
		WithArgs(userID, postID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id"}).
			AddRow(ID, userID, postID))

	item, err := repo.GetByID(userID, postID)
	assert.NotEqual(t, err, sql.ErrNoRows)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.ID, ID)
	assert.Equal(t, item.UserID, userID)
	assert.Equal(t, item.PostID, postID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSavedPostRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewSavedPostRepository(sqlxDB)

	gofakeit.Seed(0)
	userID := gofakeit.UUID()
	mock.ExpectQuery(`SELECT p.* FROM posts p INNER JOIN saved_posts sp ON sp.post_id = p.id WHERE sp.user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
			AddRow("1", userID, "Doe1", "John1", true, "2024-01-01", "2024-01-01").
			AddRow("2", userID, "Doe2", "John2", false, "2024-01-01", "2024-01-01"))

	list, err := repo.GetAll(userID)
	assert.NotEqual(t, err, sql.ErrNoRows)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list, 2)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
