package db

import (
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
	return nil, nil
}

func (bbdb *bboltDatabase) SaveShortenedURL(short short.ShortURL) (err error) {
	return nil
}

func (bbdb *bboltDatabase) BreakCache(id string) (found bool) {
	_, found = bbdb.cache.Get(id)
	bbdb.cache.Delete(id)
	return
}
