package models

type SavedPost struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	PostID    string `json:"post_id" db:"post_id"`
}

type SavedPostAction struct {
	UserID    string `json:"user_id" db:"user_id" validate:"required"`
	PostID string `json:"post_id" db:"post_id" validate:"required"`
}

func (spa *SavedPostAction) IsValid() error {
	return validate.Struct(spa)
}
