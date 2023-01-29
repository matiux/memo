package eventstore

import (
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"reflect"
)

type EventListenerError struct {
	EventListener string
	aggregate.DomainMessage
	OriginalError error
}

func (e EventListenerError) Error() string {
	return fmt.Sprintf(
		"Error in Event Listener `%v` with Message `%v`. Original error: %v",
		e.EventListener,
		e.DomainMessage.Event.Kind(),
		e.OriginalError,
	)
}

type EventListener interface {
	handle(message aggregate.DomainMessage) error
}

type EventBus interface {
	subscribe(eventListener EventListener)
	publish(domainMessages aggregate.DomainEventStream) error
}

type SimpleEventBus struct {
	eventListeners []EventListener
	queue          []aggregate.DomainMessage
	isPublishing   bool
}

func (eb *SimpleEventBus) subscribe(eventListener EventListener) {
	eb.eventListeners = append(eb.eventListeners, eventListener)
}

func (eb *SimpleEventBus) publish(domainMessages aggregate.DomainEventStream) error {
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
