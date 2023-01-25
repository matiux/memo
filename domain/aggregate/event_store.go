package aggregate

import "fmt"

var EventStreamNotFound = fmt.Errorf("event stream not found")

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

	for _, domainMessage := range eventStream {
		e.stream[stringId][domainMessage.Playhead] = domainMessage
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
