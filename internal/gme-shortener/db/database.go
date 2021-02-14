package db

import (
	"log"
	"time"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
)

// Database -> Database Interface
type Database interface {
	FindShortenedURL(id string) (res *short.ShortURL, err error)
	SaveShortenedURL(url short.ShortURL) (err error)
	ShortURLAvailable(id string) bool
	SaveShortenedURLWithExpiration(url short.ShortURL, expireAfter time.Duration) (err error)
	BreakCache(id string) (found bool)
}

// Must -> Don't use database, if some error occurred
func Must(db Database, err error) Database {
	if err != nil {
		log.Fatalln("🚨 Error creating database:", err)
		return nil
	}
	return db
}

func shortURLAvailable(db Database, id string) bool {
	if url, err := db.FindShortenedURL(id); url != nil || err == nil {
		return false
	}
	return true
}
