package web

import (
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/db"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type WebServer struct {
	db.PersistentDatabase
	db.TemporaryDatabase
	config *config.Config
}

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

func NewWebServer(persistent db.PersistentDatabase, temporary db.TemporaryDatabase, cfg *config.Config) *WebServer {
	return &WebServer{
		persistent,
		temporary,
		cfg,
	}
}

func (ws *WebServer) FindShort(id *short.ShortID) (url *short.ShortURL, err error) {
	if ws.TemporaryDatabase != nil {
		url, err = ws.TemporaryDatabase.FindShortenedURL(id)
	}
	if url == nil || err != nil {
		url, err = ws.PersistentDatabase.FindShortenedURL(id)
	}
	return
}

func (ws *WebServer) DeleteShort(id *short.ShortID) (persError error, tempError error) {
	if ws.TemporaryDatabase != nil {
		tempError = ws.TemporaryDatabase.DeleteShortenedURL(id)
	}
	persError = ws.PersistentDatabase.DeleteShortenedURL(id)
	return
}

func (ws *WebServer) ShortAvailable(id *short.ShortID, temp bool) bool {
	if temp && ws.TemporaryDatabase != nil {
		return ws.TemporaryDatabase.ShortURLAvailable(id)
	} else {
		return ws.PersistentDatabase.ShortURLAvailable(id)
	}
}
