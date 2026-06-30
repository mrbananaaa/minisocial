package logger

import (
	"log/slog"
	"os"
)

func New(cfg *Config) *slog.Logger {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: parseLevel(cfg.Level),
		},
	)

	return slog.New(handler)
}
