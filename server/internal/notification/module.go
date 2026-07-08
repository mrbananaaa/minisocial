package notification

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
)

type Module struct {
}

func New(
	db *pgxpool.Pool,
	broker messaging.Broker,
	logger *slog.Logger,
) *Module {
	log := logger.With("[INFRA]", "notifications_module")

	log.Info("notifications module cretead!")

	return &Module{}
}
