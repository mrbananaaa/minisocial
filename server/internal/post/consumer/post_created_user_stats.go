package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type UserStatsHandler struct {
	logger *slog.Logger
}

func NewUserStatsHandler(log *slog.Logger) *UserStatsHandler {
	return &UserStatsHandler{
		logger: log,
	}
}

func (h *UserStatsHandler) EventType() string {
	return string(domain.EvPostCreated)
}

func (h *UserStatsHandler) Handle(
	ctx context.Context,
	event events.EventContext,
) error {
	var payload domain.PostCreated
	if err := json.Unmarshal(event.Envelope.Payload, &payload); err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.logger.
		With("[CONSUMER]", "post.created_user-stats").
		Info("new post created",
			"message", "updating user stats",
			"author_id", payload.AuthorID,
			"post_id", payload.PostID,
		)

	return nil
}
