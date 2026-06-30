package nats

import (
	"fmt"

	"github.com/mrbananaaa/minisocial/internal/platform/config"
	"github.com/nats-io/nats.go"
)

func New(cfg config.NATSConfig) (*nats.Conn, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to nats server: %v", err)
	}

	if !conn.IsConnected() {
		return nil, fmt.Errorf("Nats disconnected")
	}

	return conn, nil
}
