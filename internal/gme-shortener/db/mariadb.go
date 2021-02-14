package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/config"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	_ "github.com/go-sql-driver/mysql" // mysql database driver
	"github.com/patrickmn/go-cache"
)

// NewMariaDB -> Create a new MariaDB connection
func NewMariaDB(config config.MariaConfig) (Database, error) {
	conn := fmt.Sprintf("%s:%s@%s", config.User, config.Password, config.DBName)
	database, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	query, err := ioutil.ReadFile("./setup/setup_mariadb.sql")

	if err != nil {
		return nil, err
	}

	database.Query(string(query))

	return &mariaDB{
		db:    database,
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

func (sql *mariaDB) ShortURLAvailable(id string) bool {
	return shortURLAvailable(sql, id)
}

func (sql *mariaDB) SaveShortenedURLWithExpiration(short short.ShortURL, expireAfter time.Duration) (err error) {
	return errors.New("Not implemented")
}
