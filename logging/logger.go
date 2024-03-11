package logging

import (
	"log/slog"
	"os"
)

type LoggerGetter interface {
	Logger() *slog.Logger
}

func GetLogger(logLevel slog.Level, json bool) *slog.Logger {
	handlerOpts := slog.HandlerOptions{Level: logLevel}

	var handler slog.Handler
	handler = slog.NewTextHandler(os.Stdout, &handlerOpts)
	if json {
		handler = slog.NewJSONHandler(os.Stdout, &handlerOpts)
	}
	return slog.New(handler)
}
