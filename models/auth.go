package models

type Login struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=6"`
}

type Token struct {
	Token string `json:"token"`
}

func (l *Login) IsValid() error {
	return validate.Struct(l)
}
