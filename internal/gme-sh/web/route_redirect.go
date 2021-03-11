package web

import (
	"fmt"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/shortreq"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (ws *WebServer) fiberRouteRedirect(ctx *fiber.Ctx) (err error) {
	id := short.ShortID(ctx.Params("id"))
	if id.Empty() {
		// TODO: redirect to 404?!
		return shortreq.ResponseErrEmptyID.Send(ctx)
	}
	// check if requested file
	if strings.Contains(id.String(), ".") {
		return shortreq.ResponseErrRequestedFile.Send(ctx)
	}
	// find short url
	var sh *short.ShortURL
	sh, err = ws.persistentDB.FindShortenedURL(&id)
	if sh == nil || err != nil {
		// TODO: redirect to 404?!
		return shortreq.ResponseErrURLNotFound.SendWithMessage(ctx,
			fmt.Sprintf("short url [%s] not found", id.String()))
	}
	// check if expired
	if sh.IsExpired() {
		// delete
		err = ws.persistentDB.DeleteShortenedURL(&id)
		if err != nil {
			return
		}
		return shortreq.ResponseErrExpired.SendWithData(ctx, sh)
	}
	// add stats
	if !sh.IsTemporary() {
		go func() {
			_ = ws.statsDB.AddStats(&id)
		}()
	}
	// dry redirect (debug)
	if ws.config.DryRedirect {
		return shortreq.ResponseOkRedirectDry.SendWithMessage(ctx,
			fmt.Sprintf("would redirect to [%s]", sh.FullURL))
	}

	// redirect
	return ctx.Redirect(sh.FullURL, 302)
}
