package aggregate

import "github.com/google/uuid"

type EntityId interface {
	Equals(otherEntityId EntityId) bool
}

type UUIDv4 struct {
	val string
}

func NewUUIDv4() UUIDv4 {
	return UUIDv4{val: uuid.NewString()}
}

func (id UUIDv4) Equals(otherEntityId EntityId) bool {

	otherUUIDv4, ok := otherEntityId.(UUIDv4)
	return ok && otherUUIDv4.val == id.val
}
