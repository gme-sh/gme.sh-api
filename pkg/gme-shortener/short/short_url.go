package short

import "time"

type ShortID string

// ShortURL -> Structure for shortened urls
type ShortURL struct {
	ID           ShortID   `json:"id" bson:"id"`
	FullURL      string    `json:"full_url" bson:"full_url"`
	CreationDate time.Time `json:"creation_date" bson:"creation_date"`
	Secret       string    `json:"secret" bson:"secret"`
}

func (id *ShortID) RedisKey() string {
	return "gme::short::" + string(*id)
}
func (id *ShortID) String() string {
	return string(*id)
}
