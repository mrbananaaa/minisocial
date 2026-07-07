package jetstream

import (
	"fmt"
	"strings"
)

func (b *JetStreamBroker) withPrefix(subject string) string {
	return fmt.Sprintf("%s.%s", b.prefix, subject)
}

func (b *JetStreamBroker) fromPrefix(subject string) string {
	return strings.TrimPrefix(subject, b.prefix+".")
}
