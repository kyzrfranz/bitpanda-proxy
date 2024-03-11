package main

import (
	"context"
	"github.com/kyzrlabs/bitpanda-proxy/intern/config"
	v1 "github.com/kyzrlabs/bitpanda-proxy/intern/pkg/bitpanda/v1"
	"github.com/kyzrlabs/bitpanda-proxy/intern/pkg/handlers"
	"github.com/kyzrlabs/bitpanda-proxy/intern/transport"
	"github.com/kyzrlabs/bitpanda-proxy/logging"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port, host, logLevel, logJson, err := config.Load()

	if err != nil {
		log.Fatalf("cannot get cli %v", err)
	}

	logger := logging.GetLogger(logLevel, logJson)

	httpConf := &transport.HttpConfig{
		Port:   port,
		Host:   host,
		Logger: logger,
	}

	if err != nil {
		logger.Error("cannot get cli", "err", err)
		os.Exit(1)
	}
	bitpandaService := v1.NewService(logger)
	txHandler := handlers.NewTransactionsHandler(bitpandaService)

	server := transport.NewHttpServer(httpConf, transport.WithFunc(txHandler.Path(), txHandler.HandlerFunc))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	slog.Info("server gracefully stopped")
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("could not shut down server", err)
	}
}
