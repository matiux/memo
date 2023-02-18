package domain

import (
	"fmt"
)

var DuplicatePlayhead = fmt.Errorf("duplicate playhead not allowed")

type EventStore interface {
	Append(id EntityId, eventStream EventStream) error
	Load(id EntityId) (EventStream, error)
	//LoadFromPlayhead(id EntityId, playhead Playhead) EventStream
}

type InMemoryEventStore struct {
	Stream map[string]map[Playhead]Message
}

func (e *InMemoryEventStore) Append(id EntityId, eventStream EventStream) error {

	stringId := id.(UUIDv4).Val

	if _, exists := e.Stream[stringId]; !exists {
		e.Stream[stringId] = make(map[Playhead]Message)
	}

	e.assertStream(e.Stream[stringId], eventStream)

	for _, domainMessage := range eventStream {
		e.Stream[stringId][domainMessage.Playhead] = domainMessage
	}

	return nil
}

func (e *InMemoryEventStore) assertStream(events map[Playhead]Message, eventsToAppend EventStream) {

	for _, event := range eventsToAppend {
		if _, exists := events[event.Playhead]; exists {
			panic(DuplicatePlayhead)
		}
	}
}

func (e *InMemoryEventStore) Load(id EntityId) (EventStream, error) {

	stringId := id.(UUIDv4).Val

	if _, exists := e.Stream[stringId]; !exists {
		return nil, fmt.Errorf("aggregate with id '%v' not found", id)
	}

	domainEventStream := EventStream{}

	for _, domainMessage := range e.Stream[stringId] {
		domainEventStream = append(domainEventStream, domainMessage)
	}

	return domainEventStream, nil
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		Stream: make(map[string]map[Playhead]Message),
	}
}

type TraceableEventStore struct {
	EventStore
	tracing  bool
	recorded EventStream
}

func (e *TraceableEventStore) Append(id EntityId, eventStream EventStream) error {
	e.EventStore.Append(id, eventStream)

	if !e.tracing {
		return nil
	}

	for _, event := range eventStream {
		e.recorded = append(e.recorded, event)
	}

	return nil
}

func (e *TraceableEventStore) GetEvents() (events []Event) {

	for _, event := range e.recorded {
		events = append(events, event.Payload)
	}

	return
}

func (e *TraceableEventStore) Trace() {
	e.tracing = true
}

func (e *TraceableEventStore) ClearEvents() {
	e.recorded = EventStream{}
}

func NewTraceableEventStore(eventStore EventStore) *TraceableEventStore {
	return &TraceableEventStore{
		EventStore: eventStore,
		tracing:    false,
	}
}
