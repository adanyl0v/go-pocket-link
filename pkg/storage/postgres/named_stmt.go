package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type NamedStmt struct {
	stmt *sqlx.NamedStmt
}

func newNamedStmt(stmt *sqlx.NamedStmt) *NamedStmt {
	return &NamedStmt{stmt: stmt}
}

func (s *NamedStmt) Close() error {
	err := s.stmt.Close()
	if err != nil {
		return errFailedToCloseStmt(err)
	}
	return nil
}

func (s *NamedStmt) Exec(arg any) (sql.Result, error) {
	res, err := s.stmt.Exec(arg)
	if err != nil {
		return nil, errFailedToExecStmt(err)
	}
	return res, nil
}

func (s *NamedStmt) ExecContext(ctx context.Context, arg any) (sql.Result, error) {
	res, err := s.stmt.ExecContext(ctx, arg)
	if err != nil {
		return nil, errFailedToExecStmt(err)
	}
	return res, nil
}

func (s *NamedStmt) Query(arg any) (*sqlx.Rows, error) {
	rows, err := s.stmt.Queryx(arg)
	if err != nil {
		return nil, errFailedToQueryStmt(err)
	}
	return rows, nil
}

func (s *NamedStmt) QueryContext(ctx context.Context, arg any) (*sqlx.Rows, error) {
	rows, err := s.stmt.QueryxContext(ctx, arg)
	if err != nil {
		return nil, errFailedToQueryStmt(err)
	}
	return rows, nil
}

func (s *NamedStmt) QueryRow(arg any) *sqlx.Row {
	return s.stmt.QueryRowx(arg)
}

func (s *NamedStmt) QueryRowContext(ctx context.Context, arg any) *sqlx.Row {
	return s.stmt.QueryRowxContext(ctx, arg)
}
