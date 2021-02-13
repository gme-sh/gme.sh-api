package db

import (
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/config"
	"github.com/go-redis/redis/v8"
)

// NewRedisClient -> Create a new Redis client
func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
