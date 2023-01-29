package eventstore

import (
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
)

var EventStreamNotFound = fmt.Errorf("event stream not found")
var DuplicatePlayhead = fmt.Errorf("duplicate playhead not allowed")

type EventStore interface {
	Append(id aggregate.EntityId, eventStream aggregate.DomainEventStream)

	Load(id aggregate.EntityId) (aggregate.DomainEventStream, error)

	//LoadFromPlayhead(id EntityId, playhead Playhead) DomainEventStream
}

type InMemoryEventStore struct {
	stream map[string]map[aggregate.Playhead]aggregate.DomainMessage
}

func (e *InMemoryEventStore) Append(id aggregate.EntityId, eventStream aggregate.DomainEventStream) {

	stringId := id.(aggregate.UUIDv4).Val

	if _, exists := e.stream[stringId]; !exists {
		e.stream[stringId] = make(map[aggregate.Playhead]aggregate.DomainMessage)
	}

	e.assertStream(e.stream[stringId], eventStream)

	for _, domainMessage := range eventStream {
		e.stream[stringId][domainMessage.Playhead] = domainMessage
	}
}

func (e *InMemoryEventStore) assertStream(events map[aggregate.Playhead]aggregate.DomainMessage, eventsToAppend aggregate.DomainEventStream) {

	for _, event := range eventsToAppend {
		if _, exists := events[event.Playhead]; exists {
			panic(DuplicatePlayhead)
		}
	}
}

func (e *InMemoryEventStore) Load(id aggregate.EntityId) (aggregate.DomainEventStream, error) {

	stringId := id.(aggregate.UUIDv4).Val

	if _, exists := e.stream[stringId]; !exists {
		return nil, EventStreamNotFound
	}

	domainEventStream := aggregate.DomainEventStream{}

	for _, domainMessage := range e.stream[stringId] {
		domainEventStream = append(domainEventStream, domainMessage)
	}

	return domainEventStream, nil
}
