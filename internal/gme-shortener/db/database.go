package db

import (
	"log"
	"time"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
)

type Database interface {
	FindShortenedURL(id short.ShortID) (res *short.ShortURL, err error)
	BreakCache(id short.ShortID) (found bool)
	ShortURLAvailable(id short.ShortID) (available bool)
}

// PersistentDatabase -> PersistentDatabase Interface
type PersistentDatabase interface /* implements Database */ {
	// PersistentDatabase Functions
	SaveShortenedURL(url *short.ShortURL) (err error)
	DeleteShortenedURL(id *short.ShortID) (err error)

	// Database Functions
	FindShortenedURL(id short.ShortID) (res *short.ShortURL, err error)
	BreakCache(id short.ShortID) (found bool)
	ShortURLAvailable(id short.ShortID) (available bool)
}

type TemporaryDatabase interface /* implements Database */ {
	// TemporaryDatabase Functions
	SaveShortenedURLWithExpiration(url *short.ShortURL, expireAfter time.Duration) (err error)
	Heartbeat() (err error)
	FindStats(id short.ShortID) (stats *short.Stats, err error)
	AddStats(id short.ShortID) (err error)
	DeleteStats(id short.ShortID) (err error)

	// Database Functions
	FindShortenedURL(id short.ShortID) (res *short.ShortURL, err error)
	BreakCache(id short.ShortID) (found bool)
	ShortURLAvailable(id short.ShortID) (available bool)
}

// Must -> Don't use database, if some error occurred
func Must(db Database, err error) Database {
	if err != nil {
		log.Fatalln("🚨 Error creating database:", err)
		return nil
	}
	return db
}

func shortURLAvailable(db Database, id short.ShortID) bool {
	if url, err := db.FindShortenedURL(id); url != nil || err == nil {
		return false
	}
	return true
}
