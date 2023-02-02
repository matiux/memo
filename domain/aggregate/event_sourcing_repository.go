package aggregate

import (
	"fmt"
	"reflect"
)

type EventSourcingRepository struct {
	EventStore
	EventBus
	aggregateClass reflect.Type
	AggregateFactory
}

func (esr *EventSourcingRepository) Save(aggregate Root) error {

	if reflect.TypeOf(aggregate) != esr.aggregateClass {
		return fmt.Errorf("aggregate type mismatch. Expected %v, but got %v", esr.aggregateClass, reflect.TypeOf(aggregate))
	}

	fmt.Println(aggregate)

	return nil
}
