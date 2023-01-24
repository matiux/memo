package aggregate

import (
	"time"
)

type Memo struct {
	EventSourcedAggregateRoot
	id           UUIDv4
	body         string
	creationDate time.Time
}

func (m *Memo) getAggregateRootId() string {
	return m.id.val
}

func (m *Memo) create(id UUIDv4, body string, creationDate time.Time) {
	event := NewMemoCreated(id, body, creationDate)
	m.apply(event, m)
}

func (m *Memo) ApplyMemoCreated(event MemoCreated) {
	m.id = event.id
	m.body = event.body
	m.creationDate = event.occurredAt
}

func NewMemo(id UUIDv4, body string, creationDate time.Time) Memo {
	memo := Memo{}
	memo.create(id, body, creationDate)
	return memo
}
