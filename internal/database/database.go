// Package database is a wrapper for sqlx database connection.
// Here we can see a base interface of connection that is implemented by connection pool and transaction object.
// Feature packages implement own storages which are using the connection interface. So, we can switch between base
// connection and transaction dynamically.
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
	ExecWithID(ctx context.Context, query string, args ...any) (uint64, error)
	Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
	Commit() error
	Rollback() error
}
