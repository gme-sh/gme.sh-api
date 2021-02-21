package web

import (
	"errors"
	"fmt"
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

	return b.String()
}

func getLoopStatus(domain string) (status int, err error) {
	if res, found := loopCache.Get(domain); found {
		log.Println("-> Loop from cache")
		return res.(int), nil
	}

	b := domain + "/gme-sh-block"
	log.Println("👉 Checking loop:", b)

	var resp *http.Response
	resp, err = http.Get(b)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	status = resp.StatusCode
	loopCache.Set(domain, status, cache.DefaultExpiration)
	return
}

func (ws *WebServer) getBlockedHostLocation(u *url.URL) (int, bool) {
	// make input lower case
	host := strings.ToLower(u.Host)
	for index, block := range ws.config.BlockedHosts {
		if strings.ToLower(block) == host {
			return index, true
		}
	}
	return -1, false
}

var errorLoop = errors.New("loop detected")

func (ws *WebServer) checkDomain(u *url.URL) error {
	// block detection
	if index, blocked := ws.getBlockedHostLocation(u); blocked {
		return fmt.Errorf("host (%s) is blocked (%d)", u.Host, index)
	}

	// loop detection
	status, err := getLoopStatus(extractToDomain(u))
	if err != nil {
		return err
	}
	if status == http.StatusLoopDetected {
		return errorLoop
	}

	return nil
}
