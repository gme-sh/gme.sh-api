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
	client  *redis.Client
	cache   *cache.Cache
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
		cache:   cache.New(15*time.Minute, 10*time.Minute),
		context: ctx,
	}, nil
}

func (rdb *redisDB) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	var data *redis.StringCmd
	data = rdb.client.Get(rdb.context, "short::"+id)
	err = data.Err()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data.Val()), &res)
	return
}

func (rdb *redisDB) SaveShortenedURL(short *short.ShortURL) (err error) {
	return nil
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

func (rdb *redisDB) BreakCache(id string) (found bool) {
	_, found = rdb.cache.Get(id)
	rdb.cache.Delete(id)
	return
}

func (rdb *redisDB) ShortURLAvailable(id string) bool {
	return shortURLAvailable(rdb, id)
}

func (rdb *redisDB) Heartbeat() (err error) {
	err = rdb.client.Set(rdb.context, "heartbeat", 1, 1*time.Second).Err()
	return
}
