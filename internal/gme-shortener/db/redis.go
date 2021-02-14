package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/config"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
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
	*redis.Client
	*cache.Cache
}

// NewRedisDatabase -> Use Redis as backend
func NewRedisDatabase(cfg config.RedisConfig) (Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if res := client.Set(context.TODO(), "heartbeat", 1, 0); res.Err() != nil {
		log.Fatalln("Error connecting to Redis:", res.Err())
		return nil, res.Err()
	}

	return &redisDB{
		Client: client,
		Cache:  cache.New(15*time.Minute, 10*time.Minute),
	}, nil
}

func (rdb *redisDB) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	var data *redis.StringCmd
	data = rdb.Client.Get(context.TODO(), "short::"+id)
	if data.Err() != nil {
		return nil, data.Err()
	}
	json.Unmarshal([]byte(data.Val()), &res)
	return
}

func (rdb *redisDB) SaveShortenedURL(short short.ShortURL) (err error) {
	return nil
}

func (rdb *redisDB) SaveShortenedURLWithExpiration(url short.ShortURL, expireAfter time.Duration) (err error) {
	var data []byte
	data, err = json.Marshal(url)
	if err != nil {
		return
	}
	rdb.Client.Set(context.TODO(), "short::"+url.ID, string(data), expireAfter)
	return
}

func (rdb *redisDB) BreakCache(id string) (found bool) {
	_, found = rdb.Cache.Get(id)
	rdb.Cache.Delete(id)
	return
}
