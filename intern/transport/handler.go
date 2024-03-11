package transport

import "net/http"

type HandlerOption struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request, <-chan struct{})
}

func WithFunc(path string, handlerFunc func(http.ResponseWriter, *http.Request, <-chan struct{})) HandlerOption {
	return HandlerOption{
		Path:    path,
		Handler: handlerFunc,
	}
}
