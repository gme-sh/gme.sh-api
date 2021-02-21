package config

import (
	"log"
	"strings"
)

type BlockedHosts struct {
	Hosts []string
}

func (b *BlockedHosts) Set(val string) error {
	b.Hosts = strings.Split(val, ",")
	log.Println("BlockedHosts :: New val =", b.Hosts)
	return nil
}
