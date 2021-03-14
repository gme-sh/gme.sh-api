package db

import (
	"context"
	"encoding/json"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/tpl"
	"go.etcd.io/bbolt"
	"log"
	"time"
)

// PersistentDatabase
type bboltDatabase struct {
	database              *bbolt.DB
	cache                 DBCache
	shortedURLsBucketName []byte
	metaBucketName        []byte
	tplBucketName         []byte
}

// NewBBoltDatabase -> Create new BBoltDatabase
func NewBBoltDatabase(cfg *config.BBoltConfig, cache DBCache) (bbdb PersistentDatabase, err error) {
	// Open file {path} with permission-mode 0666
	// 0666 = All users can read/write, but cannot execute
	// 666 = 110 (u) 110 (g) 110 (o)
	//       rwx     rwx     rwx
	db, err := bbolt.Open(cfg.Path, cfg.FileMode, nil)
	if err != nil {
		return nil, err
	}
	bbdb = &bboltDatabase{
		database:              db,
		cache:                 cache,
		shortedURLsBucketName: []byte(cfg.ShortedURLsBucketName),
		metaBucketName:        []byte(cfg.MetaBucketName),
		tplBucketName:         []byte(cfg.TplBucketName),
	}
	return
}

/*
 * ==================================================================================================
 *                            P E R M A N E N T  D A T A B A S E
 * ==================================================================================================
 */

func (*bboltDatabase) ServiceName() string {
	return "BBolt"
}

func (bdb *bboltDatabase) HealthCheck(context.Context) (err error) {
	err = bdb.database.Update(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists([]byte("health_check")); err != nil {
			return
		}
		err = bucket.Put([]byte("ping"), []byte("pong"))
		return
	})
	return
}

///

func (bdb *bboltDatabase) FindShortenedURL(id *short.ShortID) (res *short.ShortURL, err error) {
	// check cache
	if u := bdb.cache.GetShortURL(id); u != nil {
		return u, nil
	}
	// load from bbolt
	var content []byte
	err = bdb.database.View(func(tx *bbolt.Tx) (err error) {
		bucket := tx.Bucket(bdb.shortedURLsBucketName)
		if bucket == nil {
			return
		}
		content = bucket.Get(id.Bytes())
		log.Println("BBolt :: Content =", string(content))
		return
	})
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &res)
	if err == nil {
		err = bdb.cache.UpdateCache(res)
	}

	return
}

func (bdb *bboltDatabase) SaveShortenedURL(short *short.ShortURL) (err error) {
	var shortAsJson []byte
	shortAsJson, err = json.Marshal(short)
	if err != nil {
		return
	}
	err = bdb.database.Update(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists(bdb.shortedURLsBucketName); err != nil {
			return
		}
		err = bucket.Put(short.ID.Bytes(), shortAsJson)
		log.Println("Saving Short URL", short, ":", err)
		return
	})
	if err == nil {
		err = bdb.cache.UpdateCache(short)
	}
	return
}

func (bdb *bboltDatabase) DeleteShortenedURL(id *short.ShortID) (err error) {
	err = bdb.database.Update(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists(bdb.shortedURLsBucketName); err != nil {
			return
		}
		err = bucket.Delete(id.Bytes())
		log.Println("Deleting Short URL #", id, ":", err)
		return
	})
	if err == nil {
		err = bdb.cache.BreakCache(id)
	}
	return
}

func (bdb *bboltDatabase) ShortURLAvailable(id *short.ShortID) bool {
	if u := bdb.cache.GetShortURL(id); u != nil {
		return false
	}
	return shortURLAvailable(bdb, id)
}

func (bdb *bboltDatabase) FindExpiredURLs() (res []*short.ShortURL, err error) {
	log.Println("CALLED FindExpiredURLs with time")
	now := time.Now()

	err = bdb.database.View(func(tx *bbolt.Tx) (err error) {
		bucket := tx.Bucket(bdb.shortedURLsBucketName)
		if bucket == nil {
			return
		}

		err = bucket.ForEach(func(_, v []byte) (err error) {
			var sh *short.ShortURL
			err = json.Unmarshal(v, &sh)
			if err != nil {
				log.Println("error #2")
				return
			}
			log.Println("Found:", sh.String())
			// check if expired
			if sh.ExpirationDate != nil && sh.ExpirationDate.Before(now) {
				// add to result
				res = append(res, sh)
			}
			return
		})
		if err != nil {
			log.Println("error #3")
		}
		return
	})

	return
}

func (bdb *bboltDatabase) GetLastExpirationCheck() (m *LastExpirationCheckMeta) {
	log.Println("CALLED GetLastExpirationCheck with time")
	m = &LastExpirationCheckMeta{
		LastCheck: time.Unix(5, 0),
	}

	var content []byte
	if err := bdb.database.View(func(tx *bbolt.Tx) (err error) {
		bucket := tx.Bucket(bdb.metaBucketName)
		if bucket == nil {
			return
		}

		content = bucket.Get([]byte("last_expired"))
		log.Println("BBolt :: Content =", string(content))
		return
	}); err != nil {
		return
	}

	_ = json.Unmarshal(content, &m)
	return
}

func (bdb *bboltDatabase) UpdateLastExpirationCheck(t time.Time) {
	log.Println("CALLED UpdateLastExpirationCheck with time", t)

	m := &LastExpirationCheckMeta{
		LastCheck: t,
	}

	if err := bdb.database.Update(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists(bdb.metaBucketName); err != nil {
			return
		}

		var by []byte
		by, err = json.Marshal(m)
		if err != nil {
			return
		}

		err = bucket.Put([]byte("last_expired"), by)
		return
	}); err != nil {
		log.Println("ERROR updating last expiration for bbolt:", err)
		return
	} else {
		log.Println("OK updating last expiration for bbolt:", err)
	}
}

func (bdb *bboltDatabase) FindTemplates() (templates []*tpl.Template, err error) {
	templates = []*tpl.Template{}
	err = bdb.database.View(func(tx *bbolt.Tx) (e error) {
		bucket := tx.Bucket(bdb.tplBucketName)
		if bucket == nil {
			return nil
		}
		e = bucket.ForEach(func(_, v []byte) error {
			var t = new(tpl.Template)
			if err := json.Unmarshal(v, t); err != nil {
				return err
			}
			templates = append(templates, t)
			return nil
		})
		return
	})
	return
}

func (bdb *bboltDatabase) SaveTemplate(t *tpl.Template) (err error) {
	err = bdb.database.Update(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		bucket, err = tx.CreateBucketIfNotExists(bdb.tplBucketName)
		if err != nil {
			return
		}
		var data []byte
		data, err = json.Marshal(t)
		if err != nil {
			return
		}
		err = bucket.Put([]byte(t.TemplateURL), data)
		return
	})
	return
}
