package repository

import (
	"log"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment models.Comment) (*models.Comment, error) {
	query := `
		INSERT INTO comments (id, user_id, post_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, post_id, body, created_at 
	`

	err := r.db.QueryRow(query, comment.ID, comment.UserID, comment.PostID, comment.Body).Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Body, &comment.CreatedAt)
	if err != nil {
		log.Println("repository_comment: create - query error: ", err)
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) GetAll(post_id string) ([]*models.Comment, error) {
	query := `
		SELECT  id, user_id, post_id, body, created_at 
		FROM comments
		WHERE post_id = $1	
	`

	var list []*models.Comment
	err := r.db.Select(&list, query, post_id)
	if err != nil {
		log.Println("repository_comment: GetAll - query error: ", err)
		return nil, err
	}

	return list, nil
}

func (r *CommentRepository) GetByID(id string) (*models.Comment, error) {
	query := `
		SELECT  id, user_id, post_id, body, created_at 
		FROM comments
		WHERE id = $1	
	`

	var item models.Comment
	err := r.db.Get(&item, query, id)
	if err != nil {
		log.Println("repository_comment: GetByID - query error: ", err)
		return nil, err
	}

	return &item, nil
}

func (r *CommentRepository) Update(id, body string) (*models.Comment, error) {
	query := `
		UPDATE comments SET
			body = $1
		WHERE id = $2
		RETURNING id, user_id, post_id, body, created_at
	`

	var comment models.Comment
	err := r.db.QueryRow(query, body, id).Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Body, &comment.CreatedAt)
	if err != nil {
		log.Println("repository_comment: update - query error: ", err)
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) Delete(id string) error {
	query := `DELETE FROM comments WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("repository_comment: delete - exec error: ", err)
		return err
	}

	return nil
}
