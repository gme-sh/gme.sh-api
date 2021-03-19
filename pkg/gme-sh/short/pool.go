package short

import "time"

type PoolID string

func (id *PoolID) String() string {
	return string(*id)
}

func (id *PoolID) Bytes() []byte {
	return []byte(*id)
}

type Pool struct {
	ID      PoolID                  `bson:"id" json:"id"`
	Created time.Time               `bson:"created" json:"created"`
	Secret  string                  `bson:"secret" json:"secret"`
	Entries map[string][]*PoolEntry `bson:"entries" json:"entries"`
}

type PoolEntry struct {
	URL  string    `bson:"url" json:"url"`
	Time time.Time `bson:"time" json:"time"`
}

/*
{
  "id": "ihawd78q898q89f",
  "secret": ij8!of9a.O194-983kdks!.4348md2kd9gm",
  "entries": {
    "mbp": [
       { "url": "https://github.com", "timestamp": 12387512323 },
       { "url": "https://github.com", "timestamp": 12387512323 },
       { "url": "https://github.com", "timestamp": 12387512323 }
    ]
  }
}
*/
