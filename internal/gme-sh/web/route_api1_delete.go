package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type deleteShortURLResponse struct {
	Success bool
	Message string
}

func dieDelete(w http.ResponseWriter, o interface{}) {
	var res *deleteShortURLResponse

	switch v := o.(type) {
	case error:
		res = &deleteShortURLResponse{
			Success: false,
			Message: v.Error(),
		}
		break
	case string:
		res = &deleteShortURLResponse{
			Success: false,
			Message: v,
		}
		break
	case *deleteShortURLResponse:
		res = v
		break
	}

	if res == nil {
		res = &deleteShortURLResponse{
			Success: false,
			Message: "an unknown error occurred.",
		}
	}

	var msg []byte
	var err error
	msg, err = json.Marshal(res)
	if err != nil {
		_, _ = fmt.Fprintln(w, res.Message)
		return
	}

	if res.Success {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

	_, _ = fmt.Fprintln(w, string(msg))
}

// mux
// DELETE /api/v1/{id}/{secret}
func (ws *WebServer) handleApiV1Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := short.ShortID(vars["id"])
	secret64 := vars["secret64"]

	s, err := base64.StdEncoding.DecodeString(secret64)
	if err != nil {
		dieDelete(w, err)
		return
	}
	secret := string(s)

	// console output
	log.Println("🚀", r.RemoteAddr, "requested to DELETE", id, "with secret", secret)

	// find short url
	sh, err := ws.PersistentDatabase.FindShortenedURL(&id)
	if err != nil {
		log.Println("    🤬 But", id, "was not found")
		dieDelete(w, err)
		return
	}

	if sh.Secret == "" {
		log.Println("    🤬 But", id, "was locked")
		dieDelete(w, "url is locked.")
		return
	}

	// compare secrets
	if secret != sh.Secret {
		log.Println("    🤬 But", id, "has a different secret than provided")
		dieDelete(w, "invalid secret.")
		return
	}

	// delete secret
	if err := ws.PersistentDatabase.DeleteShortenedURL(&id); err != nil {
		dieDelete(w, err)
		return
	}
	log.Println("    ✅ That shit (probably) worked")

	dieDelete(w, &deleteShortURLResponse{
		Success: true,
		Message: "success",
	})
}
