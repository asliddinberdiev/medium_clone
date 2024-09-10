package postgres

import (
	"context"
	"database/sql"

	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, req *repo.User) (*repo.User, error) {
	query := `
		INSERT INTO users (
			id, first_name, last_name,
			email, password
		) VALUES($1, $2, $3, $4, $5) RETURNING id
	`

	err := u.db.QueryRow(query, req.ID, req.FirstName, req.LastName, req.Email, req.Password).Scan(&req.ID)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *userRepo) Get(ctx context.Context, id string) (*repo.User, error) {
	query := `
		SELECT 
			id, first_name, last_name,
			email, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user repo.User
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) Update(ctx context.Context, req *repo.UpdateUser) error {
	tsx, err := u.db.Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE users SET 
			first_name = $1,
			last_name = $2,
			password = $3
		WHERE id = $4
	`

	res, err := tsx.Exec(query, req.FirstName, req.LastName, req.Password, req.ID)
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

func (u *userRepo) Delete(ctx context.Context, id string) error {
	tsx, err := u.db.Begin()
	if err != nil {
		return err
	}

	res, err := tsx.Exec("DELETE FROM users WHERE id = $1", id)
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
