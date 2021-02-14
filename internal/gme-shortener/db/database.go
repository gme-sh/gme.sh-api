package db

import (
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
)

// Database -> Database Interface
type Database interface {
	FindShortenedURL(id string) (res *short.ShortURL, err error)
	SaveShortenedURL(url short.ShortURL) (err error)
	BreakCache(id string) (found bool)
}
