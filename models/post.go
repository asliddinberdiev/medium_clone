package models

type Post struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	Title     string `json:"title" db:"title"`
	Body      string `json:"body" db:"body"`
	Published bool   `json:"published" db:"published"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type PersonalPost struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Body      string `json:"body" db:"body"`
	Published bool   `json:"published" db:"published"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type CreatePost struct {
	Title     string `json:"title" db:"title" validate:"required,min=2"`
	Body      string `json:"body" db:"body" validate:"required,min=10"`
	Published bool   `json:"published" db:"published"`
}

type UpdatePost struct {
	Title     string `json:"title" db:"title" validate:"omitempty,min=2"`
	Body      string `json:"body" db:"body" validate:"omitempty,min=10"`
	Published *bool   `json:"published" db:"published" validate:"omitempty"`
}

func (pp *PersonalPost) IsValid() error {
	return validate.Struct(pp)
}

func (cp *CreatePost) IsValid() error {
	return validate.Struct(cp)
}

func (up *UpdatePost) IsValid() error {
	return validate.Struct(up)
}