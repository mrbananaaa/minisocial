package events

import (
	"context"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
)

type Worker struct {
	consumer   Consumer
	dispatcher *Dispatcher
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
		logger:     logger.With("[INFRA]", "event_worker"),
	}
}

func (w *Worker) Run(ctx context.Context) {
	sub, err := w.consumer.Subscribe(ctx, ">")
	if err != nil {
		w.logger.Error("failed to subscribe",
			"err", err.Error(),
		)
		return
	}

	w.logger.Info("event worker is running ✨")
	for {
		select {
		case evt, ok := <-sub.Message:
			if !ok {
				return
			}
			evtMsg := toEventMessage(evt)

			handlers := w.dispatcher.Handlers(evt.Topic)
			// WARN: concurrent later
			for _, handler := range handlers {
				if err := handler.Handle(ctx, evtMsg); err != nil {
					w.logger.Error("event handler error",
						"err", err.Error(),
					)
					continue
				}
			}

			evt.AckFunc()

		case err := <-sub.Errors:
			w.logger.Info("Event error",
				"err", err.Error(),
			)

		case <-ctx.Done():
			w.logger.Info("event worker closed")
			return
		}
	}
}

func toEventMessage(m messaging.Message) EventMessage {
	return EventMessage{
		Topic:   m.Topic,
		Payload: m.Payload,
	}
}
