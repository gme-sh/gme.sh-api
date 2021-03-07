package web

import (
	"fmt"
	"net/url"
	"strings"
)

func (ws *WebServer) getBlockedHostLocation(u *url.URL) (int, bool) {
	hosts := ws.config.BlockedHosts
	if hosts == nil || hosts.Hosts == nil {
		return -1, false
	}

	// make input lower case
	host := strings.ToLower(u.Host)

	// look for blocked host
	for index, block := range hosts.Hosts {
		if strings.ToLower(block) == host {
			return index, true
		}
	}
	return -1, false
}

func (ws *WebServer) checkDomain(u *url.URL) error {
	// block detection
	if index, blocked := ws.getBlockedHostLocation(u); blocked {
		return fmt.Errorf("host (%s) is blocked (%d)", u.Host, index)
	}

	return nil
}
