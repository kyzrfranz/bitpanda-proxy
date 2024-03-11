package handlers

import (
	v1 "github.com/kyzrlabs/bitpanda-proxy/intern/pkg/bitpanda/v1"
	"github.com/kyzrlabs/bitpanda-proxy/intern/transport"
	"net/http"
)

const TransactionsPath = "/transactions"

type transactionsHandler struct {
	service v1.Service
}

func NewTransactionsHandler(service v1.Service) Handler {
	return &transactionsHandler{
		service: service,
	}
}

func (h *transactionsHandler) HandlerFunc(w http.ResponseWriter, r *http.Request, done <-chan struct{}) {
	select {
	case <-done:
		break
	default:
		options := v1.FilterFromReq(r)
		apiKey := r.Header.Get("X-API-KEY")
		tx, err := h.service.Transactions(apiKey, options...)
		if err != nil {
			transport.Error(w, err)
			return
		} else if len(tx.Errors) > 0 {
			transport.Error(w, tx.Errors[0])
		} else {
			transport.JSON(w)
			transport.WriteResponse(w, tx)
		}
	}
}

func (h *transactionsHandler) Path() string {
	return TransactionsPath
}
