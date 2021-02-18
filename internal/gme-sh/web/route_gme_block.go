package web

import (
	"fmt"
	"net/http"
)

// mux
// GET /gme-sh-block
func (ws *WebServer) handleGMEBlock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusLoopDetected)
	_, _ = fmt.Fprintln(w, "block")
}
