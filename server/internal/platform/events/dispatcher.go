package events

type Dispatcher struct {
	handlers map[string][]EventHandler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *Dispatcher) Register(handlers ...EventHandler) {
	for _, handler := range handlers {
		d.handlers[handler.EventType()] = append(
			d.handlers[handler.EventType()],
			handler,
		)
	}
}

func (d *Dispatcher) Handlers(
	eventType string,
) []EventHandler {
	return d.handlers[eventType]
}
