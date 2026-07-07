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
) (*messaging.Subscription, error) {
	stream, err := b.stream(ctx)
	if err != nil {
		return nil, err
	}

	consumer, err := b.consumer(ctx, stream, topic)
	if err != nil {
		return nil, err
	}

	return b.subscribe(ctx, consumer)
}

func (b *JetStreamBroker) stream(
	ctx context.Context,
) (jetstream.Stream, error) {

	stream, err := b.js.Stream(ctx, b.streamName)
	if err != nil {
		return nil, fmt.Errorf(
			"find stream %q: %w",
			b.streamName,
			err,
		)
	}

	return stream, nil
}

func (b *JetStreamBroker) consumer(
	ctx context.Context,
	stream jetstream.Stream,
	topic string,
) (jetstream.Consumer, error) {

	cfg := jetstream.ConsumerConfig{
		Durable:       "EVENT_PROCESSOR_MINISOCIAL",
		FilterSubject: b.withPrefix(topic),
		AckPolicy:     jetstream.AckExplicitPolicy,
	}

	consumer, err := stream.CreateOrUpdateConsumer(
		ctx,
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create consumer: %w",
			err,
		)
	}

	return consumer, nil
}

func (b *JetStreamBroker) subscribe(
	ctx context.Context,
	consumer jetstream.Consumer,
) (*messaging.Subscription, error) {

	const buffer = 64

	messages := make(chan messaging.Envelope, buffer)
	errs := make(chan error, buffer)

	var once sync.Once

	consumeCtx, err := consumer.Consume(
		b.handler(ctx, messages, errs),
	)
	if err != nil {
		close(messages)
		close(errs)

		return nil, err
	}

	closeFn := func() {
		once.Do(func() {
			consumeCtx.Stop()

			close(messages)
			close(errs)
		})
	}

	go func() {
		<-ctx.Done()
		closeFn()
	}()

	return messaging.NewSubscription(
		messages,
		errs,
		closeFn,
	), nil
}

func (b *JetStreamBroker) handler(
	ctx context.Context,
	messages chan<- messaging.Envelope,
	errs chan<- error,
) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		env, err := b.envelope(msg)
		if err != nil {
			select {
			case errs <- err:
			case <-ctx.Done():
			}
			return
		}

		select {
		case messages <- env:
		case <-ctx.Done():
			return
		}
	}
}

func (b *JetStreamBroker) envelope(
	msg jetstream.Msg,
) (messaging.Envelope, error) {

	meta, err := msg.Metadata()
	if err != nil {
		return messaging.Envelope{}, err
	}

	return messaging.NewEnvelope(
		strconv.FormatUint(
			meta.Sequence.Stream,
			10,
		),
		b.fromPrefix(msg.Subject()),
		msg.Data(),
		msg.Ack,
	), nil
}

func (b *JetStreamBroker) withPrefix(subject string) string {
	return fmt.Sprintf("%s.%s", b.prefix, subject)
}

func (b *JetStreamBroker) fromPrefix(subject string) string {
	return strings.TrimPrefix(subject, b.prefix+".")
}
