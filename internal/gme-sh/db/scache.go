package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/patrickmn/go-cache"
	"log"
	"strings"
	"time"
)

const (
	SCacheChannelBreak  = "gme.sh-scache:break"
	SCacheChannelUpdate = "gme.sh-scache:update"
)

type SharedCache struct {
	NodeID  string
	Cache   *cache.Cache
	temp    TemporaryDatabase
	enabled bool
	share   bool
}

func NewSharedCache(tempDB TemporaryDatabase, enabled, share bool) *SharedCache {
	return &SharedCache{
		NodeID:  string(short.GenerateID(6, short.AlwaysTrue, 0)),
		Cache:   cache.New(15*time.Minute, 20*time.Minute),
		temp:    tempDB,
		enabled: enabled,
		share:   share,
	}
}

func (s *SharedCache) UpdateCache(u *short.ShortURL) (err error) {
	if !s.enabled {
		return errors.New("shared cache is not enabled")
	}
	s.Cache.Set(u.ID.String(), u, cache.DefaultExpiration)

	if s.share {
		log.Println("Publishing update for #", u.ID.String(), "...")
		err = s.temp.Publish(s.createSCacheUpdatePayload(u))
	}
	return
}

func (s *SharedCache) BreakCache(id *short.ShortID) (found bool, err error) {
	if !s.enabled {
		return false, nil
	}
	_, found = s.Cache.Get(id.String())
	s.Cache.Delete(id.String())

	if s.share {
		log.Println("Publishing break for #", id.String(), "...")
		err = s.temp.Publish(s.createSCacheBreakPayload(id))
	}
	return
}

func (s *SharedCache) Get(key string) (interface{}, bool) {
	if !s.enabled {
		return nil, false
	}
	return s.Cache.Get(key)
}

func extractID(in *string) (id string) {
	// strip
	*in = strings.TrimSpace(*in)
	// find space
	space := strings.Index(*in, " ")
	// extract id
	id = (*in)[:space]
	// update in
	*in = strings.TrimSpace((*in)[space:])
	return
}

func (s *SharedCache) sameID(id string) bool {
	if id == "" {
		log.Println("DEBUG :: Skipped scache update because the node-id was empty")
		return true
	}
	if id == s.NodeID {
		log.Println("DEBUG :: Skipped scache update because the node-id was the same")
		return true
	}
	return false
}

// publish gme.sh-scache:update <nodeid> <json>
func (s *SharedCache) createSCacheUpdatePayload(u *short.ShortURL) (string, string) {
	js, _ := json.Marshal(u)
	return SCacheChannelUpdate, fmt.Sprintf("%s %s", s.NodeID, string(js))
}

// publish gme.sh-scache:break <nodeid> <id>
func (s *SharedCache) createSCacheBreakPayload(i *short.ShortID) (string, string) {
	return SCacheChannelBreak, fmt.Sprintf("%s %s", s.NodeID, i.String())
}

func (s *SharedCache) Subscribe() (err error) {
	err = s.temp.Subscribe(func(channel, payload string) {
		switch channel {
		case SCacheChannelUpdate:
			// publish gme.sh-scache:update <nodeid> <json>
			log.Println("DEBUG x SCacheChannelUpdate :: Subscribe channel, payload (", channel, payload, ")")

			// get node id
			nodeID := extractID(&payload)
			log.Println("DEBUG x SCacheChannelUpdate :: NodeID:", nodeID)
			if s.sameID(nodeID) {
				return
			}

			log.Println("DEBUG x SCacheChannelUpdate :: JSON:", payload)

			// decode json to shortURL object
			var sh *short.ShortURL
			if err := json.Unmarshal([]byte(payload), &sh); err != nil {
				log.Println("DEBUG x SCacheChannelUpdate :: S-Cache: WARN: Invalid JSON for short object received")
				return
			}

			// save to PersistentDatabase
			s.Cache.Set(sh.ID.String(), sh, cache.DefaultExpiration)
			log.Println("DEBUG x SCacheChannelUpdate :: Cached short-url (by subscribe):", payload)

			break
		case SCacheChannelBreak:
			// publish gme.sh-scache:break <nodeid> <id>

			// get node id
			nodeID := extractID(&payload)
			log.Println("DEBUG x SCacheChannelBreak :: NodeID:", nodeID)
			if s.sameID(nodeID) {
				return
			}

			log.Println("DEBUG x SCacheChannelBreak :: ShortID:", payload)

			// remove from cache
			s.Cache.Delete(payload)

			break
		default:
			log.Println("WARN: Subscibed to a channl we don't know")
		}
	}, SCacheChannelBreak, SCacheChannelUpdate)
	return
}
