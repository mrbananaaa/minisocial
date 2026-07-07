package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type LoggerHandler struct {
	logger *slog.Logger
}

func NewLoggerHandler(log *slog.Logger) *LoggerHandler {
	return &LoggerHandler{
		logger: log,
	}
}

func (h *LoggerHandler) EventType() string {
	return string(domain.EvPostCreated)
}

func (h *LoggerHandler) Handle(
	ctx context.Context,
	event events.EventContext,
) error {
	var payload domain.PostCreated
	if err := json.Unmarshal(event.Envelope.Payload, &payload); err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.logger.
		With("[CONSUMER]", "post.created_logger").
		Info("new post created",
			"topic", event.Envelope.Topic,
			"id", payload.PostID,
			"title", payload.Title,
		)

	return nil
}
