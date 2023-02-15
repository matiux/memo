package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"log"
	"os"
)

type MySQLEventStore struct {
	conn      *sql.DB
	tableName string
}

func (e *MySQLEventStore) Append(id aggregate.EntityId, eventStream aggregate.DomainEventStream) error {

	_ = id.(aggregate.UUIDv4).Val

	ctx := context.Background()
	tx, err := e.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, domainMessage := range eventStream {
		stmt, err := tx.Prepare("INSERT INTO memo_db.events(uuid, playhead, payload, metadata, recorded_on, type) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		marshaledPayload, _ := json.Marshal(domainMessage.Payload)
		recOn, _ := json.Marshal(domainMessage.RecordedOn)

		if _, err = stmt.ExecContext(
			ctx,
			domainMessage.AggregateId.(aggregate.UUIDv4).Val,
			int(domainMessage.Playhead),
			string(marshaledPayload),
			"",
			string(recOn),
			domainMessage.EventType,
		); err != nil {
			tx.Rollback()
			return err
		}

	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (e *MySQLEventStore) Load(id aggregate.EntityId) (aggregate.DomainEventStream, error) {

	stringId := id.(aggregate.UUIDv4).Val

	var payload string

	statement := e.prepareLoadStatement()
	defer statement.Close()
	err := statement.QueryRow(stringId, 0).Scan(payload)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return nil, nil
}

func (e *MySQLEventStore) prepareLoadStatement() *sql.Stmt {

	query := fmt.Sprintf(
		"SELECT uuid, playhead, metadata, payload, recorded_on "+
			"FROM %v "+
			"WHERE uuid = ? "+
			"AND playhead >= ? "+
			"ORDER BY playhead ASC;",
		e.tableName,
	)
	statement, err := e.conn.Prepare(query)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return statement
}

func NewMySQLEventStore(conn *sql.DB, tableName string) *MySQLEventStore {
	return &MySQLEventStore{conn, tableName}
}
