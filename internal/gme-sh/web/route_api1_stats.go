package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/shortreq"
	"github.com/gofiber/fiber/v2"
)

// GET /api/v1/stats/{id}
func (ws *WebServer) fiberRouteStats(ctx *fiber.Ctx) (err error) {
	id := short.ShortID(ctx.Params("id"))
	if id.Empty() {
		return shortreq.ResponseErrEmptyID.Send(ctx)
	}

	var stats *short.Stats
	stats, err = ws.statsDB.FindStats(&id)
	if err != nil {
		return
	}

	return shortreq.ResponseOkStats.SendWithData(ctx, stats)
}
