package transport

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type HttpConfig struct {
	Port   int
	Host   string
	Logger *slog.Logger
}

func NewHttpServer(config *HttpConfig, opts ...HandlerOption) *http.Server {
	requestMux := http.NewServeMux()
	requestMux.HandleFunc("/status", WithCorsContext(WithContext(statusHandler)))

	for _, opt := range opts {
		requestMux.HandleFunc(opt.Path, WithCorsContext(WithContext(opt.Handler)))
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      requestMux,
	}

	go func() {
		config.Logger.Info(fmt.Sprintf("started listening on %s:%d", config.Host, config.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			config.Logger.Error(fmt.Sprintf("cannot listen on: %s", err))
			os.Exit(1)
		}
	}()

	return server
}
