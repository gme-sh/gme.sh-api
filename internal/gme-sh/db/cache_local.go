package db

import (
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/patrickmn/go-cache"
	"time"
)

type LocalCache struct {
	Cache *cache.Cache
}

func NewLocalCache() *LocalCache {
	return &LocalCache{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (l *LocalCache) UpdateCache(u *short.ShortURL) (_ error) {
	l.Cache.Set(u.ID.String(), u, cache.DefaultExpiration)
	return
}

func (l *LocalCache) BreakCache(id *short.ShortID) (_ error) {
	l.Cache.Delete(id.String())
	return
}

func (l *LocalCache) Get(key string) (interface{}, bool) {
	return l.Cache.Get(key)
}
