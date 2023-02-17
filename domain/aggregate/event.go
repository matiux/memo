package aggregate

import (
	"encoding/json"
	"fmt"
	"time"
)

var ErrEventNotRegistered = fmt.Errorf("event not registered")
var EventDateFormat = "2006-01-02\\T15:04:05.000000Z07:00"

// DomainEvent -------------------------------------
type DomainEvent interface {
	GetOccurredAt() time.Time
	Kind() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

// BasicEvent -------------------------------------
type BasicEvent struct {
	occurredAt time.Time
}

func (e *BasicEvent) GetOccurredAt() time.Time {
	return e.occurredAt
}

// MemoCreated -------------------------------------
type MemoCreated struct {
	Id   UUIDv4
	Body string
	BasicEvent
}

func (e *MemoCreated) Kind() string {
	return "MemoCreated"
}

func (e *MemoCreated) MarshalJSON() ([]byte, error) {
	//occurredAt, _ := json.Marshal(e.occurredAt)
	return json.Marshal(&struct {
		Id         string `json:"id"`
		Body       string `json:"body"`
		OccurredAt string `json:"occurred_at"`
	}{
		Id:         e.Id.Val,
		Body:       e.Body,
		OccurredAt: e.occurredAt.Format(EventDateFormat), //"Y-m-d\\TH:i:s.uP"
		//OccurredAt: string(occurredAt),
	})
}

func (e *MemoCreated) UnmarshalJSON(b []byte) error {
	var aux struct {
		Id         string `json:"id"`
		Body       string `json:"body"`
		OccurredAt string `json:"occurred_at"`
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	e.Id = NewUUIDv4From(aux.Id)
	e.Body = aux.Body
	t, err := time.Parse(EventDateFormat, aux.OccurredAt)
	if err != nil {
		return err
	}
	e.occurredAt = t

	return nil
}

func NewMemoCreated(id UUIDv4, body string, occurredAt time.Time) *MemoCreated {
	return &MemoCreated{id, body, BasicEvent{occurredAt}}
}

// MemoBodyUpdated -------------------------------------
type MemoBodyUpdated struct {
	Id   UUIDv4
	Body string
	BasicEvent
}

func (e *MemoBodyUpdated) Kind() string {
	return "MemoBodyUpdated"
}

func (e *MemoBodyUpdated) MarshalJSON() ([]byte, error) {
	//occurredAt, _ := json.Marshal(e.occurredAt)
	return json.Marshal(&struct {
		Id         string `json:"id"`
		Body       string `json:"body"`
		OccurredAt string `json:"occurred_at"`
	}{
		Id:         e.Id.Val,
		Body:       e.Body,
		OccurredAt: e.occurredAt.Format(EventDateFormat), //"Y-m-d\\TH:i:s.uP"
		//OccurredAt: string(occurredAt),
	})
}

func (e *MemoBodyUpdated) UnmarshalJSON(b []byte) error {
	var aux struct {
		Id         string `json:"id"`
		Body       string `json:"body"`
		OccurredAt string `json:"occurred_at"`
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	e.Id = NewUUIDv4From(aux.Id)
	e.Body = aux.Body
	t, err := time.Parse(EventDateFormat, aux.OccurredAt)
	if err != nil {
		return err
	}
	e.occurredAt = t

	return nil
}

func NewMemoBodyUpdated(id UUIDv4, body string, updatedAd time.Time) *MemoBodyUpdated {
	return &MemoBodyUpdated{id, body, BasicEvent{updatedAd}}
}

// EventDeserializerRegistry is a registry to deserialize json events
func EventDeserializerRegistry(eventType, payload string) (*DomainEvent, error) {
	switch eventType {
	case "MemoCreated":
		var memoCreated MemoCreated
		_ = json.Unmarshal([]byte(payload), &memoCreated)
		event := DomainEvent(&memoCreated)
		return &event, nil
	case "MemoBodyUpdated":
		var memoBodyUpdated MemoBodyUpdated
		_ = json.Unmarshal([]byte(payload), &memoBodyUpdated)
		event := DomainEvent(&memoBodyUpdated)
		return &event, nil
	}

	return nil, fmt.Errorf("invalid event type %v", eventType)
}
