package sqx

import (
	"context"
	"database/sql"
)

// Queryable is an interface wrapping common database access methods.
//
// This is useful in cases where it doesn't matter whether the database handle is the root handle
// (db.DBConnector) or an already-open transaction (*sql.Tx).
//
// This does not contain BeginTx because we don't (yet) support nested transactions or savepoints.
// Implementing that is also non-trivial. See Go issue 7898 for discussion.
type Queryable interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
