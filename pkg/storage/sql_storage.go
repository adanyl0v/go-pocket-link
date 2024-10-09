package storage

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type SQLStorage interface {
	sqlCloser
	sqlPreparer
	sqlExecer
	sqlContextExecer
	sqlQuerier
	sqlContextQuerier
	sqlRowQuerier
	sqlContextRowQuerier
	sqlNamedExecer
	sqlNamedContextExecer
	sqlNamedQuerier
	sqlNamedContextQuerier
	DriverName() string
	Begin() (SQLTx, error)
	BeginTx(ctx context.Context, options *sql.TxOptions) (SQLTx, error)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
}

type SQLTx interface {
	sqlPreparer
	sqlExecer
	sqlContextExecer
	sqlQuerier
	sqlContextQuerier
	sqlRowQuerier
	sqlContextRowQuerier
	sqlNamedExecer
	sqlNamedContextExecer
	sqlNamedQuerier
	Commit() error
	Rollback() error
}

type SQLStmt interface {
	sqlCloser
	Exec(args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, args ...any) (sql.Result, error)
	Query(args ...any) (*sqlx.Rows, error)
	QueryContext(ctx context.Context, args ...any) (*sqlx.Rows, error)
	QueryRow(args ...any) *sqlx.Row
	QueryRowContext(ctx context.Context, args ...any) *sqlx.Row
}

type SQLNamedStmt interface {
	sqlCloser
	Exec(arg any) (sql.Result, error)
	ExecContext(ctx context.Context, arg any) (sql.Result, error)
	Query(arg any) (*sqlx.Rows, error)
	QueryContext(ctx context.Context, arg any) (*sqlx.Rows, error)
	QueryRow(arg any) *sqlx.Row
	QueryRowContext(ctx context.Context, arg any) *sqlx.Row
}

type sqlPreparer interface {
	Prepare(query string) (SQLStmt, error)
	PrepareContext(ctx context.Context, query string) (SQLStmt, error)
	PrepareNamed(query string) (SQLNamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (SQLNamedStmt, error)
}

type sqlCloser interface {
	Close() error
}

type sqlExecer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

type sqlNamedExecer interface {
	ExecNamed(query string, arg any) (sql.Result, error)
}

type sqlContextExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type sqlNamedContextExecer interface {
	ExecNamedContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

type sqlQuerier interface {
	Query(query string, args ...any) (*sqlx.Rows, error)
}

type sqlNamedQuerier interface {
	QueryNamed(query string, arg any) (*sqlx.Rows, error)
}

type sqlContextQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
}

type sqlNamedContextQuerier interface {
	QueryNamedContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error)
}

type sqlRowQuerier interface {
	QueryRow(query string, args ...any) *sqlx.Row
}

type sqlNamedRowQuerier interface {
	QueryNamedRow(query string, arg any) *sqlx.Row
}

type sqlContextRowQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sqlx.Row
}

type sqlNamedContextRowQuerier interface {
	QueryNamedRowContext(ctx context.Context, query string, arg any) *sqlx.Row
}
