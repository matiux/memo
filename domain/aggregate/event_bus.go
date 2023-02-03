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
	handle(message DomainMessage) error
}

type EventBus interface {
	subscribe(eventListener EventListener)
	publish(domainMessages DomainEventStream) error
}

type SimpleEventBus struct {
	eventListeners []EventListener
	queue          []DomainMessage
	isPublishing   bool
}

func (eb *SimpleEventBus) subscribe(eventListener EventListener) {
	eb.eventListeners = append(eb.eventListeners, eventListener)
}

func (eb *SimpleEventBus) publish(domainMessages DomainEventStream) error {
	for _, domainMessage := range domainMessages {
		eb.queue = append(eb.queue, domainMessage)
	}

	defer func() {
		eb.isPublishing = false
	}()

	if !eb.isPublishing {
		eb.isPublishing = true

		for len(eb.queue) > 0 {

			domainMessage := eb.queue[0]
			eb.queue = eb.queue[1:]

			for _, eventListener := range eb.eventListeners {
				if err := eventListener.handle(domainMessage); err != nil {
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
		isPublishing: false,
	}
}

type TraceableEventBus struct {
	EventBus
	tracing  bool
	recorded DomainEventStream
}

func (eb *TraceableEventBus) publish(domainMessages DomainEventStream) error {
	if err := eb.EventBus.publish(domainMessages); err != nil {
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
