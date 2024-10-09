package postgres

import "go-pocket-link/pkg/errb"

var b = errb.Default()

func errFailedToConnectStorage(dsn string, err error) error {
	return b.Errorf("failed to connect to '%s': %v", dsn, err)
}

func errFailedToCloseStorage(err error) error {
	return b.Errorf("failed to close storage: %v", err)
}

func errFailedToPrepareQuery(query string, err error) error {
	return b.Errorf("failed to prepare '%s': %v", query, err)
}

func errFailedToExecQuery(query string, err error) error {
	return b.Errorf("failed to execute '%s': %v", query, err)
}

func errFailedToBeginTx(err error) error {
	return b.Errorf("failed to begin transaction: %v", err)
}

func errFailedToPrepareTx(query string, err error) error {
	return b.Errorf("failed to prepare transaction '%s': %v", query, err)
}

func errFailedToExecTx(query string, err error) error {
	return b.Errorf("failed to execute transaction '%s': %v", query, err)
}

func errFailedToQueryTx(query string, err error) error {
	return b.Errorf("failed to query transaction '%s': %v", query, err)
}

func errFailedToCommitTx(err error) error {
	return b.Errorf("failed to commit transaction: %v", err)
}

func errFailedToRollbackTx(err error) error {
	return b.Errorf("failed to rollback transaction: %v", err)
}

func errFailedToCloseStmt(err error) error {
	return b.Errorf("failed to close statement: %v", err)
}

func errFailedToExecStmt(err error) error {
	return b.Errorf("failed to execute statement: %v", err)
}

func errFailedToQueryStmt(err error) error {
	return b.Errorf("failed to query statement: %v", err)
}
