package aggregate_test

import (
	"github.com/matiux/memo/domain/aggregate"
	"time"
)

var memoId = aggregate.NewUUIDv4()
var body = "Vegetables are good"
var creationDate = time.Now()

func createMemo() *aggregate.Memo {

	return aggregate.NewMemo(memoId, body, creationDate)
}

func createEvents() (aggregate.DomainMessage, aggregate.DomainMessage) {

	memoCreatedDomainMessage := aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(1),
		EventType:   "MemoCreated",
		Payload:     aggregate.NewMemoCreated(memoId, body, creationDate),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	memoBodyUpdatedDomainMessage := aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(2),
		EventType:   "MemoBodyUpdated",
		Payload:     aggregate.NewMemoBodyUpdated(memoId, "Vegetables and fruits are good", time.Now()),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	return memoCreatedDomainMessage, memoBodyUpdatedDomainMessage
}

func setupTestEventSourcingRepository() (
	*aggregate.TraceableEventStore,
	*aggregate.TraceableEventBus,
	aggregate.EventSourcingRepository,
) {

	eventStore := aggregate.NewTraceableEventStore(aggregate.NewInMemoryEventStore())
	eventStore.Trace()

	eventBus := aggregate.NewTraceableEventBus(
		&aggregate.SimpleEventBus{
			EventListeners: nil,
			Queue:          nil,
			IsPublishing:   false,
		},
	)
	eventBus.Trace()

	eventSourcingRepository := aggregate.NewEventSourcingRepository(
		eventStore,
		eventBus,
		&aggregate.PublicConstructorAggregateFactory{},
	)

	return eventStore, eventBus, eventSourcingRepository
}
