package aggregate

import (
	"sync"
	"time"
)

type Playhead int
type DomainEventStream []DomainMessage

// DomainMessage represents an important change in the domain.
type DomainMessage struct {
	Playhead
	EventType   string
	Event       DomainEvent
	AggregateId EntityId
	RecordedOn  time.Time
}

// Root represents an AggregateRoot
type Root interface {
	getAggregateRootId() EntityId
	Apply(event DomainEvent) (err error)
}

// EventSourcedAggregateRoot is the basic struct for an AggregateRoot
type EventSourcedAggregateRoot struct {
	UncommittedEvents []DomainMessage
	Playhead
	mutex sync.Mutex
}

func (e *EventSourcedAggregateRoot) Record(event DomainEvent, aggregate Root) error {

	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.handleRecursively(event, aggregate); err != nil {
		return err
	}

	e.Playhead++
	e.UncommittedEvents = append(
		e.UncommittedEvents,
		DomainMessage{
			Playhead: e.Playhead,
			//EventType:   reflect.ValueOf(event).Kind().String(),
			EventType:   event.Kind(),
			Event:       event,
			AggregateId: aggregate.getAggregateRootId(),
			RecordedOn:  event.GetOccurredAt(),
		},
	)

	return nil
}

func (e *EventSourcedAggregateRoot) InitializeState(stream DomainEventStream, aggregate Root) error {

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, message := range stream {
		e.Playhead++
		if err := e.handleRecursively(message.Event, aggregate); err != nil {
			return err
		}
	}

	return nil
}

func (e *EventSourcedAggregateRoot) handleRecursively(event DomainEvent, aggregate Root) error {

	if err := e.handle(event, aggregate); err != nil {
		return err
	}

	for _, entity := range e.getChildEntities() {
		entity.registerAggregateRoot(aggregate)
		entity.handleRecursively(event)
	}

	return nil
}

func (e *EventSourcedAggregateRoot) handle(event DomainEvent, aggregate Root) error {
	if err := aggregate.Apply(event); err != nil {
		return err
	}

	return nil
}

func (e *EventSourcedAggregateRoot) getChildEntities() []EventSourcedEntity {

	return []EventSourcedEntity{}
}
