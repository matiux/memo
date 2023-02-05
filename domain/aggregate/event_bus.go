package aggregate

import (
	"fmt"
	"reflect"
)

type EventListenerError struct {
	EventListener string
	DomainMessage
	OriginalError error
}

func (e EventListenerError) Error() string {
	return fmt.Sprintf(
		"Error in Payload Listener `%v` with Message `%v`. Original error: %v",
		e.EventListener,
		e.DomainMessage.Payload.Kind(),
		e.OriginalError,
	)
}

type EventListener interface {
	Handle(message DomainMessage) error
}

type EventBus interface {
	Subscribe(eventListener EventListener)
	Publish(domainMessages DomainEventStream) error
}

type SimpleEventBus struct {
	EventListeners []EventListener
	Queue          []DomainMessage
	IsPublishing   bool
}

func (eb *SimpleEventBus) Subscribe(eventListener EventListener) {
	eb.EventListeners = append(eb.EventListeners, eventListener)
}

func (eb *SimpleEventBus) Publish(domainMessages DomainEventStream) error {
	for _, domainMessage := range domainMessages {
		eb.Queue = append(eb.Queue, domainMessage)
	}

	defer func() {
		eb.IsPublishing = false
	}()

	if !eb.IsPublishing {
		eb.IsPublishing = true

		for len(eb.Queue) > 0 {

			domainMessage := eb.Queue[0]
			eb.Queue = eb.Queue[1:]

			for _, eventListener := range eb.EventListeners {
				if err := eventListener.Handle(domainMessage); err != nil {
					return EventListenerError{
						EventListener: reflect.TypeOf(eventListener).Elem().Name(),
						DomainMessage: domainMessage,
						OriginalError: err,
					}
				}
			}
		}
	}

	return nil
}

func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		IsPublishing: false,
	}
}

type TraceableEventBus struct {
	EventBus
	tracing  bool
	recorded DomainEventStream
}

func (eb *TraceableEventBus) Publish(domainMessages DomainEventStream) error {
	if err := eb.EventBus.Publish(domainMessages); err != nil {
		return err
	}

	if !eb.tracing {
		return nil
	}

	for _, event := range domainMessages {
		eb.recorded = append(eb.recorded, event)
	}

	return nil
}

func (eb *TraceableEventBus) GetEvents() (events []DomainEvent) {

	for _, event := range eb.recorded {
		events = append(events, event.Payload)
	}

	return
}

func (eb *TraceableEventBus) Trace() {
	eb.tracing = true
}

func NewTraceableEventBus(eventBus EventBus) *TraceableEventBus {
	return &TraceableEventBus{
		EventBus: eventBus,
		tracing:  false,
	}
}
