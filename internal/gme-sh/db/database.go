package db

import (
	"context"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/tpl"
	"log"
	"time"
)

// PersistentDatabase functions
type PersistentDatabase interface {
	// HealthChecked
	ServiceName() string
	HealthCheck(context.Context) error

	// PersistentDatabase Functions
	// ShortURL
	SaveShortenedURL(*short.ShortURL) error
	DeleteShortenedURL(*short.ShortID) error
	FindShortenedURL(*short.ShortID) (*short.ShortURL, error)
	ShortURLAvailable(*short.ShortID) bool

	// Expiration
	FindExpiredURLs() ([]*short.ShortURL, error)
	GetLastExpirationCheck() *LastExpirationCheckMeta
	UpdateLastExpirationCheck(time.Time)

	// Template
	FindTemplates() ([]*tpl.Template, error)
	SaveTemplate(*tpl.Template) error
}

// StatsDatabase functions
type StatsDatabase interface {
	// HealthChecked
	ServiceName() string
	HealthCheck(context.Context) error

	// StatsDatabase Functions
	FindStats(*short.ShortID) (*short.Stats, error)
	AddStats(*short.ShortID) error
	DeleteStats(*short.ShortID) error
}

type PubSub interface {
	// HealthChecked
	ServiceName() string
	HealthCheck(context.Context) error

	Publish(string, string) error
	Subscribe(func(string, string), ...string) error
	Close() error
}

type HealthChecked interface {
	ServiceName() string
	HealthCheck(context.Context) error
}

// Must -> Don't use database, if some error occurred
func MustPersistent(db PersistentDatabase, err error) PersistentDatabase {
	if err != nil {
		log.Fatalln("🚨 Error creating persistent-database:", err)
		return nil
	}
	return db
}

func MustStats(db StatsDatabase, err error) StatsDatabase {
	if err != nil {
		log.Fatalln("🚨 Error creating stats-database:", err)
		return nil
	}
	return db
}

func MustPubSub(db PubSub, err error) PubSub {
	if err != nil {
		log.Fatalln("🚨 Error creating pubsub-database:", err)
		return nil
	}
	return db
}

func shortURLAvailable(db PersistentDatabase, id *short.ShortID) bool {
	if url, _ := db.FindShortenedURL(id); url != nil && !url.IsExpired() {
		return false
	}
	return true
}
