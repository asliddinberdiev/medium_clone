package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewAuthRepository(db *sqlx.DB, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{db: db, rdb: rdb}
}

func (r *AuthRepository) AddBlack(tokenID, token string, exp time.Duration) error {
	ctx := context.Background()
	return r.rdb.Set(ctx, tokenID, token, exp).Err()
}

func (r *AuthRepository) GetBlackToken(tokenID string) (string, error) {
	ctx := context.Background()
	return r.rdb.Get(ctx, tokenID).Result()
}
