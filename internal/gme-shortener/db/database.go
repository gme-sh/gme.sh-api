package db

import (
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"log"
)

// Database -> Database Interface
type Database interface {
	FindShortenedURL(id string) (res *short.ShortURL, err error)
	SaveShortenedURL(url short.ShortURL) (err error)
	ShortURLAvailable(id string) bool
	BreakCache(id string) (found bool)
}

func Must(db Database, err error) Database {
	if err != nil {
		log.Fatalln("🚨 Error creating database:", err)
		return nil
	} else {
		return db
	}
}

func shortURLAvailable(db Database, id string) bool {
	if url, err := db.FindShortenedURL(id); url != nil || err == nil {
		return false
	}
	return true
}
