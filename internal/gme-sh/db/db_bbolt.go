package db

import (
	"encoding/json"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"go.etcd.io/bbolt"
	"log"
)

// PersistentDatabase
type bboltDatabase struct {
	database              *bbolt.DB
	cache                 DBCache
	shortedURLsBucketName []byte
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
	}
	return
}

/*
 * ==================================================================================================
 *                            P E R M A N E N T  D A T A B A S E
 * ==================================================================================================
 */

func (bdb *bboltDatabase) FindShortenedURL(id *short.ShortID) (res *short.ShortURL, err error) {
	// check cache
	if u := bdb.cache.GetShortURL(id); u != nil {
		return u, nil
	}
	// load from bbolt
	var content []byte
	err = bdb.database.View(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists(bdb.shortedURLsBucketName); err != nil {
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
