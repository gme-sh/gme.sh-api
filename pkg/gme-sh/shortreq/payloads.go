package shortreq

import "github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"

type CreateShortURLPayload struct {
	FullURL            string        `json:"full_url"`
	PreferredAlias     short.ShortID `json:"preferred_alias"`
	ExpireAfterSeconds int           `json:"expire_after_seconds"`
}

type UpdatePoolPayload struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
