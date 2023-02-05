package aggregate

type EventSourcingRepository struct {
	eventStore       EventStore
	eventBus         EventBus
	aggregateClass   Root
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

func (esr *EventSourcingRepository) Load(id EntityId) (Root, error) {

	domainEventStream, err := esr.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	aggregateRoot, err := esr.aggregateFactory.create(esr.aggregateClass, domainEventStream)

	if err != nil {
		return nil, err
	}

	return aggregateRoot, nil
}

func NewEventSourcingRepository(
	store EventStore,
	bus EventBus,
	aggregateClass Root,
	aggregateFactory AggregateFactory,
) EventSourcingRepository {
	return EventSourcingRepository{
		eventStore:       store,
		eventBus:         bus,
		aggregateClass:   aggregateClass,
		aggregateFactory: aggregateFactory,
	}
}
