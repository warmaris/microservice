package database

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Pool struct {
	db *sqlx.DB
}

func NewConnection(dsn string) (*Pool, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Pool{db: db}, nil
}

func (p *Pool) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	// todo: add some metrics
	err := p.db.GetContext(ctx, dest, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRows
	}

	return err
}

func (p *Pool) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.db.SelectContext(ctx, dest, query, args...)
}

func (p *Pool) Exec(ctx context.Context, query string, args ...any) error {
	res, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if num, _ := res.RowsAffected(); num == 0 {
		return ErrNoAffectedRows
	}

	return nil
}

func (p *Pool) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := p.db.BeginTxx(ctx, opts)

	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

func (p *Pool) Commit() error {
	return nil
}

func (p *Pool) Rollback() error {
	return nil
}
