package repository

import (
	"fmt"

	"log"

	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

func NewPostgresDB(cfg PostgresConfig) (*sqlx.DB, error) {
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
		log.Fatalf("postgres ping error: %v\n", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("postgres ping error: %v\n", err)
		return nil, err
	}

	log.Println("initialize postgres")
	return db, nil
}
