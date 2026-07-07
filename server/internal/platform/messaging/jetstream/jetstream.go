package jetstream

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type JetStreamBroker struct {
	url        string
	nc         *nats.Conn
	js         jetstream.JetStream
	streamName string
	prefix     string
}

func New(url string) *JetStreamBroker {
	return &JetStreamBroker{
		url:        url,
		streamName: "EVENTS",
		prefix:     "events",
	}
}

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

func (b *JetStreamBroker) Publish(ctx context.Context, topic string, payload []byte) error {
	natsMsg := nats.NewMsg(b.withPrefix(topic))
	natsMsg.Data = payload

	_, err := b.js.PublishMsg(ctx, natsMsg)
	if err != nil {
		return fmt.Errorf("failed to publish to JetStream: %w", err)
	}

	return nil
}

func (b *JetStreamBroker) Subscribe(
	ctx context.Context,
	topic string,
) (messaging.SubscribtionPayload, error) {
	subject := b.withPrefix(topic)

	stream, err := b.js.Stream(ctx, b.streamName)
	if err != nil {
		return messaging.SubscribtionPayload{}, fmt.Errorf("failed to find stream for subject %s: %w", b.streamName, err)
	}

	// WARN: Hardcoder consumer config
	// TODO: Pass durable name
	// durableName := fmt.Sprintf("EVENT_PROCESSOR_%s", strings.ReplaceAll(topic, ".", "_"))
	durableName := fmt.Sprintf("EVENT_PROCESSOR_%s", "minisocial")
	consumerCfg := jetstream.ConsumerConfig{
		FilterSubject: subject,
		Durable:       durableName,
		AckPolicy:     jetstream.AckExplicitPolicy,
	}
	cons, err := stream.CreateOrUpdateConsumer(ctx, consumerCfg)
	if err != nil {
		return messaging.SubscribtionPayload{}, fmt.Errorf("failed to create or update consumer: %w", err)
	}

	messageChan := make(chan messaging.Message, 64)
	errChan := make(chan error, 64)

	var closeOnce sync.Once

	consCtx, err := cons.Consume(func(msg jetstream.Msg) {
		meta, err := msg.Metadata()
		if err != nil {
			select {
			case errChan <- err:
			case <-ctx.Done():
			}
			return
		}

		m := messaging.Message{
			ID:      strconv.FormatUint(meta.Sequence.Stream, 10),
			Topic:   b.fromPrefix(msg.Subject()),
			Payload: msg.Data(),
			AckFunc: func() error { return msg.Ack() },
		}

		select {
		case messageChan <- m:
		case <-ctx.Done():
			return
		}
	})
	if err != nil {
		close(messageChan)
		close(errChan)

		return messaging.SubscribtionPayload{}, fmt.Errorf("failed to start consume: %w", err)
	}

	closeFunc := func() {
		closeOnce.Do(func() {
			consCtx.Stop()
		})
	}

	go func() {
		<-ctx.Done()
		closeFunc()

		// wait for nats stopping all active task then close the channel
		time.Sleep(200 * time.Millisecond)
		close(messageChan)
		close(errChan)
	}()

	return messaging.SubscribtionPayload{
		Message: messageChan,
		Errors:  errChan,
		Close:   closeFunc,
	}, nil
}

func (b *JetStreamBroker) withPrefix(subject string) string {
	return fmt.Sprintf("%s.%s", b.prefix, subject)
}

func (b *JetStreamBroker) fromPrefix(subject string) string {
	return strings.TrimPrefix(subject, b.prefix+".")
}
