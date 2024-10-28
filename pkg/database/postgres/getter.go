package postgres

import (
	"context"
)

func (db *DB) Get(ctx context.Context, dest any, query string, args ...any) error {
	err := db.db.GetContext(ctx, dest, query, args...)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}

func (db *DB) GetNamed(ctx context.Context, dest any, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return errPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.GetContext(ctx, dest, arg)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}

func (db *DB) GetPrepared(ctx context.Context, dest any, query string, args ...any) error {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return errPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.GetContext(ctx, dest, args...)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}
