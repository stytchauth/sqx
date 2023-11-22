package sqx

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sq "github.com/stytchauth/squirrel"
)

type runCtx struct {
	logger    Logger
	queryable Queryable
	ctx       context.Context
}

// WithQueryable configures a Queryable for this ctx instance
func (rc runCtx) WithQueryable(queryable Queryable) runCtx {
	return runCtx{queryable: queryable, logger: rc.logger, ctx: rc.ctx}
}

// WithLogger configures a Logger for this ctx instance
func (rc runCtx) WithLogger(logger Logger) runCtx {
	return runCtx{queryable: rc.queryable, logger: logger, ctx: rc.ctx}
}

// typedRunCtx wraps a generic type + a runCtx, it can be used to create typed Select builders
type typedRunCtx[T any] struct {
	runCtx
}

// WithQueryable configures a Queryable for this ctx instance
func (rc typedRunCtx[T]) WithQueryable(queryable Queryable) typedRunCtx[T] {
	return typedRunCtx[T]{runCtx{queryable: queryable, logger: rc.logger, ctx: rc.ctx}}
}

// WithLogger configures a Logger for this ctx instance
func (rc typedRunCtx[T]) WithLogger(logger Logger) typedRunCtx[T] {
	return typedRunCtx[T]{runCtx{queryable: rc.queryable, logger: logger, ctx: rc.ctx}}
}

// Read is the entrypoint for creating generic Select builders
func Read[T any](ctx context.Context) typedRunCtx[T] {
	return typedRunCtx[T]{Write(ctx)}
}

// Write is the entrypoint for creating sql-extra builders that call ExecCtx
// and its variants - it does not have a generic b/c Exec cannot return arbitrary data
func Write(ctx context.Context) runCtx {
	return runCtx{
		ctx:       ctx,
		logger:    defaultLogger,
		queryable: defaultQueryable,
	}
}

// Select constructs a new SelectBuilder for the given columns for this typedRunCtx.
func (rc typedRunCtx[T]) Select(columns ...string) SelectBuilder[T] {
	builder := sq.Select(columns...).PlaceholderFormat(defaultPlaceholderFormat)
	return SelectBuilder[T]{builder: builder, queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

// Update constructs a new UpdateBuilder for the given table for this typedRunCtx.
func (rc runCtx) Update(table string) UpdateBuilder {
	builder := sq.Update(table).PlaceholderFormat(defaultPlaceholderFormat)
	return UpdateBuilder{builder: builder, queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

// Insert constructs a new InsertBuilder for the given table for this typedRunCtx.
func (rc runCtx) Insert(table string) InsertBuilder {
	builder := sq.Insert(table).PlaceholderFormat(defaultPlaceholderFormat)
	return InsertBuilder{builder: builder, queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

// Delete constructs a new DeleteBuilder for the given table for this typedRunCtx.
func (rc runCtx) Delete(table string) DeleteBuilder {
	builder := sq.Delete(table).PlaceholderFormat(defaultPlaceholderFormat)
	return DeleteBuilder{builder: builder, queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

// runShim maps a Queryable to the squirrel.BaseRunner interface
// by patching unused Exec and Query methods
type runShim struct {
	Queryable
}

// Exec is required by squirrel.BaseRunner so runShim patches it here. This should never be called internally since we
// instead use ExecContext.
func (t runShim) Exec(_ string, _ ...interface{}) (sql.Result, error) {
	return nil, fmt.Errorf("exec is not implemented, please use ExecCtx")
}

// Query is required by squirrel.BaseRunner so runShim patches it here. This should never be called internally since we
// instead use QueryContext.
func (t runShim) Query(_ string, _ ...interface{}) (*sql.Rows, error) {
	return nil, fmt.Errorf("query is not implemented, please use QueryCtx")
}

// debug prints the query and args to the logger if it is non-nil, or using log.Printf otherwise.
func debug(logger Logger, builder interface{ ToSql() (string, []any, error) }) {
	query, args, err := builder.ToSql()
	if logger != nil {
		logger.Printf("[DEBUG] %+v\n", map[string]any{"sql": query, "args": args, "error": err})
	} else {
		log.Printf("missing default logger in SQX")
		log.Printf("[DEBUG] %+v\n", map[string]any{"sql": query, "args": args, "error": err})
	}
}
