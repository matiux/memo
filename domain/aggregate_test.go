package domain_test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/matiux/memo/application"
	"github.com/matiux/memo/domain"
	"github.com/matiux/memo/infrastructure"
	"os"
	"time"
)

var memoId = domain.NewUUIDv4()
var body = "Vegetables are good"
var creationDate = time.Now()

func createMemo() *domain.Memo {

	return domain.NewMemo(memoId, body, creationDate)
}

func createEvents() (domain.DomainMessage, domain.DomainMessage) {

	memoCreatedDomainMessage := domain.DomainMessage{
		Playhead:    domain.Playhead(1),
		EventType:   "MemoCreated",
		Payload:     domain.NewMemoCreated(memoId, body, creationDate),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	memoBodyUpdatedDomainMessage := domain.DomainMessage{
		Playhead:    domain.Playhead(2),
		EventType:   "MemoBodyUpdated",
		Payload:     domain.NewMemoBodyUpdated(memoId, "Vegetables and fruits are good", time.Now()),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	return memoCreatedDomainMessage, memoBodyUpdatedDomainMessage
}

func createTraceableEventBus() *domain.TraceableEventBus {
	eventBus := domain.NewTraceableEventBus(
		&domain.SimpleEventBus{
			EventListeners: nil,
			Queue:          nil,
			IsPublishing:   false,
		},
	)
	eventBus.Trace()

	return eventBus
}

func setupInMemoryEventSourcingRepository() (
	*domain.TraceableEventStore,
	*domain.TraceableEventBus,
	domain.EventSourcingRepository,
) {

	eventStore := domain.NewTraceableEventStore(domain.NewInMemoryEventStore())
	eventStore.Trace()

	eventBus := createTraceableEventBus()

	eventSourcingRepository := domain.NewEventSourcingRepository(
		eventStore,
		eventBus,
		&domain.PublicConstructorAggregateFactory{},
	)

	return eventStore, eventBus, eventSourcingRepository
}

func setupMySqlEventSourcingRepository() (
	*domain.TraceableEventStore,
	*domain.TraceableEventBus,
	domain.EventSourcingRepository,
) {

	application.LoadEnv()
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	eventStore := domain.NewTraceableEventStore(infrastructure.NewMySQLEventStore(db, "events"))
	eventStore.Trace()

	eventBus := createTraceableEventBus()

	eventSourcingRepository := domain.NewEventSourcingRepository(
		eventStore,
		eventBus,
		&domain.PublicConstructorAggregateFactory{},
	)

	return eventStore, eventBus, eventSourcingRepository
}
