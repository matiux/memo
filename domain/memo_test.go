package domain_test

import (
	"github.com/matiux/memo/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemo_it_should_be_create_new_memo(t *testing.T) {

	memo := createMemo()

	assert.True(t, memo.Id.Equals(memoId))
	assert.Equal(t, body, memo.Body)
	assert.Equal(t, creationDate, memo.CreationDate)

	assert.Len(t, memo.GetUncommittedEvents(), 1)

	domainMessage := memo.GetUncommittedEvents()[0]
	memoCreated := domainMessage.Payload.(*domain.MemoCreated)

	assert.IsType(t, &domain.MemoCreated{}, memoCreated)
	assert.Equal(t, domain.Playhead(1), domainMessage.Playhead)
	assert.Equal(t, domain.Playhead(1), memo.Playhead)

	assert.Equal(t, creationDate, memoCreated.GetOccurredAt())
	assert.True(t, memoCreated.Id.Equals(memoId))
	assert.Equal(t, body, memoCreated.Body)
}

func TestMemo_it_should_be_update_memo(t *testing.T) {

	newBody := "Vegetables and fruits are good"
	updatingDate := time.Now()

	memo := createMemo()
	memo.UpdateBody(newBody, updatingDate)

	assert.True(t, memo.Id.Equals(memoId))
	assert.Equal(t, newBody, memo.Body)
	assert.Equal(t, creationDate, memo.CreationDate)

	assert.Len(t, memo.GetUncommittedEvents(), 2)

	memoCreated := (memo.GetUncommittedEvents()[0].Payload).(*domain.MemoCreated)
	memoBodyUpdated := (memo.GetUncommittedEvents()[1].Payload).(*domain.MemoBodyUpdated)

	assert.IsType(t, &domain.MemoCreated{}, memoCreated)
	assert.IsType(t, &domain.MemoBodyUpdated{}, memoBodyUpdated)
	assert.Equal(t, domain.Playhead(1), memo.GetUncommittedEvents()[0].Playhead)
	assert.Equal(t, domain.Playhead(2), memo.GetUncommittedEvents()[1].Playhead)

	assert.Equal(t, creationDate, memoCreated.GetOccurredAt())
	assert.Equal(t, updatingDate, memoBodyUpdated.GetOccurredAt())
}
