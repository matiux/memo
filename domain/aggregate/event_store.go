package aggregate

import (
	"fmt"
)

var EventStreamNotFound = fmt.Errorf("event stream not found")
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
		return nil, EventStreamNotFound
	}

	domainEventStream := DomainEventStream{}

	for _, domainMessage := range e.stream[stringId] {
		domainEventStream = append(domainEventStream, domainMessage)
	}

	return domainEventStream, nil
}
