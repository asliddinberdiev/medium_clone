package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
	"log"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, first_name, last_name, email, password, role) 
		VALUES($1, $2, $3, $4, $5, $6) 
		RETURNING id, first_name, last_name, email, role, created_at, updated_at 
	`

	err := r.db.QueryRow(query, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Println("repository: user create: this email already used")
			return nil, errors.New("unique")
		}
		log.Printf("repository: user create error: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	return []models.User{}, nil
}

func (u *UserRepository) Get(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT 
			id, first_name, last_name,
			email, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Update(ctx context.Context, req *models.User) error {
	tsx, err := u.db.Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE users SET 
			first_name = $1,
			last_name = $2,
			password = $3,
			updated_at = $4
		WHERE id = $5
	`

	res, err := tsx.Exec(query, req.FirstName, req.LastName, req.Password, time.Now(), req.ID)
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

func (u *UserRepository) Delete(ctx context.Context, id string) error {
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
