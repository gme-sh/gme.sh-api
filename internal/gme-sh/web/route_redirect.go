package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (ws *WebServer) fiberRouteRedirect(ctx *fiber.Ctx) (err error) {
	id := short.ShortID(ctx.Params("id"))
	if id.Empty() {
		// TODO: redirect to 404?!
		return UserErrorResponse(ctx, "empty short id")
	}
	// check if requested file
	if strings.Contains(id.String(), ".") {
		return UserErrorResponse(ctx, "requested file")
	}
	// find short url
	var sh *short.ShortURL
	sh, err = ws.persistentDB.FindShortenedURL(&id)
	if sh == nil || err != nil {
		// TODO: redirect to 404?!
		return UserErrorResponse(ctx, "short url ["+id.String()+"] not found")
	}
	// check if expired
	if sh.IsExpired() {
		// delete
		err = ws.persistentDB.DeleteShortenedURL(&id)
		if err != nil {
			return
		}
		return UserErrorResponse(ctx, "expired")
	}
	// add stats
	if !sh.IsTemporary() {
		go func() {
			_ = ws.statsDB.AddStats(&id)
		}()
	}
	// dry redirect (debug)
	if ws.config.DryRedirect {
		return SuccessMessageResponse(ctx, "would redirect to ["+sh.FullURL+"]")
	}

	// redirect
	return ctx.Redirect(sh.FullURL, 302)
}
