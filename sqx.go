package sqx

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/stytchauth/squirrel"
	"log"
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

func (rc typedRunCtx[T]) Select(columns ...string) SelectBuilder[T] {
	return SelectBuilder[T]{builder: sq.Select(columns...), queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

func (rc typedRunCtx[T]) FromSquirrelSelect(sel sq.SelectBuilder) SelectBuilder[T] {
	return SelectBuilder[T]{builder: sel, queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

func (rc runCtx) Update(table string) UpdateBuilder {
	return UpdateBuilder{builder: sq.Update(table), queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

func (rc runCtx) Insert(table string) InsertBuilder {
	return InsertBuilder{builder: sq.Insert(table), queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

func (rc runCtx) Delete(table string) DeleteBuilder {
	return DeleteBuilder{builder: sq.Delete(table), queryable: rc.queryable, logger: rc.logger, ctx: rc.ctx}
}

// runShim maps a Queryable to the squirrel.BaseRunner interface
// by patching unused Exec and Query methods
type runShim struct {
	Queryable
}

func (t runShim) Exec(_ string, _ ...interface{}) (sql.Result, error) {
	return nil, fmt.Errorf("exec is not implemented, please use ExecCtx")
}

func (t runShim) Query(_ string, _ ...interface{}) (*sql.Rows, error) {
	return nil, fmt.Errorf("query is not implemented, please use QueryCtx")
}

func debug(logger Logger, builder interface{ ToSql() (string, []any, error) }) {
	query, args, err := builder.ToSql()
	if logger != nil {
		logger.Printf("[DEBUG] %+v\n", map[string]any{"sql": query, "args": args, "error": err})
	} else {
		log.Printf("missing default logger in SQX")
		log.Printf("[DEBUG] %+v\n", map[string]any{"sql": query, "args": args, "error": err})
	}

}
