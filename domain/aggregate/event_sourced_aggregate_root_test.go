package aggregate_test

import (
	"github.com/matiux/memo/domain/aggregate"
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
		assert.Equal(t, aggregate.Playhead(i), eventStream[i-1].Playhead)
	}

	assert.Len(t, eventStream, 2)
}

func TestEventSourcedAggregateRoot_it_sets_internal_playhead_when_initializing(t *testing.T) {
	// TODO
}
