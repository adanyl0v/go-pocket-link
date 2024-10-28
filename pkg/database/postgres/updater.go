package postgres

import "context"

func (db *DB) Update(ctx context.Context, query string, args ...any) error {
	_, err := db.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}

func (db *DB) UpdateNamed(ctx context.Context, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return errPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, arg)
	if err != nil {
		return errExecutingQuery(query, err)
	}
	return nil
}
