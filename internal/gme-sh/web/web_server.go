package web

import (
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// WebServer struct that holds databases and configs
type WebServer struct {
	db.PersistentDatabase
	db.StatsDatabase
	config *config.Config
}

// Start starts the WebServer and listens on the specified port
func (ws *WebServer) Start() {
	router := mux.NewRouter()

	// GET /gme-sh-block
	// Used to check if a web page blocks gme.sh (e.g. for loop prevention).
	router.HandleFunc("/gme-sh-block", ws.handleGMEBlock).Methods(http.MethodGet)

	// GET /heartbeat
	// This page returns the status 200 if the backend is running. otherwise a 5xx error
	router.HandleFunc("/heartbeat", ws.handleApiV1Heartbeat).Methods(http.MethodGet)

	// POST /create
	// Used to create new short URLs
	router.HandleFunc("/create", ws.handleApiV1Create).Methods(http.MethodPost)

	// DELETE /{id}/{secret}
	// Used to delete short URLs
	router.HandleFunc("/{id}/{secret64}", ws.handleApiV1Delete).Methods(http.MethodDelete)

	// GET /stats/{id}
	// Used to retrieve stats for a short url
	router.HandleFunc("/stats/{id}", ws.handleApiV1Stats).Methods(http.MethodGet)

	// GET /{id}
	// Used for redirection to long url
	router.HandleFunc("/{id}", ws.handleRedirect).Methods(http.MethodGet)

	log.Println("🌎 Binding", ws.config.WebServer.Addr, "...")
	if err := http.ListenAndServe(ws.config.WebServer.Addr, router); err != nil {
		log.Fatalln("    └ ❌ FAILED:", err)
	}
}

// NewWebServer returns a new WebServer object (reference)
func NewWebServer(persistent db.PersistentDatabase, temporary db.StatsDatabase, cfg *config.Config) *WebServer {
	return &WebServer{
		persistent,
		temporary,
		cfg,
	}
}
