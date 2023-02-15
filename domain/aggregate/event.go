package aggregate

import (
	"encoding/json"
	"fmt"
	"time"
)

var ErrEventNotRegistered = fmt.Errorf("event not registered")

type DomainEvent interface {
	GetOccurredAt() time.Time
	Kind() string
	MarshalJSON() ([]byte, error)
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

func (e MemoCreated) MarshalJSON() ([]byte, error) {
	//occurredAt, _ := json.Marshal(e.occurredAt)
	return json.Marshal(&struct {
		Id         string `json:"id"`
		Body       string `json:"body"`
		OccurredAt string `json:"occurred_at"`
	}{
		Id:         e.Id.Val,
		Body:       e.Body,
		OccurredAt: e.occurredAt.Format("2006-01-02\\T15:04:05.000000Z07:00"), //"Y-m-d\\TH:i:s.uP"
		//OccurredAt: string(occurredAt),
	})
}

type MemoBodyUpdated struct {
	Id   UUIDv4
	Body string
	BasicEvent
}

func (e MemoBodyUpdated) Kind() string {
	return "MemoBodyUpdated"
}

func (e MemoBodyUpdated) MarshalJSON() ([]byte, error) {
	//occurredAt, _ := json.Marshal(e.occurredAt)
	return json.Marshal(&struct {
		Id         string `json:"id"`
		Body       string `json:"body"`
		OccurredAt string `json:"occurred_at"`
	}{
		Id:         e.Id.Val,
		Body:       e.Body,
		OccurredAt: e.occurredAt.Format("2006-01-02\\T15:04:05.000000Z07:00"), //"Y-m-d\\TH:i:s.uP"
		//OccurredAt: string(occurredAt),
	})
}
