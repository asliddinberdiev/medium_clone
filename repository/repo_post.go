package repository

import (
	"context"
	"database/sql"
	"time"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, req *models.Post) (*models.Post, error) {
	query := `
		INSERT INTO posts (
			id, user_id, title,
			body, published) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING * 
	`

	err := r.db.QueryRow(query, req.ID, req.UserID, req.Title, req.Body, req.Published).Scan(&req.ID, &req.UserID, &req.Title, &req.Body, &req.Published, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *PostRepository) GetAll(ctx context.Context, limit int, offset int) ([]*models.Post, error) {
	query := `
		SELECT * FROM posts
		ORDER BY id LIMIT $1 OFFSET $2
		`

	var posts []*models.Post
	if err := r.db.Select(&posts, query, limit, offset); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetAllPersonal(ctx context.Context, userID string, limit int, offset int) ([]*models.PersonalPost, error) {
	query := `
		SELECT id, title, body, published, created_at, updated_at 
		FROM posts
		WHERE user_id = $1
		ORDER BY id LIMIT $2 OFFSET $3
		`

	var posts []*models.PersonalPost
	if err := r.db.Select(&posts, query, userID, limit, offset); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) Get(ctx context.Context, id string) (*models.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`

	var post models.Post
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Published, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) Update(ctx context.Context, req *models.UpdatePost) error {
	tsx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE posts SET 
			title = $1,
			body = $2,
			published = $3,
			updated_at = $4
		WHERE id = $5 
		RETURNING id, user_id, title, body, published, created_at, updated_at
	`

	res, err := tsx.Exec(query, req.Title, req.Body, req.Published, time.Now(), req.ID)
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			err = errRoll
		}
		return err
	}

	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			err = errRoll
		}
		return err
	}

	if data == 0 {
		tsx.Commit()
		return sql.ErrNoRows
	}

	return tsx.Commit()
}

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	tsx, err := r.db.Begin()
	if err != nil {
		return err
	}

	res, err := tsx.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}

	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			err = errRoll
		}
		return err
	}

	if data == 0 {
		tsx.Commit()
		return sql.ErrNoRows
	}

	return tsx.Commit()
}
