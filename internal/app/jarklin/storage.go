package jarklin

import (
	"context"
	"database/sql"
	"encoding/json"
	"microservice/internal/database"
)

type mysqlStorage struct {
	conn database.Connection
}

func NewStorage(conn database.Connection) Storage {
	return &mysqlStorage{
		conn: conn,
	}
}

func (m *mysqlStorage) getNewEvents(ctx context.Context) ([]NotifyRequest, error) {
	const query = `SELECT payload FROM jarklin_events WHERE status = ? ORDER BY created_at ASC LIMIT 100`

	var rows []json.RawMessage
	err := m.conn.Select(ctx, &rows, query, StatusPrepared)
	if err != nil {
		return nil, err
	}

	res := make([]NotifyRequest, 0, len(rows))
	for _, row := range rows {
		event := NotifyRequest{}
		err = json.Unmarshal(row, &event)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}

	return res, nil
}

func (m *mysqlStorage) save(ctx context.Context, entity *Jarklin) error {
	const queryEntity = `INSERT INTO jarklin(name, created_at) VALUES (?, ?)`
	const queryEvent = `INSERT INTO jarklin_events(message_id, entity_id, payload, created_at, status) VALUES (?, ?, ?, ?, ?)`

	tx, err := m.conn.Begin(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	id, err := tx.ExecWithID(ctx, queryEntity, entity.Name, entity.CreatedAt)
	if err != nil {
		return err
	}

	entity.ID = id
	payload, err := json.Marshal(entity.GetEvent())
	if err != nil {
		return err
	}

	err = tx.Exec(ctx, queryEvent, 
		entity.GetEvent().MessageID,
		entity.GetEvent().EntityID,
		payload,
		entity.GetEvent().CreatedAt,
		StatusPrepared,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *mysqlStorage) updateStatus(ctx context.Context, msgID string, status NotifyStatus, statusInfo string) error {
	const query = `UPDATE jarklin_events SET status = ?, status_info = ? WHERE message_id = ?`

	return m.conn.Exec(ctx, query, status, statusInfo, msgID)
}
