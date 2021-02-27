package web

import (
	"fmt"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/db"
	"net/http"
)

// mux
// GET /api/v1/stats/{id}
func (ws *WebServer) handleApiV1Heartbeat(w http.ResponseWriter, r *http.Request) {
	err := db.LastHeartbeatError
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintln(w, err.Error())
	} else {
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, "")
	}
}
