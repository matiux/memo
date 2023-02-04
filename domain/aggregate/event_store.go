package aggregate

import (
	"fmt"
)

var DuplicatePlayhead = fmt.Errorf("duplicate playhead not allowed")

type EventStore interface {
	Append(id EntityId, eventStream DomainEventStream)
	Load(id EntityId) (DomainEventStream, error)
	//LoadFromPlayhead(id EntityId, playhead Playhead) DomainEventStream
}

type InMemoryEventStore struct {
	stream map[string]map[Playhead]DomainMessage
}

func (e *InMemoryEventStore) Append(id EntityId, eventStream DomainEventStream) {

	stringId := id.(UUIDv4).Val

	if _, exists := e.stream[stringId]; !exists {
		e.stream[stringId] = make(map[Playhead]DomainMessage)
	}

	e.assertStream(e.stream[stringId], eventStream)

	for _, domainMessage := range eventStream {
		e.stream[stringId][domainMessage.Playhead] = domainMessage
	}
}

func (e *InMemoryEventStore) assertStream(events map[Playhead]DomainMessage, eventsToAppend DomainEventStream) {

	for _, event := range eventsToAppend {
		if _, exists := events[event.Playhead]; exists {
			panic(DuplicatePlayhead)
		}
	}
}

func (e *InMemoryEventStore) Load(id EntityId) (DomainEventStream, error) {

	stringId := id.(UUIDv4).Val

	if _, exists := e.stream[stringId]; !exists {
		return nil, fmt.Errorf("aggregate with id '%v' not found", id)
	}

	domainEventStream := DomainEventStream{}

	for _, domainMessage := range e.stream[stringId] {
		domainEventStream = append(domainEventStream, domainMessage)
	}

	return domainEventStream, nil
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		stream: make(map[string]map[Playhead]DomainMessage),
	}
}

type TraceableEventStore struct {
	EventStore
	tracing  bool
	recorded DomainEventStream
}

func (e *TraceableEventStore) Append(id EntityId, eventStream DomainEventStream) {
	e.EventStore.Append(id, eventStream)

	if !e.tracing {
		return
	}

	for _, event := range eventStream {
		e.recorded = append(e.recorded, event)
	}
}

func (e *TraceableEventStore) GetEvents() (events []DomainEvent) {

	for _, event := range e.recorded {
		events = append(events, event.Payload)
	}

	return
}

func (e *TraceableEventStore) Trace() {
	e.tracing = true
}

func (e *TraceableEventStore) ClearEvents() {
	e.recorded = DomainEventStream{}
}

func NewTraceableEventStore(eventStore EventStore) *TraceableEventStore {
	return &TraceableEventStore{
		EventStore: eventStore,
		tracing:    false,
	}
}
