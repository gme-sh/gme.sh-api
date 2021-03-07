package web

import (
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/db"
	"github.com/gofiber/fiber/v2"
	"log"
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

	// POST /create
	// Used to create new short URLs
	app.Post("/create", ws.fiberRouteCreate)

	// DELETE /{id}/{secret}
	// Used to delete short URLs
	app.Delete("/:id/:secret", ws.fiberRouteDelete)

	// GET /stats/{id}
	// Used to retrieve stats for a short url
	app.Get("/stats/:id", ws.fiberRouteStats)

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
	app := fiber.New()
	return &WebServer{
		persistentDB: persistentDB,
		statsDB:      statsDB,
		config:       cfg,
		App:          app,
	}
}
