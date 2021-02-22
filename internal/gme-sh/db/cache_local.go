package db

import (
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/patrickmn/go-cache"
	"time"
)

// LocalCache is a purely local cache and only while the backend is running.
// If the application is terminated, the cache is also flushed.
type LocalCache struct {
	Cache *cache.Cache
}

// NewLocalCache creates and returns a new LocalCache with a *hodl* time of 5 minutes and a clear time of 10 minutes.
func NewLocalCache() *LocalCache {
	return &LocalCache{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// UpdateCache adds a new ShortURL object to the cache.
// Since no error can occur here, nil is always returned.
func (l *LocalCache) UpdateCache(u *short.ShortURL) (_ error) {
	l.Cache.Set(u.ID.String(), u, cache.DefaultExpiration)
	return
}

// BreakCache removes the ShortURL object from the cache that matches the ShortID.
// No further check is made whether it was already in the cache or not.
// Since no error can occur here, nil is always returned.
func (l *LocalCache) BreakCache(id *short.ShortID) (_ error) {
	l.Cache.Delete(id.String())
	return
}

// Get returns an interface from the cache if it exists.
// Otherwise the interface is nil and the return value is false.
func (l *LocalCache) Get(key string) (interface{}, bool) {
	return l.Cache.Get(key)
}

func (l *LocalCache) GetShortURL(id *short.ShortID) *short.ShortURL {
	i, found := l.Cache.Get(id.String())
	if !found {
		return nil
	}
	u, ok := i.(*short.ShortURL)
	if !ok {
		return nil
	}
	if u.IsExpired() {
		_ = l.BreakCache(id)
		return nil
	}
	return u
}
