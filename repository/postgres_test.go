package repository_test

import (
	"log"
	"testing"

	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var cfgP *viper.Viper

func init() {
	if err := godotenv.Load("./../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfgP = viper.New()
	cfgP.AutomaticEnv()
}

func confSetupTestDB() repository.PostgresConfig {
	return repository.PostgresConfig{
		User:     cfgP.GetString("POSTGRES_USER"),
		Password: cfgP.GetString("POSTGRES_PASSWORD"),
		Host:     cfgP.GetString("POSTGRES_HOST"),
		Port:     cfgP.GetString("POSTGRES_PORT"),
		Database: cfgP.GetString("POSTGRES_DB"),
		SSLMode:  cfgP.GetString("POSTGRES_SSLMODE"),
	}
}

func invalidConfSetupTestDB() repository.PostgresConfig {
	return repository.PostgresConfig{
		User:     "invalid_user",
		Password: "invalid_password",
		Host:     "localhost",
		Port:     "5431",
		Database: "invalid_db",
		SSLMode:  "disable",
	}
}

func TestNewPostgresDB(t *testing.T) {
	cfg := confSetupTestDB()
	invalidCfg := invalidConfSetupTestDB()

	t.Run("ValidConnection", func(t *testing.T) {
		db, err := repository.NewPostgresDB(cfg)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			t.Fatalf("could not ping the database: %v", err)
		}
	})

	t.Run("InvalidConnection", func(t *testing.T) {
		db, err := repository.NewPostgresDB(invalidCfg)
		if err == nil {
			t.Fatalf("expected an error due to invalid credentials, but got none")
		}

		if db != nil {
			t.Fatalf("expected db to be nil on error, but got %v", db)
		}
	})
}
