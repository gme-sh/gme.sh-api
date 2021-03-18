package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/shortreq"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var urlRegex *regexp.Regexp

func init() {
	var err error
	urlRegex, err = regexp.Compile(`^(https?://)?((([\dA-Za-z.-]+)\.([a-z.]{2,6}))|[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(:[0-9]+)?/?(.*)$`)
	if err != nil {
		log.Fatalln("error compiling regex:", err)
	}
}

func (ws *WebServer) fiberRouteCreate(ctx *fiber.Ctx) (err error) {
	req := new(shortreq.CreateShortURLPayload)
	err = ctx.BodyParser(req)
	if err != nil {
		return
	}
	// check url
	if !urlRegex.MatchString(req.FullURL) {
		log.Println("    └ 🤬 But the URL didn't match the regex")
		return shortreq.ResponseErrInvalidURL.Send(ctx)
	}
	// parse given url
	u, err := url.Parse(req.FullURL)
	if err != nil {
		return shortreq.ResponseErrInvalidURL.Send(ctx)
	}
	// check if url is blacklisted
	if i, b := ws.getBlockedHostLocation(u); b {
		return shortreq.ResponseErrDomainBlocked.SendWithMessage(ctx,
			"domain is blocked (i#"+strconv.Itoa(i)+")")
	}
	// no custom alias set?
	// -> generate alias
	if req.PreferredAlias == "" {
		if generated := short.GenerateShortID(ws.persistentDB.ShortURLAvailable); !generated.IsEmpty() {
			req.PreferredAlias = generated
		} else {
			return shortreq.ResponseErrGeneratedAliasNotAvailable.Send(ctx)
		}
	} else {
		if available := ws.persistentDB.ShortURLAvailable(&req.PreferredAlias); !available {
			return shortreq.ResponseErrAliasOccupied.Send(ctx)
		}
	}

	// check short id
	if !req.PreferredAlias.IsValid() {
		return shortreq.ResponseErrInvalidID.Send(ctx)
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
		return shortreq.ResponseErrDatabaseSave.SendWithMessage(ctx, err.Error())
	}

	log.Println("    └ 💚 Looks like it worked out")
	return shortreq.ResponseOkCreate.SendWithData(ctx, sh)
}
