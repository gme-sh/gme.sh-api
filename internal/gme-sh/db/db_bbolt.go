package db

import (
	"encoding/json"
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"
	"time"

	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/patrickmn/go-cache"
	"go.etcd.io/bbolt"
)

// PersistentDatabase
type bboltDatabase struct {
	database              *bbolt.DB
	cache                 *cache.Cache
	shortedURLsBucketName []byte
}

// NewBBoltDatabase -> Create new BBoltDatabase
func NewBBoltDatabase(cfg *config.BBoltConfig) (bbdb PersistentDatabase, err error) {
	// Open file {path} with permission-mode 0666
	// 0666 = All users can read/write, but cannot execute
	// 666 = 110 (u) 110 (g) 110 (o)
	//       rwx     rwx     rwx
	db, err := bbolt.Open(cfg.Path, cfg.FileMode, nil)
	if err != nil {
		return nil, err
	}
	// create a cache that holds its objects for 10 minutes and deletes them after 15 minutes
	c := cache.New(10*time.Minute, 15*time.Minute)
	bbdb = &bboltDatabase{
		database:              db,
		cache:                 c,
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
	if o, found := bdb.cache.Get(id.String()); found {
		return o.(*short.ShortURL), nil
	}
	var content []byte
	err = bdb.database.View(func(tx *bbolt.Tx) (err error) {
		var bucket *bbolt.Bucket
		if bucket, err = tx.CreateBucketIfNotExists(bdb.shortedURLsBucketName); err != nil {
			return
		}
		content = bucket.Get(id.Bytes())
		return
	})
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &res)
	if err == nil {
		bdb.cache.Set(id.String(), res, cache.DefaultExpiration)
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
		bdb.cache.Set(short.ID.String(), short, cache.DefaultExpiration)
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
		bdb.cache.Delete(id.String())
	}
	return
}

func (bdb *bboltDatabase) BreakCache(id *short.ShortID) (found bool) {
	_, found = bdb.cache.Get(id.String())
	bdb.cache.Delete(id.String())
	return
}

func (bdb *bboltDatabase) ShortURLAvailable(id *short.ShortID) bool {
	if _, found := bdb.cache.Get(id.String()); found {
		return false
	}

	return shortURLAvailable(bdb, id)
}
