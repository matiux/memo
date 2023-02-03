package aggregate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

var aggregateId UUIDv4
var memoCreatedDomainMessage DomainMessage
var memoBodyUpdatedDomainMessage DomainMessage

func setupTestEventStore() {
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

	memoBodyUpdatedDomainMessage = DomainMessage{
		Playhead:    Playhead(2),
		EventType:   "MemoBodyUpdated",
		Payload:     NewMemoBodyUpdated(aggregateId, "Vegetables and fruits are good", time.Now()),
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}

}

func TestEventStore_Append(t *testing.T) {

	setupTestEventStore()

	eventStore := NewInMemoryEventStore()

	eventStream := DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	}

	eventStore.Append(aggregateId, eventStream)

	assert.Len(t, eventStore.stream, 1)
	assert.Contains(t, eventStore.stream, aggregateId.Val)
	assert.Len(t, eventStore.stream[aggregateId.Val], 2)
	assert.Contains(t, eventStore.stream[aggregateId.Val], Playhead(1))
	assert.Contains(t, eventStore.stream[aggregateId.Val], Playhead(2))
	assert.True(t, reflect.DeepEqual(memoCreatedDomainMessage, eventStore.stream[aggregateId.Val][Playhead(1)]))
	assert.True(t, reflect.DeepEqual(memoBodyUpdatedDomainMessage, eventStore.stream[aggregateId.Val][Playhead(2)]))
}

func TestEventStore_Load(t *testing.T) {

	setupTestEventStore()

	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[Playhead]DomainMessage),
	}

	eventStore.Append(aggregateId, DomainEventStream{
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
		stream: make(map[string]map[Playhead]DomainMessage),
	}

	eventStore.Append(aggregateId, DomainEventStream{
		memoCreatedDomainMessage,
	})

	memoBodyUpdatedDomainMessage.Playhead = Playhead(1)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, fmt.Sprintf("%v", r), "duplicate playhead not allowed")
		} else {
			t.Error("??")
		}
	}()

	eventStore.Append(aggregateId, DomainEventStream{
		memoBodyUpdatedDomainMessage,
	})
}
