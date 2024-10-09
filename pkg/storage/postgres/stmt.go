package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Stmt struct {
	stmt *sqlx.Stmt
}

func newStmt(stmt *sqlx.Stmt) *Stmt {
	return &Stmt{stmt: stmt}
}

func (s *Stmt) Close() error {
	err := s.stmt.Close()
	if err != nil {
		return errFailedToCloseStmt(err)
	}
	return nil
}

func (s *Stmt) Exec(args ...any) (sql.Result, error) {
	res, err := s.stmt.Exec(args...)
	if err != nil {
		return nil, errFailedToExecStmt(err)
	}
	return res, nil
}

func (s *Stmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errFailedToExecStmt(err)
	}
	return res, nil
}

func (s *Stmt) Query(args ...any) (*sqlx.Rows, error) {
	rows, err := s.stmt.Queryx(args...)
	if err != nil {
		return nil, errFailedToQueryStmt(err)
	}
	return rows, nil
}

func (s *Stmt) QueryContext(ctx context.Context, args ...any) (*sqlx.Rows, error) {
	rows, err := s.stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, errFailedToQueryStmt(err)
	}
	return rows, nil
}

func (s *Stmt) QueryRow(args ...any) *sqlx.Row {
	return s.stmt.QueryRowx(args...)
}

func (s *Stmt) QueryRowContext(ctx context.Context, args ...any) *sqlx.Row {
	return s.stmt.QueryRowxContext(ctx, args...)
}
