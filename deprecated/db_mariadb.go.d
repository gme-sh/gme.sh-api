import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"io/ioutil"

	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	_ "github.com/go-sql-driver/mysql" // mysql database driver
)

// PersistentDatabase
type mariaDB struct {
	db    *sql.DB
	cache db2.DBCache
}

// NewMariaDB -> Create a new MariaDB connection
func NewMariaDB(cfg *config.MariaConfig, cache db2.DBCache) (db2.PersistentDatabase, error) {
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
	return db2.shortURLAvailable(sql, id)
}
