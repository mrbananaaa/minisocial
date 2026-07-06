package messaging

import "context"

type Broker interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(ctx context.Context, topic string) (SubscribtionPayload, error)
}
