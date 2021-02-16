package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"
	"io/ioutil"

	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	_ "github.com/go-sql-driver/mysql" // mysql database driver
)

// PersistentDatabase
type mariaDB struct {
	db    *sql.DB
	cache DBCache
}

// NewMariaDB -> Create a new MariaDB connection
func NewMariaDB(cfg *config.MariaConfig, cache DBCache) (PersistentDatabase, error) {
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
		cache: cache,
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

func (sql *mariaDB) ShortURLAvailable(id *short.ShortID) bool {
	return shortURLAvailable(sql, id)
}
