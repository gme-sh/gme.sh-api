package web

import (
	"encoding/base64"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (ws *WebServer) handleRedirect(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := short.ShortID(vars["id"])

	log.Println("🚀", request.RemoteAddr, "requested to GET redirect to", id)

	// look for redirection
	url, err := ws.PersistentDatabase.FindShortenedURL(&id)
	log.Println("url, err :=", url, err)
	if url == nil || err != nil {
		log.Println("    🤬 But it was not found")
		b64id := base64.StdEncoding.EncodeToString([]byte(id))
		http.Redirect(writer, request, "/404/"+b64id, 302)
		return
	}

	http.Redirect(writer, request, url.FullURL, 302)

	// add stats async
	go func() {
		if err = ws.TemporaryDatabase.AddStats(&id); err != nil {
			log.Println("    ⏱ Stats could not be stored:", err)
		}
	}()
}
