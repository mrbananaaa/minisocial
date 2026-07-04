package events

type Collector struct {
	records []EventSource
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Track(r EventSource) {
	c.records = append(c.records, r)
}

func (c *Collector) Flush() []Event {
	var events []Event

	for _, r := range c.records {
		events = append(events, r.PullEvents()...)
	}

	c.records = nil

	return events
}
