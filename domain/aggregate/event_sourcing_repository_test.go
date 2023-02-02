package aggregate

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var eventStore EventStore
var eventSourcingRepository EventSourcingRepository

func setupEventSourcingRepository() {

	eventStore = &InMemoryEventStore{}

	eventBus := &SimpleEventBus{
		eventListeners: nil,
		queue:          nil,
		isPublishing:   false,
	}

	eventSourcingRepository = EventSourcingRepository{
		EventStore:       eventStore,
		EventBus:         eventBus,
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

	fmt.Println(err)
}
