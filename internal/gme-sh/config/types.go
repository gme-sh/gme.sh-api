package config

import (
	"log"
	"strings"
)

// BlockedHosts ...
type BlockedHosts struct {
	Hosts []string
}

// Set -> Set BlockedHosts.Hosts from a string
func (b *BlockedHosts) Set(val string) error {
	b.Hosts = strings.Split(val, ",")
	log.Println("BlockedHosts :: New val =", b.Hosts)
	return nil
}
