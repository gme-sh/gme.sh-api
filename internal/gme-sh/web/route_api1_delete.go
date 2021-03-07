package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gofiber/fiber/v2"
)

// DELETE /:id/:secret
func (ws *WebServer) fiberRouteDelete(ctx *fiber.Ctx) (err error) {
	id := short.ShortID(ctx.Params("id"))
	if id.Empty() {
		return UserErrorResponse(ctx, "empty short-id")
	}

	secret := ctx.Params("secret")

	// find short url
	sh, err := ws.persistentDB.FindShortenedURL(&id)
	if err != nil {
		return UserErrorResponse(ctx, err)
	}

	// check if locked
	if sh.IsLocked() {
		return UserErrorResponse(ctx, "url is locked")
	}

	// compare secrets
	if sh.Secret != secret {
		return UserErrorResponse(ctx, "secret mismatch")
	}

	// delete
	err = ws.persistentDB.DeleteShortenedURL(&id)
	if err != nil {
		return
	}

	return SuccessResponse(ctx)
}
