package config

import (
	"log"
	"strings"
	"time"
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

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
