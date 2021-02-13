package short

import "time"

type ShortURL struct {
	ID           string    `json:"id" bson:"id"`
	FullURL      string    `json:"full_url" bson:"full_url"`
	CreationDate time.Time `json:"creation_date" bson:"creation_date"`
}
