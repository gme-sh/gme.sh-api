package web

import (
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var loopCache = cache.New(time.Hour, 75*time.Minute)

// extractToDomain ...
func extractToDomain(u *url.URL) string {
	var b strings.Builder

	// https
	b.WriteString(u.Scheme)
	// ://
	b.WriteString("://")

	// user info
	if u.User != nil {
		s := u.User.String()
		b.WriteString(s)
		if len(s) > 0 {
			b.WriteByte('@')
		}
	}

	// host
	b.WriteString(u.Host)
	b.WriteByte('/')

	return b.String()
}

func getLoopStatus(u *url.URL) (status int, err error) {
	e := extractToDomain(u)
	if res, found := loopCache.Get(e); found {
		log.Println("-> Loop from cache")
		return res.(int), nil
	}

	b := e + "gme-sh-block"
	log.Println("👉 Checking loop:", b)

	var resp *http.Response
	resp, err = http.Get(b)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	status = resp.StatusCode
	loopCache.Set(e, status, cache.DefaultExpiration)
	return
}
