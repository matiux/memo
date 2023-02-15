package aggregate_test

import (
	"encoding/json"
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemoCreated_it_should_marshal_memo_created_event(t *testing.T) {

	id := aggregate.NewUUIDv4From("1750c0c3-06b2-46cf-b140-b36cdc215474")
	occurredAt, _ := time.Parse("2006-01-02\\T15:04:05.000000Z07:00", "2023-02-15\\T10:19:52.642901+01:00")
	memoCreate := aggregate.NewMemoCreated(id, body, occurredAt)
	marshaledMemoCreated, err := json.Marshal(memoCreate)
	if err != nil {
		fmt.Println(err)
		return
	}

	expected := "{\"id\":\"1750c0c3-06b2-46cf-b140-b36cdc215474\",\"body\":\"Vegetables are good\",\"occurred_at\":\"2023-02-15\\\\T10:19:52.642901+01:00\"}"
	assert.Equal(t, expected, string(marshaledMemoCreated))

}

func TestMemoCreated_it_should_marshal_memo_updated_event(t *testing.T) {

	id := aggregate.NewUUIDv4From("1750c0c3-06b2-46cf-b140-b36cdc215474")
	occurredAt, _ := time.Parse("2006-01-02\\T15:04:05.000000Z07:00", "2023-02-15\\T10:19:52.642901+01:00")
	memoBodyUpdated := aggregate.NewMemoBodyUpdated(id, body, occurredAt)
	marshaledMemoBodyUpdated, err := json.Marshal(memoBodyUpdated)
	if err != nil {
		fmt.Println(err)
		return
	}

	expected := "{\"id\":\"1750c0c3-06b2-46cf-b140-b36cdc215474\",\"body\":\"Vegetables are good\",\"occurred_at\":\"2023-02-15\\\\T10:19:52.642901+01:00\"}"
	assert.Equal(t, expected, string(marshaledMemoBodyUpdated))
}
