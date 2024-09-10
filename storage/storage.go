package storage

import (
	"github.com/asliddinberdiev/medium_clone/storage/postgres"
	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
}

type storagePg struct {
	userRepo repo.UserStorageI
}

func NewStorage(psqlConn *sqlx.DB) StorageI {
	return &storagePg{
		userRepo: postgres.NewUserStorage(psqlConn),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}
