package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/config"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
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

type redisDB struct {
	client  *redis.Client
	context context.Context
}

// NewRedisDatabase -> Use Redis as backend
func NewRedisDatabase(cfg config.RedisConfig) (Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.TODO()
	if res := client.Set(ctx, "heartbeat", 1, 0); res.Err() != nil {
		log.Fatalln("Error connecting to Redis:", res.Err())
		return nil, res.Err()
	}

	return &redisDB{
		client:  client,
		context: ctx,
	}, nil
}

func (rdb *redisDB) FindShortenedURL(id short.ShortID) (res *short.ShortURL, err error) {
	var data *redis.StringCmd
	data = rdb.client.Get(rdb.context, id.RedisKey())
	err = data.Err()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data.Val()), &res)

	return
}

func (rdb *redisDB) SaveShortenedURL(short *short.ShortURL) (err error) {
	var data []byte
	data, err = json.Marshal(short)
	if err != nil {
		return
	}
	err = rdb.client.Set(rdb.context, short.ID.RedisKey(), string(data), redis.KeepTTL).Err()
	return
}

func (rdb *redisDB) SaveShortenedURLWithExpiration(url *short.ShortURL, expireAfter time.Duration) (err error) {
	var data []byte
	data, err = json.Marshal(url)
	if err != nil {
		return
	}

	err = rdb.client.Set(rdb.context, url.ID.RedisKey(), string(data), expireAfter).Err()
	return
}

func (rdb *redisDB) BreakCache(_ short.ShortID) (found bool) {
	return false
}

func (rdb *redisDB) ShortURLAvailable(id short.ShortID) bool {
	return shortURLAvailable(rdb, id)
}

func (rdb *redisDB) Heartbeat() (err error) {
	err = rdb.client.Set(rdb.context, "heartbeat", 1, 1*time.Second).Err()
	return
}

func (rdb *redisDB) FindStats(id short.ShortID) *short.Stats {
	return nil
}
