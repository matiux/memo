package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"log"
	"os"
	"time"
)

type eventRow struct {
	uuid       string
	playhead   int
	payload    string
	recordedOn string
	eventType  string
}

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

	var statement = e.prepareLoadStatement()
	defer statement.Close()
	rows, err := statement.Query(stringId, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventStream := aggregate.DomainEventStream{}

	for rows.Next() {

		var event eventRow
		if err = rows.Scan(
			&event.uuid,
			&event.playhead,
			&event.payload,
			&event.recordedOn,
			&event.eventType,
		); err != nil {
			return nil, err
		}

		eventStream = append(eventStream, e.deserializeEvent(event))
	}

	return eventStream, nil
}

func (e *MySQLEventStore) deserializeEvent(row eventRow) aggregate.DomainMessage {

	payload, err := aggregate.EventDeserializerRegistry(row.eventType, row.payload)
	if err != nil {
		panic(err)
	}

	t, _ := time.Parse("2006-01-02\\T15:04:05.000000Z07:00", row.recordedOn)

	domainMessage := aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(row.playhead),
		EventType:   row.eventType,
		Payload:     *payload,
		AggregateId: aggregate.NewUUIDv4From(row.uuid),
		RecordedOn:  t,
	}

	return domainMessage
}

func (e *MySQLEventStore) prepareLoadStatement() *sql.Stmt {

	query := fmt.Sprintf(
		"SELECT uuid, playhead, payload, recorded_on, type "+
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
