package web

import (
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/db"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type WebServer struct {
	db.PersistentDatabase
	db.TemporaryDatabase
	redis *redis.Client
}

func (ws *WebServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/{id}", ws.handleRedirect)
	router.HandleFunc("/api/v1/create", ws.handleApiV1Create).Methods("POST")
	router.HandleFunc("/api/v1/stats/{id}", ws.handleApiV1Stats).Methods("GET")
	router.HandleFunc("/api/v1/heartbeat", ws.handleApiV1Heartbeat).Methods("GET")
	router.HandleFunc("/api/v1/{id}/{secret64}", ws.handleApiV1Delete).Methods("DELETE")

	if err := http.ListenAndServe(":1336", router); err != nil {
		log.Fatalln("Error listening and serving:", err)
	}
}

func NewWebServer(persistent db.PersistentDatabase, temporary db.TemporaryDatabase, red *redis.Client) *WebServer {
	return &WebServer{
		persistent,
		temporary,
		red,
	}
}
