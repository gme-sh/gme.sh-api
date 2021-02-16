package web

import (
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type WebServer struct {
	db.PersistentDatabase
	db.TemporaryDatabase
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

func NewWebServer(persistent db.PersistentDatabase, temporary db.TemporaryDatabase) *WebServer {
	return &WebServer{
		persistent,
		temporary,
	}
}
