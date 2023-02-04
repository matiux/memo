package aggregate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var eventStore *TraceableEventStore
var traceableEventBus *TraceableEventBus
var eventSourcingRepository EventSourcingRepository

func setupEventSourcingRepository() {

	eventStore = NewTraceableEventStore(NewInMemoryEventStore())
	eventStore.Trace()

	traceableEventBus = NewTraceableEventBus(
		&SimpleEventBus{
			eventListeners: nil,
			queue:          nil,
			isPublishing:   false,
		},
	)
	traceableEventBus.Trace()

	eventSourcingRepository = EventSourcingRepository{
		EventStore:       eventStore,
		EventBus:         traceableEventBus,
		aggregateClass:   &Memo{},
		AggregateFactory: &PublicConstructorAggregateFactory{},
	}
}

func TestEventSourcingRepository_it_adds_an_aggregate_root(t *testing.T) {

	setupEventSourcingRepository()

	memoId := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memo := NewMemo(memoId, body, creationDate)

	err := eventSourcingRepository.Save(memo)

	assert.Nil(t, err)
	assert.Len(t, eventStore.GetEvents(), 1)
	assert.Len(t, traceableEventBus.GetEvents(), 1)
}

func TestEventSourcingRepository_it_loads_an_aggregate(t *testing.T) {

	setupEventSourcingRepository()

	aggregateId = NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memoCreatedDomainMessage = DomainMessage{
		Playhead:    Playhead(1),
		EventType:   "MemoCreated",
		Payload:     NewMemoCreated(aggregateId, body, creationDate),
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}

	eventStream := DomainEventStream{
		memoCreatedDomainMessage,
	}

	eventStore.Append(aggregateId, eventStream)

	aggregate, err := eventSourcingRepository.Load(aggregateId)
	expectedMemo := NewMemo(aggregateId, body, creationDate)

	assert.Nil(t, err)

	var actualMemo = aggregate.(*Memo)

	assert.True(t, expectedMemo.id.Equals(aggregate.getAggregateRootId()))
	assert.Equal(t, expectedMemo.body, actualMemo.body)
	assert.Equal(t, expectedMemo.creationDate, actualMemo.creationDate)
}

func TestEventSourcingRepository_it_return_an_error_if_aggregate_was_not_found(t *testing.T) {

	setupEventSourcingRepository()

	aggregateId = NewUUIDv4()

	_, err := eventSourcingRepository.Load(aggregateId)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("aggregate with id '%v' not found", aggregateId))
}
