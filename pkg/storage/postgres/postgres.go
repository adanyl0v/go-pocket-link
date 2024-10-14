package postgres

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go-pocket-link/pkg/errb"
	"time"
)

var b = errb.Default()

const driverName = "pgx"

type Options struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type DB struct {
	db *sqlx.DB
}

func Connect(dsn string, opts *Options) (*DB, error) {
	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, b.Errorf("connect to %s: %v", dsn, err)
	}
	if opts != nil {
		db.SetMaxOpenConns(opts.MaxOpenConns)
		db.SetMaxIdleConns(opts.MaxIdleConns)
		db.SetConnMaxLifetime(opts.ConnMaxLifetime)
		db.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
	}
	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	err := db.db.Close()
	if err != nil {
		return b.Errorf("close db connection: %v", err)
	}
	return nil
}

func (db *DB) DriverName() string {
	return driverName
}

func (db *DB) Save(ctx context.Context, dest any, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return b.Errorf("prepare `%s`: %v", query, err)
	}
	defer func() { _ = stmt.Close() }()

	row := stmt.QueryRowxContext(ctx, arg)
	if dest == nil {
		err = row.Scan()
	} else {
		err = row.Scan(dest)
	}
	if err != nil {
		return b.Errorf("exec `%s`: %v", query, err)
	}
	return nil
}

func (db *DB) Get(ctx context.Context, dest any, query string, args ...any) error {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return b.Errorf("prepare `%s`: %v", query, err)
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.GetContext(ctx, dest, args...)
	if err != nil {
		return b.Errorf("exec `%s`: %v", query, err)
	}
	return nil
}

func (db *DB) GetNamed(ctx context.Context, dest any, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return b.Errorf("prepare `%s`: %v", query, err)
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.GetContext(ctx, dest, arg)
	if err != nil {
		return b.Errorf("exec `%s`: %v", query, err)
	}
	return nil
}

func (db *DB) GetAll(ctx context.Context, dest any, query string, args ...any) error {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return b.Errorf("prepare `%s`: %v", query, err)
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.SelectContext(ctx, dest, args...)
	if err != nil {
		return b.Errorf("exec `%s`: %v", query, err)
	}
	return nil
}

func (db *DB) Update(ctx context.Context, query string, arg any) error {
	stmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return b.Errorf("prepare `%s`: %v", query, err)
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, arg)
	if err != nil {
		return b.Errorf("exec `%s`: %v", query, err)
	}
	return nil
}

func (db *DB) Delete(ctx context.Context, query string, args ...any) error {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return b.Errorf("prepare `%s`: %v", query, err)
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return b.Errorf("exec `%s`: %v", query, err)
	}
	return nil
}
