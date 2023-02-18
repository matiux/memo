package domain_test

import (
	"fmt"
	"github.com/matiux/memo/domain"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEventStore_Append(t *testing.T) {

	memoCreatedDomainMessage, memoBodyUpdatedDomainMessage := createEvents()

	eventStore := domain.NewInMemoryEventStore()

	eventStream := domain.EventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	}

	eventStore.Append(memoId, eventStream)

	assert.Len(t, eventStore.Stream, 1)
	assert.Contains(t, eventStore.Stream, memoId.Val)
	assert.Len(t, eventStore.Stream[memoId.Val], 2)
	assert.Contains(t, eventStore.Stream[memoId.Val], domain.Playhead(1))
	assert.Contains(t, eventStore.Stream[memoId.Val], domain.Playhead(2))
	assert.True(t, reflect.DeepEqual(memoCreatedDomainMessage, eventStore.Stream[memoId.Val][domain.Playhead(1)]))
	assert.True(t, reflect.DeepEqual(memoBodyUpdatedDomainMessage, eventStore.Stream[memoId.Val][domain.Playhead(2)]))
}

func TestEventStore_Load(t *testing.T) {

	memoCreatedDomainMessage, memoBodyUpdatedDomainMessage := createEvents()

	eventStore := &domain.InMemoryEventStore{
		Stream: make(map[string]map[domain.Playhead]domain.Message),
	}

	eventStore.Append(memoId, domain.EventStream{
		memoCreatedDomainMessage,
		memoBodyUpdatedDomainMessage,
	})

	domainEventStream, _ := eventStore.Load(memoId)

	assert.Len(t, domainEventStream, 2)
	assert.True(t, domainEventStream[0].AggregateId.Equals(memoId))
	assert.True(t, domainEventStream[1].AggregateId.Equals(memoId))
}

func TestEventStore_DuplicatedPlayhead(t *testing.T) {

	memoCreatedDomainMessage, memoBodyUpdatedDomainMessage := createEvents()

	eventStore := &domain.InMemoryEventStore{
		Stream: make(map[string]map[domain.Playhead]domain.Message),
	}

	eventStore.Append(memoId, domain.EventStream{
		memoCreatedDomainMessage,
	})

	memoBodyUpdatedDomainMessage.Playhead = domain.Playhead(1)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, fmt.Sprintf("%v", r), "duplicate playhead not allowed")
		} else {
			t.Error("??")
		}
	}()

	eventStore.Append(memoId, domain.EventStream{
		memoBodyUpdatedDomainMessage,
	})
}
