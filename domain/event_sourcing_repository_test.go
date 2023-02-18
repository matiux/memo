package domain_test

import (
	"fmt"
	"github.com/matiux/memo/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventSourcingRepository_it_adds_an_aggregate_root(t *testing.T) {

	eventStore, eventBus, eventSourcingRepository := setupInMemoryEventSourcingRepository()
	memo := createMemo()

	err := eventSourcingRepository.Save(memo)

	assert.Nil(t, err)
	assert.Len(t, eventStore.GetEvents(), 1)
	assert.Len(t, eventBus.GetEvents(), 1)
}

func TestEventSourcingRepository_it_loads_an_aggregate(t *testing.T) {

	eventStore, _, eventSourcingRepository := setupInMemoryEventSourcingRepository()

	memoCreatedDomainMessage := domain.DomainMessage{
		Playhead:    domain.Playhead(1),
		EventType:   "MemoCreated",
		Payload:     domain.NewMemoCreated(memoId, body, creationDate),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	eventStream := domain.EventStream{
		memoCreatedDomainMessage,
	}

	eventStore.Append(memoId, eventStream)

	aggregate1, err := eventSourcingRepository.Load(memoId, &domain.Memo{})
	expectedMemo := domain.NewMemo(memoId, body, creationDate)

	assert.Nil(t, err)

	var actualMemo = aggregate1.(*domain.Memo)

	assert.True(t, expectedMemo.Id.Equals(actualMemo.GetAggregateRootId()))
	assert.Equal(t, expectedMemo.Body, actualMemo.Body)
	assert.Equal(t, expectedMemo.CreationDate, actualMemo.CreationDate)
}

func TestEventSourcingRepository_it_return_an_error_if_aggregate_was_not_found(t *testing.T) {

	_, _, eventSourcingRepository := setupInMemoryEventSourcingRepository()

	aggregateId := domain.NewUUIDv4()

	_, err := eventSourcingRepository.Load(aggregateId, &domain.Memo{})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("aggregate with id '%v' not found", aggregateId))
}
