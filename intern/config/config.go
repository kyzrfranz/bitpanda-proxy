package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func Load() (port int, host string, logLevel slog.Level, logJson bool, err error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err = strconv.Atoi(portStr)
	if err != nil {
		return 8080, "", slog.LevelInfo, false, nil
	}

	host = os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	logJsonStr := os.Getenv("LOG_JSON")
	logJson, err = strconv.ParseBool(logJsonStr)
	if err != nil {
		return 8080, "", slog.LevelInfo, false, nil
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
