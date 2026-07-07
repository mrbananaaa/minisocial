package events

import (
	"context"
	"log/slog"
)

type Worker struct {
	consumer   Consumer
	dispatcher *Dispatcher
	executor   *Executor
	logger     *slog.Logger
}

func NewWorker(
	consumer Consumer,
	dispatcher *Dispatcher,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		consumer:   consumer,
		dispatcher: dispatcher,
		executor:   NewExecutor(dispatcher),
		logger:     logger.With("[INFRA]", "event_worker"),
	}
}

func (w *Worker) Run(ctx context.Context) {
	sub, err := w.consumer.Subscribe(ctx, "*")
	if err != nil {
		w.logger.Error("failed to subscribe",
			"err", err.Error(),
		)
		return
	}
	defer func() {
		sub.Close()

		w.logger.Info("Event worker is closed ⚠️")
	}()

	w.logger.Info("Event worker is listening ✨")
	for {
		select {
		case env, ok := <-sub.Messages():
			if !ok {
				return
			}

			if err := w.executor.Execute(ctx, env); err != nil && len(err.Errors) > 0 {
				w.logger.Error("failed to execute event",
					"err", err.Errors,
				)

				continue
			}

			env.Ack()

		case err, ok := <-sub.Errors():
			if !ok {
				return
			}
			w.logger.Info("Event error",
				"err", err.Error(),
			)

		case <-ctx.Done():
			return
		}
	}
}
