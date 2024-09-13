package storage

import (
	"github.com/asliddinberdiev/medium_clone/storage/postgres"
	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Post() repo.PostStorageI
}

type storagePg struct {
	userRepo repo.UserStorageI
	postRepo repo.PostStorageI
}

func NewStorage(psqlConn *sqlx.DB) StorageI {
	return &storagePg{
		userRepo: postgres.NewUserStorage(psqlConn),
		postRepo: postgres.NewPostStorage(psqlConn),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *storagePg) Post() repo.PostStorageI {
	return s.postRepo
}
