package repository_test

import (
	"io"
	"log"
	"testing"

	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
	log.SetOutput(io.Discard)

	cfg := confSetupTestDB()
	invalidCfg := invalidConfSetupTestDB()

	t.Run("correct", func(t *testing.T) {
		db, err := repository.NewPostgresDB(cfg)
		if err == nil {
			defer db.Close()
		}

		assert.NoError(t, err)
		assert.NotNil(t, db)

		pingErr := db.Ping()
		assert.NoError(t, pingErr)
	})

	t.Run("incorrect", func(t *testing.T) {
		db, err := repository.NewPostgresDB(invalidCfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
