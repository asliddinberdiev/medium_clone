package repository

import (
	"context"
	"fmt"

	"log"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func NewRedisDB(cfg RedisConfig) (*redis.Client, error) {
	rsUrl := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     rsUrl,
		Password: cfg.Password,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("repository_redis: NewRedisDB - ping error: ", err)
		return nil, err
	}

	log.Println("repository_redis: initialize")
	return rdb, nil
}
