package aggregate_test

import (
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

var memoCreatedDomainMessage aggregate.DomainMessage
var memoBodyUpdatedDomainMessage aggregate.DomainMessage

func setupTestEventStore() {

	memoCreatedDomainMessage = aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(1),
		EventType:   "MemoCreated",
		Payload:     aggregate.NewMemoCreated(memoId, body, creationDate),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	memoBodyUpdatedDomainMessage = aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(2),
		EventType:   "MemoBodyUpdated",
		Payload:     aggregate.NewMemoBodyUpdated(memoId, "Vegetables and fruits are good", time.Now()),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

}

func TestEventStore_Append(t *testing.T) {

	setupTestEventStore()

	eventStore := aggregate.NewInMemoryEventStore()

	eventStream := aggregate.DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	}

	eventStore.Append(memoId, eventStream)

	assert.Len(t, eventStore.Stream, 1)
	assert.Contains(t, eventStore.Stream, memoId.Val)
	assert.Len(t, eventStore.Stream[memoId.Val], 2)
	assert.Contains(t, eventStore.Stream[memoId.Val], aggregate.Playhead(1))
	assert.Contains(t, eventStore.Stream[memoId.Val], aggregate.Playhead(2))
	assert.True(t, reflect.DeepEqual(memoCreatedDomainMessage, eventStore.Stream[memoId.Val][aggregate.Playhead(1)]))
	assert.True(t, reflect.DeepEqual(memoBodyUpdatedDomainMessage, eventStore.Stream[memoId.Val][aggregate.Playhead(2)]))
}

func TestEventStore_Load(t *testing.T) {

	setupTestEventStore()

	eventStore := &aggregate.InMemoryEventStore{
		Stream: make(map[string]map[aggregate.Playhead]aggregate.DomainMessage),
	}

	eventStore.Append(memoId, aggregate.DomainEventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	})

	domainEventStream, _ := eventStore.Load(memoId)

	assert.Len(t, domainEventStream, 2)
	assert.True(t, domainEventStream[0].AggregateId.Equals(memoId))
	assert.True(t, domainEventStream[1].AggregateId.Equals(memoId))
}

func TestEventStore_DuplicatedPlayhead(t *testing.T) {

	setupTestEventStore()

	eventStore := &aggregate.InMemoryEventStore{
		Stream: make(map[string]map[aggregate.Playhead]aggregate.DomainMessage),
	}

	eventStore.Append(memoId, aggregate.DomainEventStream{
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

	eventStore.Append(memoId, aggregate.DomainEventStream{
		memoBodyUpdatedDomainMessage,
	})
}
