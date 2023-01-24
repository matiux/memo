package aggregate

import (
	"time"
)

type Event interface {
	GetOccurredAt() time.Time
}

type BasicEvent struct {
	occurredAt time.Time
}

func (e BasicEvent) GetOccurredAt() time.Time {
	return e.occurredAt
}

type MemoCreated struct {
	id   UUIDv4
	body string
	BasicEvent
}

func NewMemoCreated(id UUIDv4, body string, occurredAt time.Time) MemoCreated {
	return MemoCreated{id, body, BasicEvent{occurredAt}}
}
