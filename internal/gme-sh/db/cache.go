package db

import "github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"

type DBCache interface {
	UpdateCache(u *short.ShortURL) (err error)
	BreakCache(id *short.ShortID) (err error)
	Get(key string) (interface{}, bool)
}
