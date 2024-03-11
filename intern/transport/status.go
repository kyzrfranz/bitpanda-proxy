package transport

import (
	"net/http"
)

type Status struct {
	Healthy bool `json:"healthy"`
}

func statusHandler(w http.ResponseWriter, r *http.Request, done <-chan struct{}) {
	select {
	case <-done:
		break
	default:
		JSON(w)
		WriteResponse(w, Status{Healthy: true})
	}
}
