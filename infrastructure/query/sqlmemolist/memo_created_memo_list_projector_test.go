package sqlmemolist_test

import (
	"database/sql"
	"github.com/matiux/memo/application"
	"github.com/matiux/memo/domain"
	"github.com/matiux/memo/infrastructure/query/sqlmemolist"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func prepareReadmodelMemoListConnection() (*sql.DB, string) {
	tableName := "memo"
	application.LoadEnv()
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	return db, tableName
}

func TestMemoCreatedMemoListProjector_it_should_project_memo_created_for_memo_list_readmodel(t *testing.T) {

	db, tableName := prepareReadmodelMemoListConnection()
	db.Exec("TRUNCATE TABLE " + tableName)
	projector := sqlmemolist.NewMemoCreatedMemoListProjector(db)

	memoId := domain.NewUUIDv4From("1750c0c3-06b2-46cf-b140-b36cdc215474")
	body := "Vegetables are good"
	creationDate, _ := time.Parse(domain.EventDateFormat, "2023-02-15\\T10:19:52.642901+01:00")

	memoCreatedDomainMessage := domain.Message{
		Playhead:    domain.Playhead(1),
		EventType:   "MemoCreated",
		Payload:     domain.NewMemoCreated(memoId, body, creationDate),
		AggregateId: memoId,
		RecordedOn:  time.Now(),
	}

	err := projector.Handle(memoCreatedDomainMessage)

	assert.Nil(t, err)

	rows, _ := db.Query("SELECT * FROM " + tableName)
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	assert.Equal(t, 1, counter)
	db.Exec("TRUNCATE TABLE " + tableName)
}
