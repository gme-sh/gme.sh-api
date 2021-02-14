package db

import (
	"encoding/json"
	"time"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"github.com/patrickmn/go-cache"
	"go.etcd.io/bbolt"
)

// NewBBoltDatabase -> Create new BBoltDatabase
func NewBBoltDatabase(path string) (Database, error) {
	db, err := bbolt.Open(path, 0666, nil)

	if err != nil {
		return nil, err
	}

	return &bboltDatabase{db, cache.New(10*time.Minute, 15*time.Minute)}, nil
}

type bboltDatabase struct {
	DB    *bbolt.DB
	cache *cache.Cache
}

func (bbdb *bboltDatabase) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	var ser []byte
	bbdb.DB.View(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte("ShortenedUrls"))
		ser = b.Get([]byte(id))
		return
	})

	err = json.Unmarshal(ser, &res)

	return
}

func (bbdb *bboltDatabase) SaveShortenedURL(short short.ShortURL) (err error) {
	ser, err := json.Marshal(short)

	if err != nil {
		return
	}

	err = bbdb.DB.Update(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte("ShortenedUrls"))
		err = b.Put([]byte(short.ID), []byte(ser))
		return
	})

	return
}

func (bbdb *bboltDatabase) BreakCache(id string) (found bool) {
	_, found = bbdb.cache.Get(id)
	bbdb.cache.Delete(id)
	return
}
