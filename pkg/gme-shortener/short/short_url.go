package short

import "time"

// ShortURL -> Structure for shortened urls
type ShortURL struct {
	ID           string    `json:"id" bson:"id"`
	FullURL      string    `json:"full_url" bson:"full_url"`
	CreationDate time.Time `json:"creation_date" bson:"creation_date"`
	Secret       string    `json:"secret" bson:"secret"`
}
