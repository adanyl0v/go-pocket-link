package postgres

import "context"

func (db *DB) Delete(ctx context.Context, query string, args ...any) error {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return errorPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errorExecutingQuery(query, err)
	}
	return nil
}

func (db *DB) DeleteNamed(ctx context.Context, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return errorPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, arg)
	if err != nil {
		return errorExecutingQuery(query, err)
	}
	return nil
}
