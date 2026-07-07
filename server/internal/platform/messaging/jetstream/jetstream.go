// Package jetstream are implementation of messaging.Broker interface
package jetstream

import (
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
