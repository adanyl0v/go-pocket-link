package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-pocket-link/pkg/storage"
)

type Tx struct {
	tx *sqlx.Tx
}

func newTx(tx *sqlx.Tx) *Tx {
	return &Tx{tx: tx}
}

func (t *Tx) Prepare(query string) (storage.SQLStmt, error) {
	stmt, err := t.tx.Preparex(query)
	if err != nil {
		return nil, errFailedToPrepareTx(query, err)
	}
	return newStmt(stmt), nil
}

func (t *Tx) PrepareContext(ctx context.Context, query string) (storage.SQLStmt, error) {
	stmt, err := t.tx.PreparexContext(ctx, query)
	if err != nil {
		return nil, errFailedToPrepareTx(query, err)
	}
	return newStmt(stmt), nil
}

func (t *Tx) PrepareNamed(query string) (storage.SQLNamedStmt, error) {
	stmt, err := t.tx.PrepareNamed(query)
	if err != nil {
		return nil, errFailedToPrepareTx(query, err)
	}
	return newNamedStmt(stmt), nil
}

func (t *Tx) PrepareNamedContext(ctx context.Context, query string) (storage.SQLNamedStmt, error) {
	stmt, err := t.tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errFailedToPrepareTx(query, err)
	}
	return newNamedStmt(stmt), nil
}

func (t *Tx) Exec(query string, args ...any) (sql.Result, error) {
	res, err := t.tx.Exec(query, args...)
	if err != nil {
		return nil, errFailedToExecTx(query, err)
	}
	return res, nil
}

func (t *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errFailedToExecTx(query, err)
	}
	return res, nil
}

func (t *Tx) Query(query string, args ...any) (*sqlx.Rows, error) {
	rows, err := t.tx.Queryx(query, args...)
	if err != nil {
		return nil, errFailedToQueryTx(query, err)
	}
	return rows, nil
}

func (t *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	rows, err := t.tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errFailedToQueryTx(query, err)
	}
	return rows, nil
}

func (t *Tx) QueryRow(query string, args ...any) *sqlx.Row {
	return t.tx.QueryRowx(query, args...)
}

func (t *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	return t.tx.QueryRowxContext(ctx, query, args...)
}

func (t *Tx) ExecNamed(query string, arg any) (sql.Result, error) {
	res, err := t.tx.NamedExec(query, arg)
	if err != nil {
		return nil, errFailedToExecTx(query, err)
	}
	return res, nil
}

func (t *Tx) ExecNamedContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	res, err := t.tx.NamedExecContext(ctx, query, arg)
	if err != nil {
		return nil, errFailedToExecTx(query, err)
	}
	return res, nil
}

func (t *Tx) QueryNamed(query string, arg any) (*sqlx.Rows, error) {
	rows, err := t.tx.NamedQuery(query, arg)
	if err != nil {
		return nil, errFailedToQueryTx(query, err)
	}
	return rows, nil
}

func (t *Tx) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		return errFailedToCommitTx(err)
	}
	return nil
}

func (t *Tx) Rollback() error {
	err := t.tx.Rollback()
	if err != nil {
		return errFailedToRollbackTx(err)
	}
	return nil
}
