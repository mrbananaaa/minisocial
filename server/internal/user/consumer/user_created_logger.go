package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

type UserCreatedLoggerHandler struct {
	logger *slog.Logger
}

func NewUserCreatedLoggerHandler(log *slog.Logger) *UserCreatedLoggerHandler {
	return &UserCreatedLoggerHandler{
		logger: log,
	}
}

func (h *UserCreatedLoggerHandler) EventType() string {
	return string(domain.EvUserCreated)
}

func (h *UserCreatedLoggerHandler) Handle(ctx context.Context, evtMsg events.EventMessage) error {
	var payload domain.UserCreated
	if err := json.Unmarshal(evtMsg.Payload, &payload); err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.logger.Info("new user created",
		"id", payload.UserID,
		"email", payload.Email,
	)

	return nil
}
