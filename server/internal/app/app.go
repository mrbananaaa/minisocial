package app

import (
	"context"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/config"
	"github.com/mrbananaaa/minisocial/internal/platform/db"
	"github.com/mrbananaaa/minisocial/internal/platform/logger"
	"github.com/mrbananaaa/minisocial/internal/platform/nats"
)

type App struct {
	Log *slog.Logger
	DB  *db.Database
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	log := logger.New(&logger.Config{Level: cfg.Logger.Level})

	db, err := db.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}
	_, err = nats.New(cfg.NATS)
	if err != nil {
		return nil, err
	}

	return &App{
		Log: log,
		DB:  db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.Log.Info("App is up and running!")

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Log.Warn("App is shuting down...")

	// var errs []error
	// append each cleaning methods error
	// then return it with erros.Join

	return nil
}
