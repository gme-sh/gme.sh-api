package web

import (
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	// GET /heartbeat
	// This page returns the status 200 if the backend is running. otherwise a 5xx error
	// TODO: Fiber
	ws.Router.HandleFunc("/heartbeat", ws.handleApiV1Heartbeat).Methods(http.MethodGet)

	// POST /create
	// Used to create new short URLs
	app.Post("/create", ws.fiberRouteCreate)

	// DELETE /{id}/{secret}
	// Used to delete short URLs
	// TODO: Fiber
	ws.Router.HandleFunc("/{id}/{secret64}", ws.handleApiV1Delete).Methods(http.MethodDelete)

	// GET /stats/{id}
	// Used to retrieve stats for a short url
	// TODO: Fiber
	ws.Router.HandleFunc("/stats/{id}", ws.handleApiV1Stats).Methods(http.MethodGet)

	// GET /{id}
	// Used for redirection to long url
	// TODO: Fiber
	ws.Router.HandleFunc("/{id}", ws.handleRedirect).Methods(http.MethodGet)

	log.Println("🌎 Binding", ws.config.WebServer.Addr, "...")
	if err := http.ListenAndServe(ws.config.WebServer.Addr, ws.Router); err != nil {
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
