package aggregate

import (
	"fmt"
	"time"
)

var ErrEventNotRegistered = fmt.Errorf("event not registered")

type DomainEvent interface {
	GetOccurredAt() time.Time
	Kind() string
}

type BasicEvent struct {
	occurredAt time.Time
}

func (e BasicEvent) GetOccurredAt() time.Time {
	return e.occurredAt
}

type MemoCreated struct {
	Id   UUIDv4
	Body string
	BasicEvent
}

func (e MemoCreated) Kind() string {
	return "MemoCreated"
}

type MemoBodyUpdated struct {
	id   UUIDv4
	body string
	BasicEvent
}

func (e MemoBodyUpdated) Kind() string {
	return "MemoBodyUpdated"
}
