package models

type Post struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreatePost struct {
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
}

type UpdatePost struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
}
