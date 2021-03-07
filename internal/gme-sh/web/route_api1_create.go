package web

import (
	"encoding/json"
	"fmt"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// CreateResponseOK is returned when everything worked fine
const CreateResponseOK = "success"

type createShortURLPayload struct {
	FullURL            string        `json:"full_url"`
	PreferredAlias     short.ShortID `json:"preferred_alias"`
	ExpireAfterSeconds int           `json:"expire_after_seconds"`
}

type createShortURLResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Short   *short.ShortURL `json:"short"`
}

var urlRegex *regexp.Regexp

func init() {
	var err error
	urlRegex, err = regexp.Compile("^(https?://)?((([\\da-z.-]+)\\.([a-z.]{2,6}))|[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})(:[0-9]+)?([/\\w .-]*)/?([/\\w .-]*)/?(([?&]).+?(=.+?)?)*$")
	if err != nil {
		log.Fatalln("error compiling regex:", err)
	}
}

func dieCreate(w http.ResponseWriter, o interface{}) {
	var res *createShortURLResponse

	switch v := o.(type) {
	case error:
		res = &createShortURLResponse{
			Success: false,
			Message: v.Error(),
			Short:   nil,
		}
		break
	case string:
		res = &createShortURLResponse{
			Success: false,
			Message: v,
			Short:   nil,
		}
		break
	case *createShortURLResponse:
		res = v
		break
	}

	if res == nil {
		res = &createShortURLResponse{
			Success: false,
			Message: "an unknown error occurred.",
			Short:   nil,
		}
	}

	var msg []byte
	var err error
	msg, err = json.Marshal(res)
	if err != nil {
		_, _ = fmt.Fprintln(w, res.Message)
		return
	}

	if res.Success {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

	_, _ = fmt.Fprintln(w, string(msg))
}

// fiber
// TODO: Add status
// TODO: return error
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
		return ctx.JSON(&createShortURLResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	// checks if url is blacklisted
	if err := ws.checkDomain(u); err != nil {
		return ctx.JSON(&createShortURLResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// no custom alias set?
	// -> generate alias
	if req.PreferredAlias == "" {
		if generated := short.GenerateShortID(ws.persistentDB.ShortURLAvailable); !generated.Empty() {
			req.PreferredAlias = generated
		} else {
			return ctx.JSON(&createShortURLResponse{
				Success: false,
				Message: "no generated alias available",
			})
		}
	} else {
		if available := ws.persistentDB.ShortURLAvailable(&req.PreferredAlias); !available {
			return ctx.JSON(&createShortURLResponse{
				Success: false,
				Message: "alias not available",
			})
		}
	}
	log.Println("    └ 👉 Preferred alias:", req.PreferredAlias)

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
		return ctx.JSON(&createShortURLResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	log.Println("    └ 💚 Looks like it worked out")

	return ctx.JSON(&createShortURLResponse{
		Success: true,
		Message: CreateResponseOK,
		Short:   sh,
	})
}
