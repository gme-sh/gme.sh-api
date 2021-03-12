package short

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"regexp"
	"strings"
)

// ShortID -> ID of a ShortURL object of type string
type ShortID string

// String converts the ShortID to a string
func (id *ShortID) String() string {
	return string(*id)
}

// BsonFilter returns a bson map (bson.M) with the search option: "id": ShortID
func (id *ShortID) BsonFilter() bson.M {
	return bson.M{
		"id": id.String(),
	}
}

// Bytes converts the ShortID to a byte array (splice)
func (id *ShortID) Bytes() []byte {
	return []byte(id.String())
}

func (id *ShortID) IsEmpty() bool {
	return len(strings.TrimSpace(id.String())) <= 0
}

var (
	shortIDPattern = `^[\w-]{1,32}$`
	shortIDRegex   *regexp.Regexp
)

func init() {
	var err error
	shortIDRegex, err = regexp.Compile(shortIDPattern)
	if err != nil {
		log.Fatalln("error compiling pattern for short ids:", err)
		return
	}
}

func (id *ShortID) IsValid() bool {
	return !id.IsEmpty() && shortIDRegex.MatchString(id.String())
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
