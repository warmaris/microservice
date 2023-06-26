package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx *sqlx.Tx
}

func (t *Tx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := t.tx.GetContext(ctx, dest, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRows
	}

	return err
}

func (t *Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

func (t *Tx) Exec(ctx context.Context, query string, args ...any) error {
	res, err := t.tx.ExecContext(ctx, query, args...)

	if num, _ := res.RowsAffected(); num == 0 {
		return ErrNoAffectedRows
	}

	return err
}

func (t *Tx) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	return nil, ErrNestedTxNotSupported
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
