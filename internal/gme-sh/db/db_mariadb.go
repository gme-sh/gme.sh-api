package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"

	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	_ "github.com/go-sql-driver/mysql" // mysql database driver
	"github.com/patrickmn/go-cache"
)

// PersistentDatabase
type mariaDB struct {
	db    *sql.DB
	cache *cache.Cache
}

// NewMariaDB -> Create a new MariaDB connection
func NewMariaDB(cfg *config.MariaConfig) (PersistentDatabase, error) {
	conn := fmt.Sprintf("%s:%s@%s", cfg.User, cfg.Password, cfg.DBName)
	database, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	query, err := ioutil.ReadFile("./setup/setup_mariadb.sql")
	if err != nil {
		return nil, err
	}

	if _, err := database.Query(string(query)); err != nil {
		return nil, err
	}

	return &mariaDB{
		db:    database,
		cache: cache.New(15*time.Minute, 10*time.Minute),
	}, nil
}

func (sql *mariaDB) FindShortenedURL(_ *short.ShortID) (res *short.ShortURL, err error) {
	return nil, errors.New("not implemented")
}

func (sql *mariaDB) SaveShortenedURL(_ *short.ShortURL) (err error) {
	return errors.New("not implemented")
}

func (sql *mariaDB) DeleteShortenedURL(_ *short.ShortID) (err error) {
	return errors.New("not implemented")
}

func (sql *mariaDB) BreakCache(id *short.ShortID) (found bool) {
	_, found = sql.cache.Get(id.String())
	sql.cache.Delete(id.String())
	return
}

func (sql *mariaDB) ShortURLAvailable(id *short.ShortID) bool {
	return shortURLAvailable(sql, id)
}
