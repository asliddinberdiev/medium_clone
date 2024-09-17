package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      App
	Postgres Postgres
}

type App struct {
	Port        string
	Version     string
	TokenKey    string
	AccessTime  string
	RefreshTime string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func Load(path string) Config {
	err := godotenv.Load(path + "/.env")
	if err != nil {
		log.Fatalln("loading .env file error: ", err)
	}

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		App: App{
			Port:        conf.GetString("APP_PORT"),
			Version:     conf.GetString("APP_VERSION"),
			TokenKey:    conf.GetString("APP_TOKEN_KEY"),
			AccessTime:  conf.GetString("APP_ACCESS_TIME"),
			RefreshTime: conf.GetString("APP_REFRESH_TIME"),
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
