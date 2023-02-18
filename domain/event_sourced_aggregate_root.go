package domain

import (
	"sync"
	"time"
)

type Playhead int
type EventStream []DomainMessage

type EventSourcedEntity interface {
	handleRecursively(event Event)
	registerAggregateRoot(aggregate Root)
}

// DomainMessage represents an important change in the domain.
type DomainMessage struct {
	Playhead
	EventType   string
	Payload     Event
	AggregateId EntityId
	RecordedOn  time.Time
}

// Root represents an AggregateRoot
type Root interface {
	GetAggregateRootId() EntityId
	Apply(event Event) (err error)
	GetUncommittedEvents() []DomainMessage
	InitializeState(stream EventStream, aggregate Root) error
}

// EventSourcedAggregateRoot is the basic struct for an AggregateRoot
type EventSourcedAggregateRoot struct {
	uncommittedEvents []DomainMessage
	Playhead
	mutex sync.Mutex
}

func (e *EventSourcedAggregateRoot) Record(event Event, aggregate Root) error {

	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.handleRecursively(event, aggregate); err != nil {
		return err
	}

	e.Playhead++
	e.uncommittedEvents = append(
		e.uncommittedEvents,
		DomainMessage{
			Playhead: e.Playhead,
			//EventType:   reflect.ValueOf(event).Kind().String(),
			EventType:   event.Kind(),
			Payload:     event,
			AggregateId: aggregate.GetAggregateRootId(),
			RecordedOn:  event.GetOccurredAt(),
		},
	)

	return nil
}

func (e *EventSourcedAggregateRoot) InitializeState(stream EventStream, aggregate Root) error {

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, message := range stream {
		e.Playhead++
		if err := e.handleRecursively(message.Payload, aggregate); err != nil {
			return err
		}
	}

	return nil
}

func (e *EventSourcedAggregateRoot) handleRecursively(event Event, aggregate Root) error {

	if err := e.handle(event, aggregate); err != nil {
		return err
	}

	for _, entity := range e.getChildEntities() {
		entity.registerAggregateRoot(aggregate)
		entity.handleRecursively(event)
	}

	return nil
}

func (e *EventSourcedAggregateRoot) handle(event Event, aggregate Root) error {
	if err := aggregate.Apply(event); err != nil {
		return err
	}

	return nil
}

func (e *EventSourcedAggregateRoot) getChildEntities() []EventSourcedEntity {

	return []EventSourcedEntity{}
}

func (e *EventSourcedAggregateRoot) GetUncommittedEvents() []DomainMessage {
	return e.uncommittedEvents
}
