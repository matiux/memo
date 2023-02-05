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
	Id           UUIDv4
	Body         string
	CreationDate time.Time
}

func (m *Memo) GetAggregateRootId() EntityId {
	return m.Id
}

func (m *Memo) Apply(event DomainEvent) (err error) {

	switch t := event.(type) {
	case MemoCreated:
		m.Id = t.Id
		m.Body = t.Body
		m.CreationDate = t.GetOccurredAt()
	case MemoBodyUpdated:
		m.Body = t.body
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

func (m *Memo) UpdateBody(body string, updatedAd time.Time) {
	event := NewMemoBodyUpdated(m.Id, body, updatedAd)
	if err := m.Record(event, m); err != nil {
		panic(err)
	}
}

func NewMemo(id UUIDv4, body string, creationDate time.Time) *Memo {
	memo := &Memo{}
	memo.create(id, body, creationDate)
	return memo
}
