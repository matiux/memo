package aggregate

import "reflect"

type AggregateFactory interface {
	create(aggregateClass reflect.Value, domainEventStream DomainEventStream) (Root, error)
}

type PublicConstructorAggregateFactory struct {
}

func (pc PublicConstructorAggregateFactory) create(aggregateClass reflect.Value, domainEventStream DomainEventStream) (Root, error) {

	inputs := make([]reflect.Value, 2)
	inputs[0] = reflect.ValueOf(domainEventStream)
	inputs[1] = aggregateClass.Elem().Addr().Convert(reflect.TypeOf((*Root)(nil)).Elem())

	aggregateClass.MethodByName("InitializeState").Call(inputs)

	return aggregateClass.Interface().(Root), nil

	// ------------- OK ------------------
	//if err := aggregateClass.InitializeState(domainEventStream, aggregateClass); err != nil {
	//	return nil, err
	//}
	//
	//return aggregateClass, nil
	// --------------- OK ------------------
}
