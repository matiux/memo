package domain

import (
	"fmt"
	"reflect"
)

type EventListenerError struct {
	EventListener string
	Message
	OriginalError error
}

func (e EventListenerError) Error() string {

	payload := e.Message.Payload

	return fmt.Sprintf(
		"Error in Payload Listener `%v` with Message `%v`. Original error: %v",
		e.EventListener,
		payload.Kind(),
		e.OriginalError,
	)
}

type EventListener interface {
	Handle(message Message) error
	Support(message Message) bool
}

type EventBus interface {
	Subscribe(eventListener EventListener)
	Publish(domainMessages EventStream) error
}

type SimpleEventBus struct {
	EventListeners []EventListener
	Queue          []Message
	IsPublishing   bool
}

func (eventBus *SimpleEventBus) Subscribe(eventListener EventListener) {
	eventBus.EventListeners = append(eventBus.EventListeners, eventListener)
}

func (eventBus *SimpleEventBus) Publish(domainMessages EventStream) error {

	for _, domainMessage := range domainMessages {
		eventBus.Queue = append(eventBus.Queue, domainMessage)
	}

	if !eventBus.IsPublishing {

		defer func() {
			eventBus.IsPublishing = false
		}()

		eventBus.IsPublishing = true

		for len(eventBus.Queue) > 0 {

			domainMessage := eventBus.Queue[0]
			eventBus.Queue = eventBus.Queue[1:]

			for _, eventListener := range eventBus.EventListeners {
				if !eventListener.Support(domainMessage) {
					continue
				}
				if err := eventListener.Handle(domainMessage); err != nil {
					return EventListenerError{
						EventListener: reflect.TypeOf(eventListener).Elem().Name(),
						Message:       domainMessage,
						OriginalError: err,
					}
				}
			}
		}
	}

	return nil
}

func NewSimpleEventBus(eventListeners []EventListener) *SimpleEventBus {
	return &SimpleEventBus{
		EventListeners: eventListeners,
		IsPublishing:   false,
	}
}

type TraceableEventBus struct {
	EventBus
	tracing  bool
	recorded EventStream
}

func (eb *TraceableEventBus) Publish(domainMessages EventStream) error {
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

func (eb *TraceableEventBus) GetEvents() (events []Event) {

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
