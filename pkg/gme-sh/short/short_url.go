package short

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"regexp"
	"time"
)

// ShortURL -> Structure for shortened urls
type ShortURL struct {
	ID             ShortID    `json:"id" bson:"id"`
	FullURL        string     `json:"full_url" bson:"full_url"`
	CreationDate   time.Time  `json:"creation_date" bson:"creation_date"`
	ExpirationDate *time.Time `json:"expiration_date" bson:"expiration_date"`
	Secret         string     `json:"secret" bson:"secret"`
}

func (u *ShortURL) String() string {
	return fmt.Sprintf("ShortURL #%s (short) :: Long = %s | Created: %s", u.ID.String(), u.FullURL, u.CreationDate.String())
}

var (
	protocolPattern = `^https?:\/\/((.){1,})$`
	protocolRegex   *regexp.Regexp
)

func init() {
	var err error
	protocolRegex, err = regexp.Compile(protocolPattern)
	if err != nil {
		log.Fatalln("Error compiling protocol pattern:", err)
	}
}

func (u *ShortURL) GetRedirectURL() (url string) {
	url = u.FullURL
	if !protocolRegex.MatchString(url) {
		url = "//" + url
	}
	return
}

///////////////////////////////////////////////////////////////////////

// BsonUpdate returns a bson map (bson.M) with the field "$set": ShortURL
func (u *ShortURL) BsonUpdate() bson.M {
	return bson.M{
		"$set": u,
	}
}

func (u *ShortURL) IsExpired() bool {
	return u.IsTemporary() && time.Now().After(*u.ExpirationDate)
}

func (u *ShortURL) IsTemporary() bool {
	return u.ExpirationDate != nil
}

func (u *ShortURL) IsLocked() bool {
	return u.Secret == ""
}
