package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/go-redis/redis/v8"
)

// PersistentDatabase
// TemporaryDatabase
type redisDB struct {
	client  *redis.Client
	context context.Context
	ps      *redis.PubSub
}

// NewRedisDatabase -> Use Redis as backend
func NewRedisDatabase(cfg *config.RedisConfig) (Database, error) {
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

/*
 * ==================================================================================================
 *                            P E R M A N E N T  D A T A B A S E
 * ==================================================================================================
 */

func (rdb *redisDB) FindShortenedURL(id *short.ShortID) (res *short.ShortURL, err error) {
	data := rdb.client.Get(rdb.context, id.RedisKey())
	err = data.Err()
	if err != nil {
		return
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

func (rdb *redisDB) DeleteShortenedURL(id *short.ShortID) (err error) {
	err = rdb.client.Del(rdb.context, id.RedisKey()).Err()
	return
}

func (rdb *redisDB) SaveShortenedURLWithExpiration(u *short.ShortURL, e time.Duration) (err error) {
	var data []byte
	data, err = json.Marshal(u)
	if err != nil {
		return
	}
	err = rdb.client.Set(rdb.context, u.ID.RedisKey(), string(data), e).Err()
	return
}

func (rdb *redisDB) ShortURLAvailable(id *short.ShortID) bool {
	return shortURLAvailable(rdb, id)
}

/*
 * ==================================================================================================
 *                            T E M P O R A R Y  D A T A B A S E
 * ==================================================================================================
 */

func (rdb *redisDB) Heartbeat() (err error) {
	err = rdb.client.Set(rdb.context, "heartbeat", 1, 1*time.Second).Err()
	return
}

func (rdb *redisDB) FindStats(id *short.ShortID) (stats *short.Stats, err error) {
	var calls, calls60 uint64
	calls, err = rdb.client.Get(rdb.context, id.RedisKeyf(short.RedisKeyCountGlobal)).Uint64()
	if err != nil {
		return
	}
	calls60, err = rdb.client.Get(rdb.context, id.RedisKeyf(short.RedisKeyCount60)).Uint64()
	if err != nil {
		return
	}
	stats = &short.Stats{
		Calls:   calls,
		Calls60: calls60,
	}
	return
}

func (rdb *redisDB) AddStats(id *short.ShortID) (err error) {
	err = rdb.client.Incr(rdb.context, id.RedisKeyf(short.RedisKeyCountGlobal)).Err()
	if err != nil {
		return
	}
	count60Key := id.RedisKeyf(short.RedisKeyCount60)
	resultExists := rdb.client.Exists(rdb.context, count60Key)
	err = resultExists.Err()
	expire := resultExists.Val() == 0
	if err != nil {
		return
	}
	err = rdb.client.Incr(rdb.context, count60Key).Err()
	if err != nil {
		return
	}
	if !expire {
		err = rdb.client.Expire(rdb.context, count60Key, time.Hour).Err()
	}
	return
}

func (rdb *redisDB) DeleteStats(id *short.ShortID) (err error) {
	err = rdb.client.Del(
		rdb.context,
		id.RedisKeyf(short.RedisKeyCountGlobal),
		id.RedisKeyf(short.RedisKeyCount60),
	).Err()
	return
}

func (rdb *redisDB) Close() (err error) {
	if rdb.ps == nil {
		return
	}
	err = rdb.ps.Close()
	return
}

func (rdb *redisDB) Publish(channel, msg string) (err error) {
	err = rdb.client.Publish(rdb.context, channel, msg).Err()
	return
}

func (rdb *redisDB) Subscribe(c func(channel, payload string), channels ...string) (err error) {
	rdb.ps = rdb.client.Subscribe(rdb.context, channels...)
	// wait for confirmation
	_, err = rdb.ps.Receive(rdb.context)
	if err != nil {
		return
	}
	for msg := range rdb.ps.Channel() {
		c(msg.Channel, msg.Payload)
	}
	return
}
