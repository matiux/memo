package domain_test

import (
	"github.com/matiux/memo/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventSourcedAggregateRoot_it_applies_using_an_incrementing_playhead(t *testing.T) {

	updateTime := time.Now()

	memo := createMemo()
	memo.UpdateBody("Vegetables and fruits are good", updateTime)

	eventStream := memo.GetUncommittedEvents()

	for i := 1; i < len(eventStream); i++ {
		assert.Equal(t, domain.Playhead(i), eventStream[i-1].Playhead)
	}

	assert.Len(t, eventStream, 2)
}

func TestEventSourcedAggregateRoot_it_sets_internal_playhead_when_initializing(t *testing.T) {

	memoCreatedDomainMessage, _ := createEvents()

	memo := &domain.Memo{}
	_ = memo.InitializeState(
		domain.DomainEventStream{
			memoCreatedDomainMessage,
		},
		memo,
	)

	_ = memo.Record(memoCreatedDomainMessage.Payload, memo)

	eventStream := memo.GetUncommittedEvents()

	assert.Len(t, eventStream, 1)
	assert.Equal(t, domain.Playhead(2), eventStream[0].Playhead)

}
