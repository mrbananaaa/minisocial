// Package nats wrapping nats client implementation
package nats

import (
	"fmt"

	"github.com/mrbananaaa/minisocial/internal/platform/config"
	"github.com/nats-io/nats.go"
)

func New(cfg config.NATSConfig) (*nats.Conn, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to nats server: %v", err)
	}

	if !conn.IsConnected() {
		return nil, fmt.Errorf("nats disconnected")
	}

	return conn, nil
}
