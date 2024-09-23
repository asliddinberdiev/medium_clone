package models

type Comment struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	PostID    string `json:"post_id" db:"post_id"`
	Body      string `json:"body" db:"body"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type CreateComment struct {
	PostID string `json:"post_id" db:"post_id" validate:"required"`
	Body   string `json:"body" db:"body" validate:"required,min=3"`
}

type UpdateComment struct {
	Body string `json:"body" db:"body" validate:"omitempty,min=3"`
}

func (cc *CreateComment) IsValid() error {
	return validate.Struct(cc)
}

func (uc *UpdateComment) IsValid() error {
	return validate.Struct(uc)
}
