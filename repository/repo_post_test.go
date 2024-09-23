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

func TestPostRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewPostRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		gofakeit.Seed(0)
		post := models.Post{
			ID:        gofakeit.UUID(),
			UserID:    gofakeit.UUID(),
			Title:     gofakeit.BookTitle(),
			Body:      gofakeit.ProductDescription(),
			Published: gofakeit.Bool(),
		}

		mock.ExpectQuery(`INSERT INTO posts`).WithArgs(post.ID, post.UserID, post.Title, post.Body, post.Published).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow(post.ID, post.UserID, post.Title, post.Body, post.Published, "2024-01-01", "2024-01-01"))

		newPost, err := repo.Create(post)
		assert.NoError(t, err)
		assert.NotNil(t, newPost)
		assert.Equal(t, newPost.ID, post.ID)
		assert.Equal(t, newPost.UserID, post.UserID)
		assert.Equal(t, newPost.Title, post.Title)
		assert.Equal(t, newPost.Body, post.Body)
		assert.Equal(t, newPost.Published, post.Published)
	})

	t.Run("incorrect", func(t *testing.T) {
		gofakeit.Seed(0)
		post := models.Post{
			ID:        gofakeit.UUID(),
			UserID:    gofakeit.UUID(),
			Title:     gofakeit.BookTitle(),
			Body:      gofakeit.ProductDescription(),
			Published: true,
		}

		mock.ExpectQuery(`INSERT INTO posts`).WithArgs(post.ID, post.UserID, post.Title, post.Body, post.Published).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow("12345", "123456", "fakeError", "fakeError", false, "2024-01-01", "2024-01-01"))

		newPost, err := repo.Create(post)

		assert.NoError(t, err)
		assert.NotNil(t, newPost)
	})

	t.Run("exec_error", func(t *testing.T) {
		post := models.Post{
			ID:        gofakeit.UUID(),
			UserID:    gofakeit.UUID(),
			Title:     gofakeit.BookTitle(),
			Body:      gofakeit.ProductDescription(),
			Published: gofakeit.Bool(),
		}

		mock.ExpectQuery(`INSERT INTO posts`).WithArgs(post.ID, post.UserID, post.Title, post.Body, post.Published).
			WillReturnError(sql.ErrTxDone)

		newPost, err := repo.Create(post)
		assert.Error(t, err)
		assert.Nil(t, newPost)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewPostRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, title, body, published, created_at, updated_at FROM posts`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow("1", "11", "Doe1", "John1", true, "2024-01-01", "2024-01-01").
				AddRow("2", "22", "Doe2", "John2", false, "2024-01-01", "2024-01-01"))

		list, err := repo.GetAll()
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Len(t, list, 2)
	})

	t.Run("incorrect", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, title, body, published, created_at, updated_at FROM posts`).
			WillReturnError(sql.ErrNoRows)

		list, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, list)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostRepository_GetPersonal(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewPostRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		user_id := "11"
		mock.ExpectQuery(`SELECT id, user_id, title, body, published, created_at, updated_at FROM posts WHERE user_id = \$1`).
			WithArgs(user_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow("1", user_id, "Doe1", "John1", true, "2024-01-01", "2024-01-01").
				AddRow("2", "22", "Doe2", "John2", false, "2024-01-01", "2024-01-01"))

		list, err := repo.GetPersonal(user_id)
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Len(t, list, 2)
	})

	t.Run("incorrect", func(t *testing.T) {
		user_id := "11"
		mock.ExpectQuery(`SELECT id, user_id, title, body, published, created_at, updated_at FROM posts WHERE user_id = \$1`).
			WithArgs(user_id).
			WillReturnError(sql.ErrNoRows)

		list, err := repo.GetPersonal(user_id)
		assert.Error(t, err)
		assert.Nil(t, list)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewPostRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		post_id := "11"
		mock.ExpectQuery(`SELECT id, user_id, title, body, published, created_at, updated_at FROM posts WHERE id = \$1`).
			WithArgs(post_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow(post_id, "11", "Doe1", "John1", true, "2024-01-01", "2024-01-01").
				AddRow("2", "22", "Doe2", "John2", false, "2024-01-01", "2024-01-01"))

		item, err := repo.GetByID(post_id)
		assert.NotEqual(t, err, sql.ErrNoRows)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, item.ID, post_id)
	})

	t.Run("incorrect", func(t *testing.T) {
		post_id := "11"
		mock.ExpectQuery(`SELECT id, user_id, title, body, published, created_at, updated_at FROM posts WHERE id = \$1`).
			WithArgs(post_id).
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetByID(post_id)
		assert.Equal(t, err, sql.ErrNoRows)
		assert.Error(t, err)
		assert.Nil(t, item)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewPostRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		gofakeit.Seed(0)
		published := gofakeit.Bool()
		post := models.UpdatePost{
			Title:     gofakeit.BookTitle(),
			Body:      gofakeit.ProductDescription(),
			Published: &published,
		}

		post_id := "1"
		user_id := "11"
		mock.ExpectQuery(`UPDATE posts SET title = \$1, body = \$2, published = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, user_id, title, body, published, created_at, updated_at`).
			WithArgs(post.Title, post.Body, *post.Published, sqlmock.AnyArg(), post_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow(post_id, user_id, post.Title, post.Body, published, "2024-01-01", "2024-01-01"))

		newPost, err := repo.Update(post_id, post)
		assert.NoError(t, err)
		assert.NotNil(t, newPost)
		assert.Equal(t, post.Title, newPost.Title)
		assert.Equal(t, post.Body, newPost.Body)
		assert.Equal(t, *post.Published, newPost.Published)
	})

	t.Run("incorrect", func(t *testing.T) {
		gofakeit.Seed(0)
		published := gofakeit.Bool()
		post := models.UpdatePost{
			Title:     gofakeit.BookTitle(),
			Body:      gofakeit.ProductDescription(),
			Published: &published,
		}

		post_id := "1"
		user_id := "11"
		mock.ExpectQuery(`UPDATE posts SET title = \$1, body = \$2, published = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, user_id, title, body, published, created_at, updated_at`).
			WithArgs(post.Title, post.Body, *post.Published, sqlmock.AnyArg(), post_id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "body", "published", "created_at", "updated_at"}).
				AddRow(post_id, user_id, "fake", "fake", post.Published, "2024-01-01", "2024-01-01"))

		newPost, err := repo.Update(post_id, post)
		assert.NoError(t, err)
		assert.NotNil(t, newPost)
		assert.NotEqual(t, post.Title, newPost.Title)
		assert.NotEqual(t, post.Body, newPost.Body)
		assert.NotEqual(t, *post.Published, !newPost.Published)
	})

	t.Run("not_found", func(t *testing.T) {
		gofakeit.Seed(0)
		published := gofakeit.Bool()
		post := models.UpdatePost{
			Title:     gofakeit.BookTitle(),
			Body:      gofakeit.ProductDescription(),
			Published: &published,
		}

		post_id := "1"
		mock.ExpectQuery(`UPDATE posts SET title = \$1, body = \$2, published = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, user_id, title, body, published, created_at, updated_at`).
			WithArgs(post.Title, post.Body, *post.Published, sqlmock.AnyArg(), post_id).
			WillReturnError(sql.ErrNoRows)

		newPost, err := repo.Update(post_id, post)
		assert.Error(t, err)
		assert.Nil(t, newPost)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewPostRepository(sqlxDB)

	t.Run("correct", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
			WithArgs("1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete("1")
		assert.NoError(t, err)
	})

	t.Run("incorrect", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
			WithArgs("1").
			WillReturnError(sql.ErrNoRows)

		err := repo.Delete("1")
		assert.Error(t, err)
		assert.Equal(t, err, sql.ErrNoRows)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
