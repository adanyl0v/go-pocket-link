package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

const driverName = "pgx"

type DB struct {
	db *sqlx.DB
}

type ConnOptions struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Connect(dataSourceName string, opts *ConnOptions) (*DB, error) {
	sqlxDB, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, errConnecting(err)
	}
	if opts != nil {
		sqlxDB.SetMaxOpenConns(opts.MaxOpenConns)
		sqlxDB.SetMaxIdleConns(opts.MaxIdleConns)
		sqlxDB.SetConnMaxLifetime(opts.ConnMaxLifetime)
		sqlxDB.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
	}
	return &DB{db: sqlxDB}, nil
}

func DriverName() string {
	return driverName
}

func (db *DB) Close() error {
	if err := db.db.Close(); err != nil {
		return errClosing(err)
	}
	return nil
}
