// Package app composition root
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
	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/platform/logger"
	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
	"github.com/mrbananaaa/minisocial/internal/platform/messaging/jetstream"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/post"
	"github.com/mrbananaaa/minisocial/internal/user"
	createpost "github.com/mrbananaaa/minisocial/internal/workflows/create_post"
	createuser "github.com/mrbananaaa/minisocial/internal/workflows/create_user"

	postapi "github.com/mrbananaaa/minisocial/internal/post/api"
	postConsumer "github.com/mrbananaaa/minisocial/internal/post/consumer"
	userapi "github.com/mrbananaaa/minisocial/internal/user/api"
	userConsumer "github.com/mrbananaaa/minisocial/internal/user/consumer"
)

type App struct {
	Log *slog.Logger
	DB  *db.Database
	cfg *config.Config

	messageBroker messaging.Broker
	eventWorker   *events.Worker
	httpServer    *http.Server
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	log := logger.New(&logger.Config{Level: cfg.Logger.Level})

	database, err := db.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}
	dbPool := database.GetPool()

	messageBroker := jetstream.New(cfg.NATS.URL)
	if err := messageBroker.Connect(ctx); err != nil {
		return nil, err
	}

	dispatcher := events.NewDispatcher()

	dispatcher.Register(
		userConsumer.NewLoggerHandler(log),
		postConsumer.NewLoggerHandler(log),
		postConsumer.NewUserStatsHandler(log),
	)

	eventWorker := events.NewWorker(
		messageBroker,
		dispatcher,
		log,
	)

	txManager := db.NewTxManager(dbPool)

	userModule := user.New(dbPool)
	postModule := post.New(dbPool)

	createUserWorkflow := createuser.New(
		userModule.Service(),
		txManager,
	)
	createPostWorkflow := createpost.New(
		userModule.Service(),
		postModule.Service(),
		txManager,
	)

	validator := validation.New()

	userHandler := userapi.New(userModule.Service(), createUserWorkflow, validator)
	posthandler := postapi.New(postModule.Service(), createPostWorkflow, validator)

	r := NewRouter(Routes{
		userHandler: userHandler,
		postHandler: posthandler,
	})

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &App{
		Log:           log,
		DB:            database,
		cfg:           cfg,
		messageBroker: messageBroker,
		eventWorker:   eventWorker,
		httpServer:    s,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.Log.Info("App is up and running ✨",
		"port", a.cfg.HTTP.Port,
	)

	go a.eventWorker.Run(ctx)

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Log.Warn("App is shuting down...")
	a.DB.Close()

	var errs []error

	if err := a.messageBroker.Close(ctx); err != nil {
		errs = append(errs, err)
	}

	if err := a.httpServer.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
