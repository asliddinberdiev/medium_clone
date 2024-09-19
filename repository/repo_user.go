package repository

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user models.User) (*models.User, error) {
	tsx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (id, first_name, last_name, email, password, role) 
		VALUES($1, $2, $3, $4, $5, $6) 
		RETURNING id, first_name, last_name, email, role, created_at, updated_at 
	`

	err = tsx.QueryRow(query, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Println("repository_user: create - this email already used")
			return nil, errors.New("unique")
		}
		errRoll := tsx.Rollback()
		if errRoll != nil {
			log.Println("repository_user: create - rollback error: ", errRoll)
			err = errRoll
		}
		log.Println("repository_user: create - query error: ", err)
		return nil, err
	}

	if err := tsx.Commit(); err != nil {
		log.Println("repository_user: create - commit error: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	var input []*models.User
	query := `
		SELECT 
			id, first_name, last_name,
			email, role, created_at, updated_at
		FROM users 
	`
	err := r.db.Select(&input, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.User{}, nil
		}
		log.Println("repository_user: getAll - query error: ", err)
		return nil, err
	}

	return input, nil
}

func (u *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT 
			id, first_name, last_name,
			email, role, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("repository_user: getByID - not found")
			return nil, err
		}
		log.Println("repository_user: getByID - query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT 
			id, first_name, last_name,
			email, password, role, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := u.db.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("repository_user: getByEmail - not found")
			return nil, err
		}
		log.Println("repository_user: getByEmail - query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Update(id string, req models.UpdateUser) (*models.User, error) {
	var input models.User
	tsx, err := u.db.Begin()
	if err != nil {
		log.Println("repository_user: update - begin error: ", err)
		return nil, err
	}

	query := `
		UPDATE users SET 
			first_name = $1,
			last_name = $2,
			role = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, first_name, last_name, email, role, created_at, updated_at
	`

	err = tsx.QueryRow(query, req.FirstName, req.LastName, req.Role, time.Now(), id).
		Scan(&input.ID, &input.FirstName, &input.LastName, &input.Email, &input.Role, &input.CreatedAt, &input.UpdatedAt)
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			log.Println("repository_user: update - rollback error: ", errRoll)
			err = errRoll
		}
		log.Println("repository_user: update - exec error: ", err)
		return nil, err
	}

	err = tsx.Commit()
	if err != nil {
		log.Println("repository_user: update - commit error: ", err)
		return nil, err
	}

	return &input, nil
}

func (u *UserRepository) Delete(id string) error {
	tsx, err := u.db.Begin()
	if err != nil {
		log.Println("repository_user: delete - begin error: ", err)
		return err
	}

	res, err := tsx.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("repository_user: delete - exec error: ", err)
		return err
	}

	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			log.Println("repository_user: delete - rowsaffected rollback error: ", err)
			err = errRoll
		}
		log.Println("repository_user: delete - rowsaffected error: ", err)
		return err
	}

	if data == 0 {
		tsx.Commit()
		log.Println("repository_user: delete - not found error: ", err)
		return sql.ErrNoRows
	}

	return tsx.Commit()
}
