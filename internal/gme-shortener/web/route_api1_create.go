package web

import (
	"encoding/json"
	"fmt"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"io/ioutil"
	"net/http"
	"time"
)

type createShortURLPayload struct {
	FullURL            string
	PreferredAlias     string
	ExpireAfterSeconds int
}

type createShortURLResponse struct {
	Success bool
	Message string
	Short   *short.ShortURL
}

func dieCreate(w http.ResponseWriter, o interface{}) {
	var res *createShortURLResponse

	switch v := o.(type) {
	case error:
		res = &createShortURLResponse{
			Success: false,
			Message: v.Error(),
			Short:   nil,
		}
		break
	case string:
		res = &createShortURLResponse{
			Success: false,
			Message: v,
			Short:   nil,
		}
		break
	case *createShortURLResponse:
		res = v
		break
	}

	if res == nil {
		res = &createShortURLResponse{
			Success: false,
			Message: "an unknown error occurred.",
			Short:   nil,
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
// POST application/json /api/v1/create
func (ws *WebServer) handleApiV1Create(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		dieCreate(w, err)
		return
	}

	// parse body
	var req *createShortURLPayload
	err = json.Unmarshal(body, &req)
	if err != nil {
		dieCreate(w, err)
		return
	}

	/// Generate ID and check if alias already exists
	if req.PreferredAlias == "" {
		if generated := short.GenerateShortID(ws.PersistentDatabase.ShortURLAvailable); generated != "" {
			req.PreferredAlias = generated
		} else {
			dieCreate(w, "generated id not available")
			return
		}
	}

	// check if alias already exists
	if available := ws.PersistentDatabase.ShortURLAvailable(req.PreferredAlias); !available {
		dieCreate(w, "preferred alias is not available")
		return
	}
	///

	// create short id
	sh := short.ShortURL{
		ID:           short.ShortID(req.PreferredAlias),
		FullURL:      req.FullURL,
		CreationDate: time.Now(),
	}

	// try to save shorted url
	if err := ws.PersistentDatabase.SaveShortenedURL(sh); err != nil {
		dieCreate(w, err)
		return
	}

	dieCreate(w, &createShortURLResponse{
		Success: true,
		Message: "success",
		Short:   &sh,
	})
}
