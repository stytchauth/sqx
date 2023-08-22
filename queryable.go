package sqx

import (
	"context"
	"database/sql"
)

var defaultQueryable Queryable = nil

// SetDefaultQueryable sets the DB query handler that should be used to run requests.
// If you need to change the DB query handler for a specific request, use WithQueryable
func SetDefaultQueryable(queryable Queryable) {
	defaultQueryable = queryable
}

// Queryable is an interface wrapping common database access methods.
//
// This is useful in cases where it doesn't matter whether the database handle is the root handle
// (*sql.DB) or an already-open transaction (*sql.Tx).
type Queryable interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
