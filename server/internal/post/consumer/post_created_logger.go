package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type PostCreatedLoggerHandler struct {
	logger *slog.Logger
}

func NewPostCreatedLoggerHandler(log *slog.Logger) *PostCreatedLoggerHandler {
	return &PostCreatedLoggerHandler{
		logger: log,
	}
}

func (h *PostCreatedLoggerHandler) EventType() string {
	return string(domain.EvPostCreated)
}

func (h *PostCreatedLoggerHandler) Handle(ctx context.Context, evtMsg events.EventMessage) error {
	var payload domain.PostCreated
	if err := json.Unmarshal(evtMsg.Payload, &payload); err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.logger.Info("new post created",
		"id", payload.PostID,
		"title", payload.Title,
	)

	return nil
}
