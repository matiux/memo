package aggregate

type EventSourcingRepository struct {
	EventStore
	EventBus
	AggregateClass Root
	AggregateFactory
}

func (esr *EventSourcingRepository) Save(aggregate Root) error {

	domainEventStream := aggregate.GetUncommittedEvents()
	esr.EventStore.Append(aggregate.GetAggregateRootId(), domainEventStream)
	if err := esr.EventBus.Publish(domainEventStream); err != nil {
		return err
	}

	return nil
}

func (esr *EventSourcingRepository) Load(id EntityId) (Root, error) {

	domainEventStream, err := esr.EventStore.Load(id)
	if err != nil {
		return nil, err
	}

	aggregateRoot, err := esr.AggregateFactory.create(esr.AggregateClass, domainEventStream)

	if err != nil {
		return nil, err
	}

	return aggregateRoot, nil
}
