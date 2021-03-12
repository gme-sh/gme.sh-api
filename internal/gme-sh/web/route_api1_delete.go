package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/shortreq"
	"github.com/gofiber/fiber/v2"
)

// DELETE /:id/:secret
func (ws *WebServer) fiberRouteDelete(ctx *fiber.Ctx) (err error) {
	id := short.ShortID(ctx.Params("id"))
	if id.IsEmpty() {
		return shortreq.ResponseErrEmptyID.Send(ctx)
	}

	secret := ctx.Params("secret")

	// find short url
	sh, err := ws.persistentDB.FindShortenedURL(&id)
	if err != nil {
		return shortreq.ResponseErrURLNotFound.SendWithMessage(ctx, err.Error())
	}

	// check if locked
	if sh.IsLocked() {
		return shortreq.ResponseErrLocked.Send(ctx)
	}

	// compare secrets
	if sh.Secret != secret {
		return shortreq.ResponseErrSecretMismatch.Send(ctx)
	}

	// delete
	err = ws.persistentDB.DeleteShortenedURL(&id)
	if err != nil {
		return
	}

	return shortreq.ResponseOkDeleted.SendWithData(ctx, sh)
}
