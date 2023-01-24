package aggregate

import "github.com/google/uuid"

type EntityId interface {
	Equals(otherEntityId EntityId) bool
}

type UUIDv4 struct {
	Val string
}

func NewUUIDv4() UUIDv4 {
	return UUIDv4{Val: uuid.NewString()}
}

func NewUUIDv4From(id string) UUIDv4 {
	return UUIDv4{Val: id}
}

func (id UUIDv4) Equals(otherEntityId EntityId) bool {

	otherUUIDv4, ok := otherEntityId.(UUIDv4)
	return ok && otherUUIDv4.Val == id.Val
}
