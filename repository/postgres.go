package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database     string
	SSLMode  string
}

func NewPostgresDB(cfg PostgresConfig, log *zap.Logger) (*sqlx.DB, error) {
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Database,
		cfg.Password,
		cfg.SSLMode,
	)

	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("postgres ping error", zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("postgres ping error", zap.Error(err))
		return nil, err
	}

	log.Info("initialize postgres")
	return db, nil
}
