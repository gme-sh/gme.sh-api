package short

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
	"time"
)

type ShortID string

// ShortURL -> Structure for shortened urls
type ShortURL struct {
	ID           ShortID   `json:"id" bson:"id"`
	FullURL      string    `json:"full_url" bson:"full_url"`
	CreationDate time.Time `json:"creation_date" bson:"creation_date"`
	Secret       string    `json:"secret" bson:"secret"`
}

///////////////////////////////////////////////////////////////////////

func (u *ShortURL) BsonUpdate() bson.M {
	return bson.M{
		"$set": u,
	}
}

func (id *ShortID) BsonFilter() bson.M {
	return bson.M{
		"id": id.String(),
	}
}

///////////////////////////////////////////////////////////////////////

type RedisKey uint64

const (
	RedisKeyCountGlobal RedisKey = iota
	RedisKeyCount60
)

func (id *ShortID) RedisKey() string {
	return "gme::short::" + string(*id)
}
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
func (id *ShortID) String() string {
	return string(*id)
}
func (id *ShortID) Bytes() []byte {
	return []byte(id.String())
}
