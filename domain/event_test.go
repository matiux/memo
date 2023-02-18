package domain_test

import (
	"encoding/json"
	"fmt"
	"github.com/matiux/memo/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemoCreated_it_should_marshal_memo_created_event(t *testing.T) {

	id := domain.NewUUIDv4From("1750c0c3-06b2-46cf-b140-b36cdc215474")
	occurredAt, _ := time.Parse(domain.EventDateFormat, "2023-02-15\\T10:19:52.642901+01:00")
	memoCreate := domain.NewMemoCreated(id, body, occurredAt)
	marshaledMemoCreated, err := json.Marshal(memoCreate)
	if err != nil {
		fmt.Println(err)
		return
	}

	expected := "{\"id\":\"1750c0c3-06b2-46cf-b140-b36cdc215474\",\"body\":\"Vegetables are good\",\"occurred_at\":\"2023-02-15\\\\T10:19:52.642901+01:00\"}"
	assert.Equal(t, expected, string(marshaledMemoCreated))

}

func TestMemoCreated_it_should_unmarshal_memo_created_event(t *testing.T) {

	jsonStr := `{
		"id": "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2",
		"body": "Vegetables are good",
		"occurred_at": "2023-02-16\\T10:30:22.695498+01:00"
	}`

	var memoCreated domain.MemoCreated
	err := json.Unmarshal([]byte(jsonStr), &memoCreated)

	assert.Nil(t, err)
	assert.Equal(t, "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2", memoCreated.Id.Val)
	assert.Equal(t, "Vegetables are good", memoCreated.Body)
	assert.Equal(t, time.Date(2023, time.February, 16, 10, 30, 22, 695498000, time.Local), memoCreated.GetOccurredAt())
}

func TestMemoBodyUpdated_it_should_marshal_memo_body_updated_event(t *testing.T) {

	id := domain.NewUUIDv4From("1750c0c3-06b2-46cf-b140-b36cdc215474")
	occurredAt, _ := time.Parse(domain.EventDateFormat, "2023-02-15\\T10:19:52.642901+01:00")
	memoBodyUpdated := domain.NewMemoBodyUpdated(id, body, occurredAt)
	marshaledMemoBodyUpdated, err := json.Marshal(memoBodyUpdated)
	if err != nil {
		fmt.Println(err)
		return
	}

	expected := "{\"id\":\"1750c0c3-06b2-46cf-b140-b36cdc215474\",\"body\":\"Vegetables are good\",\"occurred_at\":\"2023-02-15\\\\T10:19:52.642901+01:00\"}"
	assert.Equal(t, expected, string(marshaledMemoBodyUpdated))
}

func TestMemoBodyUpdated_it_should_unmarshal_memo_body_updated_event(t *testing.T) {

	jsonStr := `{
		"id": "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2",
		"body": "Vegetables are good",
		"occurred_at": "2023-02-16\\T10:30:22.695498+01:00"
	}`

	var memoCreated domain.MemoBodyUpdated
	err := json.Unmarshal([]byte(jsonStr), &memoCreated)

	assert.Nil(t, err)
	assert.Equal(t, "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2", memoCreated.Id.Val)
	assert.Equal(t, "Vegetables are good", memoCreated.Body)
	assert.Equal(t, time.Date(2023, time.February, 16, 10, 30, 22, 695498000, time.Local), memoCreated.GetOccurredAt())
}

func TestEventDeserializerRegistry_it_should_unmarshal_memo_created_event(t *testing.T) {

	jsonStr := `{
		"id": "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2",
		"body": "Vegetables are good",
		"occurred_at": "2023-02-16\\T10:30:22.695498+01:00"
	}`

	event, err := domain.EventDeserializerRegistry("MemoCreated", jsonStr)
	memoCreated := (*event).(*domain.MemoCreated)
	assert.Nil(t, err)
	assert.Equal(t, "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2", memoCreated.Id.Val)
	assert.Equal(t, "Vegetables are good", memoCreated.Body)
	assert.Equal(t, time.Date(2023, time.February, 16, 10, 30, 22, 695498000, time.Local), memoCreated.GetOccurredAt())
}

func TestEventDeserializerRegistry_it_should_unmarshal_memo_body_updated_event(t *testing.T) {

	jsonStr := `{
		"id": "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2",
		"body": "Vegetables are good",
		"occurred_at": "2023-02-16\\T10:30:22.695498+01:00"
	}`

	event, err := domain.EventDeserializerRegistry("MemoBodyUpdated", jsonStr)
	memoBodyUpdated := (*event).(*domain.MemoBodyUpdated)
	assert.Nil(t, err)
	assert.Equal(t, "ce567a4f-1d9e-4b15-bcf3-f78f7e0340b2", memoBodyUpdated.Id.Val)
	assert.Equal(t, "Vegetables are good", memoBodyUpdated.Body)
	assert.Equal(t, time.Date(2023, time.February, 16, 10, 30, 22, 695498000, time.Local), memoBodyUpdated.GetOccurredAt())
}
