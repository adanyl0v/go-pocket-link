package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-pocket-link/pkg/storage"
	"time"
)

const (
	driverName = "postgres"
)

type Storage struct {
	db   *sqlx.DB
	opts *Options
}

func New(dsn string, opts *Options) (*Storage, error) {
	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, errFailedToConnectStorage(dsn, err)
	}
	if opts != nil {
		db.SetMaxOpenConns(opts.MaxOpenConns)
		db.SetMaxIdleConns(opts.MaxIdleConns)
		db.SetConnMaxLifetime(opts.ConnMaxLifetime)
		db.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
	} else {
		opts = &Options{}
	}
	return &Storage{
		db:   db,
		opts: opts,
	}, nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return errFailedToCloseStorage(err)
	}
	return nil
}

func (s *Storage) Prepare(query string) (storage.SQLStmt, error) {
	stmt, err := s.db.Preparex(query)
	if err != nil {
		return nil, errFailedToPrepareQuery(query, err)
	}
	return newStmt(stmt), nil
}

func (s *Storage) PrepareContext(ctx context.Context, query string) (storage.SQLStmt, error) {
	stmt, err := s.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, errFailedToPrepareQuery(query, err)
	}
	return newStmt(stmt), nil
}

func (s *Storage) PrepareNamed(query string) (storage.SQLNamedStmt, error) {
	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return nil, errFailedToPrepareQuery(query, err)
	}
	return newNamedStmt(stmt), nil
}

func (s *Storage) PrepareNamedContext(ctx context.Context, query string) (storage.SQLNamedStmt, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errFailedToPrepareQuery(query, err)
	}
	return newNamedStmt(stmt), nil
}

func (s *Storage) Exec(query string, args ...any) (sql.Result, error) {
	res, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return res, nil
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return res, nil
}

func (s *Storage) Query(query string, args ...any) (*sqlx.Rows, error) {
	rows, err := s.db.Queryx(query, args...)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return rows, nil
}

func (s *Storage) QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return rows, nil
}

func (s *Storage) QueryRow(query string, args ...any) *sqlx.Row {
	return s.db.QueryRowx(query, args...)
}

func (s *Storage) QueryRowContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	return s.db.QueryRowxContext(ctx, query, args...)
}

func (s *Storage) ExecNamed(query string, arg any) (sql.Result, error) {
	res, err := s.db.NamedExec(query, arg)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return res, nil
}

func (s *Storage) ExecNamedContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	res, err := s.db.NamedExecContext(ctx, query, arg)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return res, nil
}

func (s *Storage) QueryNamed(query string, arg any) (*sqlx.Rows, error) {
	rows, err := s.db.NamedQuery(query, arg)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return rows, nil
}

func (s *Storage) QueryNamedContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error) {
	rows, err := s.db.NamedQueryContext(ctx, query, arg)
	if err != nil {
		return nil, errFailedToExecQuery(query, err)
	}
	return rows, nil
}

func (s *Storage) DriverName() string {
	return driverName
}

func (s *Storage) Options() *Options {
	return s.opts
}

func (s *Storage) Begin() (storage.SQLTx, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errFailedToBeginTx(err)
	}
	return newTx(tx), nil
}

func (s *Storage) BeginTx(ctx context.Context, options *sql.TxOptions) (storage.SQLTx, error) {
	tx, err := s.db.BeginTxx(ctx, options)
	if err != nil {
		return nil, errFailedToBeginTx(err)
	}
	return newTx(tx), nil
}

func (s *Storage) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}
func (s *Storage) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}
func (s *Storage) SetConnMaxLifetime(d time.Duration) {
	s.db.SetConnMaxLifetime(d)
}
func (s *Storage) SetConnMaxIdleTime(d time.Duration) {
	s.db.SetConnMaxIdleTime(d)
}
