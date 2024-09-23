package repository

import (
	"log"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type SavedPostRepository struct {
	db *sqlx.DB
}

func NewSavedPostRepository(db *sqlx.DB) *SavedPostRepository {
	return &SavedPostRepository{db: db}
}

func (r *SavedPostRepository) Add(savedPost models.SavedPost) error {
	query := `
		INSERT INTO saved_posts (id, user_id, post_id)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, savedPost.ID, savedPost.UserID, savedPost.PostID)
	if err != nil {
		log.Println("repository_saved_post: Add - exec error: ", err)
		return err
	}

	return nil
}

func (r *SavedPostRepository) Remove(post_id string) error {
	query := `DELETE FROM saved_posts WHERE post_id = $1`

	_, err := r.db.Exec(query, post_id)
	if err != nil {
		log.Println("repository_saved_post: Remove - exec error: ", err)
		return err
	}

	return nil
}

func (r *SavedPostRepository) GetByID(user_id, post_id string) (*models.SavedPost, error) {
	query := `
		SELECT id, user_id, post_id 
		FROM saved_posts
		WHERE user_id = $1 AND post_id = $2 
	`

	var input models.SavedPost
	err := r.db.Get(&input, query, user_id, post_id)
	if err != nil {
		log.Println("repository_saved_post: GetByID - query error: ", err)
		return nil, err
	}

	return &input, nil
}

func (r *SavedPostRepository) GetAll(user_id string) ([]*models.Post, error) {
	query := `
        SELECT p.*
        FROM posts p
        INNER JOIN saved_posts sp ON sp.post_id = p.id
        WHERE sp.user_id = $1
    `

	var list []*models.Post
	err := r.db.Select(&list, query, user_id)
	if err != nil {
		log.Println("repository_saved_posts: GetAll - query error: ", err)
		return nil, err
	}

	return list, nil
}
