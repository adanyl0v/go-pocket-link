package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNoRowsInResultSet = sql.ErrNoRows
)

const (
	ErrCodeIntegrityConstraintViolation = "23000"
	ErrCodeRestrictViolation            = "23001"
	ErrCodeNotNullViolation             = "23502"
	ErrCodeForeignKeyViolation          = "23503"
	ErrCodeUniqueViolation              = "23505"
	ErrCodeCheckViolation               = "23514"
	ErrCodeExclusionViolation           = "23P01"
	ErrCodeNoDataFound                  = "P0002"
	ErrCodeInvalidSyntax                = "42601"
)

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func errConnecting(err error) error {
	return fmt.Errorf("%w (connecting to postgres)", err)
}

func errClosing(err error) error {
	return fmt.Errorf("%w (closing postgres connection)", err)
}

func errPreparingQuery(query string, err error) error {
	return fmt.Errorf("%w (preparing '%s')", err, query)
}

func errExecutingQuery(query string, err error) error {
	return fmt.Errorf("%w (executing '%s')", err, query)
}
