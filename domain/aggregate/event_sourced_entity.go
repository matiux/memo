package aggregate

type EventSourcedEntity interface {
	handleRecursively(event DomainEvent)
	registerAggregateRoot(aggregate Root)
}
