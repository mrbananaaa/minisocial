package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
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
	return string(domain.EvUserCreated)
}

func (h *LoggerHandler) Handle(
	ctx context.Context,
	event events.EventContext,
) error {
	var payload domain.UserCreated
	if err := json.Unmarshal(event.Envelope.Payload, &payload); err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.logger.
		With("[CONSUMER]", "user.created_logger").
		Info("new user created",
			"topic", event.Envelope.Topic,
			"id", payload.UserID,
			"email", payload.Email,
		)

	return nil
}
