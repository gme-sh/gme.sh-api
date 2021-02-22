package web

import (
	"encoding/base64"
	"fmt"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (ws *WebServer) handleRedirect(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := short.ShortID(vars["id"])

	if strings.Contains(string(id), ".") {
		log.Println("👋 Rejected", request.RemoteAddr, "because he/she/it requested file", id)
		_, _ = fmt.Fprintln(writer, "requested file. but this isn't a file server, got that?!")
		return
	}

	log.Println("🚀", request.RemoteAddr, "requested to GET redirect to", id)

	url, err := ws.PersistentDatabase.FindShortenedURL(&id)
	log.Println("url, err :=", url, err)
	if url == nil || err != nil {
		log.Println("    🤬 But it was not found:", err)
		b64id := base64.StdEncoding.EncodeToString([]byte(id))

		if ws.config.DryRedirect {
			_, _ = fmt.Fprintln(writer, "would redirect to /404/"+b64id, "with code 302 (disabled because DryRedirect = True)")
		} else {
			http.Redirect(writer, request, "/404/"+b64id, 302)
		}

		return
	}

	if ws.config.DryRedirect {
		_, _ = fmt.Fprintln(writer, "would redirect to", url.FullURL, "with code 302 (disabled because DryRedirect = True)")
	} else {
		http.Redirect(writer, request, url.FullURL, 302)
	}

	// add stats async
	if !url.IsTemporary() {
		log.Println("  📊 Add stats for", id.String())
		go func() {
			if err = ws.StatsDatabase.AddStats(&id); err != nil {
				log.Println("    ⏱ Stats could not be stored:", err)
			}
		}()
	} else {
		log.Println("  📊 Skipped stats for", id.String(), "because url was temporary")
	}
}
