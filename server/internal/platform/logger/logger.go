// Package logger contain configuration and logger wrapper
package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/Marlliton/slogpretty"
)

func New(cfg *Config) *slog.Logger {
	var handler slog.Handler

	handler = slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: parseLevel(cfg.Level),
		},
	)

	// TODO: refactor cfg
	if os.Getenv("APP_ENV") == "development" {
		handler = slogpretty.New(os.Stdout, &slogpretty.Options{
			Level:      parseLevel(cfg.Level),
			AddSource:  true,
			Colorful:   true,
			Multiline:  true,
			TimeFormat: time.Kitchen,
		})
	}

	return slog.New(handler)
}
