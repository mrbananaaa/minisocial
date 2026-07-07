package jetstream

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
	"github.com/nats-io/nats.go/jetstream"
)

func (b *JetStreamBroker) Subscribe(
	ctx context.Context,
	topic string,
) (*messaging.Subscription, error) {
	stream, err := b.stream(ctx)
	if err != nil {
		return nil, err
	}

	consumer, err := b.createConsumer(ctx, stream, topic)
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

func (b *JetStreamBroker) createConsumer(
	ctx context.Context,
	stream jetstream.Stream,
	topic string,
) (jetstream.Consumer, error) {
	if topic == "*" {
		topic = ">"
	}

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
		b.consumeHandler(ctx, messages, errs),
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

func (b *JetStreamBroker) consumeHandler(
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
