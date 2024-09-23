package repository

import (
	"log"
	"time"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post models.Post) (*models.Post, error) {
	query := `
		INSERT INTO posts (id, user_id, title, body, published)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, body, published, created_at, updated_at
	`

	err := r.db.QueryRow(query, post.ID, post.UserID, post.Title, post.Body, post.Published).Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Published, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		log.Println("repository_post: create - query error: ", err)
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetByID(id string) (*models.Post, error) {
	query := `SELECT  id, user_id, title, body, published, created_at, updated_at FROM posts WHERE id = $1`

	var post models.Post
	err := r.db.Get(&post, query, id)
	if err != nil {
		log.Println("repository_post: GetByID - query error: ", err)
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPersonal(user_id string) ([]*models.Post, error) {
	query := `SELECT  id, user_id, title, body, published, created_at, updated_at FROM posts WHERE user_id = $1`

	var list []*models.Post
	err := r.db.Select(&list, query, user_id)
	if err != nil {
		log.Println("repository_post: GetPersonal - query error: ", err)
		return nil, err
	}

	return list, nil
}

func (r *PostRepository) GetAll() ([]*models.Post, error) {
	query := `SELECT  id, user_id, title, body, published, created_at, updated_at FROM posts`

	var list []*models.Post
	err := r.db.Select(&list, query)
	if err != nil {
		log.Println("repository_post: GetPersonal - query error: ", err)
		return nil, err
	}

	return list, nil
}

func (r *PostRepository) Update(id string, post models.UpdatePost) (*models.Post, error) {
	query := `
		UPDATE posts SET
			title = $1,
			body = $2,
			published = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, user_id, title, body, published, created_at, updated_at
	`

	var updatePost models.Post
	err := r.db.QueryRow(query, post.Title, post.Body, post.Published, time.Now(), id).Scan(&updatePost.ID, &updatePost.UserID, &updatePost.Title, &updatePost.Body, &updatePost.Published, &updatePost.CreatedAt, &updatePost.UpdatedAt)
	if err != nil {
		log.Println("repository_post: update - query error: ", err)
		return nil, err
	}

	return &updatePost, nil
}

func (r *PostRepository) Delete(id string) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("repository_post: delete - exec error: ", err)
		return err
	}

	return nil
}
