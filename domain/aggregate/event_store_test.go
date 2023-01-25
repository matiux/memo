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

	eventStream := DomainEventStream{
		memoCreatedDomainMessage,
	}

	eventStore.Append(id, eventStream)

	assert.Len(t, eventStore.stream, 1)
	assert.Contains(t, eventStore.stream, id.Val)
	assert.Len(t, eventStore.stream[id.Val], 1)
	assert.Contains(t, eventStore.stream[id.Val], Playhead(1))
	assert.True(t, reflect.DeepEqual(memoCreatedDomainMessage, eventStore.stream[id.Val][Playhead(1)]))
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

	eventStore.Append(id, DomainEventStream{
		memoCreatedDomainMessage,
	})

	domainEventStream, _ := eventStore.Load(id)

	assert.Len(t, domainEventStream, 1)
	assert.True(t, domainEventStream[0].AggregateId.Equals(id))
}
