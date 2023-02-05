package aggregate_test

import (
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemos_it_should_add_new_memo(t *testing.T) {

	eventStore, eventBus, _ := setupTestEventSourcingRepository()
	memo := createMemo()

	memos := aggregate.NewMemos(eventStore, eventBus)
	_ = memos.Add(memo)

	byIdMemo, _ := memos.ById(memoId)

	assert.True(t, byIdMemo.Id.Equals(memoId))
}
