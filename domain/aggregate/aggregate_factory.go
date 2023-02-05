package aggregate

type AggregateFactory interface {
	create(aggregateClass Root, domainEventStream DomainEventStream) (Root, error)
}

type PublicConstructorAggregateFactory struct {
}

func (pc *PublicConstructorAggregateFactory) create(aggregateClass Root, domainEventStream DomainEventStream) (Root, error) {

	if err := aggregateClass.InitializeState(domainEventStream, aggregateClass); err != nil {
		return nil, err
	}

	return aggregateClass, nil
}
