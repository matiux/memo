package aggregate

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventSourcedAggregateRoot_it_applies_using_an_incrementing_playhead(t *testing.T) {

	memoId := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()
	updateTime := time.Now()

	memo := NewMemo(memoId, body, creationDate)
	memo.updateBody("Vegetables and fruits are good", updateTime)
	eventStream := memo.GetUncommittedEvents()

	for i := 1; i < len(eventStream); i++ {
		assert.Equal(t, Playhead(i), eventStream[i-1].Playhead)
	}

	assert.Len(t, eventStream, 2)
}

func TestEventSourcedAggregateRoot_it_sets_internal_playhead_when_initializing(t *testing.T) {
	// TODO
}
