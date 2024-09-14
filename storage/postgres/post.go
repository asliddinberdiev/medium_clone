package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostStorage(db *sqlx.DB) repo.PostStorageI {
	return &postRepo{
		db: db,
	}
}

func (u *postRepo) Create(ctx context.Context, req *repo.Post) (*repo.Post, error) {
	query := `
		INSERT INTO posts (
			id, user_id, title,
			body, published) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING * 
	`

	err := u.db.QueryRow(query, req.ID, req.UserID, req.Title, req.Body, req.Published).Scan(&req.ID, &req.UserID, &req.Title, &req.Body, &req.Published, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *postRepo) GetAll(ctx context.Context) ([]*repo.Post, error) {
	query := `SELECT * FROM posts`

	var posts []*repo.Post
	if err := u.db.Select(&posts, query); err != nil {
		return nil, err
	}

	if posts == nil {
		return []*repo.Post{}, nil
	}

	return posts, nil
}

func (u *postRepo) Get(ctx context.Context, id string) (*repo.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`

	var post repo.Post
	err := u.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Published, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (u *postRepo) Update(ctx context.Context, req *repo.UpdatePost) error {
	tsx, err := u.db.Begin()
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

func (u *postRepo) Delete(ctx context.Context, id string) error {
	tsx, err := u.db.Begin()
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
