package infrastructure_test

import (
	"database/sql"
	"github.com/matiux/memo/application"
	"github.com/matiux/memo/domain"
	"github.com/matiux/memo/infrastructure"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func prepareConnection() (*sql.DB, string) {
	tableName := "events"
	application.LoadEnv()
	dsn := os.Getenv("DATABASE_URL")
	db, _ := sql.Open("mysql", dsn)

	return db, tableName
}

func TestMySqlEventStore_it_should_append_event_stream(t *testing.T) {

	db, tableName := prepareConnection()
	db.Exec("TRUNCATE TABLE " + tableName)

	eventStore := infrastructure.NewMySQLEventStore(db, tableName)

	var memoId = domain.NewUUIDv4()
	var body = "Vegetables are good"
	var creationDate = time.Now()

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

	eventStore.Append(memoId, domain.DomainEventStream{memoCreatedDomainMessage, memoBodyUpdatedDomainMessage})

	rows, _ := db.Query("SELECT * FROM " + tableName)
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	assert.Equal(t, 2, counter)
	db.Exec("TRUNCATE TABLE " + tableName)
}

func TestMySqlEventStore_it_should_load_event_stream(t *testing.T) {

	db, tableName := prepareConnection()
	db.Exec("TRUNCATE TABLE " + tableName)

	eventStore := infrastructure.NewMySQLEventStore(db, tableName)

	var memoId = domain.NewUUIDv4()
	var body = "Vegetables are good"
	var creationDate = time.Now()

	eventStore.Append(
		memoId,
		domain.DomainEventStream{
			domain.DomainMessage{
				Playhead:    domain.Playhead(1),
				EventType:   "MemoCreated",
				Payload:     domain.NewMemoCreated(memoId, body, creationDate),
				AggregateId: memoId,
				RecordedOn:  time.Now()},
			domain.DomainMessage{
				Playhead:    domain.Playhead(2),
				EventType:   "MemoBodyUpdated",
				Payload:     domain.NewMemoBodyUpdated(memoId, "Vegetables and fruits are good", time.Now()),
				AggregateId: memoId,
				RecordedOn:  time.Now(),
			},
		},
	)

	domainEventStream, err := eventStore.Load(memoId)
	assert.Nil(t, err)
	assert.Len(t, domainEventStream, 2)
}
