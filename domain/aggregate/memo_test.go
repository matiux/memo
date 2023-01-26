package aggregate

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_it_should_be_create_new_memo(t *testing.T) {

	memoId := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memo := NewMemo(memoId, body, creationDate)

	assert.True(t, memo.id.Equals(memoId))
	assert.Equal(t, body, memo.body)
	assert.Equal(t, creationDate, memo.creationDate)

	assert.Len(t, memo.UncommittedEvents, 1)

	domainMessage := memo.UncommittedEvents[0]
	memoCreated := domainMessage.Event.(MemoCreated)

	assert.IsType(t, MemoCreated{}, memoCreated)
	assert.Equal(t, Playhead(1), domainMessage.Playhead)
	assert.Equal(t, Playhead(1), memo.Playhead)

	assert.Equal(t, creationDate, memoCreated.GetOccurredAt())
	assert.True(t, memoCreated.id.Equals(memoId))
	assert.Equal(t, body, memoCreated.body)
}

func Test_it_should_be_update_memo(t *testing.T) {

	memoId := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	newBody := "Vegetables and fruits are good"
	updatingDate := time.Now()

	memo := NewMemo(memoId, body, creationDate)
	memo.updateBody(newBody, updatingDate)

	assert.True(t, memo.id.Equals(memoId))
	assert.Equal(t, newBody, memo.body)
	assert.Equal(t, creationDate, memo.creationDate)

	assert.Len(t, memo.UncommittedEvents, 2)

	memoCreated := memo.UncommittedEvents[0].Event.(MemoCreated)
	memoBodyUpdated := memo.UncommittedEvents[1].Event.(MemoBodyUpdated)

	assert.IsType(t, MemoCreated{}, memoCreated)
	assert.IsType(t, MemoBodyUpdated{}, memoBodyUpdated)
	assert.Equal(t, Playhead(1), memo.UncommittedEvents[0].Playhead)
	assert.Equal(t, Playhead(2), memo.UncommittedEvents[1].Playhead)

	assert.Equal(t, creationDate, memoCreated.GetOccurredAt())
	assert.Equal(t, updatingDate, memoBodyUpdated.GetOccurredAt())
}
