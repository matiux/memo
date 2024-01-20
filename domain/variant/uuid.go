package variant

type Id[T comparable] interface {
	Create(id T) Id[T]
	Value() T
	Equals(id Id[T]) bool
}

type BasicId[T comparable] struct{ id T }

func (e BasicId[T]) Value() T { return e.id }

func (e BasicId[T]) Equals(id Id[T]) bool {
	return e.Value() == id.Value()
}

type UuidV4 struct{ BasicId[string] }

func (uuidV4 UuidV4) Create(id string) Id[string] {
	return UuidV4{BasicId[string]{id}}
}

func NewUuidV4(id string) UuidV4 {
	return UuidV4{BasicId[string]{id}}
}

type IntId struct{ BasicId[int] }

func (intId IntId) Create(id int) Id[int] {
	return IntId{BasicId[int]{id}}
}

func NewIntId(id int) IntId {
	return IntId{BasicId[int]{id}}
}
