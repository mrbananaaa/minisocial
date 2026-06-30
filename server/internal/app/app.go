package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrbananaaa/minisocial/internal/platform/config"
	"github.com/mrbananaaa/minisocial/internal/platform/db"
	"github.com/mrbananaaa/minisocial/internal/platform/logger"
	"github.com/mrbananaaa/minisocial/internal/platform/nats"
	"github.com/mrbananaaa/minisocial/internal/post"
	"github.com/mrbananaaa/minisocial/internal/user"
	"github.com/mrbananaaa/minisocial/internal/workflows"
)

type App struct {
	Log        *slog.Logger
	DB         *db.Database
	httpServer *http.Server
	cfg        *config.Config
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	log := logger.New(&logger.Config{Level: cfg.Logger.Level})

	db, err := db.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}
	dbPool := db.GetPool()

	_, err = nats.New(cfg.NATS)
	if err != nil {
		return nil, err
	}

	userModule := user.New(dbPool, log)
	postModule := post.New(dbPool, log)

	workflowsModule := workflows.New(
		userModule.Service(),
		postModule.Service(),
	)

	r := NewRouter()

	userModule.RegisterRoutes(r)
	postModule.RegisterRoutes(r)
	workflowsModule.RegisterRoutes(r)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &App{
		Log:        log,
		DB:         db,
		httpServer: s,
		cfg:        cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.Log.Info("App is up and running ✨",
		"port", a.cfg.HTTP.Port,
	)

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Log.Warn("App is shuting down...")
	a.DB.Close()

	// append each cleaning methods error
	// then return it with erros.Join
	var errs []error

	if err := a.httpServer.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
