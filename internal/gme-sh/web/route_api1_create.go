package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type createShortURLPayload struct {
	FullURL            string        `json:"full_url"`
	PreferredAlias     short.ShortID `json:"preferred_alias"`
	ExpireAfterSeconds int           `json:"expire_after_seconds"`
}

var urlRegex *regexp.Regexp

func init() {
	var err error
	urlRegex, err = regexp.Compile("^(https?://)?((([\\da-z.-]+)\\.([a-z.]{2,6}))|[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})(:[0-9]+)?([/\\w .-]*)/?([/\\w .-]*)/?(([?&]).+?(=.+?)?)*$")
	if err != nil {
		log.Fatalln("error compiling regex:", err)
	}
}

func (ws *WebServer) fiberRouteCreate(ctx *fiber.Ctx) (err error) {
	req := new(createShortURLPayload)
	err = ctx.BodyParser(req)
	if err != nil {
		return
	}
	// check url
	if !urlRegex.MatchString(req.FullURL) {
		log.Println("    └ 🤬 But the URL didn't match the regex")

	}
	// parse given url
	u, err := url.Parse(req.FullURL)
	if err != nil {
		return UserErrorResponse(ctx, err)
	}
	// check if url is blacklisted
	if i, b := ws.getBlockedHostLocation(u); b {
		return UserErrorResponse(ctx, "domain is blocked (i#"+strconv.Itoa(i)+")")
	}
	// no custom alias set?
	// -> generate alias
	if req.PreferredAlias == "" {
		if generated := short.GenerateShortID(ws.persistentDB.ShortURLAvailable); !generated.Empty() {
			req.PreferredAlias = generated
		} else {
			return ServerErrorResponse(ctx, "no generated alias available")
		}
	} else {
		if available := ws.persistentDB.ShortURLAvailable(&req.PreferredAlias); !available {
			return UserErrorResponse(ctx, "alias not available")
		}
	}

	// expiration
	duration := time.Duration(req.ExpireAfterSeconds) * time.Second
	var expiration *time.Time
	if duration > 0 {
		v := time.Now().Add(duration)
		expiration = &v
	}

	// generate secret
	secret := short.GenerateID(32, short.AlwaysTrue, 0)

	// create short url object
	sh := &short.ShortURL{
		ID:             req.PreferredAlias,
		FullURL:        req.FullURL,
		CreationDate:   time.Now(),
		ExpirationDate: expiration,
		Secret:         secret.String(),
	}

	// save to database
	if err := ws.persistentDB.SaveShortenedURL(sh); err != nil {
		return ServerErrorResponse(ctx, err)
	}

	log.Println("    └ 💚 Looks like it worked out")
	return SuccessDataResponse(ctx, sh)
}
