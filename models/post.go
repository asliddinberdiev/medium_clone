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
	UserID    string `json:"user_id" db:"user_id"`
	Title     string `json:"title" db:"title"`
	Body      string `json:"body" db:"body"`
	Published bool   `json:"published" db:"published"`
}

type UpdatePost struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Body      string `json:"body" db:"body"`
	Published bool   `json:"published" db:"published"`
}
