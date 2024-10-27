package postgres

import (
	"context"
)

func (db *DB) Save(ctx context.Context, dest any, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return errorPreparingQuery(query, err)
	}
	defer func() { _ = stmt.Close() }()

	if dest == nil {
		_, err = stmt.ExecContext(ctx, arg)
	} else {
		row := stmt.QueryRowxContext(ctx, arg)
		err = row.Scan(dest)
	}
	if err != nil {
		return errorExecutingQuery(query, err)
	}
	return nil
}
