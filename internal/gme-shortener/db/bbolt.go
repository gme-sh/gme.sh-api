package db

import (
	"encoding/json"
	"time"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"github.com/patrickmn/go-cache"
	"go.etcd.io/bbolt"
)

const (
	BBoltBucketName = "stonks-urls"
)

// NewBBoltDatabase -> Create new BBoltDatabase
func NewBBoltDatabase(path string) (bbdb Database, err error) {
	db, err := bbolt.Open(path, 0666, nil)

	if err != nil {
		return nil, err
	}

	bbdb = &bboltDatabase{
		database: db,
		cache:    cache.New(10*time.Minute, 15*time.Minute),
	}

	return
}

type bboltDatabase struct {
	database *bbolt.DB
	cache    *cache.Cache
}

func (bbdb *bboltDatabase) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	var ser []byte

	if err = bbdb.database.View(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists([]byte(BBoltBucketName)); err != nil {
			return
		}
		ser = bucket.Get([]byte(id))
		return
	}); err != nil {
		return
	}

	err = json.Unmarshal(ser, &res)
	return
}

func (bbdb *bboltDatabase) SaveShortenedURL(short short.ShortURL) (err error) {
	var ser []byte
	ser, err = json.Marshal(short)

	if err != nil {
		return
	}

	err = bbdb.database.Update(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists([]byte(BBoltBucketName)); err != nil {
			return
		}
		err = bucket.Put([]byte(short.ID), ser)
		return
	})

	return
}

func (bbdb *bboltDatabase) BreakCache(id string) (found bool) {
	_, found = bbdb.cache.Get(id)
	bbdb.cache.Delete(id)
	return
}
