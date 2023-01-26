package aggregate

import (
	"time"
)

func NewMemoCreated(id UUIDv4, body string, occurredAt time.Time) MemoCreated {
	return MemoCreated{id, body, BasicEvent{occurredAt}}
}

func NewMemoBodyUpdated(id UUIDv4, body string, updatedAd time.Time) MemoBodyUpdated {
	return MemoBodyUpdated{id, body, BasicEvent{updatedAd}}
}

type Memo struct {
	EventSourcedAggregateRoot
	id           UUIDv4
	body         string
	creationDate time.Time
}

func (m *Memo) getAggregateRootId() EntityId {
	return m.id
}

func (m *Memo) Apply(event DomainEvent) (err error) {

	switch t := event.(type) {
	case MemoCreated:
		m.id = t.id
		m.body = t.body
		m.creationDate = t.GetOccurredAt()
	case MemoBodyUpdated:
		m.body = t.body
	default:
		err = ErrEventNotRegistered
	}

	return
}

func (m *Memo) create(id UUIDv4, body string, creationDate time.Time) {
	event := NewMemoCreated(id, body, creationDate)
	if err := m.Record(event, m); err != nil {
		panic(err)
	}
}

func (m *Memo) updateBody(body string, updatedAd time.Time) {
	event := NewMemoBodyUpdated(m.id, body, updatedAd)
	if err := m.Record(event, m); err != nil {
		panic(err)
	}
}

func NewMemo(id UUIDv4, body string, creationDate time.Time) *Memo {
	memo := &Memo{}
	memo.create(id, body, creationDate)
	return memo
}
