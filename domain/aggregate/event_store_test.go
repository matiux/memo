package aggregate

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestInMemoryEventStore_Append(t *testing.T) {
	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[Playhead]DomainMessage),
	}

	id := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memoCreatedDomainMessage := DomainMessage{
		Playhead:    Playhead(1),
		EventType:   "MemoCreated",
		Event:       NewMemoCreated(id, body, creationDate),
		AggregateId: id,
		RecordedOn:  time.Now(),
	}

	newBody := "Vegetables and fruits are good"
	updatingDate := time.Now()

	memoBodyUpdatedDomainMessage := DomainMessage{
		Playhead:    Playhead(2),
		EventType:   "MemoBodyUpdated",
		Event:       NewMemoBodyUpdated(id, newBody, updatingDate),
		AggregateId: id,
		RecordedOn:  time.Now(),
	}

	eventStream := DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	}

	eventStore.Append(id, eventStream)

	assert.Len(t, eventStore.stream, 1)
	assert.Contains(t, eventStore.stream, id.Val)
	assert.Len(t, eventStore.stream[id.Val], 2)
	assert.Contains(t, eventStore.stream[id.Val], Playhead(1))
	assert.Contains(t, eventStore.stream[id.Val], Playhead(2))
	assert.True(t, reflect.DeepEqual(memoCreatedDomainMessage, eventStore.stream[id.Val][Playhead(1)]))
	assert.True(t, reflect.DeepEqual(memoBodyUpdatedDomainMessage, eventStore.stream[id.Val][Playhead(2)]))
}

func TestInMemoryEventStore_Load(t *testing.T) {

	eventStore := &InMemoryEventStore{
		stream: make(map[string]map[Playhead]DomainMessage),
	}

	id := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memoCreatedDomainMessage := DomainMessage{
		Playhead:    Playhead(1),
		EventType:   "MemoCreated",
		Event:       NewMemoCreated(id, body, creationDate),
		AggregateId: id,
		RecordedOn:  time.Now(),
	}

	newBody := "Vegetables and fruits are good"
	updatingDate := time.Now()

	memoBodyUpdatedDomainMessage := DomainMessage{
		Playhead:    Playhead(2),
		EventType:   "MemoBodyUpdated",
		Event:       NewMemoBodyUpdated(id, newBody, updatingDate),
		AggregateId: id,
		RecordedOn:  time.Now(),
	}

	eventStore.Append(id, DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	})

	domainEventStream, _ := eventStore.Load(id)

	assert.Len(t, domainEventStream, 2)
	assert.True(t, domainEventStream[0].AggregateId.Equals(id))
	assert.True(t, domainEventStream[1].AggregateId.Equals(id))
}
