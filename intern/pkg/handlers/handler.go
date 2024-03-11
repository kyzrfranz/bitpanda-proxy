package handlers

import (
	"net/http"
)

type Handler interface {
	HandlerFunc(w http.ResponseWriter, r *http.Request, done <-chan struct{})
	Path() string
}
