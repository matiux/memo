package eventstore

import (
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

var aggregateId aggregate.UUIDv4
var memoCreatedDomainMessage aggregate.DomainMessage
var memoBodyUpdatedDomainMessage aggregate.DomainMessage

func setupTestEventStore() {
	aggregateId = aggregate.NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memoCreatedDomainMessage = aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(1),
		EventType:   "MemoCreated",
		Event:       aggregate.NewMemoCreated(aggregateId, body, creationDate),
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}

	memoBodyUpdatedDomainMessage = aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(2),
		EventType:   "MemoBodyUpdated",
		Event:       aggregate.NewMemoBodyUpdated(aggregateId, "Vegetables and fruits are good", time.Now()),
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}

}

func TestEventStore_Append(t *testing.T) {

	setupTestEventStore()

	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[aggregate.Playhead]aggregate.DomainMessage),
	}

	eventStream := aggregate.DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	}

	eventStore.Append(aggregateId, eventStream)

	assert.Len(t, eventStore.stream, 1)
	assert.Contains(t, eventStore.stream, aggregateId.Val)
	assert.Len(t, eventStore.stream[aggregateId.Val], 2)
	assert.Contains(t, eventStore.stream[aggregateId.Val], aggregate.Playhead(1))
	assert.Contains(t, eventStore.stream[aggregateId.Val], aggregate.Playhead(2))
	assert.True(t, reflect.DeepEqual(memoCreatedDomainMessage, eventStore.stream[aggregateId.Val][aggregate.Playhead(1)]))
	assert.True(t, reflect.DeepEqual(memoBodyUpdatedDomainMessage, eventStore.stream[aggregateId.Val][aggregate.Playhead(2)]))
}

func TestEventStore_Load(t *testing.T) {

	setupTestEventStore()

	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[aggregate.Playhead]aggregate.DomainMessage),
	}

	eventStore.Append(aggregateId, aggregate.DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	})

	domainEventStream, _ := eventStore.Load(aggregateId)

	assert.Len(t, domainEventStream, 2)
	assert.True(t, domainEventStream[0].AggregateId.Equals(aggregateId))
	assert.True(t, domainEventStream[1].AggregateId.Equals(aggregateId))
}

func TestEventStore_DuplicatedPlayhead(t *testing.T) {

	setupTestEventStore()

	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[aggregate.Playhead]aggregate.DomainMessage),
	}

	eventStore.Append(aggregateId, aggregate.DomainEventStream{
		memoCreatedDomainMessage,
	})

	memoBodyUpdatedDomainMessage.Playhead = aggregate.Playhead(1)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, fmt.Sprintf("%v", r), "duplicate playhead not allowed")
		} else {
			t.Error("??")
		}
	}()

	eventStore.Append(aggregateId, aggregate.DomainEventStream{
		memoBodyUpdatedDomainMessage,
	})
}
