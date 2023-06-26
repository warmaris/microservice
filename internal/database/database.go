package database

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoRows               = errors.New("no rows in result")
	ErrNoAffectedRows       = errors.New("no rows affected")
	ErrNestedTxNotSupported = errors.New("nested transactions not supported")
)

type Connection interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...any) error
	Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
	Commit() error
	Rollback() error
}
