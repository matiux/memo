package infrastructure_test

import (
	"database/sql"
	"github.com/matiux/memo/application"
	"github.com/matiux/memo/domain/aggregate"
	"github.com/matiux/memo/infrastructure"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySqlEventStore_it_should_append_event_stream(t *testing.T) {

	tableName := "events"
	application.LoadEnv()
	dsn := os.Getenv("DATABASE_URL")
	db, _ := sql.Open("mysql", dsn)
	db.Exec("TRUNCATE TABLE " + tableName)

	eventStore := infrastructure.NewMySQLEventStore(db, tableName)

	var memoId = aggregate.NewUUIDv4()
	var body = "Vegetables are good"
	var creationDate = time.Now()

	memoCreatedDomainMessage := aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(1),
		EventType:   "MemoCreated",
		Payload:     aggregate.NewMemoCreated(memoId, body, creationDate),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	memoBodyUpdatedDomainMessage := aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(2),
		EventType:   "MemoBodyUpdated",
		Payload:     aggregate.NewMemoBodyUpdated(memoId, "Vegetables and fruits are good", time.Now()),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	eventStore.Append(memoId, aggregate.DomainEventStream{memoCreatedDomainMessage, memoBodyUpdatedDomainMessage})

	rows, _ := db.Query("SELECT * FROM " + tableName)
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	assert.Equal(t, 2, counter)
	db.Exec("TRUNCATE TABLE " + tableName)
}
