package repository_test

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

	t.Run("correct", func(t *testing.T) {
		rdb, err := repository.NewRedisDB(validCfg)
		assert.NoError(t, err)
		assert.NotNil(t, rdb)
		defer rdb.Close()

		_, err = rdb.Ping(context.Background()).Result()
		assert.NoError(t, err)
	})

	t.Run("InvalidConnection", func(t *testing.T) {
		rdb, err := repository.NewRedisDB(invalidCfg)
		assert.Error(t, err)
		assert.Nil(t, rdb)
	})
}
