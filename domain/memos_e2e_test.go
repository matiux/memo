package domain_test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/matiux/memo/application"
	"github.com/matiux/memo/domain"
	"github.com/matiux/memo/infrastructure"
	"github.com/matiux/memo/infrastructure/query/sqlmemolist"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMemos_it_should_add_new_memo_and_project_events(t *testing.T) {

	application.LoadEnv()
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.Exec("TRUNCATE TABLE memo_db.events;")
	db.Exec("TRUNCATE TABLE memo_db.memo;")

	eventStore := infrastructure.NewMySQLEventStore(db, "events")

	eventListeners := []domain.EventListener{sqlmemolist.NewMemoCreatedMemoListProjector(db)}

	memos := domain.NewMemos(eventStore, domain.NewSimpleEventBus(eventListeners))

	memoId := domain.NewUUIDv4()
	body := "Vegetables are good"

	memo := domain.NewMemo(memoId, body, time.Now())
	memos.Add(memo)

	myMemo, _ := memos.ById(memoId)
	myMemo.UpdateBody("Vegetables amd fruits are good", time.Now())
	memos.Update(myMemo)

	events, _ := db.Query("SELECT * FROM memo_db.events;")
	memoList, _ := db.Query("SELECT * FROM memo_db.memo;")

	var eventsCounter int
	for events.Next() {
		eventsCounter++
	}

	var memoCounter int
	for memoList.Next() {
		memoCounter++
	}

	assert.Equal(t, 2, eventsCounter)
	assert.Equal(t, 1, memoCounter)
}
