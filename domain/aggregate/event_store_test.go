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

func TestMain(m *testing.M) {
	aggregateId = NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memoCreatedDomainMessage = DomainMessage{
		Playhead:    Playhead(1),
		EventType:   "MemoCreated",
		Event:       NewMemoCreated(aggregateId, body, creationDate),
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}

	memoBodyUpdatedDomainMessage = DomainMessage{
		Playhead:    Playhead(2),
		EventType:   "MemoBodyUpdated",
		Event:       NewMemoBodyUpdated(aggregateId, "Vegetables and fruits are good", time.Now()),
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}

	m.Run()
}

func TestInMemoryEventStore_Append(t *testing.T) {

	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[Playhead]DomainMessage),
	}

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

func TestInMemoryEventStore_Load(t *testing.T) {

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

func TestInMemoryEventStore_DuplicatedPlayhead(t *testing.T) {
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
