package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID        string `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	Role      string `json:"role" db:"role"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type UserCreate struct {
	FirstName string `json:"first_name" db:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" db:"last_name" validate:"omitempty,min=2"`
	Email     string `json:"email" db:"email" validate:"required,email"`
	Password  string `json:"password" db:"password" validate:"required,min=6"`
	Role      string `json:"role" db:"role" validate:"omitempty,oneof=admin user"`
}

type UpdateUser struct {
	FirstName string `json:"first_name" db:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name" db:"last_name" validate:"omitempty,min=2"`
	Password  string `json:"password" db:"password" validate:"omitempty,min=6"`
	Role      string `json:"role" db:"role" validate:"omitempty,oneof=admin user"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (uc *UserCreate) IsValid() error {
	return validate.Struct(uc)
}

func (uu *UpdateUser) IsValid() error {
	return validate.Struct(uu)
}
