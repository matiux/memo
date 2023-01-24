package aggregate

import (
	"fmt"
	"time"
)

var ErrEventNotRegistered = fmt.Errorf("event not registered")

type DomainEvent interface {
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
