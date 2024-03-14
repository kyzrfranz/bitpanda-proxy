package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

func Load() (port int, host string, logLevel slog.Level, logJson bool, cacheDuration time.Duration, err error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	durationStr := os.Getenv("CACHE_DURATION_MINUTES")
	durationMinutes, parseErr := strconv.Atoi(durationStr)
	if parseErr != nil {
		durationMinutes = 30
	}

	cacheDuration = time.Minute * time.Duration(durationMinutes)

	port, parseErr = strconv.Atoi(portStr)
	if parseErr != nil {
		port = 8080
	}

	host = os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	logJsonStr := os.Getenv("LOG_JSON")
	logJson, parseErr = strconv.ParseBool(logJsonStr)
	if parseErr != nil {
		logJson = false
	}

	loglevelFromEnv := os.Getenv("LOG_LEVEL")
	if loglevelFromEnv != "" {
		switch strings.ToUpper(loglevelFromEnv) {
		case "ERROR":
			logLevel = slog.LevelError
			break
		case "WARN", "warning":
			logLevel = slog.LevelWarn
			break
		case "INFO":
			logLevel = slog.LevelInfo
			break
		case "DEBUG":
			logLevel = slog.LevelDebug
			break
		default:
			logLevel = slog.LevelInfo
			break
		}
	} else {
		logLevel = slog.LevelInfo
	}

	return
}
