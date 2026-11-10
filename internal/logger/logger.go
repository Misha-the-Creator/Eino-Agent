package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetupLogger(isLocal bool) *slog.Logger {
	var log *slog.Logger

	switch isLocal {
	case true:
		log = slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		}))
	case false:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}