package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

func (b *JetStreamBroker) Publish(ctx context.Context, topic string, payload []byte) error {
	natsMsg := nats.NewMsg(b.withPrefix(topic))
	natsMsg.Data = payload

	_, err := b.js.PublishMsg(ctx, natsMsg)
	if err != nil {
		return fmt.Errorf("failed to publish to JetStream: %w", err)
	}

	return nil
}
