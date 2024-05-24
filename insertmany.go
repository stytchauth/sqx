package sqx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blockloop/scan/v2"
	sq "github.com/stytchauth/squirrel"
)

// InsertManyBuilder wraps squirrel.InsertBuilder and adds syntactic sugar for common usage patterns.
// This is a special case that includes a type constraint since generic methods are not supported in Go.
// In order to implement FromItems, we need a type constraint, but this is only possible if the struct itself
// is also generic. As such, the InsertManyBuilder is more constrained than InsertBuilder, but this is by
// design since *most* use cases should prefer the InsertBuilder unless they explicitly want to use FromItems.
type InsertManyBuilder[T any] struct {
	builder   sq.InsertBuilder
	queryable Queryable
	ctx       context.Context
	err       error
	logger    Logger
}

// ============================================
// BEGIN: squirrel-InsertBuilder parity section
// ============================================

// Columns adds insert columns to the query.
func (b InsertManyBuilder[T]) Columns(columns ...string) InsertManyBuilder[T] {
	return b.withBuilder(b.builder.Columns(columns...))
}

// Values adds a single row's values to the query.
func (b InsertManyBuilder[T]) Values(values ...any) InsertManyBuilder[T] {
	return b.withBuilder(b.builder.Values(values...))
}

// ==========================================
// END: squirrel-InsertBuilder parity section
// ==========================================

// FromItems generates an InsertManyBuilder from a slice of items. The first item in the slice is used to determine the
// columns for the insert statement. If excluded columns are provided, they will be removed from the list of columns.
// All items should be of the same type.
func (b InsertManyBuilder[T]) FromItems(items []T, excluded ...string) InsertManyBuilder[T] {
	if len(items) == 0 {
		return b
	}

	cols, err := scan.ColumnsStrict(&items[0], excluded...)
	if err != nil {
		return b.withError(err)
	}

	b = b.Columns(cols...)
	for _, item := range items {
		vals, err := scan.Values(cols, &item)
		if err != nil {
			return b.withError(err)
		}
		b = b.Values(vals...)
	}
	return b
}

// Do executes the InsertManyBuilder
func (b InsertManyBuilder[T]) Do() error {
	_, err := b.DoResult()
	return err
}

// DoResult executes the InsertManyBuilder and also returns the sql.Result for a successful query. This is useful if you
// wish to check the value of the LastInsertId() or RowsAffected() methods since Do() will discard this information.
func (b InsertManyBuilder[T]) DoResult() (sql.Result, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.queryable == nil {
		return nil, fmt.Errorf("missing queryable - call SetDefaultQueryable or WithQueryable to set it")
	}
	return b.builder.RunWith(runShim{b.queryable}).ExecContext(b.ctx)
}

// Debug prints the InsertManyBuilder state out to the provided logger
func (b InsertManyBuilder[T]) Debug() InsertManyBuilder[T] {
	debug(b.logger, b.builder)
	return b
}

// WithQueryable configures a Queryable for this InsertManyBuilder instance
func (b InsertManyBuilder[T]) WithQueryable(queryable Queryable) InsertManyBuilder[T] {
	return InsertManyBuilder[T]{builder: b.builder, queryable: queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}

// WithLogger configures a Queryable for this InsertManyBuilder instance
func (b InsertManyBuilder[T]) WithLogger(logger Logger) InsertManyBuilder[T] {
	return InsertManyBuilder[T]{builder: b.builder, queryable: b.queryable, logger: logger, ctx: b.ctx, err: b.err}
}

func (b InsertManyBuilder[T]) withError(err error) InsertManyBuilder[T] {
	if b.err != nil {
		return b
	}
	return InsertManyBuilder[T]{builder: b.builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: err}
}

func (b InsertManyBuilder[T]) withBuilder(builder sq.InsertBuilder) InsertManyBuilder[T] {
	return InsertManyBuilder[T]{builder: builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}
