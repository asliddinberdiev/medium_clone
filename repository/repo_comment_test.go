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

func TestCommentRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewCommentRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		gofakeit.Seed(0)
		comment := models.Comment{
			ID:     gofakeit.UUID(),
			UserID: gofakeit.UUID(),
			PostID: gofakeit.UUID(),
			Body:   gofakeit.ProductDescription(),
		}

		mock.ExpectQuery(`INSERT INTO comments`).WithArgs(comment.ID, comment.UserID, comment.PostID, comment.Body).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "body", "created_at"}).
				AddRow(comment.ID, comment.UserID, comment.PostID, comment.Body, "2024-01-01"))

		newComment, err := repo.Create(comment)
		assert.NoError(t, err)
		assert.NotNil(t, newComment)
		assert.Equal(t, newComment.ID, comment.ID)
		assert.Equal(t, newComment.UserID, comment.UserID)
		assert.Equal(t, newComment.PostID, comment.PostID)
		assert.Equal(t, newComment.Body, comment.Body)
	})

	t.Run("incorrect", func(t *testing.T) {
		gofakeit.Seed(0)
		comment := models.Comment{
			ID:     gofakeit.UUID(),
			UserID: gofakeit.UUID(),
			PostID: gofakeit.UUID(),
			Body:   gofakeit.ProductDescription(),
		}

		mock.ExpectQuery(`INSERT INTO comments`).WithArgs(comment.ID, comment.UserID, comment.PostID, comment.Body).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "body", "created_at"}).
				AddRow("12345", "123456", "1234567", "fakeError", "2024-01-01"))

		newComment, err := repo.Create(comment)

		assert.NoError(t, err)
		assert.NotNil(t, newComment)
	})

	t.Run("exec_error", func(t *testing.T) {
		comment := models.Comment{
			ID:     gofakeit.UUID(),
			UserID: gofakeit.UUID(),
			PostID: gofakeit.UUID(),
			Body:   gofakeit.ProductDescription(),
		}

		mock.ExpectQuery(`INSERT INTO comments`).WithArgs(comment.ID, comment.UserID, comment.PostID, comment.Body).
			WillReturnError(sql.ErrTxDone)

		newComment, err := repo.Create(comment)
		assert.Error(t, err)
		assert.Nil(t, newComment)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCommentRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewCommentRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		post_id := "99"
		mock.ExpectQuery(`SELECT id, user_id, post_id, body, created_at FROM comments WHERE post_id = \$1`).
			WithArgs(post_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "body", "created_at"}).
				AddRow("1", "11", post_id, "John1", "2024-01-01").
				AddRow("2", "22", post_id, "John2", "2024-01-01"))

		list, err := repo.GetAll(post_id)
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Len(t, list, 2)
	})

	t.Run("incorrect", func(t *testing.T) {
		post_id := "99"
		mock.ExpectQuery(`SELECT id, user_id, post_id, body, created_at FROM comments WHERE post_id = \$1`).
			WithArgs(post_id).
			WillReturnError(sql.ErrNoRows)

		list, err := repo.GetAll(post_id)
		assert.Error(t, err)
		assert.Nil(t, list)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCommentRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewCommentRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		comment_id := "11"
		mock.ExpectQuery(`SELECT id, user_id, post_id, body, created_at FROM comments WHERE id = \$1`).
			WithArgs(comment_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "body", "created_at"}).
				AddRow(comment_id, "11", "111", "John", "2024-01-01").
				AddRow("2", "22", "222", "John2", "2024-01-01"))

		item, err := repo.GetByID(comment_id)
		assert.NotEqual(t, err, sql.ErrNoRows)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, item.ID, comment_id)
	})

	t.Run("incorrect", func(t *testing.T) {
		comment_id := "11"
		mock.ExpectQuery(`SELECT id, user_id, post_id, body, created_at FROM comments WHERE id = \$1`).
			WithArgs(comment_id).
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetByID(comment_id)
		assert.Equal(t, err, sql.ErrNoRows)
		assert.Error(t, err)
		assert.Nil(t, item)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCommentRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewCommentRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		comment_id := "1"
		comment_body := "test"
		mock.ExpectQuery(`UPDATE comments SET body = \$1 WHERE id = \$2 RETURNING id, user_id, post_id, body, created_at`).
			WithArgs(comment_body, comment_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "body", "created_at"}).
				AddRow(comment_id, "123", "1234", comment_body, "2024-01-01"))

		newComment, err := repo.Update(comment_id, comment_body)
		assert.NoError(t, err)
		assert.NotNil(t, newComment)
		assert.Equal(t, comment_id, newComment.ID)
		assert.Equal(t, comment_body, newComment.Body)
	})

	t.Run("incorrect", func(t *testing.T) {
		comment_id := "1"
		comment_body := "test"
		mock.ExpectQuery(`UPDATE comments SET body = \$1 WHERE id = \$2 RETURNING id, user_id, post_id, body, created_at`).
			WithArgs(comment_body, comment_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "body", "created_at"}).
				AddRow(comment_id, "123", "1234", "fake", "2024-01-01"))

		newComment, err := repo.Update(comment_id, comment_body)
		assert.NoError(t, err)
		assert.NotNil(t, newComment)
		assert.Equal(t, comment_id, newComment.ID)
		assert.NotEqual(t, comment_body, newComment.Body)
	})

	t.Run("not_found", func(t *testing.T) {
		comment_id := "1"
		comment_body := "test"
		mock.ExpectQuery(`UPDATE comments SET body = \$1 WHERE id = \$2 RETURNING id, user_id, post_id, body, created_at`).
			WithArgs(comment_body, comment_id).
			WillReturnError(sql.ErrNoRows)

		newComment, err := repo.Update(comment_id, comment_body)
		assert.Error(t, err)
		assert.Nil(t, newComment)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCommentRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewCommentRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		comment_id := "1"
		mock.ExpectExec(`DELETE FROM comments WHERE id = \$1`).
			WithArgs(comment_id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(comment_id)
		assert.NoError(t, err)
	})

	t.Run("incorrect", func(t *testing.T) {
		comment_id := "1"
		mock.ExpectExec(`DELETE FROM comments WHERE id = \$1`).
			WithArgs(comment_id).
			WillReturnError(sql.ErrNoRows)

		err := repo.Delete(comment_id)
		assert.Error(t, err)
		assert.Equal(t, err, sql.ErrNoRows)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
