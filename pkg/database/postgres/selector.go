package postgres

import "context"

func (db *DB) Select(ctx context.Context, dest any, query string, args ...any) error {
	err := db.db.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}

func (db *DB) SelectPrepared(ctx context.Context, dest any, query string, args ...any) error {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return errPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.SelectContext(ctx, dest, args...)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}
