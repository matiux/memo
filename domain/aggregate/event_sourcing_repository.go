package aggregate

type EventSourcingRepository struct {
	eventStore       EventStore
	eventBus         EventBus
	aggregateFactory AggregateFactory
}

func (esr *EventSourcingRepository) Save(aggregate Root) error {

	domainEventStream := aggregate.GetUncommittedEvents()
	esr.eventStore.Append(aggregate.GetAggregateRootId(), domainEventStream)
	if err := esr.eventBus.Publish(domainEventStream); err != nil {
		return err
	}

	return nil
}

func (esr *EventSourcingRepository) Load(id EntityId, aggregate Root) (Root, error) {

	domainEventStream, err := esr.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	aggregateRoot, err := esr.aggregateFactory.create(aggregate, domainEventStream)

	if err != nil {
		return nil, err
	}

	return aggregateRoot, nil
}

func NewEventSourcingRepository(
	store EventStore,
	bus EventBus,
	aggregateFactory AggregateFactory,
) EventSourcingRepository {
	return EventSourcingRepository{
		eventStore:       store,
		eventBus:         bus,
		aggregateFactory: aggregateFactory,
	}
}
