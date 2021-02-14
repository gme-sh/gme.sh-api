package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	_ "github.com/go-sql-driver/mysql" // mysql database driver
	"github.com/patrickmn/go-cache"
)

// NewMariaDB -> Create a new MariaDB connection
func NewMariaDB(user, password, path string) (Database, error) {
	conn := fmt.Sprintf("%s:%s@%s", user, password, path)
	dab, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	// TODO: Create Table either in go or manually
	dab.Query("CREATE TABLE ShortenedUrls")

	return &mariaDB{
		db:    dab,
		cache: cache.New(15*time.Minute, 10*time.Minute),
	}, nil
}

type mariaDB struct {
	db    *sql.DB
	cache *cache.Cache
}

func (sql *mariaDB) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	return nil, nil
}

func (sql *mariaDB) SaveShortenedURL(url short.ShortURL) (err error) {
	return nil
}

func (sql *mariaDB) BreakCache(id string) (found bool) {
	_, found = sql.cache.Get(id)
	sql.cache.Delete(id)
	return
}
