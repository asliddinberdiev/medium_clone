package repo

import (
	"context"
	"time"
)

type UserStorageI interface {
	Create(ctx context.Context, req *User) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, req *UpdateUser) error
	Delete(ctx context.Context, id string) error
}

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUser struct {
	ID        string
	FirstName string
	LastName  string
	Password  string
}
