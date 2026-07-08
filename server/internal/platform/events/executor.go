package events

import (
	"context"
	"sync"

	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
)

type Executor struct {
	dispatcher *Dispatcher
}

func NewExecutor(
	dispatcher *Dispatcher,
) *Executor {
	return &Executor{
		dispatcher: dispatcher,
	}
}

type ExecutionResult struct {
	Errors []error
}

func (e *Executor) Execute(
	ctx context.Context,
	env messaging.Envelope,
) *ExecutionResult {
	handlers := e.dispatcher.Resolve(env.Topic)

	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		h := handler

		wg.Go(func() {
			err := h.Handle(ctx, EventContext{
				Envelope: env,
				EventID:  env.ID,
			})
			if err != nil {
				errChan <- err
			}
		})
	}

	wg.Wait()

	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	return &ExecutionResult{
		Errors: errs,
	}
}
