package transport

import "net/http"

type HandlerWithContext func(http.ResponseWriter, *http.Request, <-chan struct{})

func WithContext(h HandlerWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		done := ctx.Done()

		h(w, r, done)
	}
}
