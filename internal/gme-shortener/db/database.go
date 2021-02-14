package db

import (
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"log"
)

// Database -> Database Interface
type Database interface {
	FindShortenedURL(id string) (res *short.ShortURL, err error)
	SaveShortenedURL(url short.ShortURL) (err error)
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
