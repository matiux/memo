package aggregate

type EventSourcingRepository struct {
	EventStore
	EventBus
	aggregateClass Root
	AggregateFactory
}

func (esr *EventSourcingRepository) Save(aggregate Root) error {

	domainEventStream := aggregate.GetUncommittedEvents()
	esr.EventStore.Append(aggregate.getAggregateRootId(), domainEventStream)
	if err := esr.EventBus.publish(domainEventStream); err != nil {
		return err
	}

	return nil
}

func (esr *EventSourcingRepository) Load(id EntityId) (Root, error) {

	domainEventStream, err := esr.EventStore.Load(id)
	if err != nil {
		return nil, err
	}

	aggregateRoot, err := esr.AggregateFactory.create(esr.aggregateClass, domainEventStream)

	if err != nil {
		return nil, err
	}

	return aggregateRoot, nil
}
