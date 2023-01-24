package variant

type EntityId[T comparable] interface {
	Create() EntityId[T]
	Value() T
	Equals(id EntityId[T]) bool
}

type BasicId[T comparable] struct {
	id T
}

func (e BasicId[T]) Value() T {
	return e.id
}

func (e BasicId[T]) Equals(id EntityId[T]) bool {
	return e.Value() == id.Value()
}

type UuidV4 struct {
	BasicId[string]
}

func (uuidV4 UuidV4) Create() EntityId[string] {
	return UuidV4{BasicId[string]{id: "myid"}}
}

func NewUuidV4ID() UuidV4 {
	return UuidV4{BasicId[string]{id: "myid"}}
}
