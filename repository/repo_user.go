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
		log.Println("repository: user create query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	return []models.User{}, nil
}

func (u *UserRepository) Get(id string) (*models.User, error) {
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
			log.Println("repository: user not found")
			return nil, err
		}
		log.Println("repository: user get query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Update(req *models.User) error {
	tsx, err := u.db.Begin()
	if err != nil {
		log.Println("repository: user update begin error: ", err)
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
			log.Println("repository: user update rollback error: ", err)
			err = errRoll
		}
		log.Println("repository: user update exec error: ", err)
		return err
	}

	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			log.Println("repository: user update rowsaffected rollback error: ", err)
			err = errRoll
		}
		log.Println("repository: user update rowsaffected error: ", err)
		return err
	}

	if data == 0 {
		tsx.Commit()
		log.Println("repository: user update not found error: ", err)
		return sql.ErrNoRows
	}

	return tsx.Commit()
}

func (u *UserRepository) Delete(id string) error {
	tsx, err := u.db.Begin()
	if err != nil {
		log.Println("repository: user delete begin error: ", err)
		return err
	}

	res, err := tsx.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("repository: user delete exec error: ", err)
		return err
	}

	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			log.Println("repository: user delete rowsaffected rollback error: ", err)
			err = errRoll
		}
		log.Println("repository: user delete rowsaffected error: ", err)
		return err
	}

	if data == 0 {
		tsx.Commit()
		log.Println("repository: user delete not found error: ", err)
		return sql.ErrNoRows
	}

	return tsx.Commit()
}
