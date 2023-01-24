package aggregate

import (
	"reflect"
)

type Root interface {
	getAggregateRootId() string
	Apply(event Event) (err error)
}

type EventSourcedAggregateRoot struct {
	UncommittedEvents []Event
	Playhead          int64
}

func (e *EventSourcedAggregateRoot) apply(event Event, aggregate Root) {
	//tv := reflect.TypeOf(event)
	//fmt.Printf("\n-----\n%v\n-----\n", tv.Name())

	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(event)

	reflect.ValueOf(aggregate).MethodByName("ApplyMemoCreated").Call(inputs)

	e.UncommittedEvents = append(e.UncommittedEvents, event)
}
