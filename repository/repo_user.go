package repository

import (
	"log"
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
		log.Println("repository_user: create - query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	query := `
		SELECT 
			id, first_name, last_name,
			email, role, created_at, updated_at
		FROM users
	`
	err := r.db.Select(&users, query)
	if err != nil {
		log.Println("repository_user: getAll - query error: ", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT 
			id, first_name, last_name,
			email, role, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("repository_user: getByID - query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT 
			id, first_name, last_name,
			email, password, role, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("repository_user: getByEmail - query error: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(id string, req models.UpdateUser) (*models.User, error) {
	query := `
		UPDATE users SET 
			first_name = $1,
			last_name = $2,
			role = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, first_name, last_name, email, role, created_at, updated_at
	`

	var input models.User
	err := r.db.QueryRow(query, req.FirstName, req.LastName, req.Role, time.Now(), id).
		Scan(&input.ID, &input.FirstName, &input.LastName, &input.Email, &input.Role, &input.CreatedAt, &input.UpdatedAt)
	if err != nil {
		log.Println("repository_user: update - query error: ", err)
		return nil, err
	}

	return &input, nil
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("repository_user: delete - exec error: ", err)
		return err
	}

	return nil
}
