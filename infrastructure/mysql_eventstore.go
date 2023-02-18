package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/matiux/memo/domain"
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

func (e *MySQLEventStore) Append(id domain.EntityId, eventStream domain.EventStream) error {

	_ = id.(domain.UUIDv4).Val
	ctx := context.Background()

	tx, err := e.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, domainMessage := range eventStream {
		stmt, err := tx.PrepareContext(
			ctx, "INSERT INTO memo_db.events(uuid, playhead, payload, metadata, recorded_on, type) VALUES(?, ?, ?, ?, ?, ?)",
		)
		defer stmt.Close()
		if err != nil {
			tx.Rollback()
			return err
		}

		marshaledPayload, _ := json.Marshal(domainMessage.Payload)
		recOn, _ := json.Marshal(domainMessage.RecordedOn)

		if _, err = stmt.ExecContext(
			ctx,
			domainMessage.AggregateId.(domain.UUIDv4).Val,
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

	return tx.Commit()
}

func (e *MySQLEventStore) Load(id domain.EntityId) (domain.EventStream, error) {

	stringId := id.(domain.UUIDv4).Val

	var statement = e.prepareLoadStatement()
	defer statement.Close()
	rows, err := statement.Query(stringId, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventStream := domain.EventStream{}

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

func (e *MySQLEventStore) deserializeEvent(row eventRow) domain.Message {

	payload, err := domain.EventDeserializerRegistry(row.eventType, row.payload)
	if err != nil {
		panic(err)
	}

	t, _ := time.Parse(domain.EventDateFormat, row.recordedOn)

	domainMessage := domain.Message{
		Playhead:    domain.Playhead(row.playhead),
		EventType:   row.eventType,
		Payload:     *payload,
		AggregateId: domain.NewUUIDv4From(row.uuid),
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
