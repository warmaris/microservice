package exibillia

import (
	"context"
	"errors"
	"microservice/internal/database"
	"time"
)

var (
	errNoRows = errors.New("not found")
)

type mysqlStorage struct {
	conn database.Connection
}

func NewStorage(conn database.Connection) Storage {
	return &mysqlStorage{
		conn: conn,
	}
}

func (m *mysqlStorage) create(ctx context.Context, exibillia Exibillia) error {
	const query = `INSERT INTO exibillia (name, description, tags) VALUES (?, ?, ?)`

	return m.conn.Exec(ctx, query, exibillia.Name, exibillia.Description, tags(exibillia.Tags))
}

func (m *mysqlStorage) getByID(ctx context.Context, id uint64) (Exibillia, error) {
	const query = `SELECT id, name, description, tags, created_at, updated_at FROM exibillia WHERE id = ?`

	var res model
	err := m.conn.Get(ctx, &res, query, id)
	if errors.Is(err, database.ErrNoRows) {
		return Exibillia{}, errNoRows
	}

	return Exibillia{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Tags:        res.Tags,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt.Time,
	}, err
}

func (m *mysqlStorage) update(ctx context.Context, exibillia *Exibillia) error {
	const query = `UPDATE exibillia SET description = ?, tags = ?, updated_at = ? WHERE id = ?`

	updTime := time.Now()
	err := m.conn.Exec(ctx, query, exibillia.Description, tags(exibillia.Tags), updTime, exibillia.ID)
	if err != nil {
		return err
	}

	exibillia.UpdatedAt = updTime

	return nil
}

func (m *mysqlStorage) delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM exibillia WHERE id = ?`

	err := m.conn.Exec(ctx, query, id)
	if errors.Is(err, database.ErrNoAffectedRows) {
		return errNoRows
	}

	return err
}
