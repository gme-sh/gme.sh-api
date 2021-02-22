package db

import "github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"

// DBCache is an interface which may be implemented by a StatsDatabase to provide cache functions.
// └ LocalCache
//    └ SharedCache
type DBCache interface {
	UpdateCache(u *short.ShortURL) (err error)
	BreakCache(id *short.ShortID) (err error)
	Get(key string) (interface{}, bool)
	GetShortURL(id *short.ShortID) *short.ShortURL
}
