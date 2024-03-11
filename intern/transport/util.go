package transport

import (
	"encoding/json"
	"errors"
	v1 "github.com/kyzrlabs/bitpanda-proxy/intern/pkg/bitpanda/v1"
	"net/http"
)

func JSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func HeaderOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func WriteResponse(w http.ResponseWriter, payload any) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, err error) {
	var bitpandaError v1.Error

	switch {
	case errors.As(err, &bitpandaError):
		http.Error(w, bitpandaError.Title, bitpandaError.Status)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
