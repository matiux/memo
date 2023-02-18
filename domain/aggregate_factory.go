package domain

type AggregateFactory interface {
	create(aggregateClass Root, domainEventStream EventStream) (Root, error)
}

type PublicConstructorAggregateFactory struct {
}

func (pc PublicConstructorAggregateFactory) create(aggregateClass Root, domainEventStream EventStream) (Root, error) {

	if err := aggregateClass.InitializeState(domainEventStream, aggregateClass); err != nil {
		return nil, err
	}

	return aggregateClass, nil
}
