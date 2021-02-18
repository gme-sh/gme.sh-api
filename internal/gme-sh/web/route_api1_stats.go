package web

import (
	"encoding/json"
	"fmt"
	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/gorilla/mux"
	"log"
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

	log.Println("🚀", r.RemoteAddr, "requested to GET stats for", alias)

	// find shorted url
	if u, _ := ws.FindShort(&alias); u == nil {
		log.Println("    🤬 But", alias, "was not found")
		dieStats(w, "url not found")
		return
	}

	stats, err := ws.TemporaryDatabase.FindStats(&alias)
	if err != nil {
		log.Println("    🤬 But stats for", alias, "not found")
		dieStats(w, err)
		return
	}

	log.Println("    ✅ Stats:", stats)

	dieStats(w, &getStatsResponse{
		Success: true,
		Message: "Success",
		Stats:   stats,
	})
}
