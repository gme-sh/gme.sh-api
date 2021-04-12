package web

import (
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"net/http"
	"strings"
	"time"
)

// WebServer struct that holds databases and configs
type WebServer struct {
	persistentDB db.PersistentDatabase
	statsDB      db.StatsDatabase
	config       *config.Config
	App          *fiber.App
}

// Start starts the WebServer and listens on the specified port
func (ws *WebServer) Start() {
	app := ws.App

	// logger middleware
	app.Use(logger.New())

	// / -> redirect to github
	app.Get("/", func(ctx *fiber.Ctx) error {
		u := ws.config.WebServer.DefaultURL
		if u == "" {
			u = "https://github.com/gme-sh/gme.sh-api"
		}
		return ctx.Redirect(u)
	})

	// limiter middleware
	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		Next: func(c *fiber.Ctx) bool {
			// do not skip /create, /delete
			if c.Method() == http.MethodDelete || c.Method() == http.MethodPost {
				return false
			}
			// do not skip stats
			if strings.HasPrefix(c.Path(), "/stats") {
				return false
			}
			return true
		},
	}))

	// panic middleware
	app.Use(recover2.New(recover2.Config{
		EnableStackTrace: true,
	}))

	// monitor "middleware"
	app.Get("/dashboard", monitor.New())

	// POST /create
	// Used to create new short URLs
	app.Post("/create", ws.fiberRouteCreate)

	// DELETE /{id}/{secret}
	// Used to delete short URLs
	app.Delete("/:id/:secret", ws.fiberRouteDelete)

	// GET /stats/{id}
	// Used to retrieve stats for a short url
	app.Get("/stats/:id", ws.fiberRouteStats)

	// POOL
	app.Get("/pool/:id/:secret", ws.fiberRoutePoolGet)
	app.Post("/pool/:id/:secret", ws.fiberRoutePoolUpdate)

	// GET /{id}
	// Used for redirection to long url
	app.Get("/:id", ws.fiberRouteRedirect)

	log.Println("🌎 Binding", ws.config.WebServer.Addr, "...")
	if err := app.Listen(ws.config.WebServer.Addr); err != nil {
		log.Fatalln("    └ ❌ FAILED:", err)
	}
}

// NewWebServer returns a new WebServer object (reference)
func NewWebServer(persistentDB db.PersistentDatabase, statsDB db.StatsDatabase, cfg *config.Config) *WebServer {
	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})
	return &WebServer{
		persistentDB: persistentDB,
		statsDB:      statsDB,
		config:       cfg,
		App:          app,
	}
}
