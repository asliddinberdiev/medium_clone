package repository_test

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var cfgR *viper.Viper

func init() {
	if err := godotenv.Load("./../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfgR = viper.New()
	cfgR.AutomaticEnv()
}

func validRedisConfig() repository.RedisConfig {
	return repository.RedisConfig{
		Host:     cfgR.GetString("REDIS_HOST"),
		Port:     cfgR.GetString("REDIS_PORT"),
		Password: cfgR.GetString("REDIS_PASSWORD"),
	}
}

func invalidRedisConfig() repository.RedisConfig {
	return repository.RedisConfig{
		Host:     "invalid",
		Port:     "9999",
		Password: "adsad546",
	}
}

func TestNewRedisDB(t *testing.T) {
	log.SetOutput(io.Discard)
	validCfg := validRedisConfig()
	invalidCfg := invalidRedisConfig()

	t.Run("ValidConnection", func(t *testing.T) {
		rdb, err := repository.NewRedisDB(validCfg)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		defer rdb.Close()

		_, err = rdb.Ping(context.Background()).Result()
		if err != nil {
			t.Fatalf("could not ping the Redis database: %v", err)
		}
	})

	t.Run("InvalidConnection", func(t *testing.T) {
		rdb, err := repository.NewRedisDB(invalidCfg)
		if err == nil {
			t.Fatalf("expected an error due to invalid credentials, but got none")
		}

		if rdb != nil {
			t.Fatalf("expected rdb to be nil on error, but got %v", rdb)
		}
	})
}
