package web

import (
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/shortreq"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func (ws *WebServer) fiberRoutePoolGet(ctx *fiber.Ctx) (err error) {
	var pool *short.Pool
	if pool, err = ws.findPoolOrDie(ctx); pool == nil {
		return
	}
	log.Println("pool, err :=", pool, err)
	// return pool
	return shortreq.ResponseOkPoolGet.SendWithData(ctx, pool)
}

func (ws *WebServer) fiberRoutePoolUpdate(ctx *fiber.Ctx) (err error) {
	var pool *short.Pool
	if pool, err = ws.findPoolOrDie(ctx); pool == nil {
		return
	}
	payload := new(shortreq.UpdatePoolPayload)
	if err = ctx.BodyParser(payload); err != nil {
		return
	}
	if len(payload.URL) > 400 {
		return shortreq.ResponseErrPoolInvalidURL.Send(ctx)
	}

	if pool.Entries == nil {
		pool.Entries = make(map[string][]*short.PoolEntry)
	}
	entry := &short.PoolEntry{
		URL:  payload.URL,
		Time: time.Now(),
	}
	if _, ok := pool.Entries[payload.Name]; !ok {
		pool.Entries[payload.Name] = []*short.PoolEntry{entry}
	} else {
		pool.Entries[payload.Name] = append(pool.Entries[payload.Name], entry)
	}
	if len(pool.Entries[payload.Name]) > 3 {
		pool.Entries[payload.Name] = pool.Entries[payload.Name][len(pool.Entries[payload.Name])-3:]
	}
	if err = ws.persistentDB.SavePool(pool); err != nil {
		return shortreq.ResponseErrPoolUpdating.SendWithMessage(ctx, err.Error())
	}
	return shortreq.ResponseOkPoolUpdating.Send(ctx)
}

func (ws *WebServer) findPoolOrDie(ctx *fiber.Ctx) (pool *short.Pool, err error) {
	id := short.PoolID(ctx.Params("id"))
	secret := ctx.Params("secret")
	// find pool
	if pool, err = ws.persistentDB.FindPool(&id); err != nil {
		pool = nil
		err = shortreq.ResponseErrPoolNotFound.SendWithMessage(ctx, err.Error())
		return
	}
	// check secret
	if pool.Secret != secret {
		pool = nil
		err = shortreq.ResponseErrPoolSecretMismatch.Send(ctx)
		return
	}
	return
}
