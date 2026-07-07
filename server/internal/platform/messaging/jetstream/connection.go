package jetstream

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (b *JetStreamBroker) Connect(ctx context.Context) error {
	conn, err := nats.Connect(b.url)
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}
	b.nc = conn

	js, err := jetstream.New(conn)
	if err != nil {
		return fmt.Errorf("failed to create jetstream instance: %w", err)
	}
	b.js = js

	_, err = b.js.Stream(ctx, b.streamName)
	if err != nil && errors.Is(err, jetstream.ErrStreamNotFound) {
		_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
			Name:        b.streamName,
			Description: "Events stream",
			Subjects:    []string{b.prefix + ".>"},
			Retention:   jetstream.LimitsPolicy,
			MaxAge:      24 * time.Hour,
		})
		if err != nil {
			return fmt.Errorf("failed to create %s stream: %w", b.streamName, err)
		}
	}

	return nil
}

func (b *JetStreamBroker) Close(ctx context.Context) error {
	if b.nc == nil {
		return nil
	}

	if err := b.nc.Drain(); err != nil {
		b.nc.Close()
		return fmt.Errorf("failed to drain nats connection: %w", err)
	}

	return nil
}
