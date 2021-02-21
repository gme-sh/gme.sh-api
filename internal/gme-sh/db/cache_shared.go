package db

import (
	"encoding/json"
	"fmt"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"log"
	"strings"
)

const (
	// SCacheChannelBreak -> Channel to subscribe for cache break notifications
	SCacheChannelBreak = "gme.sh-scache:break"

	// SCacheChannelUpdate -> Channel to subscribe for cache update notifications
	SCacheChannelUpdate = "gme.sh-scache:update"
)

// SharedCache only makes sense if you want to run multiple backend shards / servers at the same time.
// If a request is then cached on one server, this cache is passed on to all other servers via PubSub,
// whereby the requests to the database are brought to a minimum.
type SharedCache struct {
	NodeID string
	tempDB TemporaryDatabase
	local  *LocalCache
}

// NewSharedCache creates a new SharedCache object and returns it
func NewSharedCache(tempDB TemporaryDatabase) *SharedCache {
	return &SharedCache{
		NodeID: string(short.GenerateID(6, short.AlwaysTrue, 0)),
		tempDB: tempDB,
		local:  NewLocalCache(),
	}
}

// UpdateCache adds a new ShortURL object to the cache.
func (s *SharedCache) UpdateCache(u *short.ShortURL) (err error) {
	err = s.local.UpdateCache(u)
	if err != nil {
		return
	}
	log.Println("Publishing update for #", u.ID.String(), "...")
	err = s.tempDB.Publish(s.createSCacheUpdatePayload(u))
	return
}

// BreakCache removes the ShortURL object from the cache that matches the ShortID.
// No further check is made whether it was already in the cache or not.
// returns an error if there was an error publishing the break notification
func (s *SharedCache) BreakCache(id *short.ShortID) (err error) {
	// since the BreakCache from LocalCache always returns nil,
	// we don't have to deal with any exception here
	err = s.local.BreakCache(id)

	log.Println("Publishing break for #", id.String(), "...")
	err = s.tempDB.Publish(s.createSCacheBreakPayload(id))
	return
}

// Get returns an interface from the cache if it exists.
// Otherwise the interface is nil and the return value is false.
// Alias for LocalCache.Get()
func (s *SharedCache) Get(key string) (interface{}, bool) {
	return s.local.Get(key)
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

// Subscribe subscribes to SCacheChannelBreak + SCacheChannelUpdate channels and processes their messages
func (s *SharedCache) Subscribe() (err error) {
	err = s.tempDB.Subscribe(func(channel, payload string) {
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
			_ = s.local.UpdateCache(sh)
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

			id := short.ShortID(payload)
			log.Println("DEBUG x SCacheChannelBreak :: ShortID:", id)

			// remove from cache
			_ = s.local.BreakCache(&id)
			break
		default:
			log.Println("WARN: Subscibed to a channel we don't know")
		}
	}, SCacheChannelBreak, SCacheChannelUpdate)
	return
}
