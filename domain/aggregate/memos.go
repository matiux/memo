package aggregate

type Memos struct {
	eventSourcingRepository EventSourcingRepository
}

func (m Memos) Add(memo *Memo) error {
	if err := m.eventSourcingRepository.Save(memo); err != nil {
		return err
	}

	return nil
}

func (m Memos) ById(idMemo UUIDv4) (*Memo, error) {

	aggregate, err := m.eventSourcingRepository.Load(idMemo)

	if err != nil {
		return nil, err
	}

	memo := aggregate.(*Memo)

	return memo, nil
}

func (m Memos) Update(memo *Memo) error {
	return m.Add(memo)
}

func NewMemos(store EventStore, bus EventBus) Memos {
	return Memos{
		EventSourcingRepository{
			eventStore:       store,
			eventBus:         bus,
			aggregateClass:   &Memo{},
			aggregateFactory: &PublicConstructorAggregateFactory{},
		},
	}
}
