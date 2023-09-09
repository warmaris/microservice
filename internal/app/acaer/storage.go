package acaer

import (
	"context"
	"database/sql"
	"fmt"
	"microservice/internal/database"
	"strings"
)

type mysqlStorage struct {
	conn database.Connection
}

func NewStorage(conn database.Connection) Storage {
	return &mysqlStorage{
		conn: conn,
	}
}

func (m *mysqlStorage) begin(ctx context.Context) (Storage, error) {
	tx, err := m.conn.Begin(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		return nil, err
	}

	return NewStorage(tx), nil
}

func (m *mysqlStorage) commit() error {
	return m.conn.Commit()
}

func (m *mysqlStorage) rollback() error {
	return m.conn.Rollback()
}

func (m *mysqlStorage) getVersions(ctx context.Context) ([]string, error) {
	const query = `SELECT version FROM acaer_versions`

	var res []string
	err := m.conn.Select(ctx, &res, query)

	return res, err
}

func (m *mysqlStorage) create(ctx context.Context, acaer Acaer) (uint64, error) {
	const query = `INSERT INTO acaer(name, version) VALUES (?, ?)`

	return m.conn.ExecWithID(ctx, query, acaer.Name, acaer.Version)
}

func (m *mysqlStorage) createLooncans(ctx context.Context, list []looncanDTO) error {
	const query = `INSERT INTO looncan (name, value, parent_id, parent_type) VALUES %s`

	placeholders := make([]string, len(list))
	args := make([]any, 0, 4*len(list))

	for i, l := range list {
		placeholders[i] = "(?, ?, ?, ?)"
		args = append(args, l.name, l.value, l.parentID, l.parentType)
	}

	return m.conn.Exec(ctx, fmt.Sprintf(query, strings.Join(placeholders, ", ")), args...)
}

func (m *mysqlStorage) createFromAggregate(ctx context.Context, acaer Acaer) error {
	const acaerInsert = `INSERT INTO acaer(name, version) VALUES (?, ?)`
	const looncansInsert = `INSERT INTO looncan (name, value, parent_id, parent_type) VALUES %s`

	tx, err := m.conn.Begin(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer func(tx *database.Tx) {
		_ = tx.Rollback()
	}(tx)

	id, err := tx.ExecWithID(ctx, acaerInsert, acaer.Name, acaer.Version)
	if err != nil {
		return err
	}

	placeholders := make([]string, len(acaer.getLooncans()))
	args := make([]any, 0, 4*len(acaer.getLooncans()))

	for i, l := range acaer.getLooncans() {
		placeholders[i] = "(?, ?, ?, ?)"
		args = append(args, l.Name, l.Value, id, parentAcaer)
	}

	err = tx.Exec(ctx, fmt.Sprintf(looncansInsert, strings.Join(placeholders, ", ")), args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
