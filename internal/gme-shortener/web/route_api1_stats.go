package web

import (
	"encoding/json"
	"fmt"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"github.com/gorilla/mux"
	"net/http"
)

type getStatsResponse struct {
	Success bool
	Message string
	Stats   *short.Stats
}

func dieStats(w http.ResponseWriter, o interface{}) {
	var res *getStatsResponse

	switch v := o.(type) {
	case error:
		res = &getStatsResponse{
			Success: false,
			Message: v.Error(),
			Stats:   nil,
		}
		break
	case string:
		res = &getStatsResponse{
			Success: false,
			Message: v,
			Stats:   nil,
		}
		break
	case *getStatsResponse:
		res = v
		break
	default:
		res = &getStatsResponse{
			Success: false,
			Message: "an unknown error occurred.",
			Stats:   nil,
		}
		break
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
// GET /api/v1/stats/{id}
func (ws *WebServer) handleApiV1Stats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := short.ShortID(vars["id"])

	// find shorted url
	if available := ws.PersistentDatabase.ShortURLAvailable(alias); available {
		dieStats(w, "url not found")
		return
	}

	// get stats
	// TODO: Get stats
	res := &getStatsResponse{
		Success: true,
		Message: "success",
		Stats: &short.Stats{
			Calls: 133742069,
		},
	}

	dieStats(w, res)
}
