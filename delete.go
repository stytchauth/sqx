package sqx

import (
	"context"
	"fmt"

	sq "github.com/stytchauth/squirrel"
)

type DeleteBuilder struct {
	builder   sq.DeleteBuilder
	queryable Queryable
	ctx       context.Context
	err       error
	logger    Logger
}

// ============================================
// BEGIN: squirrel-DeleteBuilder parity section
// ============================================

// Prefix adds an expression to the beginning of the query
func (b DeleteBuilder) Prefix(sql string, args ...interface{}) DeleteBuilder {
	return b.withBuilder(b.builder.Prefix(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b DeleteBuilder) PrefixExpr(expr Sqlizer) DeleteBuilder {
	return b.withBuilder(b.builder.PrefixExpr(expr))
}

// From sets the table to be deleted from.
func (b DeleteBuilder) From(from string) DeleteBuilder {
	return b.withBuilder(b.builder.From(from))
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b DeleteBuilder) Where(pred interface{}, rest ...interface{}) DeleteBuilder {
	return b.withBuilder(b.builder.Where(pred, rest...))
}

// OrderBy adds ORDER BY expressions to the query.
func (b DeleteBuilder) OrderBy(orderBys ...string) DeleteBuilder {
	return b.withBuilder(b.builder.OrderBy(orderBys...))
}

// Limit sets a LIMIT clause on the query.
func (b DeleteBuilder) Limit(limit uint64) DeleteBuilder {
	return b.withBuilder(b.builder.Limit(limit))
}

// Offset sets a OFFSET clause on the query.
func (b DeleteBuilder) Offset(offset uint64) DeleteBuilder {
	return b.withBuilder(b.builder.Offset(offset))
}

// Suffix adds an expression to the end of the query
func (b DeleteBuilder) Suffix(sql string, args ...interface{}) DeleteBuilder {
	return b.withBuilder(b.builder.Suffix(sql, args...))
}

// ==========================================
// END: squirrel-UpdateBuilder parity section
// ==========================================

// Do executes the DeleteBuilder
func (b DeleteBuilder) Do() error {
	if b.err != nil {
		return b.err
	}
	if b.queryable == nil {
		return fmt.Errorf("missing queryable - call SetDefaultQueryable or WithQueryable to set it")
	}
	_, err := b.builder.RunWith(runShim{b.queryable}).ExecContext(b.ctx)
	return err
}

// Debug prints the DeleteBuilder state out to the provided logger
func (b DeleteBuilder) Debug() DeleteBuilder {
	debug(b.logger, b.builder)
	return b
}

// WithQueryable configures a Queryable for this DeleteBuilder instance
func (b DeleteBuilder) WithQueryable(queryable Queryable) DeleteBuilder {
	return DeleteBuilder{builder: b.builder, queryable: queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}

// WithLogger configures a Queryable for this DeleteBuilder instance
func (b DeleteBuilder) WithLogger(logger Logger) DeleteBuilder {
	return DeleteBuilder{builder: b.builder, queryable: b.queryable, logger: logger, ctx: b.ctx, err: b.err}
}

func (b DeleteBuilder) withBuilder(builder sq.DeleteBuilder) DeleteBuilder {
	return DeleteBuilder{builder: builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}
