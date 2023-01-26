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

	if err := aggregate.Apply(event); err != nil {
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
