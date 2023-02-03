package aggregate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
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
		aggregateClass:   reflect.TypeOf(&Memo{}),
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

	assert.Len(t, eventStore.GetEvents(), 1)
	assert.Len(t, traceableEventBus.GetEvents(), 1)

	fmt.Println(err)
}
