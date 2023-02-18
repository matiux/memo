package sqlmemolist

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/matiux/memo/domain"
)

type MemoCreatedMemoListProjector struct {
	conn      *sql.DB
	tableName string
}

func (el *MemoCreatedMemoListProjector) Handle(message domain.Message) error {

	ctx := context.Background()

	stmt, err := el.conn.PrepareContext(ctx, "INSERT INTO memo(id, body, created_at) VALUES (?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	event := message.Payload.(*domain.MemoCreated)
	recOn, _ := json.Marshal(message.RecordedOn)

	if _, err = stmt.ExecContext(ctx, event.Id.Val, event.Body, string(recOn)); err != nil {
		return err
	}

	return nil
}

func (el *MemoCreatedMemoListProjector) Support(message domain.Message) bool {
	_, ok := message.Payload.(*domain.MemoCreated)

	return ok
}

func NewMemoCreatedMemoListProjector(conn *sql.DB) *MemoCreatedMemoListProjector {
	return &MemoCreatedMemoListProjector{conn, "memo"}
}
