package web

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (ws *WebServer) handleIndex(writer http.ResponseWriter, request *http.Request) {
	log.Println("handleIndex")
	// Do something
}

func (ws *WebServer) handleShortCreate(writer http.ResponseWriter, request *http.Request) {
	log.Println("handleShortCreate")
	// Do something
}

func (ws *WebServer) handleShortURLNotFound(writer http.ResponseWriter, req *http.Request) {
	log.Println("handleShortURLNotFound")

	vars := mux.Vars(req)
	b64id := vars["b64id"]
	id, err := base64.StdEncoding.DecodeString(b64id)
	if err != nil {
		id = []byte("")
	}
	_, _ = fmt.Fprintf(writer, "Short %s not found\n", string(id))
}

func (ws *WebServer) handleRedirect(writer http.ResponseWriter, request *http.Request) {
	log.Println("handleRedirect")

	vars := mux.Vars(request)
	id := vars["id"]

	// look for redirection
	url, err := ws.FindShortenedURL(id)
	log.Println("url, err :=", url, err)
	if url == nil || err != nil {
		b64id := base64.StdEncoding.EncodeToString([]byte(id))
		http.Redirect(writer, request, "/404/"+b64id, 302)
		return
	}

	// redirection found
	// TODO: Check if url is from shortener to prevent loops
	// writer.WriteHeader(200)
	http.Redirect(writer, request, url.FullURL, 302)

	// TODO: Redis stats
}
