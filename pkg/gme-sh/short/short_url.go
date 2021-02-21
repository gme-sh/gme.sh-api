package short

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
	"time"
)

// ShortID -> ID of a ShortURL object of type string
type ShortID string

// ShortURL -> Structure for shortened urls
type ShortURL struct {
	ID           ShortID   `json:"id" bson:"id"`
	FullURL      string    `json:"full_url" bson:"full_url"`
	CreationDate time.Time `json:"creation_date" bson:"creation_date"`
	Secret       string    `json:"secret" bson:"secret"`
	Temporary    bool      `json:"temporary" bson:"is_temp"`
}

///////////////////////////////////////////////////////////////////////

// BsonUpdate returns a bson map (bson.M) with the field "$set": ShortURL
func (u *ShortURL) BsonUpdate() bson.M {
	return bson.M{
		"$set": u,
	}
}

// BsonFilter returns a bson map (bson.M) with the search option: "id": ShortID
func (id *ShortID) BsonFilter() bson.M {
	return bson.M{
		"id": id.String(),
	}
}

///////////////////////////////////////////////////////////////////////

// RedisKey is used to be able to specify keys at RedisKeyf (as magic constant)
type RedisKey uint64

const (
	// RedisKeyCountGlobal -> gme::short::{id}::count:g
	RedisKeyCountGlobal RedisKey = iota

	// RedisKeyCount60 -> gme::short::{id}::count:60
	RedisKeyCount60
)

// RedisKey returns gme::short::{id}
func (id *ShortID) RedisKey() string {
	return "gme::short::" + string(*id)
}

// RedisKeyf returns gme::short::{id}::{keys}
func (id *ShortID) RedisKeyf(keys ...interface{}) string {
	var builder strings.Builder

	// write previous
	builder.WriteString(id.RedisKey())

	for _, k := range keys {
		// write separator
		builder.WriteString("::")

		switch v := k.(type) {
		case string:
			builder.WriteString(v)
			break
		case RedisKey:
			var s string
			switch v {
			case RedisKeyCountGlobal:
				s = "count:g"
				break
			case RedisKeyCount60:
				s = "count:60"
				break
			default:
				log.Println("WARNING: Invalid Redis-Key-Format:", v)
				continue
			}
			builder.WriteString(s)
			break
		}
	}

	return builder.String()
}

// String converts the ShortID to a string
func (id *ShortID) String() string {
	return string(*id)
}

// String converts the ShortID to a byte array (splice)
func (id *ShortID) Bytes() []byte {
	return []byte(id.String())
}
