package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App      App
	Postgres Postgres
}

type App struct {
	Port    string
	Version string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func Load(path string, log *zap.Logger) Config {
	// Load environment variables from the .env file
	err := godotenv.Load(path + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file", zap.Error(err))
	}

	conf := viper.New()
	conf.AutomaticEnv()

	// Populate the Config struct
	cfg := Config{
		App: App{
			Port:    conf.GetString("APP_PORT"),
			Version: conf.GetString("APP_VERSION"),
		},
		Postgres: Postgres{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DB"),
			SSLMode:  conf.GetString("POSTGRES_SSLMODE"),
		},
	}

	return cfg
}
