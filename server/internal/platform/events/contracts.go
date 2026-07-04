package events

type Event interface {
	IsEvent()
}

type EventSource interface {
	PullEvents() []Event
}
