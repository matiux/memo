package aggregate_test

import (
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createMemosRepository() aggregate.Memos {
	eventStore, eventBus, _ := setupInMemoryEventSourcingRepository()

	return aggregate.NewMemos(eventStore, eventBus)
}

func TestMemos_it_should_add_new_memo_in_memory(t *testing.T) {

	eventStore, eventBus, _ := setupInMemoryEventSourcingRepository()
	memo := createMemo()

	memos := aggregate.NewMemos(eventStore, eventBus)
	_ = memos.Add(memo)

	byIdMemo, _ := memos.ById(memoId)

	assert.True(t, byIdMemo.Id.Equals(memoId))
}

func TestMemos_it_should_update_existing_memo_in_memory(t *testing.T) {

	memos := createMemosRepository()

	memo := createMemo()
	_ = memos.Add(memo)

	updateTime := time.Now()
	toUpdateMemo, _ := memos.ById(memoId)
	toUpdateMemo.UpdateBody("Vegetables and fruits are good", updateTime)

	_ = memos.Update(toUpdateMemo)

	updatedMemo, _ := memos.ById(memoId)
	assert.Equal(t, aggregate.Playhead(2), updatedMemo.Playhead)
	assert.Equal(t, "Vegetables and fruits are good", updatedMemo.Body)
}

func TestMemos_it_should_load_memo_by_repository(t *testing.T) {

	eventStore, eventBus, _ := setupMySqlEventSourcingRepository()
	memo := createMemo()

	memos := aggregate.NewMemos(eventStore, eventBus)
	_ = memos.Add(memo)

	byIdMemo, _ := memos.ById(memoId)

	assert.True(t, byIdMemo.Id.Equals(memoId))
	assert.Equal(t, aggregate.Playhead(1), byIdMemo.Playhead)
	assert.Equal(t, "Vegetables are good", byIdMemo.Body)
}
