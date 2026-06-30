package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrbananaaa/minisocial/internal/app"
	"github.com/mrbananaaa/minisocial/internal/platform/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load env: %v", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.New(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize app: %v", err))
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- app.Run(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
}
